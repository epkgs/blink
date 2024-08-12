package blink

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/epkgs/blink/internal/cast"
	"github.com/epkgs/blink/internal/log"
	"github.com/epkgs/blink/internal/utils"
	"github.com/epkgs/blink/pkg/async"
)

const (
	JS_MB               = "__mb"
	JS_IPC              = "ipc"
	JS_JS2GO            = "__js2go"
	JS_GO2JS            = "__go2js"
	JS_REGISTER_HANDLER = "__register_handler"
)

type Callback interface{}

type resultCallback func(result interface{}, err error)

// resultCallback 用于区分无须返回值的情况
type ipcHandler func(cb resultCallback, args ...interface{})

type IPC struct {
	mb *Blink

	handlers map[string]ipcHandler

	resultWaiting map[string]resultCallback
}

type IPCMessage struct {
	ID      string        `json:"id"`      // 消息 ID
	ReplyId string        `json:"replyId"` // 回复ID
	Channel string        `json:"channel"` // 通道
	Args    []interface{} `json:"args"`    // 参数
	Result  interface{}   `json:"result"`  // 返回值，当有回复ID时，此字段有效
	Error   string        `json:"error"`   // 是否错误，当有回复ID时，此字段有效
}

func newIPC(mb *Blink) *IPC {
	ipc := &IPC{
		mb: mb,

		handlers:      make(map[string]ipcHandler),
		resultWaiting: make(map[string]resultCallback),
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
func (ipc *IPC) Invoke(channel string, args ...interface{}) (interface{}, error) {
	handler, exist := ipc.handlers[channel]
	if !exist {
		msg := fmt.Sprintf("ipc channel %s not exist", channel)
		log.Error(msg)
		return nil, errors.New(msg)
	}

	result := make(chan interface{})
	err := make(chan error)

	// 将 callback 转 chan
	handler(func(res interface{}, e error) {
		result <- res
		err <- e
	}, args...)

	return <-result, <-err
}

func (ipc *IPC) Sent(channel string, args ...interface{}) error {
	handler, exist := ipc.handlers[channel]
	if !exist {
		msg := fmt.Sprintf("ipc channel %s not exist", channel)
		log.Error(msg)
		return errors.New(msg)
	}

	handler(nil, args...)

	return nil
}

// GO 注册 Handler
//
// handler 必须为函数，如有返回值，第一个返回值为正常返回值（可省略），第二个返回值为错误（可省略）
func (ipc *IPC) Handle(channel string, handler Callback) {

	// 使用反射获取处理函数的类型
	handlerVal := reflect.ValueOf(handler)
	if handlerVal.Kind() != reflect.Func {
		panic(fmt.Sprintf("channel %s, handler must be a function", channel))
	}

	handlerType := handlerVal.Type()

	ipc.handlers[channel] = func(cb resultCallback, inputs ...interface{}) {

		inputSize := len(inputs)

		// 构造参数列表
		pCount := handlerType.NumIn()
		isVariadic := handlerType.IsVariadic()
		if isVariadic {
			pCount = pCount - 1
		}
		inVals := make([]reflect.Value, pCount)
		for i := 0; i < pCount; i++ {

			param := handlerType.In(i)

			var inputVal reflect.Value
			var err error

			if i < inputSize {
				inputVal, err = cast.Param(param, inputs[i])
				if err != nil {
					cb(nil, err)
					return
				}
			} else {
				inputVal = reflect.Zero(param)
			}

			inVals[i] = inputVal
		}

		if isVariadic {
			// 处理可变参数
			inputs = inputs[pCount:]
			inputSize := len(inputs)
			elem := handlerType.In(handlerType.NumIn() - 1).Elem()
			for i := 0; i < inputSize; i++ {
				inputVal, err := cast.Param(elem, inputs[i])
				if err != nil {
					cb(nil, err)
					log.Error(err.Error())
					return
				}
				inVals = append(inVals, inputVal)
			}
		}

		// 调用处理函数
		out := handlerVal.Call(inVals)

		if cb == nil {
			return
		}

		// 处理返回值
		if len(out) == 0 {
			// 没有返回值
			go cb(nil, nil)
		} else if len(out) == 1 {
			// 只有一个返回值
			go cb(out[0].Interface(), nil)
		} else if len(out) == 2 {
			// 有2个返回值
			go cb(out[0].Interface(), out[1].Interface().(error))
		} else {
			// 多个返回值
			go cb(nil, fmt.Errorf("multiple return values are not supported"))
		}
	}
}

func (ipc *IPC) HasChannel(channel string) (exist bool) {
	_, exist = ipc.handlers[channel]
	return
}

//go:embed ipc.js
var ipcjs []byte

func (ipc *IPC) registerBootScript() {
	script := fmt.Sprintf(
		string(ipcjs),
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

		log.Debug("JS -> GO: %s", txt)

		var msg IPCMessage
		if err := json.Unmarshal(([]byte)(txt), &msg); err != nil {
			log.Error("JS -> GO, JSON 解析出错(%s): %s", err.Error(), txt)
			return
		}

		if msg.ReplyId != "" {
			ipc.mb.AddJob(func() {
				ipc.handleJSReply(&msg)
			})
			return
		}

		if msg.Channel != "" {
			if view, exist := ipc.mb.GetViewByJsExecState(es); exist {

				ipc.mb.AddJob(func() {
					ipc.invokeByJS(view, &msg)
				})
			}
			return
		}
	})
}

// JS 调用 handler
func (ipc *IPC) invokeByJS(view *View, msg *IPCMessage) {

	// 如果 ID 为空，则无须回复返回值
	if msg.ID == "" {
		ipc.Sent(msg.Channel, msg.Args...)
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
		ID:      "",
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

	cb, exist := ipc.resultWaiting[msg.ReplyId]
	if !exist {
		return
	}

	delete(ipc.resultWaiting, msg.ReplyId) // 接收到消息就从 map 中删除

	if msg.Error != "" {
		cb(nil, errors.New(msg.Error))
	} else {
		cb(msg.Result, nil)
	}
}

// JS 注册 handler 埋点
func (ipc *IPC) registerJSHandler() {
	// 注册 JS handler
	ipc.mb.js.bindFunction(JS_REGISTER_HANDLER, 1, func(es JsExecState) {
		arg := ipc.mb.js.Arg(es, 0)
		channel := ipc.mb.js.ToString(es, arg)

		view, exist := ipc.mb.GetViewByJsExecState(es)
		if !exist {
			log.Error("JS 注册 handler, 没有找到 view")
			return
		}

		// 将 JS handler 转为 GO handler
		ipc.handlers[channel] = func(cb resultCallback, args ...interface{}) {

			if cb == nil {
				msg := IPCMessage{
					ID:      "", // ID 为空则不需要回复
					Channel: channel,
					Args:    args,
				}
				sentMsgToView(view, msg)
				return
			}

			id := utils.RandString(8) // 生成key

			msg := IPCMessage{
				ID:      id,
				Channel: channel,
				Args:    args,
			}

			ipc.resultWaiting[id] = cb // 暂存 result callback

			sentMsgToView(view, msg)

			// 删除等待结果的callback
			go func() {
				defer delete(ipc.resultWaiting, id) // 删除等待结果的callback

				time.Sleep(10 * time.Second)
				cb, exist := ipc.resultWaiting[id]
				if !exist {
					return
				}

				cb(nil, errors.New("等待 JS Handler 处理结果超时"))
			}()
		}
	})
}

func (ipc *IPC) CallJsFunc(view *View, funcName string, args ...interface{}) async.InProgress[interface{}] {

	newArgs := make([]interface{}, 0, len(args)+1)
	newArgs = append(newArgs, funcName)
	newArgs = append(newArgs, args...)

	id := utils.RandString(8) // 生成key

	msg := IPCMessage{
		ID:      id,
		Channel: "callJsFunc",
		Args:    newArgs,
	}

	resultChan := make(chan interface{}) // result 管道
	errChan := make(chan error)          // 错误管道

	progress := async.New(10*time.Second, func() (interface{}, error) {
		return <-resultChan, <-errChan
	}).Start()

	ipc.resultWaiting[id] = func(result interface{}, err error) {
		resultChan <- result
		errChan <- err
	}

	sentMsgToView(view, msg)

	return progress
}

func sentMsgToView(view *View, msg IPCMessage) {

	msgTxt, _ := json.Marshal(msg)

	script := fmt.Sprintf(`window.top['%s'](%q)`, JS_GO2JS, msgTxt)

	log.Debug("GO -> JS: %s", msgTxt)

	view.RunJS(script)
}
