package blink

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/epkgs/mini-blink/internal/log"
	"github.com/epkgs/mini-blink/internal/utils"
)

const (
	JS_MB               = "__mb"
	JS_IPC              = "ipc"
	JS_JS2GO            = "__js2go"
	JS_GO2JS            = "__go2js"
	JS_REGISTER_HANDLER = "__register_handler"
)

type IPCHandler func(args ...any) any

type IPC struct {
	mb *Blink

	handlers map[string]IPCHandler

	resultWaiting map[string]chan any
}

type IPCMessage struct {
	ID      string        `json:"id"`      // 消息 ID
	ReplyId string        `json:"replyId"` // 回复ID
	Channel string        `json:"channel"` // 通道
	Args    []interface{} `json:"args"`    // 参数
	Error   string        `json:"error"`   // 是否错误，当有回复ID时，此字段有效
	Result  interface{}   `json:"result"`  // 返回值，当有回复ID时，此字段有效
}

func newIPC(mb *Blink) *IPC {
	ipc := &IPC{
		mb: mb,

		handlers:      make(map[string]IPCHandler),
		resultWaiting: make(map[string]chan any),
	}

	ipc.registerBootScript()
	ipc.registerJS2GO()
	ipc.registerJSHandler()

	return ipc
}

// GO 调用handler
//
//	一、GO 调用 GO handler，直接调用并返回
//
//	二、GO 调用 JS handler, 和 GO 调用 GO 流程一样，唯一区别是在 `invokeJS` 里调用 `ipc.Invoke` 执行的 `handler` 是转化后的 `JS handler`
func (ipc *IPC) Invoke(channel string, args ...any) (any, error) {
	handler, exist := ipc.handlers[channel]
	if !exist {
		msg := fmt.Sprintf("ipc channel %s not exist", channel)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	result := handler(args...)

	if err, ok := result.(error); ok {
		return nil, err
	}

	return result, nil
}

func (ipc *IPC) Sent(channel string, args ...any) error {
	handler, exist := ipc.handlers[channel]
	if !exist {
		msg := fmt.Sprintf("ipc channel %s not exist", channel)
		log.Error(msg)
		return errors.New(msg)
	}

	result := handler(args...)
	if err, ok := result.(error); ok {
		return err
	}

	return nil
}

// GO 注册 Handler
func (ipc *IPC) Handle(channel string, handler IPCHandler) {
	ipc.handlers[channel] = handler
}

func (ipc *IPC) HasChannel(channel string) (exist bool) {
	_, exist = ipc.handlers[channel]
	return
}

//go:embed internal/scripts/ipc.js
var bootjs []byte

func (ipc *IPC) registerBootScript() {
	script := fmt.Sprintf(
		string(bootjs),
		JS_MB,
		JS_IPC,
		JS_JS2GO,
		JS_GO2JS,
		JS_REGISTER_HANDLER,
	)

	ipc.mb.AddBootScript(script)
}

// JS -> GO 的消息分派、处理
func (ipc *IPC) registerJS2GO() {
	ipc.mb.js.bindFunction(JS_JS2GO, 1, func(es JsExecState) {
		arg := ipc.mb.js.Arg(es, 0)
		txt := ipc.mb.js.ToString(es, arg)

		log.Info("GO -> JS: %s", txt)

		var msg IPCMessage
		if err := json.Unmarshal(([]byte)(txt), &msg); err != nil {
			log.Error("JS -> GO, JSON 解析出错(%s): %s", err.Error(), txt)
			return
		}

		if msg.ReplyId != "" {
			ipc.handleJSReply(&msg)
			return
		}

		if msg.Channel != "" {
			view := ipc.mb.GetViewByJsExecState(es)
			ipc.invokeByJS(view, &msg)
			return
		}
	})
}

// JS 调用 handler
func (ipc *IPC) invokeByJS(view *View, msg *IPCMessage) {

	// 如果 ID 为空，则无须回复返回值
	if msg.ID == "" {
		ipc.Invoke(msg.Channel, msg.Args...)
		return
	}

	// 调用 invoke 获取到结果
	result, err := ipc.Invoke(msg.Channel, msg.Args...)

	e := ""
	if err != nil {
		e = err.Error()
		result = nil
	}

	replyMsg := IPCMessage{
		ID:      utils.RandString(8),
		ReplyId: msg.ID,
		Error:   e,
		Result:  result,
	}

	sentMsgToView(view, replyMsg)
}

func (ipc *IPC) handleJSReply(msg *IPCMessage) {
	if msg.ReplyId == "" {
		return
	}

	resultChan, exist := ipc.resultWaiting[msg.ReplyId]
	if !exist {
		return
	}

	if msg.Error != "" {
		resultChan <- errors.New(msg.Error)
	} else {
		resultChan <- msg.Result
	}
}

// JS 注册 handler 埋点
func (ipc *IPC) registerJSHandler() {
	// 注册 JS handler
	ipc.mb.js.bindFunction(JS_REGISTER_HANDLER, 1, func(es JsExecState) {
		arg := ipc.mb.js.Arg(es, 0)
		channel := ipc.mb.js.ToString(es, arg)

		view := ipc.mb.GetViewByJsExecState(es)

		// 将 JS handler 转为 GO handler
		ipc.handlers[channel] = func(args ...any) any {
			id, resultChan := ipc.handleJSChannel(view, channel, args...)
			defer close(resultChan)             // 关闭 result
			defer delete(ipc.resultWaiting, id) // 接收到消息就从 map 中删除

			select {
			case result := <-resultChan:
				return result
			case <-time.After(10 * time.Second): // 10秒等待超时
				return errors.New("等待 IPC JS Handler 处理结果超时")
			}
		}
	})
}

func (ipc *IPC) handleJSChannel(view *View, channel string, args ...any) (id string, result chan any) {
	id = utils.RandString(8) // 生成key

	msg := IPCMessage{
		ID:      id, // 关键是 ID 设置为空
		Channel: channel,
		Args:    args,
	}

	result = make(chan any, 1) // result 管道

	ipc.resultWaiting[id] = result // 暂存 result channel, 等待 JS 完毕后，通过 JS_HANDLE_PROCESS_REPLY 将结果塞进来

	sentMsgToView(view, msg)

	return id, result
}

func (ipc *IPC) RunJSFunc(view *View, funcName string, args ...any) chan any {

	newArgs := make([]any, 0, len(args)+1)
	newArgs = append(newArgs, funcName)
	newArgs = append(newArgs, args...)

	_, result := ipc.handleJSChannel(view, "runFunc", newArgs...)
	return result
}

func sentMsgToView(view *View, msg IPCMessage) {

	msgTxt, _ := json.Marshal(msg)

	script := fmt.Sprintf(`window.top['%s'](%q)`, JS_GO2JS, msgTxt)

	log.Info("GO -> JS: %s", msgTxt)

	view.RunJS(script)
}
