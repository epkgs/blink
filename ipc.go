package blink

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/chebyrash/promise"
	"github.com/epkgs/blink/internal/cast"
	"github.com/epkgs/blink/internal/log"
	"github.com/epkgs/blink/pkg/utils"
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

// 封装 ipcPedding, 为 Get/Add/Del 提供锁保护
type ipcPendding struct {
	ctx       context.Context
	mu        *sync.Mutex
	callbacks map[string]resultCallback
}

func newIPCPendding(ctx context.Context) *ipcPendding {
	return &ipcPendding{
		ctx:       ctx,
		mu:        &sync.Mutex{},
		callbacks: make(map[string]resultCallback),
	}
}

func (p *ipcPendding) Add(id string, cb resultCallback) {
	// 异步，避免阻塞
	utils.Go(func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		p.callbacks[id] = cb
	}, nil)

	// 超时处理
	utils.Go(func() {

		time.Sleep(10 * time.Second)

		if p.mu.TryLock() {
			defer p.mu.Unlock()
		}

		cb, exist := p.callbacks[id]
		if !exist {
			return
		}

		delete(p.callbacks, id) // 删除等待结果的callback

		cb(nil, errors.New("等待 JS Handler 处理结果超时"))
	}, nil)
}

func (p *ipcPendding) Del(id string) {
	// 异步，避免阻塞
	utils.Go(func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		delete(p.callbacks, id)
	}, nil)
}

func (p *ipcPendding) Get(id string) (resultCallback, bool) {
	cb, exist := p.callbacks[id]
	return cb, exist
}

type IPC struct {
	mb *Blink

	handlers map[string]ipcHandler

	pendding *ipcPendding
}

type IPCMessage struct {
	ID      string        `json:"id"`               // 消息 ID
	ReplyId string        `json:"replyId"`          // 回复ID
	Channel string        `json:"channel"`          // 通道
	Args    []interface{} `json:"args"`             // 参数
	Result  interface{}   `json:"result,omitempty"` // 返回值，当有回复ID时，此字段有效
	Error   string        `json:"error,omitempty"`  // 是否错误，当有回复ID时，此字段有效
}

func newIPC(mb *Blink) *IPC {
	ipc := &IPC{
		mb: mb,

		handlers: make(map[string]ipcHandler),
	}

	ipc.pendding = newIPCPendding(mb.Ctx)

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

	result := make(chan interface{}, 1)
	err := make(chan error, 1)

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
// handler 必须为函数，参数任意，返回值最多为2个
//   - 1个返回值：会自动判断返回值是否为 error
//   - 2个返回值：第一个为 结果，第二个为 error
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

		// 异步处理 handler
		utils.Go(func() {

			select {
			case <-time.After(10 * time.Second):
				log.Debug("IPC 调用超时")
				cb(nil, errors.New("IPC 调用超时"))
				return
			default:
				// 调用处理函数
				out := handlerVal.Call(inVals)

				if cb == nil {
					return
				}

				// 处理返回值
				if len(out) == 0 {
					// 没有返回值
					cb(nil, nil)
				} else if len(out) == 1 {
					// 只有一个返回值
					result := out[0].Interface()

					switch res := result.(type) {
					case error:
						cb(nil, res)
					default:
						cb(res, nil)
					}
				} else if len(out) == 2 {
					// 有2个返回值
					res := out[0].Interface()
					var err error
					switch e := out[1].Interface().(type) {
					case error:
						err = e
					default:
						err = nil
					}
					cb(res, err)
				} else {
					// 多个返回值
					cb(nil, fmt.Errorf("more than 2 return values are not supported"))
				}
			}

		}, func(err error) {
			log.Error("panic by ipc handler[ %v ]: %v", channel, err)
			cb(nil, err)
		})
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
		_ = ipc.Sent(msg.Channel, msg.Args...)
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

	cb, exist := ipc.pendding.Get(msg.ReplyId)
	if !exist {
		return
	}

	ipc.pendding.Del(msg.ReplyId) // 接收到消息就从 map 中删除

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

			ipc.pendding.Add(id, cb) // 添加到等待结果的 map

			sentMsgToView(view, msg)
		}
	})
}

func (ipc *IPC) RunJsFunc(view *View, funcName string, args ...interface{}) *promise.Promise[any] {

	newArgs := make([]interface{}, 0, len(args)+1)
	newArgs = append(newArgs, funcName)
	newArgs = append(newArgs, args...)

	id := utils.RandString(8) // 生成key

	msg := IPCMessage{
		ID:      id,
		Channel: "callJsFunc",
		Args:    newArgs,
	}

	resolve := func(any) {}
	reject := func(error) {}
	p := promise.New(func(resolv func(any), rej func(error)) {
		resolve = resolv
		reject = rej
	})

	cb := func(result any, err error) {
		if err != nil {
			reject(err)
		} else {
			resolve(result)
		}
	}

	ipc.pendding.Add(id, cb)

	sentMsgToView(view, msg)

	return p
}

func sentMsgToView(view *View, msg IPCMessage) {

	msgTxt, _ := json.Marshal(msg)

	script := fmt.Sprintf(`window.top['%s'](%q)`, JS_GO2JS, msgTxt)

	log.Debug("GO -> JS: %s", msgTxt)

	view.RunJs(script)
}
