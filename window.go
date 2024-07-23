package blink

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"syscall"
	"unsafe"

	"github.com/epkgs/blink/internal/log"
	"github.com/epkgs/blink/internal/utils"
	"github.com/lxn/win"
)

type WM_SIZING uintptr

const (
	WMSZ_LEFT        WM_SIZING = 1 // 左边缘
	WMSZ_RIGHT       WM_SIZING = 2 // 右边缘
	WMSZ_TOP         WM_SIZING = 3 // 上边缘
	WMSZ_TOPLEFT     WM_SIZING = 4 // 左上角
	WMSZ_TOPRIGHT    WM_SIZING = 5 // 右上角
	WMSZ_BOTTOM      WM_SIZING = 6 // 下边缘
	WMSZ_BOTTOMLEFT  WM_SIZING = 7 // 左下角
	WMSZ_BOTTOMRIGHT WM_SIZING = 8 // 右下角
)

type SIZE_TYPE uintptr

const (
	SIZE_RESTORED SIZE_TYPE = iota
	SIZE_MINIMIZED
	SIZE_MAXIMIZED
	SIZE_MAXSHOW
	SIZE_MAXHIDE
)

const (
	WM_USER       = win.WM_USER
	WM_TRAYNOTIFY = WM_USER + 1

	ID_TRAY             = WM_USER + 100
	ID_TRAYMENU_RESTORE = WM_USER + 101
	ID_TRAYMENU_EXIT    = WM_USER + 102
)

var (
	user32     = syscall.NewLazyDLL("user32.dll")
	appendMenu = user32.NewProc("AppendMenuW")
)

type WindowOnSizingCallback func(WM_SIZING, *win.RECT)
type WindowOnSizeCallback func(stype SIZE_TYPE, width, height uint16)
type WindowOnCreateCallback func(*win.CREATESTRUCT)
type WindowOnActivateAppCallback func(bool, uintptr)

type Window struct {
	mb   *Blink
	view *View
	Hwnd WkeHandle

	windowType WkeWindowType

	sizing                WM_SIZING
	borderResizeEnabled   bool
	borderResizeThickness int32 // border 反应厚度

	isMaximized bool
	fixedTitle  bool

	_oldWndProc uintptr

	iconHandle win.HANDLE

	nid               win.NOTIFYICONDATA
	useSimpleTrayMenu bool

	_onSizing      *bindEvent[WindowOnSizingCallback]
	_onSize        *bindEvent[WindowOnSizeCallback]
	_onCreate      *bindEvent[WindowOnCreateCallback]
	_onActivateApp *bindEvent[WindowOnActivateAppCallback]
}

func newWindow(mb *Blink, view *View, windowType WkeWindowType) *Window {
	window := &Window{
		mb:         mb,
		view:       view,
		windowType: windowType,
		Hwnd:       view.GetWindowHandle(),

		_onSizing:      newBindEvent[WindowOnSizingCallback](),
		_onSize:        newBindEvent[WindowOnSizeCallback](),
		_onCreate:      newBindEvent[WindowOnCreateCallback](),
		_onActivateApp: newBindEvent[WindowOnActivateAppCallback](),
	}

	window._oldWndProc = win.SetWindowLongPtr(win.HWND(window.Hwnd), win.GWL_WNDPROC, CallbackToPtr(window.hookWindowProc))

	window.view.OnTitleChanged(func(title string) {
		if window.fixedTitle {
			return
		}
		window.setTitle(title)
	})

	if window.windowType == WKE_WINDOW_TYPE_TRANSPARENT {
		window.EnableBorderResize(true)
	}

	// 监听尺寸大小变化，修改 isMaximized 状态
	window.OnSize(func(stype SIZE_TYPE, width, height uint16) {
		if stype == SIZE_MAXIMIZED {
			window.isMaximized = true
		} else if stype == SIZE_RESTORED || stype == SIZE_MINIMIZED {
			window.isMaximized = false
		}
	})

	return window
}

func (w *Window) hookWindowProc(hwnd, message, wparam, lparam uintptr) uintptr {

	handled := func() bool {
		switch message {
		case win.WM_ENTERSIZEMOVE:
			w.udpateCursor()
		case win.WM_GETMINMAXINFO:
			lpmmi := (*win.MINMAXINFO)(unsafe.Pointer(lparam))

			// 修正无边框窗口，最大化时的尺寸问题，避免遮挡任务栏
			if w.isMaximized {
				if w.windowType == WKE_WINDOW_TYPE_TRANSPARENT || w.windowType == WKE_WINDOW_TYPE_HIDE_CAPTION {
					// 获取窗口所在屏幕的句柄
					hMonitor := win.MonitorFromWindow(win.HWND(w.Hwnd), win.MONITOR_DEFAULTTONEAREST)

					// 获取屏幕工作区尺寸
					var monitorInfo win.MONITORINFO
					monitorInfo.CbSize = uint32(unsafe.Sizeof(monitorInfo))
					win.GetMonitorInfo(hMonitor, &monitorInfo)

					// 根据屏幕工作区大小，设置窗口尺寸
					lpmmi.PtMaxPosition.X = monitorInfo.RcWork.Left
					lpmmi.PtMaxPosition.Y = monitorInfo.RcWork.Top
					lpmmi.PtMaxSize.X = monitorInfo.RcWork.Right - monitorInfo.RcWork.Left
					lpmmi.PtMaxSize.Y = monitorInfo.RcWork.Bottom - monitorInfo.RcWork.Top
				}
			}

		case WM_TRAYNOTIFY:

			if lparam == win.WM_LBUTTONDBLCLK {
				log.Debug("Tray icon double clicked")
				w.Restore()
				return true
			}

			if lparam == win.WM_RBUTTONUP {
				log.Debug("Right click tray icon")
				if w.useSimpleTrayMenu {
					w.showSimpleTrayMenu()
					return true
				}
			}
		case win.WM_COMMAND:
			// 处理菜单点击事件
			menuID := LOWORD(uint32(wparam))
			switch menuID {
			case ID_TRAYMENU_RESTORE:
				log.Debug("Restore menu item clicked")
				w.Restore()
				return true
			case ID_TRAYMENU_EXIT:
				log.Debug("Exit menu item clicked")
				// w.Hide()
				w.Destroy()
				return true
			}

		case win.WM_MOUSEMOVE:

			// 子控件不处理 sizing 事件，转交给父窗口处理
			if w.windowType == WKE_WINDOW_TYPE_CONTROL && w.view.parent != nil {
				w.view.parent.Window.calcSizing()
			} else {
				w.calcSizing()
			}

			return false // 可能还有其他事件

		case win.WM_LBUTTONDOWN:
			if wparam != win.MK_LBUTTON {
				return false
			}

			// 子控件不处理 sizing 事件，转交给父窗口处理
			if w.windowType == WKE_WINDOW_TYPE_CONTROL && w.view.parent != nil {
				return w.view.parent.Window.triggerSizing()
			} else {
				return w.triggerSizing()
			}

		case win.WM_SIZING:

			pos := (WM_SIZING)(wparam)

			rect := (*win.RECT)(unsafe.Pointer(lparam))

			for _, cb := range w._onSizing.Callbacks {
				cb(pos, rect)
			}

		case win.WM_SIZE:
			stype := (SIZE_TYPE)(wparam)
			width := LOWORD(uint32(lparam))
			height := HIWORD(uint32(lparam))
			for _, cb := range w._onSize.Callbacks {
				cb(stype, width, height)
			}

		case win.WM_CREATE:

			log.Debug("trigger WM_CREATE")
			created := (*win.CREATESTRUCT)(unsafe.Pointer(lparam))

			for _, cb := range w._onCreate.Callbacks {
				cb(created)
			}

		case win.WM_ACTIVATEAPP:
			actived := wparam == 1

			for _, cb := range w._onActivateApp.Callbacks {
				cb(actived, lparam)
			}

		}

		return false
	}()

	if handled {

		return 1
	}

	return win.CallWindowProc(uintptr(w._oldWndProc), win.HWND(hwnd), uint32(message), wparam, lparam)
}

func (w *Window) calcSizing() bool {

	if !w.borderResizeEnabled {
		return false
	}

	pt := &win.POINT{}
	win.GetCursorPos(pt)

	rect := &win.RECT{}
	win.GetWindowRect(win.HWND(w.Hwnd), rect)

	inLeft := pt.X >= rect.Left && pt.X <= rect.Left+w.borderResizeThickness
	inReght := pt.X <= rect.Right && pt.X >= rect.Right-w.borderResizeThickness
	inTop := pt.Y >= rect.Top && pt.Y <= rect.Top+w.borderResizeThickness
	inBottom := pt.Y <= rect.Bottom && pt.Y >= rect.Bottom-w.borderResizeThickness

	if inLeft && inTop { // 左上角
		w.sizing = WMSZ_TOPLEFT
	} else if inLeft && inBottom { // 左下角
		w.sizing = WMSZ_BOTTOMLEFT
	} else if inReght && inTop { // 右上角
		w.sizing = WMSZ_TOPRIGHT
	} else if inReght && inBottom { // 右下角
		w.sizing = WMSZ_BOTTOMRIGHT
	} else if inLeft { // 左边
		w.sizing = WMSZ_LEFT
	} else if inReght { // 右边
		w.sizing = WMSZ_RIGHT
	} else if inTop { // 上边
		w.sizing = WMSZ_TOP
	} else if inBottom { // 下边
		w.sizing = WMSZ_BOTTOM
	} else {
		// 鼠标没有在窗口边缘
		w.sizing = 0
		return false
	}

	// win.SendMessage(win.HWND(w.Hwnd), win.WM_SYSCOMMAND, uintptr(win.WM_SETCURSOR|w.sizing), win.WM_MOUSEMOVE)
	w.udpateCursor()
	return true
}

func (w *Window) triggerSizing() bool {
	if w.sizing <= 0 {
		return false
	}
	win.SendMessage(win.HWND(w.Hwnd), win.WM_SYSCOMMAND, uintptr(win.SC_SIZE|w.sizing), 0)
	return true
}

func (w *Window) udpateCursor() bool {
	switch w.sizing {
	case WMSZ_LEFT, WMSZ_RIGHT:
		win.SetCursor(win.LoadCursor(0, AssertType[uint16](win.IDC_SIZEWE)))
		return true
	case WMSZ_TOP, WMSZ_BOTTOM:
		win.SetCursor(win.LoadCursor(0, AssertType[uint16](win.IDC_SIZENS)))
		return true
	case WMSZ_TOPLEFT, WMSZ_BOTTOMRIGHT:
		win.SetCursor(win.LoadCursor(0, AssertType[uint16](win.IDC_SIZENWSE)))
		return true
	case WMSZ_BOTTOMLEFT, WMSZ_TOPRIGHT:
		win.SetCursor(win.LoadCursor(0, AssertType[uint16](win.IDC_SIZENESW)))
		return true
	}
	return false
}

func (w *Window) Show() {
	w.mb.AddJob(func() {
		// win.ShowWindow(win.HWND(w.Hwnd), win.SW_SHOW)
		win.SetWindowPos(win.HWND(w.Hwnd), win.HWND_TOP, 0, 0, 0, 0, win.SWP_NOMOVE|win.SWP_NOSIZE|win.SWP_SHOWWINDOW)
	})
	// w.view.mb.CallFunc("wkeShowWindow", uintptr(w.view.Hwnd), 1)
}

func (w *Window) Hide() {
	w.mb.AddJob(func() {
		win.ShowWindow(win.HWND(w.Hwnd), win.SW_HIDE)
	})
	// w.view.mb.CallFunc("wkeShowWindow", uintptr(w.view.Hwnd), 0)
}

func (w *Window) Close() {
	w.mb.AddJob(func() {
		win.PostMessage(win.HWND(w.Hwnd), win.WM_CLOSE, 0, 0)
	})
}

func (w *Window) Destroy() {
	// w.mb.AddJob(func() {
	// 	win.DestroyWindow(win.HWND(w.Hwnd))
	// })
	w.view.DestroyWindow()
}

func (w *Window) MinimizeToTray() {
	w.mb.AddJob(func() {
		win.ShowWindow(win.HWND(w.Hwnd), win.SW_HIDE)
	})
}

func (w *Window) Minimize() {
	w.mb.AddJob(func() {
		win.SendMessage(win.HWND(w.Hwnd), win.WM_SYSCOMMAND, uintptr(win.SC_MINIMIZE), 0)
	})
}

func (w *Window) Maximize() {
	w.mb.AddJob(func() {
		w.isMaximized = true
		win.SendMessage(win.HWND(w.Hwnd), win.WM_SYSCOMMAND, uintptr(win.SC_MAXIMIZE), 0)
	})
}

func (w *Window) IsMaximized() bool {
	// style := win.GetWindowLong(win.HWND(w.Hwnd), win.GWL_STYLE)
	// return (style & win.WS_MAXIMIZE) != 0
	return w.isMaximized
}

func (w *Window) Restore() {
	w.mb.AddJob(func() {
		w.isMaximized = false
		win.SendMessage(win.HWND(w.Hwnd), win.WM_SYSCOMMAND, uintptr(win.SC_RESTORE), 0)
	})
}

func (w *Window) EnableDragging() {
	w.mb.AddJob(func() {
		win.ReleaseCapture() // 松开鼠标控制
		go win.SendMessage(win.HWND(w.Hwnd), win.WM_SYSCOMMAND, uintptr(win.SC_MOVE|win.HTCAPTION), 0)
	})
}

func (w *Window) CloseAsHideTray() {
	w.EnableTray()
	w.view.OnClosing(func() bool {
		w.MinimizeToTray()
		return false
	})
	w.useSimpleTrayMenu = true
}

func (w *Window) EnableTray(setups ...func(*win.NOTIFYICONDATA)) {
	w.nid = win.NOTIFYICONDATA{
		HWnd:             win.HWND(w.Hwnd),
		UID:              ID_TRAY,
		UFlags:           win.NIF_ICON | win.NIF_MESSAGE | win.NIF_TIP,
		UCallbackMessage: WM_TRAYNOTIFY,
		HIcon:            win.HICON(w.iconHandle),
		CbSize:           uint32(unsafe.Sizeof(w.nid)),
	}
	copy(w.nid.SzTip[:], StringToU16Arr("双击打开窗口"))

	for _, setup := range setups {
		setup(&w.nid)
	}

	win.Shell_NotifyIcon(win.NIM_ADD, &w.nid)

	w.view.OnDestroy(func() {
		log.Debug("RemoveTray in view OnDestroy event")
		w.RemoveTray()
	})
}

func (w *Window) RemoveTray() {
	w.mb.AddJob(func() {
		win.Shell_NotifyIcon(win.NIM_DELETE, &w.nid)
	})

}

// TODO: 抽离代码，使其更通用，可以在外部添加修改菜单
func (w *Window) showSimpleTrayMenu() {

	w.mb.AddJob(func() {
		// 创建托盘菜单
		hMenu := win.CreatePopupMenu()
		if hMenu == 0 {
			return
		}
		defer win.DestroyMenu(hMenu)

		appendMenu.Call(uintptr(hMenu), win.MF_STRING, ID_TRAYMENU_RESTORE, StringToWCharPtr("显示窗口"))
		appendMenu.Call(uintptr(hMenu), win.MF_STRING, ID_TRAYMENU_EXIT, StringToWCharPtr("退出"))

		var pt win.POINT
		win.GetCursorPos(&pt)
		win.SetMenuDefaultItem(hMenu, ID_TRAYMENU_RESTORE, false)
		win.TrackPopupMenu(hMenu, win.TPM_LEFTALIGN|win.TPM_RIGHTBUTTON, pt.X, pt.Y, 0, win.HWND(w.Hwnd), nil)
	})

}

func (w *Window) SetIcon(handle win.HANDLE) error {
	if handle == 0 {
		return errors.New("获取图标句柄失败，无法设置 ICON 。")
	}

	w.iconHandle = handle

	w.mb.AddJob(func() {

		win.SendMessage(win.HWND(w.Hwnd), win.WM_SETICON, win.IMAGE_ICON, uintptr(handle))
		// win.SendMessage(win.HWND(w.Hwnd), win.WM_SETICON, win.IMAGE_BITMAP, uintptr(handle))
	})

	return nil
}

// 设置窗口图标(从图标文件中). 快捷方法
func (w *Window) SetIconFromFile(iconFilePath string) error {
	iconHandle, err := w.loadIconFromFile(iconFilePath)
	if err != nil {
		return err
	}

	return w.SetIcon(iconHandle)
}

// 设置窗口图标(从图标二进制数据中). 快捷方法
func (w *Window) SetIconFromBytes(iconData []byte) error {
	iconHandle, err := w.loadIconFromBytes(iconData)
	if err != nil {
		return err
	}
	return w.SetIcon(iconHandle)
}

// 数据hash > icon handle的缓存映射
var iconCache = sync.Map{}

// 从二进制数组中加载icon
// TODO:目前是先把ico二进制数据存到本地,再使用winapi的LoadImage加载图标,因为暂未找到直接从内存中加载ico文件的方法
func (w *Window) loadIconFromBytes(iconData []byte) (iconHandle win.HANDLE, err error) {
	//计算数据的hash
	bh := md5.Sum(iconData)
	dataHash := hex.EncodeToString(bh[:])

	//先判断缓存里面有没有
	if handle, isExist := iconCache.Load(dataHash); isExist {
		return handle.(win.HANDLE), nil
	}

	//缓存中没有,则释放到本地目录
	iconFilePath := filepath.Join(w.mb.Config.tempPath, "icon_"+dataHash+".ico")
	if _, err := os.Stat(iconFilePath); os.IsNotExist(err) {
		if err := os.WriteFile(iconFilePath, iconData, 0644); err != nil {
			return 0, errors.New("无法创建临时icon文件: " + err.Error())
		}
	}

	//从文件中加载
	handle, err := w.loadIconFromFile(iconFilePath)
	if err != nil {
		return 0, err
	}
	//存入缓存
	iconCache.Store(dataHash, handle)
	//返回结果
	return handle, nil
}

// 从文件中加载icon
// 注意：仅支持ico文件
func (w *Window) loadIconFromFile(iconFilePath string) (iconHandle win.HANDLE, err error) {
	iconFilePathW, err := syscall.UTF16PtrFromString(iconFilePath)
	if err != nil {
		return
	}
	iconHandle = win.LoadImage(
		0,
		iconFilePathW,
		win.IMAGE_ICON,
		0,
		0,
		win.LR_LOADFROMFILE,
	)
	if iconHandle == 0 {
		return 0, errors.New("加载图标文件失败," + iconFilePath)
	}
	return
}

func (w *Window) MoveToCenter() {
	w.mb.CallFunc("wkeMoveToCenter", uintptr(w.view.Hwnd))
}

func (w *Window) SetTitle(title string) {
	w.fixedTitle = true
	w.setTitle(title)
}
func (w *Window) setTitle(title string) {
	w.mb.CallFunc("wkeSetWindowTitle", uintptr(w.view.Hwnd), StringToPtr(title))
}

// 开启边缘拖动大小功能
func (w *Window) EnableBorderResize(enable bool, thicknessOpional ...int32) {
	w.borderResizeEnabled = enable
	if len(thicknessOpional) > 0 {
		w.borderResizeThickness = thicknessOpional[0]
	} else {
		w.borderResizeThickness = 5 // 默认的拖动边缘宽度
	}
}

// 隐藏标题栏
func (w *Window) HideCaption() {
	var stop func()
	stop = w.OnActivateApp(func(b bool, u uintptr) {
		stop()
		istyle := win.GetWindowLong(win.HWND(w.Hwnd), win.GWL_STYLE)
		win.SetWindowLong(win.HWND(w.Hwnd), win.GWL_STYLE, istyle&^(win.WS_CAPTION|win.WS_THICKFRAME))
		// win.SetWindowLong(win.HWND(w.Hwnd), win.GWL_STYLE, istyle&^(win.WS_CAPTION))
		win.SetWindowPos(win.HWND(w.Hwnd), 0, 0, 0, 0, 0, win.SWP_NOSIZE|win.SWP_NOMOVE|win.SWP_NOZORDER|win.SWP_NOACTIVATE|win.SWP_FRAMECHANGED)

		w.windowType = WKE_WINDOW_TYPE_HIDE_CAPTION
	})

}

// 显示标题栏
func (w *Window) ShowCaption() {
	istyle := win.GetWindowLong(win.HWND(w.Hwnd), win.GWL_STYLE)
	win.SetWindowLong(win.HWND(w.Hwnd), win.GWL_STYLE, istyle|(win.WS_CAPTION|win.WS_THICKFRAME))
	win.SetWindowPos(win.HWND(w.Hwnd), 0, 0, 0, 0, 0, win.SWP_NOSIZE|win.SWP_NOMOVE|win.SWP_NOZORDER|win.SWP_NOACTIVATE|win.SWP_FRAMECHANGED)
}

// 修改窗口大小
func (w *Window) SetWindowPos(x, y, width, height int32) {
	win.SetWindowPos(win.HWND(w.Hwnd), 0, x, y, width, height, win.SWP_NOZORDER|win.SWP_NOACTIVATE)
}

func (w *Window) Resize(width, height int32) {
	win.SetWindowPos(win.HWND(w.Hwnd), 0, 0, 0, width, height, win.SWP_NOMOVE|win.SWP_NOZORDER|win.SWP_NOACTIVATE)
}

func (w *Window) Move(x, y int32) {
	win.SetWindowPos(win.HWND(w.Hwnd), 0, x, y, 0, 0, win.SWP_NOSIZE|win.SWP_NOZORDER|win.SWP_NOACTIVATE)
}

// Sizing 事件
func (w *Window) OnSizing(callback WindowOnSizingCallback) (stop func()) {
	key := utils.RandString(6)
	w._onSizing.Callbacks[key] = callback

	return func() {
		delete(w._onSizing.Callbacks, key)
	}
}

// Size 事件
func (w *Window) OnSize(callback WindowOnSizeCallback) (stop func()) {
	key := utils.RandString(6)
	w._onSize.Callbacks[key] = callback

	return func() {
		delete(w._onSize.Callbacks, key)
	}
}

// Create 事件
func (w *Window) OnCreate(callback WindowOnCreateCallback) (stop func()) {
	key := utils.RandString(6)
	w._onCreate.Callbacks[key] = callback
	return func() {
		delete(w._onCreate.Callbacks, key)
	}
}

// Active APP 事件
func (w *Window) OnActivateApp(callback WindowOnActivateAppCallback) (stop func()) {
	key := utils.RandString(6)
	w._onActivateApp.Callbacks[key] = callback
	return func() {
		delete(w._onActivateApp.Callbacks, key)
	}
}
