package blink

import (
	"runtime"
	"sync"
	"unsafe"

	"github.com/epkgs/blink/internal/log"
	"github.com/epkgs/blink/internal/miniblink"
	"github.com/epkgs/blink/pkg/alert"
	"github.com/epkgs/blink/pkg/downloader"
	"github.com/epkgs/blink/pkg/queue"
	"github.com/epkgs/blink/pkg/resource"
	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

var locker sync.RWMutex

type BlinkJob struct {
	job  func()
	done chan bool
}

type CallFuncJob struct {
	funcName string
	args     []uintptr
	result   chan CallFuncResult
}

func (cj *CallFuncJob) Wait() (r1 uintptr, r2 uintptr, err error) {
	result := <-cj.result

	return result.R1, result.R2, result.Err
}

type CallFuncResult struct {
	R1  uintptr
	R2  uintptr
	Err error
}

type Blink struct {
	Config *Config
	IPC    *IPC

	js *JS

	Resource   *resource.Resource
	Downloader *downloader.Downloader

	dll   *windows.DLL
	procs map[string]*windows.Proc

	views   map[WkeHandle]*View
	windows map[WkeHandle]*Window

	bootScripts []string

	threadID uint32 // 调用 mb api 的线程 id

	quit     chan struct{}
	jobs     chan BlinkJob
	calls    *queue.Queue[CallFuncJob]
	jobLoops []func()
}

func NewApp(setups ...func(*Config)) *Blink {
	config, err := NewConfig(setups...)
	if err != nil {
		log.Error("NewConfig ERR: %v", err)
		alert.Error(err.Error())
		panic(err)
	}

	dll, err := miniblink.LoadDLL(config.GetDllFile(), config.GetTempPath())
	if err != nil {
		log.Error("loadDLL ERR: %v", err)
		alert.Error(err.Error())
		panic(err)
	}

	down := downloader.New(func(o *downloader.Option) {
		o.EnableSaveFileDialog = true
	})

	blink := &Blink{
		Config:     config,
		Resource:   resource.New(),
		Downloader: down,

		dll:   dll,
		procs: make(map[string]*windows.Proc),

		views:   make(map[WkeHandle]*View),
		windows: make(map[WkeHandle]*Window),

		quit:     make(chan struct{}),
		jobs:     make(chan BlinkJob, 20),
		calls:    queue.NewQueue[CallFuncJob](999),
		jobLoops: []func(){},
	}

	// 启动任务循环
	blink.loopJobLoops()

	if !blink.isInitialize() {
		blink.initialize()
	}

	blink.js = newJS(blink)

	blink.IPC = newIPC(blink)

	return blink
}

func (mb *Blink) CloseAll() {
	for _, v := range mb.views {
		v.DestroyWindow()
	}
}

func (mb *Blink) Exit() {
	mb.Free()
}

func (mb *Blink) Free() {

	mb.CloseAll()

	mb.finalize()

	close(mb.quit) // quit 退出任务循环，须等待上面任务完成才能退出任务循环

	mb.dll.Release()
	mb = nil
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

func (mb *Blink) GetViewByHandle(viewHwnd WkeHandle) (view *View, exist bool) {
	locker.Lock()
	defer locker.Unlock()

	view, exist = mb.views[viewHwnd]
	return
}

func (mb *Blink) GetWindowByHandle(windowHwnd WkeHandle) (window *Window, exist bool) {
	locker.Lock()
	defer locker.Unlock()

	window, exist = mb.windows[windowHwnd]
	return
}

var winMsgOnce sync.Once

func (mb *Blink) LoopWinMessage() {
	winMsgOnce.Do(func() {

		msg := &win.MSG{}

		mb.AddLoop(func() {

			if win.GetMessage(msg, 0, 0, 0) <= 0 {
				return
			}

			win.TranslateMessage(msg)

			win.DispatchMessage(msg)

		})
	})
}

func (mb *Blink) KeepRunning() {

	mb.LoopWinMessage()

	<-mb.quit
}

func (mb *Blink) findProc(name string) *windows.Proc {
	proc, ok := mb.procs[name]
	if !ok {
		proc = mb.dll.MustFindProc(name)
		mb.procs[name] = proc
	}
	return proc
}

func (mb *Blink) CallFunc(funcName string, args ...uintptr) (r1 uintptr, r2 uintptr, err error) {

	threadID := windows.GetCurrentThreadId()

	// 如果和调用 MB 的线程不一致，则塞入 chan 队列，等待执行
	if mb.threadID != threadID {
		return mb.CallFuncAsync(funcName, args...).Wait()
	}

	// 一致，则直接执行
	return mb.doCallFunc(funcName, args...)
}

func (mb *Blink) CallFuncFirst(funcName string, args ...uintptr) (r1 uintptr, r2 uintptr, err error) {

	threadID := windows.GetCurrentThreadId()

	// 如果和调用 MB 的线程不一致，则塞入 chan 队列，等待执行
	if mb.threadID != threadID {
		return mb.CallFuncAsyncFirst(funcName, args...).Wait()
	}

	// 一致，则直接执行
	return mb.doCallFunc(funcName, args...)
}

func (mb *Blink) CallFuncAsync(funcName string, args ...uintptr) *CallFuncJob {

	job := CallFuncJob{
		funcName: funcName,
		args:     args,
		result:   make(chan CallFuncResult, 1),
	}
	mb.calls.AddLast(job)

	return &job
}

func (mb *Blink) CallFuncAsyncFirst(funcName string, args ...uintptr) *CallFuncJob {

	job := CallFuncJob{
		funcName: funcName,
		args:     args,
		result:   make(chan CallFuncResult, 1),
	}
	mb.calls.AddFirst(job)

	return &job
}

// 将单个任务塞入队列，仅执行一次
func (mb *Blink) AddJob(job func()) chan bool {
	done := make(chan bool, 1)
	mb.jobs <- BlinkJob{
		job,
		done,
	}

	return done
}

// 增加任务到循环队列，每次循环都会执行
func (mb *Blink) AddLoop(job ...func()) *Blink {
	mb.jobLoops = append(mb.jobLoops, job...)
	return mb
}

func (mb *Blink) loopJobLoops() {
	go func() {

		runtime.LockOSThread() // ! 由于 miniblink 的线程限制，需要锁定线程

		mb.threadID = windows.GetCurrentThreadId()

		for {
			select {
			// 退出信号
			case <-mb.quit:
				return

				// 任务
			case bj := <-mb.jobs:
				bj.job()
				close(bj.done)

				// 调用 mb api 接口的异步任务
			case ch := <-mb.calls.Chan():
				job := ch.First()
				r1, r2, err := mb.doCallFunc(job.funcName, job.args...)
				job.result <- CallFuncResult{
					R1:  r1,
					R2:  r2,
					Err: err,
				}

			default:

				// 执行剩余队列
				for _, queue := range mb.jobLoops {
					queue()
				}

			}
		}
	}()
}

func (mb *Blink) doCallFunc(name string, args ...uintptr) (r1 uintptr, r2 uintptr, err error) {
	defer func() {
		if r := recover(); r != nil {

			if r == windows.NOERROR {
				err = nil
				return
			}

			err = r.(error)
			log.Error("Panic by CallFunc: %s", err.Error())
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

func (mb *Blink) initialize() {
	mb.CallFunc("wkeInitialize")
}

func (mb *Blink) finalize() {
	mb.CallFunc("wkeFinalize")
}

func (mb *Blink) isInitialize() bool {
	r1, _, _ := mb.CallFunc("wkeIsInitialize")

	return r1 != 0
}

type WebWindowConfig struct {
	WkeRect
}

type WithWebWindowConfig func(c *WebWindowConfig)

// 设置窗口大小
func WithWebWindowSize(w, h int32) WithWebWindowConfig {
	return func(config *WebWindowConfig) {
		config.W = w
		config.H = h
	}
}

// 设置窗口位置
func WithWebWindowPos(x, y int32) WithWebWindowConfig {
	return func(config *WebWindowConfig) {
		config.X = x
		config.Y = y
	}
}

func (mb *Blink) createWebWindow(winType WkeWindowType, parent *View, withConfig ...WithWebWindowConfig) *View {
	var pHwnd WkeHandle = 0
	if parent != nil {
		pHwnd = parent.Window.Hwnd
	}

	conf := WebWindowConfig{
		WkeRect{200, 200, 800, 600},
	}
	for _, set := range withConfig {
		set(&conf)
	}

	ptr, _, _ := mb.CallFunc("wkeCreateWebWindow", uintptr(winType), uintptr(pHwnd), uintptr(conf.X), uintptr(conf.Y), uintptr(conf.W), uintptr(conf.H))
	return NewView(mb, WkeHandle(ptr), winType, parent)

}

// 普通窗口
func (mb *Blink) CreateWebWindowPopup(withConfig ...WithWebWindowConfig) *View {
	return mb.createWebWindow(WKE_WINDOW_TYPE_POPUP, nil, withConfig...)
}

// 透明窗口
func (mb *Blink) CreateWebWindowTransparent(withConfig ...WithWebWindowConfig) *View {
	return mb.createWebWindow(WKE_WINDOW_TYPE_TRANSPARENT, nil, withConfig...)
}

// 嵌入在父窗口里的子窗口
func (mb *Blink) CreateWebWindowControl(parent *View, withConfig ...WithWebWindowConfig) *View {
	return mb.createWebWindow(WKE_WINDOW_TYPE_CONTROL, parent, withConfig...)
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

func (mb *Blink) NetHookRequest(job WkeNetJob) {
	mb.CallFunc("wkeNetHookRequest", uintptr(job))
}

func (mb *Blink) GetViewByJsExecState(es JsExecState) (view *View, exist bool) {
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
