package blink

type (
	WkeString         uintptr
	WkeWebFrameHandle uintptr
	WkeHandle         uintptr
	JsExecState       uintptr
	JsValue           uintptr
)

type WkeSlist struct {
	Str  uintptr
	Next uintptr
}

type JsType uint32

const (
	JsType_NUMBER JsType = iota
	JsType_STRING
	JsType_BOOLEAN
	JsType_OBJECT
	JsType_FUNCTION
	JsType_UNDEFINED
	JsType_ARRAY
	JsType_NULL
)

type JsArg interface {
	int |
		string |
		bool
}

type JsData struct {
	Name [100]byte
	PropertyGet,
	PropertySet,
	Finalize,
	CallAsFunction uintptr
}

type JsKeys struct {
	Length uint32
	First  uintptr
}

type WkeRequestType int

const (
	WkeRequestType_Unknow WkeRequestType = iota + 1
	WkeRequestType_Get
	WkeRequestType_Post
	WkeRequestType_Put
)

type WkeKeyFlags int

const (
	WkeKeyFlags_Extend WkeKeyFlags = 0x0100
	WkeKeyFlags_Repeat WkeKeyFlags = 0x4000
)

type WkeRect struct {
	X, Y, W, H int32
}

type WkeNetJob uintptr

type WkeMouseFlags int

const (
	WkeMouseFlags_None    WkeMouseFlags = 0
	WkeMouseFlags_LBUTTON WkeMouseFlags = 0x01
	WkeMouseFlags_RBUTTON WkeMouseFlags = 0x02
	WkeMouseFlags_SHIFT   WkeMouseFlags = 0x04
	WkeMouseFlags_CONTROL WkeMouseFlags = 0x08
	WkeMouseFlags_MBUTTON WkeMouseFlags = 0x10
)

type WkeConsoleLevel int

const (
	WkeConsoleLevel_Log WkeConsoleLevel = iota + 1
	WkeConsoleLevel_Warning
	WkeConsoleLevel_Error
	WkeConsoleLevel_Debug
	WkeConsoleLevel_Info
	WkeConsoleLevel_RevokedError
)

type WkeNavigationType int

const (
	WKE_NAVIGATION_TYPE_LINKCLICK WkeNavigationType = iota
	WKE_NAVIGATION_TYPE_FORMSUBMITTED
	WKE_NAVIGATION_TYPE_BACKFORWARD
	WKE_NAVIGATION_TYPE_RELOAD
	WKE_NAVIGATION_TYPE_FORMRESUBMITTED
	WKE_NAVIGATION_TYPE_OTHER
)

type WkeWindowFeatures struct {
	X      int // 窗口x坐标
	Y      int // 窗口y坐标
	Width  int // 窗口宽度
	Height int // 窗口高度

	MenuBarVisible     bool // 是否显示菜单栏
	StatusBarVisible   bool // 是否显示状态栏
	ToolBarVisible     bool // 是否显示工具栏
	LocationBarVisible bool // 是否显示地址栏
	ScrollbarsVisible  bool // 是否显示滚动条
	Resizable          bool // 是否可调整大小
	Fullscreen         bool // 是否全屏显示
}

type WkeWindowClosingCallback func(view WkeHandle, param uintptr) (boolRes uintptr)
type WkeWindowDestroyCallback func(view WkeHandle, param uintptr) (voidRes uintptr)
type WkePaintBitUpdatedCallback func(view WkeHandle, param, buf []byte, rect *WkeRect, width, height int32) (voidRes uintptr)
type WkeNetResponseCallback func(view WkeHandle, param uintptr, url string, job WkeNetJob) (boolRes uintptr)
type WkeLoadUrlBeginCallback func(view WkeHandle, param uintptr, url string, job WkeNetJob) (boolRes uintptr)
type WkeJsNativeFunction func(es JsExecState, param uintptr) (voidRes uintptr)
type WkeDidCreateScriptContextCallback func(view WkeHandle, param uintptr, frame WkeWebFrameHandle, context uintptr, exGroup, worldId int) (voidRes uintptr)
type WkeWillReleaseScriptContextCallback func(webView WkeHandle, param uintptr, frameId WkeWebFrameHandle, context uintptr, worldId int) (voidRes uintptr)
type WkeConsoleCallback func(view WkeHandle, param uintptr, level WkeConsoleLevel, message, sourceName WkeString, sourceLine uint32, stackTrace WkeString) (voidRes uintptr)
type WkeLoadUrlEndCallback func(view WkeHandle, param uintptr, url string, job WkeNetJob, buf []byte) (voidRes uintptr)
type WkeLoadUrlFailCallback func(view WkeHandle, param, url string, job WkeNetJob) (voidRes uintptr)
type WkeDocumentReady2Callback func(view WkeHandle, param uintptr, frame WkeWebFrameHandle) (voidRes uintptr)
type WkeOnShowDevtoolsCallback func(view WkeHandle, param uintptr) (voidRes uintptr)
type WkeTitleChangedCallback func(view WkeHandle, param uintptr, title WkeString) (voidRes uintptr)
type WkeDownloadCallback func(view WkeHandle, param uintptr, url uintptr) (voidRes uintptr)
type WkeCreateViewCallback func(webView WkeHandle, param uintptr, navigationType WkeNavigationType, url WkeString, windowFeatures *WkeWindowFeatures) WkeHandle

type WkeCursorType int

const (
	WkeCursorType_Pointer WkeCursorType = iota
	WkeCursorType_Cross
	WkeCursorType_Hand
	WkeCursorType_IBeam
	WkeCursorType_Wait
	WkeCursorType_Help
	WkeCursorType_EastResize
	WkeCursorType_NorthResize
	WkeCursorType_NorthEastResize
	WkeCursorType_NorthWestResize
	WkeCursorType_SouthResize
	WkeCursorType_SouthEastResize
	WkeCursorType_SouthWestResize
	WkeCursorType_WestResize
	WkeCursorType_NorthSouthResize
	WkeCursorType_EastWestResize
	WkeCursorType_NorthEastSouthWestResize
	WkeCursorType_NorthWestSouthEastResize
	WkeCursorType_ColumnResize
	WkeCursorType_RowResize
	WkeCursorType_MiddlePanning
	WkeCursorType_EastPanning
	WkeCursorType_NorthPanning
	WkeCursorType_NorthEastPanning
	WkeCursorType_NorthWestPanning
	WkeCursorType_SouthPanning
	WkeCursorType_SouthEastPanning
	WkeCursorType_SouthWestPanning
	WkeCursorType_WestPanning
	WkeCursorType_Move
	WkeCursorType_VerticalText
	WkeCursorType_Cell
	WkeCursorType_ContextMenu
	WkeCursorType_Alias
	WkeCursorType_Progress
	WkeCursorType_NoDrop
	WkeCursorType_Copy
	WkeCursorType_None
	WkeCursorType_NotAllowed
	WkeCursorType_ZoomIn
	WkeCursorType_ZoomOut
	WkeCursorType_Grab
	WkeCursorType_Grabbing
	WkeCursorType_Custom
)

type ProxyType int

const (
	ProxyType_NONE ProxyType = iota
	ProxyType_HTTP
	ProxyType_SOCKS4
	ProxyType_SOCKS4A
	ProxyType_SOCKS5
	ProxyType_SOCKS5HOSTNAME
)

type ProxyInfo struct {
	Type     ProxyType
	HostName string
	Port     int
	UserName string
	Password string
}

type WkeWindowType uintptr

const (
	// 普通窗口
	WKE_WINDOW_TYPE_POPUP WkeWindowType = iota
	// 透明窗口。mb内部通过layer window实现
	WKE_WINDOW_TYPE_TRANSPARENT
	// 嵌入在父窗口里的子窗口。此时parent需要被设置
	WKE_WINDOW_TYPE_CONTROL
	// 隐藏标题栏
	WKE_WINDOW_TYPE_HIDE_CAPTION
)

type BOOL int32

const TRUE BOOL = 1
const FALSE BOOL = 0

type wkePrintSettings struct {
	structSize               int32 // 结构体大小，每个 int 为4, 12个int为48（极个别 C 编译器的int大小为8，暂不予考虑）
	dpi                      int32 // DPI 默认600
	width                    int32 // 纸张宽度，单位 px。（600 DPI 时，转换为像素 A4 的宽度为 4961 px）
	height                   int32 // 纸张高度，单位 px。（600 DPI 时，转换为像素 A4 的宽度为 4961 px）
	marginTop                int32 // 上边距，单位 px。（600 DPI 时， 1厘米边距转换为像素：236 px）
	marginBottom             int32 // 下边距，单位 px。（600 DPI 时， 1厘米边距转换为像素：236 px）
	marginLeft               int32 // 左边距，单位 px。（600 DPI 时， 1厘米边距转换为像素：236 px）
	marginRight              int32 // 右边距，单位 px。（600 DPI 时， 1厘米边距转换为像素：236 px）
	isPrintPageHeadAndFooter BOOL  // 是否打印页眉页脚
	isPrintBackgroud         BOOL  // 是否打印背景
	isLandscape              BOOL  // 是否横向打印
	isPrintToMultiPage       BOOL  // 是否打印到多页（分开保存为多个文档）
}

type wkePdfDatas struct {
	count int
	sizes uintptr
	datas uintptr // 二进制数据
}
