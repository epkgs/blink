package blink

import (
	"runtime"
	"sync"
	"unsafe"

	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

var locker sync.RWMutex

type BlinkJob struct {
	job  func()
	done chan bool
}

type Blink struct {
	Config *Config
	IPC    *IPC

	js *JS

	Resource   *ResourceLoader
	Downloader *Downloader

	dll   *windows.DLL
	procs map[string]*windows.Proc

	views   map[WkeHandle]*View
	windows map[WkeHandle]*Window

	bootScripts []string

	quit chan bool
	jobs chan BlinkJob
}

func NewApp(setups ...func(*Config)) *Blink {

	config := NewConfig(setups...)

	dll, err := loadDLL(config)

	if err != nil {
		MessageBoxError(0, err.Error())
		panic(err)
	}

	blink := &Blink{
		Config:     config,
		Resource:   NewResourceLoader(),
		Downloader: NewDownloader(3),

		dll:   dll,
		procs: make(map[string]*windows.Proc),

		views:   make(map[WkeHandle]*View),
		windows: make(map[WkeHandle]*Window),

		quit: make(chan bool),
		jobs: make(chan BlinkJob, 20),
	}

	if !blink.IsInitialize() {
		blink.Initialize()
	}

	blink.js = newJS(blink)

	blink.IPC = newIPC(blink)

	return blink
}

func (mb *Blink) Exit() {
	close(mb.quit)
}

func (mb *Blink) Free() {

	for _, v := range mb.views {
		v.DestroyWindow()
	}

	close(mb.quit)

	mb.Finalize()
	mb.dll.Release()
}

func (mb *Blink) GetViews() []*View {
	var views []*View

	for _, v := range mb.views {
		views = append(views, v)
	}

	return views
}

func (mb *Blink) GetFirstView() (view *View) {
	for _, view = range mb.views {
		break
	}
	return
}

func (mb *Blink) GetViewByHandle(viewHwnd WkeHandle) *View {
	locker.Lock()
	defer locker.Unlock()

	view, exist := mb.views[viewHwnd]
	if !exist {
		return nil
	}
	return view
}

func (mb *Blink) GetWindowByHandle(windowHwnd WkeHandle) *Window {
	locker.Lock()
	defer locker.Unlock()

	window, exist := mb.windows[windowHwnd]
	if !exist {
		return nil
	}
	return window
}

func (mb *Blink) AddJob(job func()) chan bool {
	done := make(chan bool, 1)
	mb.jobs <- BlinkJob{
		job,
		done,
	}

	return done
}

func (mb *Blink) KeepRunning() {

	runtime.LockOSThread()

	msg := &win.MSG{}

	for {
		select {
		case <-mb.quit:
			return

		case bj := <-mb.jobs:
			logInfo("received job")
			bj.job()
			close(bj.done)

		default:

			if win.GetMessage(msg, 0, 0, 0) <= 0 {
				return
			}

			win.TranslateMessage(msg)

			win.DispatchMessage(msg)
		}
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
	handle := mb.js.GetWebView(es)
	return mb.GetViewByHandle(handle)
}

func (mb *Blink) AddBootScript(script string) {
	mb.bootScripts = append(mb.bootScripts, script)
}

func (mb *Blink) GetString(str WkeString) string {
	p, _, _ := mb.CallFunc("wkeGetString", uintptr(str))

	return PtrToString(p)
}
