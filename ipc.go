package blink

import (
	"encoding/json"
	"fmt"
)

type IPC struct {
	mb *Blink

	handlers map[string]IPCHandler

	jsHandlerReply map[string]chan any
}

type IPCMessage struct {
	Channel string `json:"channel"`
	Args    []any  `json:"args"`
}

type IPCHandler func(args ...any) any

const IPC_CHANNEL_PREFIX = "IPC_CHANNEL_"

func newIPC(mb *Blink) *IPC {
	ipc := &IPC{
		mb: mb,

		handlers:       make(map[string]IPCHandler),
		jsHandlerReply: make(map[string]chan any),
	}

	ipc.registerJSHandler()
	ipc.registerJSHandlerReply()

	return ipc
}

func (ipc *IPC) registerJSHandler() {
	// 注册 JS handler
	ipc.mb.js.bindFunction(JS_HANDLE_REGISTER, 1, func(es JsExecState) {
		arg := ipc.mb.js.Arg(es, 0)
		channel := ipc.mb.js.ToString(es, arg)

		view := ipc.mb.GetViewByJsExecState(es)

		// 将 JS handler 转为 GO handler
		ipc.handlers[channel] = func(args ...any) any {
			key := RandString(8)
			msg := JsMessage{
				Key: key,
				Data: IPCMessage{
					Channel: channel,
					Args:    args,
				},
			}
			msgTxt, _ := json.Marshal(msg)
			result := make(chan any, 1)
			ipc.jsHandlerReply[key] = result // 暂存 result channel
			script := fmt.Sprintf(`window.top['%s'](%q)`, JS_HANDLE_PROCESS, string(msgTxt))
			view.RunJS(script)
			return <-result
		}
	})
}

func (ipc *IPC) registerJSHandlerReply() {
	// 注册 JS handler 处理结果
	ipc.mb.js.bindFunction(JS_HANDLE_PROCESS_REPLY, 1, func(es JsExecState) {
		arg := ipc.mb.js.Arg(es, 0)
		txt := ipc.mb.js.ToString(es, arg)
		msg := JsMessage{}
		json.Unmarshal([]byte(txt), &msg)
		result := ipc.jsHandlerReply[msg.Key]
		result <- msg.Data
	})
}

func (ipc *IPC) Invoke(channel string, args ...any) any {
	handler, exist := ipc.handlers[channel]
	if !exist {
		logError("IPC: Invoke: handler not exist: %s", channel)
		return nil
	}

	return handler(args...)
}

func (ipc *IPC) invokeJS(view *View, key string, msg IPCMessage) {
	res := ipc.Invoke(msg.Channel, msg.Args...)

	jsmsg := JsMessage{
		Key:  key,
		Data: res,
	}

	jsmsgTxt, _ := json.Marshal(jsmsg)

	script := fmt.Sprintf(`window.top['%s'](%q)`, JS_JS2GO_REPLY, string(jsmsgTxt))

	view.RunJS(script)
}

func (ipc *IPC) Handle(channel string, handler IPCHandler) {
	ipc.handlers[channel] = handler
}
