﻿


typedef unsigned long       DWORD;
typedef DWORD   COLORREF;
typedef unsigned long HWND;

typedef struct _RECT
{
    long    left;
    long    top;
    long    right;
    long    bottom;
} RECT;

typedef struct _wkeRect {
    int x;
    int y;
    int w;
    int h;
} wkeRect;

typedef struct _wkePoint {
    int x;
    int y;
} wkePoint;

typedef struct _wkeSize {
    int w;
    int h;
} wkeSize;

typedef enum _wkeMouseFlags {
    WKE_LBUTTON = 0x01,
    WKE_RBUTTON = 0x02,
    WKE_SHIFT = 0x04,
    WKE_CONTROL = 0x08,
    WKE_MBUTTON = 0x10,
} wkeMouseFlags;

typedef enum _wkeKeyFlags {
    WKE_EXTENDED = 0x0100,
    WKE_REPEAT = 0x4000,
} wkeKeyFlags;

typedef enum _wkeMouseMsg {
    WKE_MSG_MOUSEMOVE = 0x0200,
    WKE_MSG_LBUTTONDOWN = 0x0201,
    WKE_MSG_LBUTTONUP = 0x0202,
    WKE_MSG_LBUTTONDBLCLK = 0x0203,
    WKE_MSG_RBUTTONDOWN = 0x0204,
    WKE_MSG_RBUTTONUP = 0x0205,
    WKE_MSG_RBUTTONDBLCLK = 0x0206,
    WKE_MSG_MBUTTONDOWN = 0x0207,
    WKE_MSG_MBUTTONUP = 0x0208,
    WKE_MSG_MBUTTONDBLCLK = 0x0209,
    WKE_MSG_MOUSEWHEEL = 0x020A,
} wkeMouseMsg;



typedef char utf8;
struct JsExecStateInfo;
typedef JsExecStateInfo* jsExecState;


typedef __int64 jsValue;


struct _tagWkeWebView;
typedef struct _tagWkeWebView* wkeWebView;

struct _tagWkeString;
typedef struct _tagWkeString* wkeString;


typedef struct _tagWkeMediaPlayer* wkeMediaPlayer;
typedef struct _tagWkeMediaPlayerClient* wkeMediaPlayerClient;
typedef struct _tabblinkWebURLRequestPtr* blinkWebURLRequestPtr;


typedef enum _wkeProxyType {
    WKE_PROXY_NONE,
    WKE_PROXY_HTTP,
    WKE_PROXY_SOCKS4,
    WKE_PROXY_SOCKS4A,
    WKE_PROXY_SOCKS5,
    WKE_PROXY_SOCKS5HOSTNAME
} wkeProxyType;

typedef struct _wkeProxy {
    wkeProxyType type;
    char hostname[100];
    unsigned short port;
    char username[50];
    char password[50];
} wkeProxy;

typedef enum _wkeSettingMask {
    WKE_SETTING_PROXY = 1,
    WKE_SETTING_EXTENSION = 1 << 2, // 测试功能，请勿调用
} wkeSettingMask;

typedef struct _wkeSettings {
    wkeProxy proxy;
    unsigned int mask;
    const char* extension;
} wkeSettings;

typedef struct _wkeViewSettings {
    int size;
    unsigned int bgColor;
} wkeViewSettings;

typedef void* wkeWebFrameHandle;

typedef struct _wkeGeolocationPosition{
    double timestamp;
    double latitude;
    double longitude;
    double accuracy;
    bool providesAltitude;
    double altitude;
    bool providesAltitudeAccuracy;
    double altitudeAccuracy;
    bool providesHeading;
    double heading;
    bool providesSpeed;
    double speed;
} wkeGeolocationPosition;


typedef enum _wkeMenuItemId {
    kWkeMenuSelectedAllId = 1 << 1,
    kWkeMenuSelectedTextId = 1 << 2,
    kWkeMenuUndoId = 1 << 3,
    kWkeMenuCopyImageId = 1 << 4,
    kWkeMenuInspectElementAtId = 1 << 5,
    kWkeMenuCutId = 1 << 6,
    kWkeMenuPasteId = 1 << 7,
    kWkeMenuPrintId = 1 << 8,
    kWkeMenuGoForwardId = 1 << 9,
    kWkeMenuGoBackId = 1 << 10,
    kWkeMenuReloadId = 1 << 11,
    kWkeMenuSaveImageId = 1 << 12,
} wkeMenuItemId;

typedef void* (__cdecl  *FILE_OPEN_) (const char* path);
typedef void(__cdecl  *FILE_CLOSE_) (void* handle);
typedef size_t(__cdecl  *FILE_SIZE) (void* handle);
typedef int(__cdecl  *FILE_READ) (void* handle, void* buffer, size_t size);
typedef int(__cdecl  *FILE_SEEK) (void* handle, int offset, int origin);

typedef FILE_OPEN_ WKE_FILE_OPEN;
typedef FILE_CLOSE_ WKE_FILE_CLOSE;
typedef FILE_SIZE WKE_FILE_SIZE;
typedef FILE_READ WKE_FILE_READ;
typedef FILE_SEEK WKE_FILE_SEEK;
typedef bool(__cdecl  *WKE_EXISTS_FILE)(const char * path);

struct _wkeClientHandler; // declare warning fix
typedef void(__cdecl  *ON_TITLE_CHANGED) (const struct _wkeClientHandler* clientHandler, const wkeString title);
typedef void(__cdecl  *ON_URL_CHANGED) (const struct _wkeClientHandler* clientHandler, const wkeString url);

typedef struct _wkeClientHandler {
    ON_TITLE_CHANGED onTitleChanged;
    ON_URL_CHANGED onURLChanged;
} wkeClientHandler;

typedef bool(__cdecl  * wkeCookieVisitor)(
    void* params,
    const char* name,
    const char* value,
    const char* domain,
    const char* path, // If |path| is non-empty only URLs at or below the path will get the cookie value.
    int secure, // If |secure| is true the cookie will only be sent for HTTPS requests.
    int httpOnly, // If |httponly| is true the cookie will only be sent for HTTP requests.
    int* expires // The cookie expiration date is only valid if |has_expires| is true.
    );

typedef enum _wkeCookieCommand {
    wkeCookieCommandClearAllCookies,
    wkeCookieCommandClearSessionCookies,
    wkeCookieCommandFlushCookiesToFile,
    wkeCookieCommandReloadCookiesFromFile,
} wkeCookieCommand;

typedef enum _wkeNavigationType {
    WKE_NAVIGATION_TYPE_LINKCLICK,
    WKE_NAVIGATION_TYPE_FORMSUBMITTE,
    WKE_NAVIGATION_TYPE_BACKFORWARD,
    WKE_NAVIGATION_TYPE_RELOAD,
    WKE_NAVIGATION_TYPE_FORMRESUBMITT,
    WKE_NAVIGATION_TYPE_OTHER
} wkeNavigationType;

typedef enum _WkeCursorInfoType {
    WkeCursorInfoPointer,
    WkeCursorInfoCross,
    WkeCursorInfoHand,
    WkeCursorInfoIBeam,
    WkeCursorInfoWait,
    WkeCursorInfoHelp,
    WkeCursorInfoEastResize,
    WkeCursorInfoNorthResize,
    WkeCursorInfoNorthEastResize,
    WkeCursorInfoNorthWestResize,
    WkeCursorInfoSouthResize,
    WkeCursorInfoSouthEastResize,
    WkeCursorInfoSouthWestResize,
    WkeCursorInfoWestResize,
    WkeCursorInfoNorthSouthResize,
    WkeCursorInfoEastWestResize,
    WkeCursorInfoNorthEastSouthWestResize,
    WkeCursorInfoNorthWestSouthEastResize,
    WkeCursorInfoColumnResize,
    WkeCursorInfoRowResize,
    WkeCursorInfoMiddlePanning,
    WkeCursorInfoEastPanning,
    WkeCursorInfoNorthPanning,
    WkeCursorInfoNorthEastPanning,
    WkeCursorInfoNorthWestPanning,
    WkeCursorInfoSouthPanning,
    WkeCursorInfoSouthEastPanning,
    WkeCursorInfoSouthWestPanning,
    WkeCursorInfoWestPanning,
    WkeCursorInfoMove,
    WkeCursorInfoVerticalText,
    WkeCursorInfoCell,
    WkeCursorInfoContextMenu,
    WkeCursorInfoAlias,
    WkeCursorInfoProgress,
    WkeCursorInfoNoDrop,
    WkeCursorInfoCopy,
    WkeCursorInfoNone,
    WkeCursorInfoNotAllowed,
    WkeCursorInfoZoomIn,
    WkeCursorInfoZoomOut,
    WkeCursorInfoGrab,
    WkeCursorInfoGrabbing,
    WkeCursorInfoCustom
} WkeCursorInfoType;

typedef struct _wkeWindowFeatures {
    int x;
    int y;
    int width;
    int height;

    bool menuBarVisible;
    bool statusBarVisible;
    bool toolBarVisible;
    bool locationBarVisible;
    bool scrollbarsVisible;
    bool resizable;
    bool fullscreen;
} wkeWindowFeatures;

typedef struct _wkeMemBuf {
    int unuse;
    void* data;
    size_t length;
} wkeMemBuf;

typedef enum _wkeStorageType {
            // String data with an associated MIME type. Depending on the MIME type, there may be
            // optional metadata attributes as well.
            StorageTypeString,
            // Stores the name of one file being dragged into the renderer.
            StorageTypeFilename,
            // An image being dragged out of the renderer. Contains a buffer holding the image data
            // as well as the suggested name for saving the image to.
            StorageTypeBinaryData,
            // Stores the filesystem URL of one file being dragged into the renderer.
            StorageTypeFileSystemFile,
} wkeStorageType;

struct wkeWebDragDataItem {
        wkeStorageType storageType;

        // Only valid when storageType == StorageTypeString.
        wkeMemBuf* stringType;
        wkeMemBuf* stringData;

        // Only valid when storageType == StorageTypeFilename.
        wkeMemBuf* filenameData;
        wkeMemBuf* displayNameData;

        // Only valid when storageType == StorageTypeBinaryData.
        wkeMemBuf* binaryData;

        // Title associated with a link when stringType == "text/uri-list".
        // Filename when storageType == StorageTypeBinaryData.
        wkeMemBuf* title;

        // Only valid when storageType == StorageTypeFileSystemFile.
        wkeMemBuf* fileSystemURL;
        __int64 fileSystemFileSize;

        // Only valid when stringType == "text/html".
        wkeMemBuf* baseURL;
    } ;

typedef struct _wkeWebDragData {

    wkeWebDragDataItem* m_itemList;
    int m_itemListLength;

    int m_modifierKeyState; // State of Shift/Ctrl/Alt/Meta keys.
    wkeMemBuf* m_filesystemId;
} wkeWebDragData;

typedef enum _wkeWebDragOperation {
    wkeWebDragOperationNone = 0,
    wkeWebDragOperationCopy = 1,
    wkeWebDragOperationLink = 2,
    wkeWebDragOperationGeneric = 4,
    wkeWebDragOperationPrivate = 8,
    wkeWebDragOperationMove = 16,
    wkeWebDragOperationDelete = 32,
    wkeWebDragOperationEvery = 0xffffffff
} wkeWebDragOperation;

typedef wkeWebDragOperation wkeWebDragOperationsMask;

typedef enum _wkeResourceType {
    WKE_RESOURCE_TYPE_MAIN_FRAME = 0,       // top level page
    WKE_RESOURCE_TYPE_SUB_FRAME = 1,        // frame or iframe
    WKE_RESOURCE_TYPE_STYLESHEET = 2,       // a CSS stylesheet
    WKE_RESOURCE_TYPE_SCRIPT = 3,           // an external script
    WKE_RESOURCE_TYPE_IMAGE = 4,            // an image (jpg/gif/png/etc)
    WKE_RESOURCE_TYPE_FONT_RESOURCE = 5,    // a font
    WKE_RESOURCE_TYPE_SUB_RESOURCE = 6,     // an "other" subresource.
    WKE_RESOURCE_TYPE_OBJECT = 7,           // an object (or embed) tag for a plugin,
                                            // or a resource that a plugin requested.
    WKE_RESOURCE_TYPE_MEDIA = 8,            // a media resource.
    WKE_RESOURCE_TYPE_WORKER = 9,           // the main resource of a dedicated
                                            // worker.
    WKE_RESOURCE_TYPE_SHARED_WORKER = 10,   // the main resource of a shared worker.
    WKE_RESOURCE_TYPE_PREFETCH = 11,        // an explicitly requested prefetch
    WKE_RESOURCE_TYPE_FAVICON = 12,         // a favicon
    WKE_RESOURCE_TYPE_XHR = 13,             // a XMLHttpRequest
    WKE_RESOURCE_TYPE_PING = 14,            // a ping request for <a ping>
    WKE_RESOURCE_TYPE_SERVICE_WORKER = 15,  // the main resource of a service worker.
    WKE_RESOURCE_TYPE_LAST_TYPE
} wkeResourceType;

typedef struct _wkeWillSendRequestInfo {
    wkeString url;
    wkeString newUrl;
    wkeResourceType resourceType;
    int httpResponseCode;
    wkeString method;
    wkeString referrer;
    void* headers;
} wkeWillSendRequestInfo;

typedef enum _wkeHttBodyElementType {
    wkeHttBodyElementTypeData,
    wkeHttBodyElementTypeFile,
} wkeHttBodyElementType;

typedef struct _wkePostBodyElement {
    int size;
    wkeHttBodyElementType type;
    wkeMemBuf* data;
    wkeString filePath;
    __int64 fileStart;
    __int64 fileLength; // -1 means to the end of the file.
} wkePostBodyElement;

typedef struct _wkePostBodyElements {
    int size;
    wkePostBodyElement** element;
    size_t elementSize;
    bool isDirty;
} wkePostBodyElements;

typedef void* wkeNetJob;

typedef struct _wkeSlist {
    char* data;
    struct _wkeSlist* next;
} wkeSlist;

typedef struct _wkeTempCallbackInfo {
    int size;
    wkeWebFrameHandle frame;
    wkeWillSendRequestInfo* willSendRequestInfo;
    const char* url;
    wkePostBodyElements* postBody;
    wkeNetJob job;
} wkeTempCallbackInfo;

typedef enum _wkeRequestType {
    kWkeRequestTypeInvalidation,
    kWkeRequestTypeGet,
    kWkeRequestTypePost,
    kWkeRequestTypePut,
} wkeRequestType;

typedef struct _wkePdfDatas {
    int count;
    size_t* sizes;
    const void** datas;
} wkePdfDatas;

typedef struct _wkePrintSettings {
    int structSize;
    int dpi;
    int width; // in px
    int height;
    int marginTop;
    int marginBottom;
    int marginLeft;
    int marginRight;
    bool isPrintPageHeadAndFooter;
    bool isPrintBackgroud;
    bool isLandscape;
    bool isPrintToMultiPage;
} wkePrintSettings;

typedef struct _wkeScreenshotSettings {
    int structSize;
    int width;
    int height;
} wkeScreenshotSettings;

typedef void(__cdecl *wkeCaretChangedCallback)(wkeWebView webView, void* param, const wkeRect* r);
typedef void(__cdecl *wkeTitleChangedCallback)(wkeWebView webView, void* param, const wkeString title);
typedef void(__cdecl *wkeURLChangedCallback)(wkeWebView webView, void* param, const wkeString url);
typedef void(__cdecl *wkeURLChangedCallback2)(wkeWebView webView, void* param, wkeWebFrameHandle frameId, const wkeString url);
typedef void(__cdecl *wkePaintUpdatedCallback)(wkeWebView webView, void* param, const HDC hdc, int x, int y, int cx, int cy);
typedef void(__cdecl *wkePaintBitUpdatedCallback)(wkeWebView webView, void* param, const void* buffer, const wkeRect* r, int width, int height);
typedef void(__cdecl *wkeAlertBoxCallback)(wkeWebView webView, void* param, const wkeString msg);
typedef bool(__cdecl *wkeConfirmBoxCallback)(wkeWebView webView, void* param, const wkeString msg);
typedef bool(__cdecl *wkePromptBoxCallback)(wkeWebView webView, void* param, const wkeString msg, const wkeString defaultResult, wkeString result);
typedef bool(__cdecl *wkeNavigationCallback)(wkeWebView webView, void* param, wkeNavigationType navigationType, wkeString url);
typedef wkeWebView(__cdecl *wkeCreateViewCallback)(wkeWebView webView, void* param, wkeNavigationType navigationType, const wkeString url, const wkeWindowFeatures* windowFeatures);
typedef void(__cdecl *wkeDocumentReadyCallback)(wkeWebView webView, void* param);
typedef void(__cdecl *wkeDocumentReady2Callback)(wkeWebView webView, void* param, wkeWebFrameHandle frameId);

typedef void(__cdecl *wkeOnShowDevtoolsCallback)(wkeWebView webView, void* param);

typedef void(__stdcall *wkeNodeOnCreateProcessCallback)(wkeWebView webView, void* param, const wchar_t* applicationPath, const WCHAR* arguments, STARTUPINFOW* startup);
typedef void(__cdecl *wkeOnPluginFindCallback)(wkeWebView webView, void* param, const utf8* mime, void* initializeFunc, void* getEntryPointsFunc, void* shutdownFunc);

typedef void(__cdecl *wkeOnPrintCallback)(wkeWebView webView, void* param, wkeWebFrameHandle frameId, void* printParams);
typedef void(__cdecl *wkeOnScreenshot)(wkeWebView webView, void* param, const char* data, size_t size);

typedef wkeString(__cdecl *wkeImageBufferToDataURL)(wkeWebView webView, void* param, const char* data, size_t size);

typedef struct _wkeMediaLoadInfo {
    int size;
    int width;
    int height;
    double duration;
} wkeMediaLoadInfo;

typedef void(__cdecl *wkeWillMediaLoadCallback)(wkeWebView webView, void* param, const char* url, wkeMediaLoadInfo* info);

typedef void(__cdecl *wkeStartDraggingCallback)(
    wkeWebView webView,
    void* param, 
    wkeWebFrameHandle frame,
    const wkeWebDragData* data,
    wkeWebDragOperationsMask mask, 
    const void* image, 
    const wkePoint* dragImageOffset
    );

typedef void(__cdecl *wkeUiThreadRunCallback)(HWND hWnd, void* param);
typedef int(__cdecl *wkeUiThreadPostTaskCallback)(HWND hWnd, wkeUiThreadRunCallback callback, void* param);

typedef enum _wkeOtherLoadType {
    WKE_DID_START_LOADING,
    WKE_DID_STOP_LOADING,
    WKE_DID_NAVIGATE,
    WKE_DID_NAVIGATE_IN_PAGE,
    WKE_DID_GET_RESPONSE_DETAILS,
    WKE_DID_GET_REDIRECT_REQUEST,
    WKE_DID_POST_REQUEST,
} wkeOtherLoadType;
typedef void(__cdecl *wkeOnOtherLoadCallback)(wkeWebView webView, void* param, wkeOtherLoadType type, wkeTempCallbackInfo* info);

typedef enum _wkeOnContextMenuItemClickType {
    kWkeContextMenuItemClickTypePrint = 0x01,
} wkeOnContextMenuItemClickType;

typedef enum _wkeOnContextMenuItemClickStep {
    kWkeContextMenuItemClickStepShow = 0x01,
    kWkeContextMenuItemClickStepClick = 0x02,
} wkeOnContextMenuItemClickStep;

typedef bool(__cdecl * wkeOnContextMenuItemClickCallback)(
    wkeWebView webView, 
    void* param, 
    wkeOnContextMenuItemClickType type, 
    wkeOnContextMenuItemClickStep step, 
    wkeWebFrameHandle frameId,
    void* info
    );

typedef enum _wkeLoadingResult {
    WKE_LOADING_SUCCEEDED,
    WKE_LOADING_FAILED,
    WKE_LOADING_CANCELED
} wkeLoadingResult;

typedef enum _wkeDownloadOpt {
    kWkeDownloadOptCancel,
    kWkeDownloadOptCacheData,
} wkeDownloadOpt;

typedef void(__cdecl *wkeNetJobDataRecvCallback)(void* ptr, wkeNetJob job, const char* data, int length);
typedef void(__cdecl *wkeNetJobDataFinishCallback)(void* ptr, wkeNetJob job, wkeLoadingResult result);

typedef struct _wkeNetJobDataBind {
    void* param;
    wkeNetJobDataRecvCallback recvCallback;
    wkeNetJobDataFinishCallback finishCallback;
} wkeNetJobDataBind;

typedef void(__cdecl * wkePopupDialogSaveNameCallback)(void* ptr, const wchar_t* filePath);

typedef struct _wkeDownloadBind {
    void* param;
    wkeNetJobDataRecvCallback recvCallback;
    wkeNetJobDataFinishCallback finishCallback;
    wkePopupDialogSaveNameCallback saveNameCallback;
} wkeDownloadBind;

enum wkeDialogProperties {
    kWkeDialogPropertiesOpenFile = 1 << 1, // 允许选择文件
    kWkeDialogPropertiesOpenDirectory = 1 << 2, // 允许选择文件夹
    kWkeDialogPropertiesMultiSelections = 1 << 3, // 允许多选。
    kWkeDialogPropertiesShowHiddenFiles = 1 << 4, // 显示对话框中的隐藏文件。
    kWkeDialogPropertiesCreateDirectory = 1 << 5, // macOS - 允许你通过对话框的形式创建新的目录。
    kWkeDialogPropertiesPromptToCreate = 1 << 6, // Windows - 如果输入的文件路径在对话框中不存在, 则提示创建。 这并不是真的在路径上创建一个文件，而是允许返回一些不存在的地址交由应用程序去创建。
    kWkeDialogPropertiesNoResolveAliases = 1 << 7, // macOS - 禁用自动的别名路径(符号链接) 解析。 所选别名现在将会返回别名路径而非其目标路径。
    kWkeDialogPropertiesTreatPackageAsDirectory = 1 << 8, // macOS - 将包(如.app 文件夹) 视为目录而不是文件。
    kWkeDialogPropertiesDontAddToRecent = 1 << 9, // Windows - 不要将正在打开的项目添加到最近的文档列表中。
};

typedef struct _wkeFileFilter {
  const utf8* name; // 例如"image"、"Movies"
  const utf8* extensions; // 例如"jpg|png|gif"
} wkeFileFilter;

typedef struct _wkeDialogOptions {
    int magic; // 'mbdo'
    utf8* title;
    utf8* defaultPath;
    utf8* buttonLabel;
    wkeFileFilter* filters;
    int filtersCount;
    wkeDialogProperties prop;
    utf8* message;
    BOOL securityScopedBookmarks;
} wkeDialogOptions;

typedef wkeDownloadOpt(__cdecl * wkeDownloadInBlinkThreadCallback)(
    wkeWebView webView,
    void* param,
    size_t expectedContentLength,
    const char* url,
    const char* mime,
    const char* disposition,
    wkeNetJob job,
    wkeNetJobDataBind* dataBind
    );

typedef void(__cdecl *wkeLoadingFinishCallback)(wkeWebView webView, void* param, const wkeString url, wkeLoadingResult result, const wkeString failedReason);
typedef bool(__cdecl *wkeDownloadCallback)(wkeWebView webView, void* param, const char* url);
typedef wkeDownloadOpt(__cdecl *wkeDownload2Callback)(
    wkeWebView webView, 
    void* param,
    size_t expectedContentLength,
    const char* url, 
    const char* mime, 
    const char* disposition, 
    wkeNetJob job, 
    wkeNetJobDataBind* dataBind);

typedef enum _wkeConsoleLevel {
    wkeLevelDebug = 4,
    wkeLevelLog = 1,
    wkeLevelInfo = 5,
    wkeLevelWarning = 2,
    wkeLevelError = 3,
    wkeLevelRevokedError = 6,
    wkeLevelLast = wkeLevelInfo
} wkeConsoleLevel;
typedef void(__cdecl *wkeConsoleCallback)(wkeWebView webView, void* param, wkeConsoleLevel level, const wkeString message, const wkeString sourceName, unsigned sourceLine, const wkeString stackTrace);

typedef void(__cdecl *wkeOnCallUiThread)(wkeWebView webView, void* paramOnInThread);
typedef void(__cdecl *wkeCallUiThread)(wkeWebView webView, wkeOnCallUiThread func, void* param);

typedef wkeMediaPlayer(__cdecl * wkeMediaPlayerFactory)(wkeWebView webView, wkeMediaPlayerClient client, void* npBrowserFuncs, void* npPluginFuncs);
typedef bool(__cdecl * wkeOnIsMediaPlayerSupportsMIMEType)(const utf8* mime);

//wkeNet--------------------------------------------------------------------------------------

typedef struct wkeWebUrlRequest* wkeWebUrlRequestPtr;
typedef struct wkeWebUrlResponse* wkeWebUrlResponsePtr;

typedef void(__cdecl * wkeOnUrlRequestWillRedirectCallback)(wkeWebView webView, void* param, wkeWebUrlRequestPtr oldRequest, wkeWebUrlRequestPtr request, wkeWebUrlResponsePtr redirectResponse);
typedef void(__cdecl * wkeOnUrlRequestDidReceiveResponseCallback)(wkeWebView webView, void* param, wkeWebUrlRequestPtr request, wkeWebUrlResponsePtr response);
typedef void(__cdecl * wkeOnUrlRequestDidReceiveDataCallback)(wkeWebView webView, void* param, wkeWebUrlRequestPtr request, const char* data, int dataLength);
typedef void(__cdecl * wkeOnUrlRequestDidFailCallback)(wkeWebView webView, void* param, wkeWebUrlRequestPtr request, const utf8* error);
typedef void(__cdecl * wkeOnUrlRequestDidFinishLoadingCallback)(wkeWebView webView, void* param, wkeWebUrlRequestPtr request, double finishTime);

typedef struct _wkeUrlRequestCallbacks {
    wkeOnUrlRequestWillRedirectCallback willRedirectCallback;
    wkeOnUrlRequestDidReceiveResponseCallback didReceiveResponseCallback;
    wkeOnUrlRequestDidReceiveDataCallback didReceiveDataCallback;
    wkeOnUrlRequestDidFailCallback didFailCallback;
    wkeOnUrlRequestDidFinishLoadingCallback didFinishLoadingCallback;
} wkeUrlRequestCallbacks;

typedef bool(__cdecl *wkeLoadUrlBeginCallback)(wkeWebView webView, void* param, const utf8* url, wkeNetJob job);
typedef void(__cdecl *wkeLoadUrlEndCallback)(wkeWebView webView, void* param, const utf8* url, wkeNetJob job, void* buf, int len);
typedef void(__cdecl *wkeLoadUrlFailCallback)(wkeWebView webView, void* param, const utf8* url, wkeNetJob job);
typedef void(__cdecl *wkeLoadUrlHeadersReceivedCallback)(wkeWebView webView, void* param, const utf8* url, wkeNetJob job);
typedef void(__cdecl *wkeLoadUrlFinishCallback)(wkeWebView webView, void* param, const utf8* url, wkeNetJob job, int len);

typedef void(__cdecl *wkeDidCreateScriptContextCallback)(wkeWebView webView, void* param, wkeWebFrameHandle frameId, void* context, int extensionGroup, int worldId);
typedef void(__cdecl *wkeWillReleaseScriptContextCallback)(wkeWebView webView, void* param, wkeWebFrameHandle frameId, void* context, int worldId);
typedef bool(__cdecl *wkeNetResponseCallback)(wkeWebView webView, void* param, const utf8* url, wkeNetJob job);
typedef void(__cdecl *wkeOnNetGetFaviconCallback)(wkeWebView webView, void* param, const utf8* url, wkeMemBuf* buf);

typedef void* v8ContextPtr;
typedef void* v8Isolate;

//wke window-----------------------------------------------------------------------------------
typedef enum _wkeWindowType {
    WKE_WINDOW_TYPE_POPUP,
    WKE_WINDOW_TYPE_TRANSPARENT,
    WKE_WINDOW_TYPE_CONTROL
} wkeWindowType;

typedef struct _wkeWindowCreateInfo {
    int size;
    HWND parent;
    DWORD style; 
    DWORD styleEx; 
    int x; 
    int y; 
    int width; 
    int height;
    COLORREF color;
} wkeWindowCreateInfo;

typedef struct _wkeDefaultPrinterSettings {
    int structSize; // 默认是4 * 10
    bool isLandscape; // 是否为横向打印格式
    bool isPrintHeadFooter; // 
    bool isPrintBackgroud; // 是否打印背景
    int edgeDistanceLeft; // 上边距单位：毫米
    int edgeDistanceTop;
    int edgeDistanceRight;
    int edgeDistanceBottom;
    int copies; // 默认打印份数
    int paperType; // DMPAPER_A4等
} wkeDefaultPrinterSettings;

typedef bool(__cdecl *wkeWindowClosingCallback)(wkeWebView webWindow, void* param);
typedef void(__cdecl *wkeWindowDestroyCallback)(wkeWebView webWindow, void* param);

typedef struct {
    RECT bounds;
    bool draggable;
} wkeDraggableRegion;
typedef void(__cdecl *wkeDraggableRegionsChangedCallback)(wkeWebView webView, void* param, const wkeDraggableRegion* rects, int rectCount);

//JavaScript Bind-----------------------------------------------------------------------------------

typedef jsValue(__fastcall* jsNativeFunction) (jsExecState es);

typedef jsValue(__cdecl * wkeJsNativeFunction) (jsExecState es, void* param);

typedef enum _jsType {
    JSTYPE_NUMBER,
    JSTYPE_STRING,
    JSTYPE_BOOLEAN,
    JSTYPE_OBJECT,
    JSTYPE_FUNCTION,
    JSTYPE_UNDEFINED,
    JSTYPE_ARRAY,
    JSTYPE_NULL,
} jsType;

// cexer JS对象、函数绑定支持
typedef jsValue(__cdecl *jsGetPropertyCallback)(jsExecState es, jsValue object, const char* propertyName);
typedef bool(__cdecl *jsSetPropertyCallback)(jsExecState es, jsValue object, const char* propertyName, jsValue value);
typedef jsValue(__cdecl *jsCallAsFunctionCallback)(jsExecState es, jsValue object, jsValue* args, int argCount);
struct tagjsData; // declare warning fix
typedef void(__cdecl *jsFinalizeCallback)(struct tagjsData* data);

typedef struct tagjsData {
    char typeName[100];
    jsGetPropertyCallback propertyGet;
    jsSetPropertyCallback propertySet;
    jsFinalizeCallback finalize;
    jsCallAsFunctionCallback callAsFunction;
} jsData;

typedef struct _jsExceptionInfo {
    const utf8* message; // Returns the exception message.
    const utf8* sourceLine; // Returns the line of source code that the exception occurred within.
    const utf8* scriptResourceName; // Returns the resource name for the script from where the function causing the error originates.
    int lineNumber; // Returns the 1-based number of the line where the error occurred or 0 if the line number is unknown.
    int startPosition; // Returns the index within the script of the first character where the error occurred.
    int endPosition; // Returns the index within the script of the last character where the error occurred.
    int startColumn; // Returns the index within the line of the first character where the error occurred.
    int endColumn; // Returns the index within the line of the last character where the error occurred.
    const utf8* callstackString;
} jsExceptionInfo;

typedef struct _jsKeys {
    unsigned int length;
    const char** keys;


} jsKeys;



// 以下是wke的导出函数。格式按照【返回类型】【函数名】【参数】来排列

#define WKE_FOR_EACH_DEFINE_FUNCTION(ITERATOR0, ITERATOR1, ITERATOR2, ITERATOR3, ITERATOR4, ITERATOR5, ITERATOR6, ITERATOR9, ITERATOR10, ITERATOR11) \
    ITERATOR0(void, wkeShutdown, "") \
    ITERATOR0(void, wkeShutdownForDebug, "测试使用，不了解千万别用！") \
    \
    ITERATOR0(unsigned int, wkeVersion, "") \
    ITERATOR0(const utf8*, wkeVersionString, "") \
    ITERATOR2(void, wkeGC, wkeWebView webView, long intervalSec, "") \
    ITERATOR2(void, wkeSetResourceGc, wkeWebView webView, long intervalSec, "") \
    \
    ITERATOR5(void, wkeSetFileSystem, WKE_FILE_OPEN pfnOpen, WKE_FILE_CLOSE pfnClose, WKE_FILE_SIZE pfnSize, WKE_FILE_READ pfnRead, WKE_FILE_SEEK pfnSeek, "") \
    \
    ITERATOR1(const char*, wkeWebViewName, wkeWebView webView, "") \
    ITERATOR2(void, wkeSetWebViewName, wkeWebView webView, const char* name, "") \
    \
    ITERATOR1(BOOL, wkeIsLoaded, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeIsLoadFailed, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeIsLoadComplete, wkeWebView webView, "") \
    \
    ITERATOR1(const utf8*, wkeGetSource, wkeWebView webView, "") \
    ITERATOR1(const utf8*, wkeTitle, wkeWebView webView, "") \
    ITERATOR1(const wchar_t*, wkeTitleW, wkeWebView webView, "") \
    ITERATOR1(int, wkeWidth, wkeWebView webView, "") \
    ITERATOR1(int, wkeHeight, wkeWebView webView, "") \
    ITERATOR1(int, wkeContentsWidth, wkeWebView webView, "") \
    ITERATOR1(int, wkeContentsHeight, wkeWebView webView, "") \
    \
    ITERATOR1(void, wkeSelectAll, wkeWebView webView, "") \
    ITERATOR1(void, wkeCopy, wkeWebView webView, "") \
    ITERATOR1(void, wkeCut, wkeWebView webView, "") \
    ITERATOR1(void, wkePaste, wkeWebView webView, "") \
    ITERATOR1(void, wkeDelete, wkeWebView webView, "") \
    \
    ITERATOR1(BOOL, wkeCookieEnabled, wkeWebView webView, "") \
    ITERATOR1(float, wkeMediaVolume, wkeWebView webView, "") \
    \
    ITERATOR5(BOOL, wkeMouseEvent, wkeWebView webView, unsigned int message, int x, int y, unsigned int flags, "") \
    ITERATOR4(BOOL, wkeContextMenuEvent, wkeWebView webView, int x, int y, unsigned int flags, "") \
    ITERATOR5(BOOL, wkeMouseWheel, wkeWebView webView, int x, int y, int delta, unsigned int flags, "") \
    ITERATOR4(BOOL, wkeKeyUp, wkeWebView webView, unsigned int virtualKeyCode, unsigned int flags, bool systemKey, "") \
    ITERATOR4(BOOL, wkeKeyDown, wkeWebView webView, unsigned int virtualKeyCode, unsigned int flags, bool systemKey, "") \
    ITERATOR4(BOOL, wkeKeyPress, wkeWebView webView, unsigned int virtualKeyCode, unsigned int flags, bool systemKey, "") \
    \
    ITERATOR1(void, wkeFocus, wkeWebView webView, "") \
    ITERATOR1(void, wkeUnfocus, wkeWebView webView, "") \
    \
    ITERATOR1(wkeRect, wkeGetCaret, wkeWebView webView, "") \
    \
    ITERATOR1(void, wkeAwaken, wkeWebView webView, "") \
    \
    ITERATOR1(float, wkeZoomFactor, wkeWebView webView, "") \
    \
    ITERATOR2(void, wkeSetClientHandler, wkeWebView webView, const wkeClientHandler* handler, "") \
    ITERATOR1(const wkeClientHandler*, wkeGetClientHandler, wkeWebView webView, "") \
    \
    ITERATOR1(const utf8*, wkeToString, const wkeString string, "") \
    ITERATOR1(const wchar_t*, wkeToStringW, const wkeString string, "") \
    \
    ITERATOR2(const utf8*, jsToString, jsExecState es, jsValue v, "") \
    ITERATOR2(const wchar_t*, jsToStringW, jsExecState es, jsValue v, "") \
    \
    ITERATOR1(void, wkeConfigure, const wkeSettings* settings, "") \
    ITERATOR0(BOOL, wkeIsInitialize, "") \
    \
    ITERATOR2(void, wkeSetViewSettings, wkeWebView webView, const wkeViewSettings* settings, "") \
    ITERATOR3(void, wkeSetDebugConfig, wkeWebView webView, const char* debugString, const char* param, "") \
    ITERATOR2(void *, wkeGetDebugConfig, wkeWebView webView, const char* debugString, "") \
    \
    ITERATOR0(void, wkeFinalize, "") \
    ITERATOR0(void, wkeUpdate, "") \
    ITERATOR0(unsigned int, wkeGetVersion, "") \
    ITERATOR0(const utf8*, wkeGetVersionString, "") \
    \
    ITERATOR0(wkeWebView, wkeCreateWebView, "") \
    ITERATOR1(void, wkeDestroyWebView, wkeWebView webView, "") \
    \
    ITERATOR2(void, wkeSetMemoryCacheEnable, wkeWebView webView, bool b, "") \
    ITERATOR2(void, wkeSetMouseEnabled, wkeWebView webView, bool b, "") \
    ITERATOR2(void, wkeSetTouchEnabled, wkeWebView webView, bool b, "") \
    ITERATOR2(void, wkeSetSystemTouchEnabled, wkeWebView webView, bool b, "") \
    ITERATOR2(void, wkeSetContextMenuEnabled, wkeWebView webView, bool b, "") \
    ITERATOR2(void, wkeSetNavigationToNewWindowEnable, wkeWebView webView, bool b, "") \
    ITERATOR2(void, wkeSetCspCheckEnable, wkeWebView webView, bool b, "") \
    ITERATOR2(void, wkeSetNpapiPluginsEnabled, wkeWebView webView, bool b, "") \
    ITERATOR2(void, wkeSetHeadlessEnabled, wkeWebView webView, bool b, "可以关闭渲染") \
    ITERATOR2(void, wkeSetDragEnable, wkeWebView webView, bool b, "可关闭拖拽文件加载网页") \
    ITERATOR2(void, wkeSetDragDropEnable, wkeWebView webView, bool b, "可关闭拖拽到其他进程") \
    ITERATOR3(void, wkeSetContextMenuItemShow, wkeWebView webView, wkeMenuItemId item, bool isShow, "设置某项menu是否显示") \
    ITERATOR2(void, wkeSetLanguage, wkeWebView webView, const char* language, "") \
    \
    ITERATOR2(void, wkeSetViewNetInterface, wkeWebView webView, const char* netInterface, "") \
    \
    ITERATOR1(void, wkeSetProxy, const wkeProxy* proxy, "") \
    ITERATOR2(void, wkeSetViewProxy, wkeWebView webView, wkeProxy *proxy, "") \
    \
    ITERATOR1(const char*, wkeGetName, wkeWebView webView, "") \
    ITERATOR2(void, wkeSetName, wkeWebView webView, const char* name, "") \
    \
    ITERATOR2(void, wkeSetHandle, wkeWebView webView, HWND wnd, "") \
    ITERATOR3(void, wkeSetHandleOffset, wkeWebView webView, int x, int y, "") \
    \
    ITERATOR1(BOOL, wkeIsTransparent, wkeWebView webView, "") \
    ITERATOR2(void, wkeSetTransparent, wkeWebView webView, bool transparent, "") \
    \
    ITERATOR2(void, wkeSetUserAgent, wkeWebView webView, const utf8* userAgent, "") \
    ITERATOR1(const char*, wkeGetUserAgent, wkeWebView webView, "") \
    ITERATOR2(void, wkeSetUserAgentW, wkeWebView webView, const wchar_t* userAgent, "") \
    \
    ITERATOR4(void, wkeShowDevtools, wkeWebView webView, const wchar_t* path, wkeOnShowDevtoolsCallback callback, void* param, "") \
    \
    ITERATOR2(void, wkeLoadW, wkeWebView webView, const wchar_t* url, "") \
    ITERATOR2(void, wkeLoadURL, wkeWebView webView, const utf8* url, "") \
    ITERATOR2(void, wkeLoadURLW, wkeWebView webView, const wchar_t* url, "") \
    ITERATOR4(void, wkePostURL, wkeWebView wkeView, const utf8* url, const char* postData, int postLen, "") \
    ITERATOR4(void, wkePostURLW, wkeWebView wkeView, const wchar_t* url, const char* postData, int postLen, "") \
    \
    ITERATOR2(void, wkeLoadHTML, wkeWebView webView, const utf8* html, "") \
    ITERATOR3(void, wkeLoadHtmlWithBaseUrl, wkeWebView webView, const utf8* html, const utf8* baseUrl, "") \
    ITERATOR2(void, wkeLoadHTMLW, wkeWebView webView, const wchar_t* html, "") \
    \
    ITERATOR2(void, wkeLoadFile, wkeWebView webView, const utf8* filename, "") \
    ITERATOR2(void, wkeLoadFileW, wkeWebView webView, const wchar_t* filename, "") \
    \
    ITERATOR1(const utf8*, wkeGetURL, wkeWebView webView, "") \
    ITERATOR2(const utf8*, wkeGetFrameUrl, wkeWebView webView, wkeWebFrameHandle frameId, "") \
    \
    ITERATOR1(BOOL, wkeIsLoading, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeIsLoadingSucceeded, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeIsLoadingFailed, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeIsLoadingCompleted, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeIsDocumentReady, wkeWebView webView, "") \
    ITERATOR1(void, wkeStopLoading, wkeWebView webView, "") \
    ITERATOR1(void, wkeReload, wkeWebView webView, "") \
    ITERATOR2(void, wkeGoToOffset, wkeWebView webView, int offset, "") \
    ITERATOR2(void, wkeGoToIndex, wkeWebView webView, int index, "") \
    \
    ITERATOR1(int, wkeGetWebviewId, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeIsWebviewAlive, int id, "") \
    ITERATOR1(BOOL, wkeIsWebviewValid, wkeWebView webView, "") \
    \
    ITERATOR3(const utf8*, wkeGetDocumentCompleteURL, wkeWebView webView, wkeWebFrameHandle frameId, const utf8* partialURL, "") \
    \
    ITERATOR3(wkeMemBuf*, wkeCreateMemBuf, wkeWebView webView, void* buf, size_t length, "") \
    ITERATOR1(void, wkeFreeMemBuf, wkeMemBuf* buf, "") \
    \
    ITERATOR1(const utf8*, wkeGetTitle, wkeWebView webView, "") \
    ITERATOR1(const wchar_t*, wkeGetTitleW, wkeWebView webView, "") \
    \
    ITERATOR3(void, wkeResize, wkeWebView webView, int w, int h, "") \
    ITERATOR1(int, wkeGetWidth, wkeWebView webView, "") \
    ITERATOR1(int, wkeGetHeight, wkeWebView webView, "") \
    ITERATOR1(int, wkeGetContentWidth, wkeWebView webView, "") \
    ITERATOR1(int, wkeGetContentHeight, wkeWebView webView, "") \
    \
    ITERATOR2(void, wkeSetDirty, wkeWebView webView, bool dirty, "") \
    ITERATOR1(BOOL, wkeIsDirty, wkeWebView webView, "") \
    ITERATOR5(void, wkeAddDirtyArea, wkeWebView webView, int x, int y, int w, int h, "") \
    ITERATOR1(void, wkeLayoutIfNeeded, wkeWebView webView, "") \
    ITERATOR11(void, wkePaint2, wkeWebView webView, void* bits, int bufWid, int bufHei, int xDst, int yDst, int w, int h, int xSrc, int ySrc, bool bCopyAlpha, "") \
    ITERATOR3(void, wkePaint, wkeWebView webView, void* bits, int pitch, "") \
    ITERATOR1(void, wkeRepaintIfNeeded, wkeWebView webView, "") \
    ITERATOR1(HDC, wkeGetViewDC, wkeWebView webView, "") \
    ITERATOR1(void, wkeUnlockViewDC, wkeWebView webView, "") \
    ITERATOR1(HWND, wkeGetHostHWND, wkeWebView webView, "") \
    \
    ITERATOR1(BOOL, wkeCanGoBack, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeGoBack, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeCanGoForward, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeGoForward, wkeWebView webView, "") \
    ITERATOR2(BOOL, wkeNavigateAtIndex, wkeWebView webView, int index, "") \
    ITERATOR1(int, wkeGetNavigateIndex, wkeWebView webView, "") \
    \
    ITERATOR1(void, wkeEditorSelectAll, wkeWebView webView, "") \
    ITERATOR1(void, wkeEditorUnSelect, wkeWebView webView, "") \
    ITERATOR1(void, wkeEditorCopy, wkeWebView webView, "") \
    ITERATOR1(void, wkeEditorCut, wkeWebView webView, "") \
    ITERATOR1(void, wkeEditorPaste, wkeWebView webView, "") \
    ITERATOR1(void, wkeEditorDelete, wkeWebView webView, "") \
    ITERATOR1(void, wkeEditorUndo, wkeWebView webView, "") \
    ITERATOR1(void, wkeEditorRedo, wkeWebView webView, "") \
    \
    ITERATOR1(const wchar_t*, wkeGetCookieW, wkeWebView webView, "") \
    ITERATOR1(const utf8*, wkeGetCookie, wkeWebView webView, "") \
    ITERATOR3(void, wkeSetCookie, wkeWebView webView, const utf8* url, const utf8* cookie, "cookie格式必须是类似:cna=4UvTFE12fEECAXFKf4SFW5eo; expires=Tue, 23-Jan-2029 13:17:21 GMT; path=/; domain=.youku.com") \
    ITERATOR3(void, wkeVisitAllCookie, wkeWebView webView, void* params, wkeCookieVisitor visitor, "") \
    ITERATOR2(void, wkePerformCookieCommand, wkeWebView webView, wkeCookieCommand command, "") \
    ITERATOR2(void, wkeSetCookieEnabled, wkeWebView webView, bool enable, "") \
    ITERATOR1(BOOL, wkeIsCookieEnabled, wkeWebView webView, "") \
    ITERATOR2(void, wkeSetCookieJarPath, wkeWebView webView, const WCHAR* path, "") \
    ITERATOR2(void, wkeSetCookieJarFullPath, wkeWebView webView, const WCHAR* path, "") \
    ITERATOR1(void, wkeClearCookie, wkeWebView webView, "") \
    ITERATOR2(void, wkeSetLocalStorageFullPath, wkeWebView webView, const WCHAR* path, "") \
    ITERATOR2(void, wkeAddPluginDirectory, wkeWebView webView, const WCHAR* path, "") \
    \
    ITERATOR2(void, wkeSetMediaVolume, wkeWebView webView, float volume, "") \
    ITERATOR1(float, wkeGetMediaVolume, wkeWebView webView, "") \
    \
    ITERATOR5(BOOL, wkeFireMouseEvent, wkeWebView webView, unsigned int message, int x, int y, unsigned int flags, "") \
    ITERATOR4(BOOL, wkeFireContextMenuEvent, wkeWebView webView, int x, int y, unsigned int flags, "") \
    ITERATOR5(BOOL, wkeFireMouseWheelEvent, wkeWebView webView, int x, int y, int delta, unsigned int flags, "") \
    ITERATOR4(BOOL, wkeFireKeyUpEvent, wkeWebView webView, unsigned int virtualKeyCode, unsigned int flags, bool systemKey, "") \
    ITERATOR4(BOOL, wkeFireKeyDownEvent, wkeWebView webView, unsigned int virtualKeyCode, unsigned int flags, bool systemKey, "") \
    ITERATOR4(BOOL, wkeFireKeyPressEvent, wkeWebView webView, unsigned int charCode, unsigned int flags, bool systemKey, "") \
    ITERATOR6(BOOL, wkeFireWindowsMessage, wkeWebView webView, HWND hWnd, UINT message, WPARAM wParam, LPARAM lParam, LRESULT* result, "") \
    \
    ITERATOR1(void, wkeSetFocus, wkeWebView webView, "") \
    ITERATOR1(void, wkeKillFocus, wkeWebView webView, "") \
    \
    ITERATOR1(wkeRect, wkeGetCaretRect, wkeWebView webView, "") \
    ITERATOR1(wkeRect*, wkeGetCaretRect2, wkeWebView webView, "给一些不方便获取返回结构体的语言调用") \
    \
    ITERATOR2(jsValue, wkeRunJS, wkeWebView webView, const utf8* script, "") \
    ITERATOR2(jsValue, wkeRunJSW, wkeWebView webView, const wchar_t* script, "") \
    \
    ITERATOR1(jsExecState, wkeGlobalExec, wkeWebView webView, "") \
    ITERATOR2(jsExecState, wkeGetGlobalExecByFrame, wkeWebView webView, wkeWebFrameHandle frameId, "") \
    \
    ITERATOR1(void, wkeSleep, wkeWebView webView, "") \
    ITERATOR1(void, wkeWake, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeIsAwake, wkeWebView webView, "") \
    \
    ITERATOR2(void, wkeSetZoomFactor, wkeWebView webView, float factor, "") \
    ITERATOR1(float, wkeGetZoomFactor, wkeWebView webView, "") \
    ITERATOR0(void, wkeEnableHighDPISupport, "") \
    \
    ITERATOR2(void, wkeSetEditable, wkeWebView webView, bool editable, "") \
    \
    ITERATOR1(const utf8*, wkeGetString, const wkeString string, "") \
    ITERATOR1(const wchar_t*, wkeGetStringW, const wkeString string, "") \
    \
    ITERATOR3(void, wkeSetString, wkeString string, const utf8* str, size_t len, "") \
    ITERATOR3(void, wkeSetStringWithoutNullTermination, wkeString string, const utf8* str, size_t len, "") \
    ITERATOR3(void, wkeSetStringW, wkeString string, const wchar_t* str, size_t len, "") \
    \
    ITERATOR2(wkeString, wkeCreateString, const utf8* str, size_t len, "") \
    ITERATOR2(wkeString, wkeCreateStringW, const wchar_t* str, size_t len, "") \
    ITERATOR2(wkeString, wkeCreateStringWithoutNullTermination, const utf8* str, size_t len, "") \
    ITERATOR1(size_t, wkeGetStringLen, wkeString str, "") \
    ITERATOR1(void, wkeDeleteString, wkeString str, "") \
    \
    ITERATOR0(wkeWebView, wkeGetWebViewForCurrentContext, "") \
    ITERATOR3(void, wkeSetUserKeyValue, wkeWebView webView, const char* key, void* value, "") \
    ITERATOR2(void*, wkeGetUserKeyValue, wkeWebView webView, const char* key, "") \
    \
    ITERATOR1(int, wkeGetCursorInfoType, wkeWebView webView, "") \
    ITERATOR2(void, wkeSetCursorInfoType, wkeWebView webView, int type, "") \
    ITERATOR5(void, wkeSetDragFiles, wkeWebView webView, const POINT* clintPos, const POINT* screenPos, wkeString* files, int filesCount, "") \
    \
    ITERATOR5(void, wkeSetDeviceParameter, wkeWebView webView, const char* device, const char* paramStr, int paramInt, float paramFloat, "") \
    ITERATOR1(wkeTempCallbackInfo*, wkeGetTempCallbackInfo, wkeWebView webView, "") \
    \
    ITERATOR3(void, wkeOnCaretChanged, wkeWebView webView, wkeCaretChangedCallback callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnMouseOverUrlChanged, wkeWebView webView, wkeTitleChangedCallback callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnTitleChanged, wkeWebView webView, wkeTitleChangedCallback callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnURLChanged, wkeWebView webView, wkeURLChangedCallback callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnURLChanged2, wkeWebView webView, wkeURLChangedCallback2 callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnPaintUpdated, wkeWebView webView, wkePaintUpdatedCallback callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnPaintBitUpdated, wkeWebView webView, wkePaintBitUpdatedCallback callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnAlertBox, wkeWebView webView, wkeAlertBoxCallback callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnConfirmBox, wkeWebView webView, wkeConfirmBoxCallback callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnPromptBox, wkeWebView webView, wkePromptBoxCallback callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnNavigation, wkeWebView webView, wkeNavigationCallback callback, void* param, "") \
    ITERATOR3(void, wkeOnCreateView, wkeWebView webView, wkeCreateViewCallback callback, void* param, "") \
    ITERATOR3(void, wkeOnDocumentReady, wkeWebView webView, wkeDocumentReadyCallback callback, void* param, "") \
    ITERATOR3(void, wkeOnDocumentReady2, wkeWebView webView, wkeDocumentReady2Callback callback, void* param, "") \
    ITERATOR3(void, wkeOnLoadingFinish, wkeWebView webView, wkeLoadingFinishCallback callback, void* param, "") \
    ITERATOR3(void, wkeOnDownload, wkeWebView webView, wkeDownloadCallback callback, void* param, "") \
    ITERATOR3(void, wkeOnDownload2, wkeWebView webView, wkeDownload2Callback callback, void* param, "") \
    ITERATOR3(void, wkeOnConsole, wkeWebView webView, wkeConsoleCallback callback, void* param, "") \
    ITERATOR3(void, wkeSetUIThreadCallback, wkeWebView webView, wkeCallUiThread callback, void* param, "") \
    ITERATOR3(void, wkeOnLoadUrlBegin, wkeWebView webView, wkeLoadUrlBeginCallback callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnLoadUrlEnd, wkeWebView webView, wkeLoadUrlEndCallback callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnLoadUrlHeadersReceived, wkeWebView webView, wkeLoadUrlHeadersReceivedCallback callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnLoadUrlFinish, wkeWebView webView, wkeLoadUrlFinishCallback callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnLoadUrlFail, wkeWebView webView, wkeLoadUrlFailCallback callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnDidCreateScriptContext, wkeWebView webView, wkeDidCreateScriptContextCallback callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnWillReleaseScriptContext, wkeWebView webView, wkeWillReleaseScriptContextCallback callback, void* callbackParam, "") \
    ITERATOR3(void, wkeOnWindowClosing, wkeWebView webWindow, wkeWindowClosingCallback callback, void* param, "") \
    ITERATOR3(void, wkeOnWindowDestroy, wkeWebView webWindow, wkeWindowDestroyCallback callback, void* param, "") \
    ITERATOR3(void, wkeOnDraggableRegionsChanged, wkeWebView webView, wkeDraggableRegionsChangedCallback callback, void* param, "") \
    ITERATOR3(void, wkeOnWillMediaLoad, wkeWebView webView, wkeWillMediaLoadCallback callback, void* param, "") \
    ITERATOR3(void, wkeOnStartDragging, wkeWebView webView, wkeStartDraggingCallback callback, void* param, "") \
    ITERATOR3(void, wkeOnPrint, wkeWebView webView, wkeOnPrintCallback callback, void* param, "") \
    ITERATOR4(void, wkeScreenshot, wkeWebView webView, const wkeScreenshotSettings* settings, wkeOnScreenshot callback, void* param, "") \
    \
    ITERATOR3(void, wkeOnOtherLoad, wkeWebView webView, wkeOnOtherLoadCallback callback, void* param, "") \
    ITERATOR3(void, wkeOnContextMenuItemClick, wkeWebView webView, wkeOnContextMenuItemClickCallback callback, void* param, "") \
    \
    ITERATOR1(BOOL, wkeIsProcessingUserGesture, wkeWebView webView, "") \
    \
    ITERATOR2(void, wkeNetSetMIMEType, wkeNetJob jobPtr, const char* type, "设置response的mime") \
    ITERATOR2(const char*, wkeNetGetMIMEType, wkeNetJob jobPtr, wkeString mime, "获取response的mime") \
    ITERATOR1(const char*, wkeNetGetReferrer, wkeNetJob jobPtr, "获取request的referrer") \
    ITERATOR4(void, wkeNetSetHTTPHeaderField, wkeNetJob jobPtr, const wchar_t* key, const wchar_t* value, bool response, "") \
    ITERATOR2(const char*, wkeNetGetHTTPHeaderField, wkeNetJob jobPtr, const char* key, "") \
    ITERATOR2(const char*, wkeNetGetHTTPHeaderFieldFromResponse, wkeNetJob jobPtr, const char* key, "") \
    ITERATOR3(void, wkeNetSetData, wkeNetJob jobPtr, void* buf, int len, "调用此函数后,网络层收到数据会存储在一buf内,接收数据完成后响应OnLoadUrlEnd事件.#此调用严重影响性能,慎用" \
        "此函数和wkeNetSetData的区别是，wkeNetHookRequest会在接受到真正网络数据后再调用回调，并允许回调修改网络数据。"\
        "而wkeNetSetData是在网络数据还没发送的时候修改") \
    ITERATOR1(void, wkeNetHookRequest, wkeNetJob jobPtr, "") \
    ITERATOR3(void, wkeNetOnResponse, wkeWebView webView, wkeNetResponseCallback callback, void* param, "") \
    ITERATOR1(wkeRequestType, wkeNetGetRequestMethod, wkeNetJob jobPtr, "") \
    ITERATOR3(int, wkeNetGetFavicon, wkeWebView webView, wkeOnNetGetFaviconCallback callback, void* param, "") \
    \
    ITERATOR1(void, wkeNetContinueJob, wkeNetJob jobPtr, "")\
    ITERATOR1(const char*, wkeNetGetUrlByJob, wkeNetJob jobPtr, "")\
    ITERATOR1(const wkeSlist*, wkeNetGetRawHttpHead, wkeNetJob jobPtr, "")\
    ITERATOR1(const wkeSlist*, wkeNetGetRawResponseHead, wkeNetJob jobPtr, "")\
    \
    ITERATOR1(void, wkeNetCancelRequest, wkeNetJob jobPtr, "")\
    ITERATOR1(BOOL, wkeNetHoldJobToAsynCommit, wkeNetJob jobPtr, "")\
    ITERATOR2(void, wkeNetChangeRequestUrl, wkeNetJob jobPtr, const char* url, "")\
    \
    ITERATOR3(wkeWebUrlRequestPtr, wkeNetCreateWebUrlRequest, const utf8* url, const utf8* method, const utf8* mime, "")\
    ITERATOR1(wkeWebUrlRequestPtr, wkeNetCreateWebUrlRequest2, const blinkWebURLRequestPtr request, "")\
    ITERATOR2(blinkWebURLRequestPtr, wkeNetCopyWebUrlRequest, wkeNetJob jobPtr, bool needExtraData, "")\
    ITERATOR1(void, wkeNetDeleteBlinkWebURLRequestPtr, blinkWebURLRequestPtr request, "")\
    ITERATOR3(void, wkeNetAddHTTPHeaderFieldToUrlRequest, wkeWebUrlRequestPtr request, const utf8* name, const utf8* value, "")\
    ITERATOR4(int, wkeNetStartUrlRequest, wkeWebView webView, wkeWebUrlRequestPtr request, void* param, const wkeUrlRequestCallbacks* callbacks, "")\
    ITERATOR1(int, wkeNetGetHttpStatusCode, wkeWebUrlResponsePtr response, "")\
    ITERATOR1(__int64, wkeNetGetExpectedContentLength, wkeWebUrlResponsePtr response, "")\
    ITERATOR1(const utf8*, wkeNetGetResponseUrl, wkeWebUrlResponsePtr response, "")\
    ITERATOR1(void, wkeNetCancelWebUrlRequest, int requestId, "")\
    \
    ITERATOR1(wkePostBodyElements*, wkeNetGetPostBody, wkeNetJob jobPtr, "") \
    ITERATOR2(wkePostBodyElements*, wkeNetCreatePostBodyElements, wkeWebView webView, size_t length, "") \
    ITERATOR1(void, wkeNetFreePostBodyElements, wkePostBodyElements* elements, "") \
    ITERATOR1(wkePostBodyElement*, wkeNetCreatePostBodyElement, wkeWebView webView, "") \
    ITERATOR1(void, wkeNetFreePostBodyElement, wkePostBodyElement* element, "") \
    \
    ITERATOR9(wkeDownloadOpt, wkePopupDialogAndDownload, wkeWebView webviewHandle, const wkeDialogOptions* dialogOpt, \
        size_t expectedContentLength, const char* url, const char* mime, const char* disposition, wkeNetJob job,wkeNetJobDataBind* dataBind, wkeDownloadBind* callbackBind, "") \
    ITERATOR10(wkeDownloadOpt, wkeDownloadByPath, wkeWebView webviewHandle, const wkeDialogOptions* dialogOpt, const WCHAR* path, size_t expectedContentLength,const char* url, \
        const char* mime, const char* disposition, wkeNetJob job, wkeNetJobDataBind* dataBind, wkeDownloadBind* callbackBind, "") \
    \
    ITERATOR2(BOOL, wkeIsMainFrame, wkeWebView webView, wkeWebFrameHandle frameId, "") \
    ITERATOR2(BOOL, wkeIsWebRemoteFrame, wkeWebView webView, wkeWebFrameHandle frameId, "") \
    ITERATOR1(wkeWebFrameHandle, wkeWebFrameGetMainFrame, wkeWebView webView, "") \
    ITERATOR4(jsValue, wkeRunJsByFrame, wkeWebView webView, wkeWebFrameHandle frameId, const utf8* script, bool isInClosure, "") \
    ITERATOR3(void, wkeInsertCSSByFrame, wkeWebView webView, wkeWebFrameHandle frameId, const utf8* cssText, "") \
    \
    ITERATOR3(void, wkeWebFrameGetMainWorldScriptContext, wkeWebView webView, wkeWebFrameHandle webFrameId, v8ContextPtr contextOut, "") \
    \
    ITERATOR0(v8Isolate, wkeGetBlinkMainThreadIsolate, "") \
    \
    ITERATOR6(wkeWebView, wkeCreateWebWindow, wkeWindowType type, HWND parent, int x, int y, int width, int height, "") \
    ITERATOR1(wkeWebView, wkeCreateWebCustomWindow, const wkeWindowCreateInfo* info, "") \
    ITERATOR1(void, wkeDestroyWebWindow, wkeWebView webWindow, "") \
    ITERATOR1(HWND, wkeGetWindowHandle, wkeWebView webWindow, "") \
    \
    ITERATOR2(void, wkeShowWindow, wkeWebView webWindow, bool show, "") \
    ITERATOR2(void, wkeEnableWindow, wkeWebView webWindow, bool enable, "") \
    \
    ITERATOR5(void, wkeMoveWindow, wkeWebView webWindow, int x, int y, int width, int height, "") \
    ITERATOR1(void, wkeMoveToCenter, wkeWebView webWindow, "") \
    ITERATOR3(void, wkeResizeWindow, wkeWebView webWindow, int width, int height, "") \
    \
    ITERATOR6(wkeWebDragOperation, wkeDragTargetDragEnter, wkeWebView webView, const wkeWebDragData* webDragData, const POINT* clientPoint, const POINT* screenPoint, wkeWebDragOperationsMask operationsAllowed, int modifiers, "") \
    ITERATOR5(wkeWebDragOperation, wkeDragTargetDragOver, wkeWebView webView, const POINT* clientPoint, const POINT* screenPoint, wkeWebDragOperationsMask operationsAllowed, int modifiers, "") \
    ITERATOR1(void, wkeDragTargetDragLeave, wkeWebView webView, "") \
    ITERATOR4(void, wkeDragTargetDrop, wkeWebView webView, const POINT* clientPoint, const POINT* screenPoint, int modifiers, "") \
    ITERATOR4(void, wkeDragTargetEnd, wkeWebView webView, const POINT* clientPoint, const POINT* screenPoint, wkeWebDragOperation operation, "") \
    \
    ITERATOR1(void, wkeUtilSetUiCallback, wkeUiThreadPostTaskCallback callback, "") \
    ITERATOR1(const utf8*, wkeUtilSerializeToMHTML, wkeWebView webView, "") \
    ITERATOR3(const wkePdfDatas*, wkeUtilPrintToPdf, wkeWebView webView, wkeWebFrameHandle frameId, const wkePrintSettings* settings,"") \
    ITERATOR3(const wkeMemBuf*, wkePrintToBitmap, wkeWebView webView, wkeWebFrameHandle frameId, const wkeScreenshotSettings* settings,"") \
    ITERATOR1(void, wkeUtilRelasePrintPdfDatas, const wkePdfDatas* datas,"") \
    \
    ITERATOR2(void, wkeSetWindowTitle, wkeWebView webWindow, const utf8* title, "") \
    ITERATOR2(void, wkeSetWindowTitleW, wkeWebView webWindow, const wchar_t* title, "") \
    \
    ITERATOR3(void, wkeNodeOnCreateProcess, wkeWebView webView, wkeNodeOnCreateProcessCallback callback, void* param, "") \
    \
    ITERATOR4(void, wkeOnPluginFind, wkeWebView webView, const char* mime, wkeOnPluginFindCallback callback, void* param, "") \
    ITERATOR4(void, wkeAddNpapiPlugin, wkeWebView webView, void* initializeFunc, void* getEntryPointsFunc, void* shutdownFunc, "") \
    \
    ITERATOR4(void, wkePluginListBuilderAddPlugin, void* builder, const utf8* name, const utf8* description, const utf8* fileName, "") \
    ITERATOR3(void, wkePluginListBuilderAddMediaTypeToLastPlugin, void* builder, const utf8* name, const utf8* description, "") \
    ITERATOR2(void, wkePluginListBuilderAddFileExtensionToLastMediaType, void* builder, const utf8* fileExtension, "") \
    \
    ITERATOR1(wkeWebView, wkeGetWebViewByNData, void* ndata, "") \
    \
    ITERATOR5(BOOL, wkeRegisterEmbedderCustomElement, wkeWebView webView, wkeWebFrameHandle frameId, const char* name, void* options, void* outResult, "") \
    \
    ITERATOR3(void, wkeSetMediaPlayerFactory, wkeWebView webView, wkeMediaPlayerFactory factory, wkeOnIsMediaPlayerSupportsMIMEType callback, "") \
    \
    ITERATOR3(const utf8* , wkeGetContentAsMarkup, wkeWebView webView, wkeWebFrameHandle frame, size_t* size, "") \
    \
    ITERATOR1(const utf8*, wkeUtilDecodeURLEscape, const utf8* url, "") \
    ITERATOR1(const utf8*, wkeUtilEncodeURLEscape, const utf8* url, "") \
    ITERATOR1(const utf8*, wkeUtilBase64Encode, const utf8* str, "") \
    ITERATOR1(const utf8*, wkeUtilBase64Decode, const utf8* str, "") \
    ITERATOR1(const wkeMemBuf*, wkeUtilCreateV8Snapshot, const utf8* str, "") \
    \
    ITERATOR0(void, wkeRunMessageLoop, "") \
    \
    ITERATOR1(void, wkeSaveMemoryCache, wkeWebView webView, "") \
    \
    ITERATOR3(void, jsBindFunction, const char* name, jsNativeFunction fn, unsigned int argCount, "") \
    ITERATOR2(void, jsBindGetter, const char* name, jsNativeFunction fn, "") \
    ITERATOR2(void, jsBindSetter, const char* name, jsNativeFunction fn, "") \
    \
    ITERATOR4(void, wkeJsBindFunction, const char* name, wkeJsNativeFunction fn, void* param, unsigned int argCount, "") \
    ITERATOR3(void, wkeJsBindGetter, const char* name, wkeJsNativeFunction fn, void* param, "") \
    ITERATOR3(void, wkeJsBindSetter, const char* name, wkeJsNativeFunction fn, void* param, "") \
    \
    ITERATOR1(int, jsArgCount, jsExecState es, "") \
    ITERATOR2(jsType, jsArgType, jsExecState es, int argIdx, "") \
    ITERATOR2(jsValue, jsArg, jsExecState es, int argIdx, "") \
    \
    ITERATOR1(jsType, jsTypeOf, jsValue v, "") \
    ITERATOR1(BOOL, jsIsNumber, jsValue v, "") \
    ITERATOR1(BOOL, jsIsString, jsValue v, "") \
    ITERATOR1(BOOL, jsIsBoolean, jsValue v, "") \
    ITERATOR1(BOOL, jsIsObject, jsValue v, "") \
    ITERATOR1(BOOL, jsIsFunction, jsValue v, "") \
    ITERATOR1(BOOL, jsIsUndefined, jsValue v, "") \
    ITERATOR1(BOOL, jsIsNull, jsValue v, "") \
    ITERATOR1(BOOL, jsIsArray, jsValue v, "") \
    ITERATOR1(BOOL, jsIsTrue, jsValue v, "") \
    ITERATOR1(BOOL, jsIsFalse, jsValue v, "") \
    \
    ITERATOR2(int, jsToInt, jsExecState es, jsValue v, "") \
    ITERATOR2(float, jsToFloat, jsExecState es, jsValue v, "") \
    ITERATOR2(double, jsToDouble, jsExecState es, jsValue v, "") \
    ITERATOR2(const char*, jsToDoubleString, jsExecState es, jsValue v, "") \
    ITERATOR2(BOOL, jsToBoolean, jsExecState es, jsValue v, "") \
    ITERATOR3(jsValue, jsArrayBuffer, jsExecState es, const char* buffer, size_t size, "") \
    ITERATOR2(wkeMemBuf*, jsGetArrayBuffer, jsExecState es, jsValue value, "") \
    ITERATOR2(const utf8*, jsToTempString, jsExecState es, jsValue v, "") \
    ITERATOR2(const wchar_t*, jsToTempStringW, jsExecState es, jsValue v, "") \
    ITERATOR2(void*, jsToV8Value, jsExecState es, jsValue v, "return v8::Persistent<v8::Value>*") \
    \
    ITERATOR1(jsValue, jsInt, int n, "") \
    ITERATOR1(jsValue, jsFloat, float f, "") \
    ITERATOR1(jsValue, jsDouble, double d, "") \
    ITERATOR1(jsValue, jsDoubleString, const char* str, "") \
    ITERATOR1(jsValue, jsBoolean, bool b, "") \
    \
    ITERATOR0(jsValue, jsUndefined, "") \
    ITERATOR0(jsValue, jsNull, "") \
    ITERATOR0(jsValue, jsTrue, "") \
    ITERATOR0(jsValue, jsFalse, "") \
    \
    ITERATOR2(jsValue, jsString, jsExecState es, const utf8* str, "") \
    ITERATOR2(jsValue, jsStringW, jsExecState es, const wchar_t* str, "") \
    ITERATOR1(jsValue, jsEmptyObject, jsExecState es, "") \
    ITERATOR1(jsValue, jsEmptyArray, jsExecState es, "") \
    \
    ITERATOR2(jsValue, jsObject, jsExecState es, jsData* obj, "") \
    ITERATOR2(jsValue, jsFunction, jsExecState es, jsData* obj, "") \
    ITERATOR2(jsData*, jsGetData, jsExecState es, jsValue object, "") \
    \
    ITERATOR3(jsValue, jsGet, jsExecState es, jsValue object, const char* prop, "") \
    ITERATOR4(void, jsSet, jsExecState es, jsValue object, const char* prop, jsValue v, "") \
    \
    ITERATOR3(jsValue, jsGetAt, jsExecState es, jsValue object, int index, "") \
    ITERATOR4(void, jsSetAt, jsExecState es, jsValue object, int index, jsValue v, "") \
    ITERATOR2(jsKeys*, jsGetKeys, jsExecState es, jsValue object, "") \
    ITERATOR2(BOOL, jsIsJsValueValid, jsExecState es, jsValue object, "") \
    ITERATOR1(BOOL, jsIsValidExecState, jsExecState es, "") \
    ITERATOR3(void, jsDeleteObjectProp, jsExecState es, jsValue object, const char* prop, "") \
    \
    ITERATOR2(int, jsGetLength, jsExecState es, jsValue object, "") \
    ITERATOR3(void, jsSetLength, jsExecState es, jsValue object, int length, "") \
    \
    ITERATOR1(jsValue, jsGlobalObject, jsExecState es, "") \
    ITERATOR1(wkeWebView, jsGetWebView, jsExecState es, "") \
    \
    ITERATOR2(jsValue, jsEval, jsExecState es, const utf8* str, "") \
    ITERATOR2(jsValue, jsEvalW, jsExecState es, const wchar_t* str, "") \
    ITERATOR3(jsValue, jsEvalExW, jsExecState es, const wchar_t* str, bool isInClosure, "") \
    \
    ITERATOR5(jsValue, jsCall, jsExecState es, jsValue func, jsValue thisObject, jsValue* args, int argCount, "") \
    ITERATOR4(jsValue, jsCallGlobal, jsExecState es, jsValue func, jsValue* args, int argCount, "") \
    \
    ITERATOR2(jsValue, jsGetGlobal, jsExecState es, const char* prop, "") \
    ITERATOR3(void, jsSetGlobal, jsExecState es, const char* prop, jsValue v, "") \
    \
    ITERATOR0(void, jsGC, "") \
    ITERATOR2(BOOL, jsAddRef, jsExecState es, jsValue val, "") \
    ITERATOR2(BOOL, jsReleaseRef, jsExecState es, jsValue val, "") \
    ITERATOR1(jsExceptionInfo*, jsGetLastErrorIfException, jsExecState es, "") \
    ITERATOR2(jsValue, jsThrowException, jsExecState es, const utf8* exception, "") \
    ITERATOR1(const utf8*, jsGetCallstack, jsExecState es, "")


extern "C" __declspec(dllexport) void __cdecl  wkeInit();
extern "C" __declspec(dllexport) void __cdecl  wkeInitialize();
extern "C" __declspec(dllexport) void __cdecl  wkeInitializeEx(const wkeSettings* settings);
