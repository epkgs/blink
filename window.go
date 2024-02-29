package blink

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/lxn/win"
)

type WM_SIZING uint16

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

type Window struct {
	mb   *Blink
	view *View
	Hwnd WkeHandle

	windowType WkeWindowType

	enableBorderResize bool
	isMaximized        bool
	fixedTitle         bool

	sizing WM_SIZING

	_oldWndProc uintptr
}

func newWindow(mb *Blink, view *View, windowType WkeWindowType) *Window {
	window := &Window{
		mb:         mb,
		view:       view,
		windowType: windowType,
		Hwnd:       view.GetWindowHandle(),
	}

	window._oldWndProc = win.SetWindowLongPtr(win.HWND(window.Hwnd), win.GWL_WNDPROC, uintptr(CallbackToPtr(window.hookWindowProc)))

	window.view.OnTitleChanged(func(title string) {
		if window.fixedTitle {
			return
		}
		window.setTitle(title)
	})

	return window
}

func (w *Window) hookWindowProc(hwnd, message, wparam, lparam uintptr) uintptr {

	res := win.CallWindowProc(uintptr(w._oldWndProc), win.HWND(hwnd), uint32(message), wparam, lparam)

	switch message {
	case win.WM_ENTERSIZEMOVE:
		w.udpateCursor()
	}

	return res
}

// 返回 true 标识已处理
func (w *Window) DispatchMessage(msg *win.MSG) bool {

	if !w.enableBorderResize {
		return false
	}

	if w.windowType != WKE_WINDOW_TYPE_TRANSPARENT {
		return false
	}

	switch msg.Message {
	case win.WM_MOUSEMOVE:
		return w.handleMouseMove(msg)
	case win.WM_LBUTTONDOWN:
		if msg.WParam != win.MK_LBUTTON || w.sizing <= 0 {
			return false
		}
		win.SendMessage(win.HWND(w.Hwnd), win.WM_SYSCOMMAND, uintptr(win.SC_SIZE|w.sizing), 0)
		return true
	}

	// 没有任何处理事件
	return false
}

var border int32 = 5 // 设置 5 像素的 border 反应厚度
func (w *Window) handleMouseMove(msg *win.MSG) bool {
	rect := &win.RECT{}
	win.GetWindowRect(win.HWND(w.Hwnd), rect)

	inLeft := msg.Pt.X >= rect.Left && msg.Pt.X <= rect.Left+border
	inReght := msg.Pt.X <= rect.Right && msg.Pt.X >= rect.Right-border
	inTop := msg.Pt.Y >= rect.Top && msg.Pt.Y <= rect.Top+border
	inBottom := msg.Pt.Y <= rect.Bottom && msg.Pt.Y >= rect.Bottom-border

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

func (w *Window) Minimize() {
	go win.SendMessage(win.HWND(w.Hwnd), win.WM_SYSCOMMAND, uintptr(win.SC_MINIMIZE), 0)
}

func (w *Window) Maximize() {
	go win.SendMessage(win.HWND(w.Hwnd), win.WM_SYSCOMMAND, uintptr(win.SC_MAXIMIZE), 0)
	w.isMaximized = true
}

func (w *Window) IsMaximized() bool {
	// style := win.GetWindowLong(win.HWND(w.Hwnd), win.GWL_STYLE)
	// return (style & win.WS_MAXIMIZE) != 0
	return w.isMaximized
}

func (w *Window) Restore() {
	go win.SendMessage(win.HWND(w.Hwnd), win.WM_SYSCOMMAND, uintptr(win.SC_RESTORE), 0)
	w.isMaximized = false
}

func (w *Window) Move() {
	win.ReleaseCapture()
	go win.SendMessage(win.HWND(w.Hwnd), win.WM_SYSCOMMAND, uintptr(win.SC_MOVE|win.HTCAPTION), 0)
}

func (w *Window) SetIcon(handle win.HANDLE) error {
	if handle == 0 {
		return errors.New("获取图标句柄失败，无法设置 ICON 。")
	}
	win.SendMessage(win.HWND(w.Hwnd), win.WM_SETICON, 1, uintptr(handle))
	win.SendMessage(win.HWND(w.Hwnd), win.WM_SETICON, 0, uintptr(handle))

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
	iconFilePath := filepath.Join(w.mb.Config.runtimePath, "icon_"+dataHash+".ico")
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
