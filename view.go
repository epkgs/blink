package blink

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unsafe"

	"github.com/lxn/win"
)

type OnConsoleCallback func(level int, message, sourceName string, sourceLine int, stackTrace string)
type OnDestroyCallback func()
type OnLoadUrlBeginCallback func(url string, job WkeNetJob) bool
type OnLoadUrlEndCallback func(url string, job WkeNetJob, buf []byte)
type OnDocumentReadyCallback func(frame WkeWebFrameHandle)
type OnDidCreateScriptContextCallback func(frame WkeWebFrameHandle, context uintptr, exGroup, worldId int)
type OnTitleChangedCallback func(title string)
type OnDownloadCallback func(url string)

type View struct {
	Hwnd   WkeHandle
	Window *Window

	mb     *Blink
	parent *View

	onDestroyCallbacks       []OnDestroyCallback
	onLoadUrlBeginCallbacks  []OnLoadUrlBeginCallback
	onLoadUrlEndCallbacks    []OnLoadUrlEndCallback
	onDocumentReadyCallbacks []OnDocumentReadyCallback
	onTitleChangedCallbacks  []OnTitleChangedCallback
	onDownloadCallbacks      []OnDownloadCallback
}

func NewView(mb *Blink, hwnd WkeHandle, windowType WkeWindowType, parent ...*View) *View {

	var p *View = nil

	if len(parent) >= 1 {
		p = parent[0]
	}

	view := &View{
		mb:     mb,
		Hwnd:   hwnd,
		parent: p,

		onDestroyCallbacks:       []OnDestroyCallback{},
		onLoadUrlBeginCallbacks:  []OnLoadUrlBeginCallback{},
		onLoadUrlEndCallbacks:    []OnLoadUrlEndCallback{},
		onDocumentReadyCallbacks: []OnDocumentReadyCallback{},
		onTitleChangedCallbacks:  []OnTitleChangedCallback{},
		onDownloadCallbacks:      []OnDownloadCallback{},
	}

	view.Window = newWindow(mb, view, windowType)

	view.SetLocalStorageFullPath(view.mb.Config.GetStoragePath())
	view.SetCookieJarFullPath(view.mb.Config.GetCookieFilePath())

	view.registerFileSystem()

	view.registerOnDestroy()
	view.registerOnLoadUrlBegin()
	view.registerOnLoadUrlEnd()
	view.registerOnDocumentReady()
	view.registerOnTitleChanged()
	view.registerOnDownload()

	view.listenMinBtnClick()
	view.listenMaxBtnClick()
	view.listenCloseBtnClick()
	view.listenCaptionDrag()

	view.addToPool()

	// 添加默认下载操作
	view.OnDownload(func(url string) {
		view.mb.Downloader.Download(url)
	})

	return view
}

func (v *View) addToPool() {
	locker.Lock()
	v.mb.views[v.Hwnd] = v
	v.mb.windows[v.Window.Hwnd] = v.Window
	locker.Unlock()

	v.OnDestroy(func() {
		locker.Lock()
		delete(v.mb.windows, v.Window.Hwnd)
		delete(v.mb.views, v.Hwnd)
		locker.Unlock()

		// 删除子view
		for _, child := range v.mb.views {
			if child.parent == v {
				child.Destroy()
			}
		}
	})
}

func (v *View) Show() {
	v.mb.CallFunc("wkeShowWindow", uintptr(v.Hwnd), 1)
}

func (v *View) Hide() {
	v.mb.CallFunc("wkeShowWindow", uintptr(v.Hwnd), 0)
}

func (v *View) Destroy() {
	v.mb.CallFunc("wkeDestroyWebWindow", uintptr(v.Hwnd))
}

func (v *View) Reload() bool {
	r, _, _ := v.mb.CallFunc("wkeReload", uintptr(v.Hwnd))
	return r != 0
}

func (v *View) ForceReload() {
	v.LoadURL(v.GetURL())
}

func (v *View) LoadURL(url string) {
	v.mb.CallFunc("wkeLoadURL", uintptr(v.Hwnd), StringToPtr(url))
}

func (v *View) GetURL() string {
	r, _, _ := v.mb.CallFunc("wkeGetURL", uintptr(v.Hwnd))
	return PtrToString(r)
}

// 设置local storage的全路径。如“c:\mb\LocalStorage\”
// 注意：这个接口只能接受目录。
func (v *View) SetLocalStorageFullPath(path string) {
	v.mb.CallFunc("wkeSetLocalStorageFullPath", uintptr(v.Hwnd), StringToWCharPtr(path))
}

// 设置cookie的全路径+文件名，如“c:\mb\cookie.dat”
func (v *View) SetCookieJarFullPath(path string) {
	v.mb.CallFunc("wkeSetCookieJarFullPath", uintptr(v.Hwnd), StringToWCharPtr(path))
}

func (v *View) GetWindowHandle() WkeHandle {
	ptr, _, _ := v.mb.CallFunc("wkeGetWindowHandle", uintptr(v.Hwnd))
	return WkeHandle(ptr)
}

func (v *View) Resize(width, height int32) {
	v.mb.CallFunc("wkeResize", uintptr(v.Hwnd), uintptr(width), uintptr(height))
}

// 返回 true 标识已处理
func (v *View) DispatchMessage(msg *win.MSG) bool {

	// do something...

	return false
}

func (v *View) registerFileSystem() {
	v.OnLoadUrlBegin(func(url string, job WkeNetJob) bool {

		f := v.mb.Resource.GetFile(url)

		// 找不到文件
		if f == nil {
			return false
		}

		defer f.Close()

		byt, err := io.ReadAll(f)
		// 读取文件错误
		if err != nil {
			return false
		}

		v.mb.NetSetData(job, byt)

		// 找到并读取正常，返回 true 取消后继的网络请求
		return true

	})
}

// 可以添加多个 callback，将按照加入顺序依次执行
func (v *View) OnDestroy(callback OnDestroyCallback) {
	v.onDestroyCallbacks = append(v.onDestroyCallbacks, callback)
}
func (v *View) registerOnDestroy() {
	var handler = func(view, param uintptr) uintptr {
		for _, callback := range v.onDestroyCallbacks {
			callback()
		}
		return 0
	}
	v.mb.CallFunc("wkeOnWindowDestroy", uintptr(v.Hwnd), CallbackToPtr(handler), 0)
}

func (v *View) OnLoadUrlBegin(callback OnLoadUrlBeginCallback) {
	v.onLoadUrlBeginCallbacks = append(v.onLoadUrlBeginCallbacks, callback)
}
func (v *View) registerOnLoadUrlBegin() {
	var handler = func(view, param, url, job uintptr) (boolPtr uintptr) {
		for _, callback := range v.onLoadUrlBeginCallbacks {
			// 如果返回结果为 true，则中断后面的处理，直接返回 true
			// 返回 true 则中断、阻止后面的网络请求
			if callback(PtrToString(url), WkeNetJob(job)) {
				return 1
			}
		}
		return 0
	}

	v.mb.CallFunc("wkeOnLoadUrlBegin", uintptr(v.Hwnd), CallbackToPtr(handler), 0)
}

func (v *View) OnLoadUrlEnd(callback OnLoadUrlEndCallback) {
	v.onLoadUrlEndCallbacks = append(v.onLoadUrlEndCallbacks, callback)
}
func (v *View) registerOnLoadUrlEnd() {
	var handler = func(view, param, url, job, buf, len uintptr) uintptr {

		_url := PtrToString(url)
		_job := WkeNetJob(job)
		_buf := Read[byte](buf)
		for _, callback := range v.onLoadUrlEndCallbacks {
			callback(_url, _job, _buf)
		}
		return 0
	}
	v.mb.CallFunc("wkeOnLoadUrlEnd", uintptr(v.Hwnd), CallbackToPtr(handler), 0)
}

func (v *View) OnDocumentReady(callback OnDocumentReadyCallback) {
	v.onDocumentReadyCallbacks = append(v.onDocumentReadyCallbacks, callback)
}

func (v *View) registerOnDocumentReady() {
	var cb WkeDocumentReady2Callback = func(view WkeHandle, param uintptr, frame WkeWebFrameHandle) (voidRes uintptr) {

		for _, callback := range v.onDocumentReadyCallbacks {
			func(v *View) {
				callback(frame)
			}(v)
		}

		return 0
	}
	v.mb.CallFunc("wkeOnDocumentReady2", uintptr(v.Hwnd), CallbackToPtr(cb), 0)
}

func (v *View) IsMainFrame(frameId WkeWebFrameHandle) bool {
	p, _, _ := v.mb.CallFunc("wkeIsMainFrame", uintptr(v.Hwnd), uintptr(frameId))

	return p != 0
}

func (v *View) EnableBorderResize() {
	v.Window.enableBorderResize = true
}

func (v *View) GetRect() *WkeRect {
	ptr, _, _ := v.mb.CallFunc("wkeGetCaretRect2", uintptr(v.Hwnd))
	return (*WkeRect)(unsafe.Pointer(ptr))
}

// 仅作用于 主frame
func (v *View) RunJS(script string) {

	if v.IsDocumentReady() {
		v.mb.CallFunc("wkeRunJS", uintptr(v.Hwnd), StringToPtr(script))
		return
	}

	v.OnDocumentReady(func(frame WkeWebFrameHandle) {
		if !v.IsMainFrame(frame) {
			return
		}
		v.mb.CallFunc("wkeRunJS", uintptr(v.Hwnd), StringToPtr(script))
	})

}

func (v *View) CallJsFunc(callback func(result any), funcName string, args ...any) {

	key := RandString(8)

	script := `
	const rootWin = window.top || window.parent || window;
	const tunnel = rootWin['%s'] || (function(e) { });
	const msg = %s;
	msg.data = %s(...msg.data);
	tunnel(JSON.stringify(msg));
	`

	msg := &JsMessage{
		Key:  key,
		Data: args,
	}

	if args == nil {
		msg.Data = []any{}
	}

	jsonBytes, err := json.Marshal(&msg)

	if err != nil {
		return
	}

	jsonTxt := string(jsonBytes)

	script = fmt.Sprintf(script,
		JS_MSG_FUNC,
		jsonTxt,
		funcName,
	)

	if callback != nil {

		var jsCallback JsCallback = func(result any) {
			callback(result)
			v.mb.JS.removeCallback(key)
		}

		// 注册callback
		v.mb.JS.AddCallback(key, jsCallback)
	}

	v.RunJS(script)
}

func (v *View) OnDidCreateScriptContext(callback OnDidCreateScriptContextCallback) {

	var cb WkeDidCreateScriptContextCallback = func(view WkeHandle, param uintptr, frame WkeWebFrameHandle, context uintptr, exGroup, worldId int) (voidRes uintptr) {
		callback(frame, context, exGroup, worldId)
		return 0
	}

	v.mb.CallFunc("wkeOnDidCreateScriptContext", uintptr(v.Hwnd), CallbackToPtr(cb), 0)
}

// JS.bind(".mb-minimize-btn", "click", func)
func (v *View) AddEventListener(selector, eventType string, callback func(), preScripts ...string) {

	key := strconv.FormatUint(uint64(v.Hwnd), 10) + " " + selector + " " + eventType

	msg := &JsMessage{
		Key: key,
		Data: JsEvent{
			Selector:  selector,
			EventType: eventType,
		},
	}

	msgBytes, err := json.Marshal(msg)

	if err != nil {
		return
	}

	msgTxt := string(msgBytes)

	script := `
		const rootWin = window.top || window.parent || window;
		const tunnel = rootWin['%s'] || (function(e) { });
		const args = %s
		const els = document.querySelectorAll(args.data.selector);
		const handler = function(e) { %s; e.preventDefault(); tunnel(%q); };
		for (const el of els) {
			el.removeEventListener(args.data.eventType, handler);
			el.addEventListener(args.data.eventType, handler);
		}
	`

	script = fmt.Sprintf(script,
		JS_MSG_FUNC,
		msgTxt,
		strings.Join(preScripts, ";"),
		msgTxt,
	)

	var jsCallback JsCallback = func(args any) {
		callback()
	}

	// 注册callback
	v.mb.JS.AddCallback(key, jsCallback)

	v.RunJS(script)

}

func (v *View) RemoveEventListener(selector, eventType string) {
	key := strconv.FormatUint(uint64(v.Hwnd), 10) + " " + selector + " " + eventType

	v.mb.JS.removeCallback(key)
}

func (v *View) listenMinBtnClick() {
	v.AddEventListener(".mb-btn-min", "click", func() {
		v.Window.Minimize()
	})

}

func (v *View) listenMaxBtnClick() {

	preScript := `this.classList.toggle('maximized');`

	v.AddEventListener(".mb-btn-max", "click", func() {
		if v.Window.IsMaximized() {
			v.Window.Restore()
		} else {
			v.Window.Maximize()
		}
	}, preScript)
}

func (v *View) listenCloseBtnClick() {
	v.AddEventListener(".mb-btn-close", "click", func() {
		v.Destroy()
	})
}

// 监听窗口拖动
func (v *View) listenCaptionDrag() {

	preScript := `if(e.target.closest('.mb-caption-nodrag')) return;`

	v.AddEventListener(".mb-caption-drag", "mousedown", func() {
		if v.Window.IsMaximized() {
			return
		}
		v.Window.Move()
	}, preScript)
}

func (v *View) OnConsole(callback OnConsoleCallback) {

	var cb WkeConsoleCallback = func(view WkeHandle, param uintptr, level WkeConsoleLevel, message, sourceName WkeString, sourceLine uint32, stackTrace WkeString) (voidRes uintptr) {

		callback(int(level), v.mb.GetString(message), v.mb.GetString(sourceName), int(sourceLine), v.mb.GetString(stackTrace))

		return 0
	}

	v.mb.CallFunc("wkeOnConsole", uintptr(v.Hwnd), CallbackToPtr(cb), 0)
}

func (v *View) IsDocumentReady() bool {
	p, _, _ := v.mb.CallFunc("wkeIsDocumentReady", uintptr(v.Hwnd))
	return p != 0
}

func (v *View) OnTitleChanged(callback OnTitleChangedCallback) {
	v.onTitleChangedCallbacks = append(v.onTitleChangedCallbacks, callback)
}
func (v *View) registerOnTitleChanged() {

	var cb WkeTitleChangedCallback = func(view WkeHandle, param uintptr, title WkeString) (voidRes uintptr) {
		_title := v.mb.GetString(title)

		for _, callback := range v.onTitleChangedCallbacks {
			callback(_title)
		}
		return
	}

	v.mb.CallFunc("wkeOnTitleChanged", uintptr(v.Hwnd), CallbackToPtr(cb), 0)
}

func (v *View) OnDownload(callback OnDownloadCallback) {
	v.onDownloadCallbacks = append(v.onDownloadCallbacks, callback)
}
func (v *View) registerOnDownload() {
	var cb WkeDownloadCallback = func(view WkeHandle, param uintptr, url uintptr) (voidRes uintptr) {
		link := PtrToString(url)
		for _, callback := range v.onDownloadCallbacks {
			go callback(link)
		}
		return
	}

	v.mb.CallFunc("wkeOnDownload", uintptr(v.Hwnd), CallbackToPtr(cb), 0)
}
