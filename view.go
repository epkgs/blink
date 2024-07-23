package blink

import (
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/epkgs/blink/internal/log"
	"github.com/epkgs/blink/internal/utils"
)

type OnDomEventCallback func()
type OnConsoleCallback func(level int, message, sourceName string, sourceLine int, stackTrace string)
type OnClosingCallback func() bool // 返回 false 拒绝关闭窗口
type OnDestroyCallback func()
type OnLoadUrlBeginCallback func(url string, job WkeNetJob) bool
type OnLoadUrlEndCallback func(url string, job WkeNetJob, buf []byte)
type OnDocumentReadyCallback func(frame WkeWebFrameHandle)
type OnDidCreateScriptContextCallback func(frame WkeWebFrameHandle, context uintptr, exGroup, worldId int)
type OnWillReleaseScriptContextCallback func(frameId WkeWebFrameHandle, context uintptr, worldId int)
type OnTitleChangedCallback func(title string)
type OnDownloadCallback func(url string)
type OnOtherLoadCallback func(loadType WkeOtherLoadType, info *WkeTempCallbackInfo)

type bindEvent[T any] struct {
	Callbacks map[string]T
	Register  sync.Once
}

func newBindEvent[T any]() *bindEvent[T] {
	return &bindEvent[T]{Callbacks: make(map[string]T)}
}

type View struct {
	Hwnd     WkeHandle
	Window   *Window
	DevTools *View

	mb     *Blink
	parent *View

	_didCreateScriptContext bool // 标记是否已经创建了脚本上下文

	_onDomEvent                         *bindEvent[OnDomEventCallback]
	_onConsole                          *bindEvent[OnConsoleCallback]
	_onClosing                          *bindEvent[OnClosingCallback]
	_onDestroy                          *bindEvent[OnDestroyCallback]
	_onLoadUrlBegin                     *bindEvent[OnLoadUrlBeginCallback]
	_onLoadUrlEnd                       *bindEvent[OnLoadUrlEndCallback]
	_onDocumentReady                    *bindEvent[OnDocumentReadyCallback]
	_onTitleChanged                     *bindEvent[OnTitleChangedCallback]
	_onDownload                         *bindEvent[OnDownloadCallback]
	_onDidCreateScriptContext           *bindEvent[OnDidCreateScriptContextCallback]
	_onWillReleaseScriptContextCallback *bindEvent[OnWillReleaseScriptContextCallback]
	_onOtherLoad                        *bindEvent[OnOtherLoadCallback]
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

		_onDomEvent:                         newBindEvent[OnDomEventCallback](),
		_onConsole:                          newBindEvent[OnConsoleCallback](),
		_onClosing:                          newBindEvent[OnClosingCallback](),
		_onDestroy:                          newBindEvent[OnDestroyCallback](),
		_onLoadUrlBegin:                     newBindEvent[OnLoadUrlBeginCallback](),
		_onLoadUrlEnd:                       newBindEvent[OnLoadUrlEndCallback](),
		_onDocumentReady:                    newBindEvent[OnDocumentReadyCallback](),
		_onTitleChanged:                     newBindEvent[OnTitleChangedCallback](),
		_onDownload:                         newBindEvent[OnDownloadCallback](),
		_onDidCreateScriptContext:           newBindEvent[OnDidCreateScriptContextCallback](),
		_onWillReleaseScriptContextCallback: newBindEvent[OnWillReleaseScriptContextCallback](),
		_onOtherLoad:                        newBindEvent[OnOtherLoadCallback](),
	}

	view.Window = newWindow(mb, view, windowType)

	view.SetLocalStorageFullPath(view.mb.Config.GetStoragePath())
	view.SetCookieJarFullPath(view.mb.Config.GetCookieFileABS())

	view.registerFileSystem()

	view.injectBootScripts()
	view.watchScriptContextState()
	view.bindDomEvents() // 绑定一些DOM事件

	view.addToPool()

	// 添加默认下载操作
	view.OnDownload(func(url string) {
		view.mb.Downloader.Download(url)
	})

	return view
}

func (v *View) addToPool() {

	locker.Lock()
	defer locker.Unlock()

	v.mb.views[v.Hwnd] = v
	v.mb.windows[v.Window.Hwnd] = v.Window

	log.Debug("Add view to BLINK, now SIZE: %d", len(v.mb.views))

	v.OnDestroy(func() {

		func() {
			locker.Lock()
			defer locker.Unlock()

			delete(v.mb.windows, v.Window.Hwnd)
			delete(v.mb.views, v.Hwnd)

		}()

		for _, child := range v.mb.views {
			if child.parent == v {
				child.DestroyWindow()
			}
		}

	})
}

func (v *View) injectBootScripts() {
	var script string

	for _, s := range v.mb.bootScripts {
		script += s + ";\n"
	}

	v.OnDidCreateScriptContext(func(frame WkeWebFrameHandle, context uintptr, exGroup, worldId int) {
		v.RunJS(script)
	})
}

func (v *View) ShowWindow() {
	v.Window.Show()
}

func (v *View) HideWindow() {
	v.Window.Hide()
}

func (v *View) CloseWindow() {
	v.Window.Close()
}

// 销毁wkeWebView对应的所有数据结构，包括真实窗口等
func (v *View) DestroyWindow() {
	// v.Window.Destroy()
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
//
// callback 返回 false 拒绝关闭窗口
func (v *View) OnClosing(callback OnClosingCallback) (stop func()) {

	v._onClosing.Register.Do(func() {
		var handler WkeWindowClosingCallback = func(view WkeHandle, param uintptr) (boolRes uintptr) {
			log.Debug("Trigger view.OnClosing")
			for _, callback := range v._onClosing.Callbacks {
				if ok := callback(); !ok {
					return BoolToPtr(false)
				}
			}
			return BoolToPtr(true)
		}
		v.mb.CallFunc("wkeOnWindowClosing", uintptr(v.Hwnd), CallbackToPtr(handler), 0)
	})

	key := utils.RandString(10)

	v._onClosing.Callbacks[key] = callback

	return func() {
		delete(v._onClosing.Callbacks, key)
	}
}

// 可以添加多个 callback，将按照加入顺序依次执行
func (v *View) OnDestroy(callback OnDestroyCallback) (stop func()) {

	v._onDestroy.Register.Do(func() {
		var handler WkeWindowDestroyCallback = func(view WkeHandle, param uintptr) (voidRes uintptr) {
			log.Debug("Trigger view.OnDestroy")
			for _, callback := range v._onDestroy.Callbacks {
				callback()
			}
			return
		}
		v.mb.CallFunc("wkeOnWindowDestroy", uintptr(v.Hwnd), CallbackToPtr(handler), 0)
	})

	key := utils.RandString(10)

	v._onDestroy.Callbacks[key] = callback

	return func() {
		delete(v._onDestroy.Callbacks, key)
	}
}

func (v *View) OnLoadUrlBegin(callback OnLoadUrlBeginCallback) (stop func()) {

	v._onLoadUrlBegin.Register.Do(func() {
		var handler = func(view, param, url, job uintptr) (boolPtr uintptr) {
			urlPtr := PtrToString(url)
			jobPtr := WkeNetJob(job)
			for _, callback := range v._onLoadUrlBegin.Callbacks {
				// 返回 true 则中断、阻止后面的网络请求
				if callback(urlPtr, jobPtr) {
					return 1 // 返回 true 的 uintptr
				}
			}
			return 0 // 返回 false 的 uintptr
		}

		v.mb.CallFunc("wkeOnLoadUrlBegin", uintptr(v.Hwnd), CallbackToPtr(handler), 0)
	})

	key := utils.RandString(10)

	v._onLoadUrlBegin.Callbacks[key] = callback

	return func() {
		delete(v._onLoadUrlBegin.Callbacks, key)
	}
}

func (v *View) OnLoadUrlEnd(callback OnLoadUrlEndCallback) (stop func()) {

	v._onLoadUrlEnd.Register.Do(func() {
		var handler = func(view, param, url, job, buf, len uintptr) uintptr {

			_url := PtrToString(url)
			_job := WkeNetJob(job)
			_buf := CopyBytes(buf, int(len))
			for _, callback := range v._onLoadUrlEnd.Callbacks {
				callback(_url, _job, _buf)
			}
			return 0
		}
		v.mb.CallFunc("wkeOnLoadUrlEnd", uintptr(v.Hwnd), CallbackToPtr(handler), 0)
	})

	key := utils.RandString(10)

	v._onLoadUrlEnd.Callbacks[key] = callback

	return func() {
		delete(v._onLoadUrlEnd.Callbacks, key)
	}
}

func (v *View) OnDocumentReady(callback OnDocumentReadyCallback) (stop func()) {

	v._onDocumentReady.Register.Do(func() {
		var cb WkeDocumentReady2Callback = func(view WkeHandle, param uintptr, frame WkeWebFrameHandle) (voidRes uintptr) {

			for _, callback := range v._onDocumentReady.Callbacks {
				callback(frame)
			}

			return 0
		}
		v.mb.CallFunc("wkeOnDocumentReady2", uintptr(v.Hwnd), CallbackToPtr(cb), 0)
	})

	key := utils.RandString(10)

	v._onDocumentReady.Callbacks[key] = callback

	return func() {
		delete(v._onDocumentReady.Callbacks, key)
	}
}

func (v *View) IsMainFrame(frameId WkeWebFrameHandle) bool {
	p, _, _ := v.mb.CallFunc("wkeIsMainFrame", uintptr(v.Hwnd), uintptr(frameId))

	return p != 0
}

func (v *View) GetRect() *WkeRect {
	ptr, _, _ := v.mb.CallFunc("wkeGetCaretRect2", uintptr(v.Hwnd))
	return (*WkeRect)(unsafe.Pointer(ptr))
}

// 仅作用于 主frame，会自动判断是否 document ready，仅执行一次
func (v *View) DoWhenDocumentReady(callback func()) {
	if v.IsDocumentReady() {
		callback()
	} else {
		stop := func() {}
		stop = v.OnDocumentReady(func(frame WkeWebFrameHandle) {
			if !v.IsMainFrame(frame) {
				return
			}
			stop()
			callback()
		})
	}
}

// 仅作用于 主frame，会自动判断是否已创建 Script Context，仅执行一次
func (v *View) DoWhenDidCreateScriptContext(callback func()) {
	if v.IsDidCreateScriptContext() {
		callback()
	} else {
		stop := func() {}
		stop = v.OnDidCreateScriptContext(func(frame WkeWebFrameHandle, context uintptr, exGroup, worldId int) {
			if !v.IsMainFrame(frame) {
				return
			}
			stop()
			callback()
		})
	}
}

// 仅作用于 主frame，会自动判断是否 document ready
func (v *View) RunJS(script string) JsValue {
	r1, _, _ := v.mb.CallFunc("wkeRunJS", uintptr(v.Hwnd), StringToPtr(script))

	return JsValue(r1)
}

// 可指定 frame，会自动判断是否 document ready
func (v *View) RunJsByFrame(frame WkeWebFrameHandle, script string) JsValue {

	r1, _, _ := v.mb.CallFunc("wkeRunJsByFrame", uintptr(frame), StringToPtr(script), 0)

	return JsValue(r1)
}

func (v *View) CallJsFunc(funcName string, args ...interface{}) (result chan interface{}) {

	return v.mb.IPC.CallJsFunc(v, funcName, args...)
}

func (v *View) OnDidCreateScriptContext(callback OnDidCreateScriptContextCallback) (stop func()) {

	v._onDidCreateScriptContext.Register.Do(func() {
		var cb WkeDidCreateScriptContextCallback = func(view WkeHandle, param uintptr, frame WkeWebFrameHandle, context uintptr, exGroup, worldId int) (voidRes uintptr) {

			for _, callback := range v._onDidCreateScriptContext.Callbacks {
				callback(frame, context, exGroup, worldId)
			}
			return 0
		}
		v.mb.CallFunc("wkeOnDidCreateScriptContext", uintptr(v.Hwnd), CallbackToPtr(cb), 0)
	})

	key := utils.RandString(8)
	v._onDidCreateScriptContext.Callbacks[key] = callback

	return func() {
		delete(v._onDidCreateScriptContext.Callbacks, key)
	}
}

func (v *View) OnWillReleaseScriptContext(callback OnWillReleaseScriptContextCallback) (stop func()) {
	v._onWillReleaseScriptContextCallback.Register.Do(func() {
		var cb WkeWillReleaseScriptContextCallback = func(webView WkeHandle, param uintptr, frameId WkeWebFrameHandle, context uintptr, worldId int) (voidRes uintptr) {
			for _, callback := range v._onWillReleaseScriptContextCallback.Callbacks {
				callback(frameId, context, worldId)
			}
			return 0
		}

		v.mb.CallFunc("wkeOnWillReleaseScriptContext", uintptr(v.Hwnd), CallbackToPtr(cb), 0)
	})

	key := utils.RandString(8)
	v._onWillReleaseScriptContextCallback.Callbacks[key] = callback

	return func() {
		delete(v._onWillReleaseScriptContextCallback.Callbacks, key)
	}
}

func (v *View) watchScriptContextState() {
	v.OnDidCreateScriptContext(func(frame WkeWebFrameHandle, context uintptr, exGroup, worldId int) {
		v._didCreateScriptContext = true
	})
	v.OnWillReleaseScriptContext(func(frameId WkeWebFrameHandle, context uintptr, worldId int) {
		v._didCreateScriptContext = false
	})
}

func (v *View) IsDidCreateScriptContext() bool {
	return v._didCreateScriptContext
}

// JS.bind(".mb-minimize-btn", "click", func)
func (v *View) AddEventListener(selector, eventType string, callback func(), preScripts ...string) (stop func()) {

	script := `
	(()=>{
		const VIEW_HANDLE = '%s';
		const JS_IPC = '%s';
		const selector = '%s';
		const eventType = '%s';
		
		const els = document.querySelectorAll(selector);
		
		const handler = function(e) {
			%s; // pre-event
		
			e.preventDefault();
		
			const ipc = window.top[JS_IPC]
			ipc.sent('domEvent', VIEW_HANDLE, selector, eventType)
		};
		
		for (let i = 0; i < els.length; i++) {
			els[i].removeEventListener(eventType, handler);
			els[i].addEventListener(eventType, handler);
		}
	
	})();
	`

	script = fmt.Sprintf(
		script,
		strconv.FormatUint(uint64(v.Hwnd), 10),
		JS_IPC,
		selector,
		eventType,
		strings.Join(preScripts, ";"),
	)

	v._onDomEvent.Register.Do(func() {
		v.mb.IPC.Handle("domEvent", func(hwndStr, selector, eventType string) {
			hwnd, err := strconv.ParseUint(hwndStr, 10, 64)
			if err != nil {
				log.Error("hwnd 转换失败：%s", err.Error())
				return
			}

			view, exist := v.mb.GetViewByHandle(WkeHandle(hwnd))
			if !exist {
				return
			}

			key := selector + " " + eventType

			callback, exist := view._onDomEvent.Callbacks[key]
			if !exist {
				return
			}

			callback()
		})
	})

	key := selector + " " + eventType

	v._onDomEvent.Callbacks[key] = callback // 增加 callback

	v.RunJS(script)

	return func() {
		delete(v._onDomEvent.Callbacks, key)
	}
}

func (v *View) RemoveEventListener(selector, eventType string) {

	key := selector + " " + eventType

	delete(v._onDomEvent.Callbacks, key)
}

func (v *View) bindDomEvents() {
	v.OnDocumentReady(func(frame WkeWebFrameHandle) {
		// 最小化按钮
		v.AddEventListener(".__mb_min__", "click", func() {
			v.Window.Minimize()
		})

		// 最大化按钮
		v.AddEventListener(".__mb_max__", "click", func() {
			if v.Window.IsMaximized() {
				v.Window.Restore()
			} else {
				v.Window.Maximize()
			}
		}, `this.classList.toggle('__mb_maximized');`)

		// 关闭按钮
		v.AddEventListener(".__mb_close__", "click", func() {
			v.CloseWindow()
		})

		// 监听窗口拖动
		v.AddEventListener(".__mb_drag__, .__mb_caption__", "mousedown", func() {
			if v.Window.IsMaximized() {
				return
			}
			v.Window.EnableDragging()
		},
			`if(e.target.closest('.__mb_nodrag__')) return;`, // 如果是在禁止拖动区域，则不监听
		)

		// 监听标题栏双击事件
		v.AddEventListener(".__mb_caption__", "dblclick", func() {
			if v.Window.IsMaximized() {
				v.Window.Restore()
			} else {
				v.Window.Maximize()
			}
		})
	})
}

func (v *View) OnConsole(callback OnConsoleCallback) (stop func()) {

	v._onConsole.Register.Do(func() {
		var cb WkeConsoleCallback = func(_view WkeHandle, _param uintptr, _level WkeConsoleLevel, _message, _sourceName WkeString, _sourceLine uint32, _stackTrace WkeString) (voidRes uintptr) {
			level := int(_level)
			message := v.mb.GetString(_message)
			sourceName := v.mb.GetString(_sourceName)
			sourceLine := int(_sourceLine)
			stackTrace := v.mb.GetString(_stackTrace)

			for _, callback := range v._onConsole.Callbacks {
				callback(level, message, sourceName, sourceLine, stackTrace)
			}

			return 0
		}

		v.mb.CallFunc("wkeOnConsole", uintptr(v.Hwnd), CallbackToPtr(cb), 0)
	})

	key := utils.RandString(10)
	v._onConsole.Callbacks[key] = callback

	return func() {
		delete(v._onConsole.Callbacks, key)
	}
}

func (v *View) IsDocumentReady() bool {
	p, _, _ := v.mb.CallFunc("wkeIsDocumentReady", uintptr(v.Hwnd))
	return p != 0
}

// 阻塞等待文档加载完成，仅限主frame
func (v *View) WaitUntilDocumentReady(timeout time.Duration) bool {

	if v.IsDocumentReady() {
		return true
	}

	rst := make(chan bool, 2)

	select {
	case <-time.After(timeout):
		log.Error("等待文档加载超时")
		rst <- false
	default:
		stop := func() {}
		stop = v.OnDocumentReady(func(frame WkeWebFrameHandle) {
			if v.IsMainFrame(frame) {
				stop()
				rst <- true
			}
		})
	}

	return <-rst
}

func (v *View) OnTitleChanged(callback OnTitleChangedCallback) (stop func()) {

	v._onTitleChanged.Register.Do(func() {
		var cb WkeTitleChangedCallback = func(view WkeHandle, param uintptr, title WkeString) (voidRes uintptr) {
			_title := v.mb.GetString(title)

			for _, callback := range v._onTitleChanged.Callbacks {
				callback(_title)
			}
			return
		}

		v.mb.CallFunc("wkeOnTitleChanged", uintptr(v.Hwnd), CallbackToPtr(cb), 0)
	})

	key := utils.RandString(10)

	v._onTitleChanged.Callbacks[key] = callback

	return func() {
		delete(v._onTitleChanged.Callbacks, key)
	}
}

func (v *View) OnDownload(callback OnDownloadCallback) (stop func()) {

	v._onDownload.Register.Do(func() {
		var cb WkeDownloadCallback = func(view WkeHandle, param uintptr, url uintptr) (voidRes uintptr) {
			link := PtrToString(url)
			for _, callback := range v._onDownload.Callbacks {
				callback(link)
			}
			return
		}

		v.mb.CallFunc("wkeOnDownload", uintptr(v.Hwnd), CallbackToPtr(cb), 0)
	})

	key := utils.RandString(10)

	v._onDownload.Callbacks[key] = callback

	return func() {
		delete(v._onDownload.Callbacks, key)
	}
}

func (v *View) GetMainWebFrame() WkeWebFrameHandle {
	r1, _, _ := v.mb.CallFunc("wkeWebFrameGetMainFrame", uintptr(v.Hwnd))

	return WkeWebFrameHandle(r1)
}

func mm2px(mm float64, dpi int) int {
	return int(math.Round(float64(dpi) * mm / 25.4))
}

type PrintSettings struct {
	DPI          int
	Width        int // 单位 MM
	Height       int // 单位 MM
	MarginTop    int // 单位 MM
	MarginBottom int // 单位 MM
	MarginLeft   int // 单位 MM
	MarginRight  int // 单位 MM
}

type WithPrintSettings func(s *PrintSettings)

// 保存主 WebFrame 的内容到 PDF
func (v *View) SaveToPDF(writer io.Writer, withSetting ...WithPrintSettings) error {
	frameId := v.GetMainWebFrame()

	return v.SaveWebFrameToPDF(frameId, writer, withSetting...)
}

// 保存指定 WebFrame 的内容到 PDF
func (v *View) SaveWebFrameToPDF(frameId WkeWebFrameHandle, writer io.Writer, withSetting ...WithPrintSettings) error {

	// 默认为A4纸张，每边1厘米的边距，DPI为300
	s := PrintSettings{
		DPI:          300,
		Width:        210,
		Height:       297,
		MarginTop:    10,
		MarginBottom: 10,
		MarginLeft:   10,
		MarginRight:  10,
	}

	for _, withSet := range withSetting {
		withSet(&s)
	}

	// 假设A4纸张，每边1厘米的边距，DPI为600
	setting := wkePrintSettings{
		// structSize:               48, // 结构体大小，每个 int 为4, 12个int为48（极个别 C 编译器的int大小为8，暂不予考虑）
		dpi:                      int32(s.DPI),
		width:                    int32(mm2px(float64(s.Width), s.DPI)),     // 根据 DPI 将纸张宽度 mm 转换为像素 px
		height:                   int32(mm2px(float64(s.Height), s.DPI)),    // 根据 DPI 将纸张高度 mm 转换为像素 px
		marginTop:                int32(mm2px(float64(s.MarginTop), s.DPI)), // 根据 DPI 将纸张边距 mm 转换为像素 px
		marginBottom:             int32(mm2px(float64(s.MarginBottom), s.DPI)),
		marginLeft:               int32(mm2px(float64(s.MarginLeft), s.DPI)),
		marginRight:              int32(mm2px(float64(s.MarginRight), s.DPI)),
		isPrintPageHeadAndFooter: FALSE, // 是否打印页眉页脚
		isPrintBackgroud:         TRUE,  // 是否打印背景
		isLandscape:              FALSE, // 是否横向打印
		isPrintToMultiPage:       FALSE, // 是否打印到多页
	}

	setting.structSize = int32(unsafe.Sizeof(setting)) // 使用 unsafe 获取结构体大小，避免 C 编译器的不同

	if s.Width > s.Height {
		setting.isLandscape = TRUE // 宽大于高，则横向打印
	}

	// 调用 wkeUtilPrintToPdf 生成 PDF
	r1, _, err := v.mb.CallFunc("wkeUtilPrintToPdf", uintptr(v.Hwnd), uintptr(frameId), uintptr(unsafe.Pointer(&setting)))
	if r1 == 0 && err != nil {
		// err 为windows的最后一个错误，可能与打印无关。
		return err
	}

	// 释放内存
	defer v.mb.CallFuncAsync("wkeUtilRelasePrintPdfDatas", r1)

	pd := (*wkePdfDatas)(unsafe.Pointer(r1))

	if pd.count == 0 {
		return errors.New("生成 PDF 失败")
	}

	sizes := unsafe.Slice((*uintptr)(unsafe.Pointer(pd.sizes)), pd.count)
	datasPtrs := unsafe.Slice((**byte)(unsafe.Pointer(pd.datas)), pd.count)

	dataPtr := datasPtrs[0]
	size := sizes[0]

	chunk := unsafe.Slice(dataPtr, int(size))

	if _, err := writer.Write(chunk); err != nil {
		return err
	}

	return nil
}

func (v *View) SetHeadlessEnabled(enable bool) *CallFuncJob {
	return v.mb.CallFuncAsync("wkeSetHeadlessEnabled", uintptr(v.Hwnd), BoolToPtr(enable))
}

func (v *View) SetTransparent(transparent bool) {
	v.mb.CallFunc("wkeSetTransparent", uintptr(v.Hwnd), BoolToPtr(transparent))
}

func (v *View) OnOtherLoad(callback OnOtherLoadCallback) (stop func()) {
	v._onOtherLoad.Register.Do(func() {
		var cb WkeOnOtherLoadCallback = func(webView WkeHandle, param uintptr, loadType WkeOtherLoadType, info *WkeTempCallbackInfo) (voidRes uintptr) {
			for _, callback := range v._onOtherLoad.Callbacks {
				callback(loadType, info)
			}
			return
		}

		v.mb.CallFunc("wkeOnOtherLoad", uintptr(v.Hwnd), CallbackToPtr(cb), 0)
	})

	key := utils.RandString(10)
	v._onOtherLoad.Callbacks[key] = callback

	return func() {
		delete(v._onOtherLoad.Callbacks, key)
	}
}
