package blink

import (
	"io"
	"os"
	"runtime"
	"sync"
	"unsafe"

	"github.com/epkgs/mini-blink/internal/dll"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

var locker sync.RWMutex

type Blink struct {
	Config *Config
	JS     *JS

	Resource *ResourceLoader

	dll   *windows.DLL
	procs map[string]*windows.Proc

	views   map[WkeHandle]*View
	windows map[WkeHandle]*Window

	bootScripts []string
}

func loadDLL(conf *Config) *windows.DLL {

	// 尝试在默认目录里加载 DLL
	if loaded, err := windows.LoadDLL(DLL_FILE); err == nil {
		return loaded
	}

	fullPath := conf.GetDllFilePath()

	// 放入闭包，使其可以被释放
	func() {

		file, err := dll.FS.Open(DLL_FILE)
		if err != nil {
			panic("无法从默认路径或内嵌资源里找到 blink.dll，err: " + err.Error())
		}

		data, err := io.ReadAll(file)
		if err != nil {
			panic("读取内联DLL出错，err: " + err.Error())
		}

		newFile, err := os.Create(fullPath)
		if err != nil {
			panic("无法创建dll文件，err: " + err.Error())
		}

		defer newFile.Close()
		n, err := newFile.Write(data)
		if err != nil {
			panic("写入dll文件失败，err: " + err.Error())
		}
		if n != len(data) {
			panic("写入校验失败")
		}
	}()

	return windows.MustLoadDLL(fullPath)
}

func NewApp(setups ...func(*Config)) *Blink {

	config := NewConfig(setups...)

	blink := &Blink{
		Config:   config,
		Resource: NewResourceLoader(),

		dll:   loadDLL(config),
		procs: make(map[string]*windows.Proc),

		views:   make(map[WkeHandle]*View),
		windows: make(map[WkeHandle]*Window),
	}

	if !blink.IsInitialize() {
		blink.Initialize()
	}

	blink.JS = newJS(blink)

	return blink
}

func (mb *Blink) Free() {
	mb.Finalize()

	mb.dll.Release()
}

func (mb *Blink) GetViewByHandle(viewHwnd WkeHandle) *View {
	locker.Lock()
	view, exist := mb.views[viewHwnd]
	locker.Unlock()
	if !exist {
		return nil
	}
	return view
}

func (mb *Blink) GetWindowByHandle(windowHwnd WkeHandle) *Window {
	locker.Lock()
	window, exist := mb.windows[windowHwnd]
	locker.Unlock()
	if !exist {
		return nil
	}
	return window
}

// 返回 true 则代表已处理 msg，不需要继续执行
func (mb *Blink) DispatchMessage(msg *win.MSG) bool {

	view := mb.GetViewByHandle(WkeHandle(msg.HWnd))
	if view != nil {
		return view.DispatchMessage(msg)
	}

	window := mb.GetWindowByHandle(WkeHandle(msg.HWnd))
	if window != nil {
		return window.DispatchMessage(msg)
	}

	return false

}

func (mb *Blink) KeepRunning() {

	runtime.LockOSThread()
	defer func() {
		mb.Free()
		runtime.UnlockOSThread()
	}()

	msg := &win.MSG{}

	for win.GetMessage(msg, 0, 0, 0) > 0 {

		win.TranslateMessage(msg)

		if mb.DispatchMessage(msg) {
			continue
		}

		win.DispatchMessage(msg)
	}

}

func (mb *Blink) findProc(name string) *windows.Proc {
	proc, ok := mb.procs[name]
	if !ok {
		proc = mb.dll.MustFindProc(name)
		mb.procs[name] = proc
	}
	return proc
}

// ! 注意：args 的 GC。例如传入值是另一个 callback 的入参，那么需要确保此传入值直到 callback 调用时仍未被 GC 回收
func (mb *Blink) CallFunc(name string, args ...uintptr) (r1 uintptr, r2 uintptr, err error) {
	runtime.LockOSThread() // ! 由于 miniblink 的线程限制，需要锁定线程
	defer func() {
		if r := recover(); r != nil {

			if r == windows.NOERROR {
				err = nil
				return
			}

			err = r.(error)
			logError("Panic by CallFunc: %s", err.Error())
		}
	}()

	r1, r2, err = mb.findProc(name).Call(args...)

	if err == windows.NOERROR {
		err = nil
	}

	return
}

func (mb *Blink) Version() int {
	ver, _, _ := mb.CallFunc("wkeVersion")
	return int(ver)
}

func (mb *Blink) VersionString() string {
	ver, _, _ := mb.CallFunc("wkeVersionString")
	return PtrToString(ver)
}

func (mb *Blink) Initialize() {
	mb.CallFunc("wkeInitialize")
}

func (mb *Blink) Finalize() {
	mb.CallFunc("wkeFinalize")
}

func (mb *Blink) IsInitialize() bool {
	r1, _, _ := mb.CallFunc("wkeIsInitialize")

	return r1 != 0
}

func (mb *Blink) createWebWindow(winType WkeWindowType, parent *View, rectOptional ...WkeRect) *View {
	var pHwnd WkeHandle = 0
	if parent != nil {
		pHwnd = parent.Hwnd
	}

	rect := WkeRect{200, 200, 800, 600}
	if len(rectOptional) >= 1 {
		rect = rectOptional[0]
	}

	ptr, _, _ := mb.CallFunc("wkeCreateWebWindow", uintptr(winType), uintptr(pHwnd), uintptr(rect.X), uintptr(rect.Y), uintptr(rect.W), uintptr(rect.H))
	return NewView(mb, WkeHandle(ptr), winType, parent)

}

// 普通窗口
func (mb *Blink) CreateWebWindowPopup(rectOptional ...WkeRect) *View {
	return mb.createWebWindow(WKE_WINDOW_TYPE_POPUP, nil, rectOptional...)
}

// 透明窗口
func (mb *Blink) CreateWebWindowTransparent(rectOptional ...WkeRect) *View {
	return mb.createWebWindow(WKE_WINDOW_TYPE_TRANSPARENT, nil, rectOptional...)
}

// 嵌入在父窗口里的子窗口
func (mb *Blink) CreateWebWindowControl(parent *View, rectOptional ...WkeRect) *View {
	return mb.createWebWindow(WKE_WINDOW_TYPE_CONTROL, parent, rectOptional...)
}

// 设置response的mime
func (mb *Blink) NetSetMIMEType(job WkeNetJob, mimeType string) {
	mb.CallFunc("wkeNetSetMIMEType", uintptr(job), StringToPtr(mimeType))
}

// 获取response的mime
func (mb *Blink) NetGetMIMEType(job WkeNetJob, mime string) string {
	ptr, _, _ := mb.CallFunc("wkeNetGetMIMEType", uintptr(job), StringToPtr(mime))
	return PtrToString(ptr)
}

// 调用此函数后,网络层收到数据会存储在一buf内,接收数据完成后响应OnLoadUrlEnd事件.#此调用严重影响性能,慎用。
// 此函数和wkeNetSetData的区别是，wkeNetHookRequest会在接受到真正网络数据后再调用回调，并允许回调修改网络数据。
// 而wkeNetSetData是在网络数据还没发送的时候修改。
func (mb *Blink) NetSetData(job WkeNetJob, buf []byte) {
	if len(buf) == 0 {
		buf = []byte{0}
	}

	mb.CallFunc("wkeNetSetData", uintptr(job), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
}

func (mb *Blink) GetViewByJsExecState(es JsExecState) *View {
	handle := mb.JS.GetWebView(es)
	return mb.GetViewByHandle(handle)
}

func (mb *Blink) AddBootScript(script string) {
	mb.bootScripts = append(mb.bootScripts, script)
}

func (mb *Blink) GetString(str WkeString) string {
	p, _, _ := mb.CallFunc("wkeGetString", uintptr(str))

	return PtrToString(p)
}
