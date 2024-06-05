package wke

	ITERATOR0(void, wkeShutdown, "") \
    ITERATOR0(void, wkeShutdownForDebug, "测试使用，不了解千万别用！") \

    ITERATOR0(unsigned int, wkeVersion, "") \
    ITERATOR0(const utf8*, wkeVersionString, "") \
    ITERATOR2(void, wkeGC, wkeWebView webView, long intervalSec, "") \
    ITERATOR2(void, wkeSetResourceGc, wkeWebView webView, long intervalSec, "") \

    ITERATOR5(void, wkeSetFileSystem, WKE_FILE_OPEN pfnOpen, WKE_FILE_CLOSE pfnClose, WKE_FILE_SIZE pfnSize, WKE_FILE_READ pfnRead, WKE_FILE_SEEK pfnSeek, "") \

    ITERATOR1(const char*, wkeWebViewName, wkeWebView webView, "") \
    ITERATOR2(void, wkeSetWebViewName, wkeWebView webView, const char* name, "") \

    ITERATOR1(BOOL, wkeIsLoaded, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeIsLoadFailed, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeIsLoadComplete, wkeWebView webView, "") \

    ITERATOR1(const utf8*, wkeGetSource, wkeWebView webView, "") \
    ITERATOR1(const utf8*, wkeTitle, wkeWebView webView, "") \
    ITERATOR1(const wchar_t*, wkeTitleW, wkeWebView webView, "") \
    ITERATOR1(int, wkeWidth, wkeWebView webView, "") \
    ITERATOR1(int, wkeHeight, wkeWebView webView, "") \
    ITERATOR1(int, wkeContentsWidth, wkeWebView webView, "") \
    ITERATOR1(int, wkeContentsHeight, wkeWebView webView, "") \

    ITERATOR1(void, wkeSelectAll, wkeWebView webView, "") \
    ITERATOR1(void, wkeCopy, wkeWebView webView, "") \
    ITERATOR1(void, wkeCut, wkeWebView webView, "") \
    ITERATOR1(void, wkePaste, wkeWebView webView, "") \
    ITERATOR1(void, wkeDelete, wkeWebView webView, "") \

    ITERATOR1(BOOL, wkeCookieEnabled, wkeWebView webView, "") \
    ITERATOR1(float, wkeMediaVolume, wkeWebView webView, "") \

    ITERATOR5(BOOL, wkeMouseEvent, wkeWebView webView, unsigned int message, int x, int y, unsigned int flags, "") \
    ITERATOR4(BOOL, wkeContextMenuEvent, wkeWebView webView, int x, int y, unsigned int flags, "") \
    ITERATOR5(BOOL, wkeMouseWheel, wkeWebView webView, int x, int y, int delta, unsigned int flags, "") \
    ITERATOR4(BOOL, wkeKeyUp, wkeWebView webView, unsigned int virtualKeyCode, unsigned int flags, bool systemKey, "") \
    ITERATOR4(BOOL, wkeKeyDown, wkeWebView webView, unsigned int virtualKeyCode, unsigned int flags, bool systemKey, "") \
    ITERATOR4(BOOL, wkeKeyPress, wkeWebView webView, unsigned int virtualKeyCode, unsigned int flags, bool systemKey, "") \

    ITERATOR1(void, wkeFocus, wkeWebView webView, "") \
    ITERATOR1(void, wkeUnfocus, wkeWebView webView, "") \

    ITERATOR1(wkeRect, wkeGetCaret, wkeWebView webView, "") \

    ITERATOR1(void, wkeAwaken, wkeWebView webView, "") \

    ITERATOR1(float, wkeZoomFactor, wkeWebView webView, "") \

    ITERATOR2(void, wkeSetClientHandler, wkeWebView webView, const wkeClientHandler* handler, "") \
    ITERATOR1(const wkeClientHandler*, wkeGetClientHandler, wkeWebView webView, "") \

    ITERATOR1(const utf8*, wkeToString, const wkeString string, "") \
    ITERATOR1(const wchar_t*, wkeToStringW, const wkeString string, "") \

    ITERATOR2(const utf8*, jsToString, jsExecState es, jsValue v, "") \
    ITERATOR2(const wchar_t*, jsToStringW, jsExecState es, jsValue v, "") \

    ITERATOR1(void, wkeConfigure, const wkeSettings* settings, "") \
    ITERATOR0(BOOL, wkeIsInitialize, "") \

    ITERATOR2(void, wkeSetViewSettings, wkeWebView webView, const wkeViewSettings* settings, "") \
    ITERATOR3(void, wkeSetDebugConfig, wkeWebView webView, const char* debugString, const char* param, "") \
    ITERATOR2(void *, wkeGetDebugConfig, wkeWebView webView, const char* debugString, "") \

    ITERATOR0(void, wkeFinalize, "") \
    ITERATOR0(void, wkeUpdate, "") \
    ITERATOR0(unsigned int, wkeGetVersion, "") \
    ITERATOR0(const utf8*, wkeGetVersionString, "") \

    ITERATOR0(wkeWebView, wkeCreateWebView, "") \
    ITERATOR1(void, wkeDestroyWebView, wkeWebView webView, "") \

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

    ITERATOR2(void, wkeSetViewNetInterface, wkeWebView webView, const char* netInterface, "") \

    ITERATOR1(void, wkeSetProxy, const wkeProxy* proxy, "") \
    ITERATOR2(void, wkeSetViewProxy, wkeWebView webView, wkeProxy *proxy, "") \

    ITERATOR1(const char*, wkeGetName, wkeWebView webView, "") \
    ITERATOR2(void, wkeSetName, wkeWebView webView, const char* name, "") \

    ITERATOR2(void, wkeSetHandle, wkeWebView webView, HWND wnd, "") \
    ITERATOR3(void, wkeSetHandleOffset, wkeWebView webView, int x, int y, "") \

    ITERATOR1(BOOL, wkeIsTransparent, wkeWebView webView, "") \
    ITERATOR2(void, wkeSetTransparent, wkeWebView webView, bool transparent, "") \

    ITERATOR2(void, wkeSetUserAgent, wkeWebView webView, const utf8* userAgent, "") \
    ITERATOR1(const char*, wkeGetUserAgent, wkeWebView webView, "") \
    ITERATOR2(void, wkeSetUserAgentW, wkeWebView webView, const wchar_t* userAgent, "") \

    ITERATOR4(void, wkeShowDevtools, wkeWebView webView, const wchar_t* path, wkeOnShowDevtoolsCallback callback, void* param, "") \

    ITERATOR2(void, wkeLoadW, wkeWebView webView, const wchar_t* url, "") \
    ITERATOR2(void, wkeLoadURL, wkeWebView webView, const utf8* url, "") \
    ITERATOR2(void, wkeLoadURLW, wkeWebView webView, const wchar_t* url, "") \
    ITERATOR4(void, wkePostURL, wkeWebView wkeView, const utf8* url, const char* postData, int postLen, "") \
    ITERATOR4(void, wkePostURLW, wkeWebView wkeView, const wchar_t* url, const char* postData, int postLen, "") \

    ITERATOR2(void, wkeLoadHTML, wkeWebView webView, const utf8* html, "") \
    ITERATOR3(void, wkeLoadHtmlWithBaseUrl, wkeWebView webView, const utf8* html, const utf8* baseUrl, "") \
    ITERATOR2(void, wkeLoadHTMLW, wkeWebView webView, const wchar_t* html, "") \

    ITERATOR2(void, wkeLoadFile, wkeWebView webView, const utf8* filename, "") \
    ITERATOR2(void, wkeLoadFileW, wkeWebView webView, const wchar_t* filename, "") \

    ITERATOR1(const utf8*, wkeGetURL, wkeWebView webView, "") \
    ITERATOR2(const utf8*, wkeGetFrameUrl, wkeWebView webView, wkeWebFrameHandle frameId, "") \

    ITERATOR1(BOOL, wkeIsLoading, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeIsLoadingSucceeded, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeIsLoadingFailed, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeIsLoadingCompleted, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeIsDocumentReady, wkeWebView webView, "") \
    ITERATOR1(void, wkeStopLoading, wkeWebView webView, "") \
    ITERATOR1(void, wkeReload, wkeWebView webView, "") \
    ITERATOR2(void, wkeGoToOffset, wkeWebView webView, int offset, "") \
    ITERATOR2(void, wkeGoToIndex, wkeWebView webView, int index, "") \

    ITERATOR1(int, wkeGetWebviewId, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeIsWebviewAlive, int id, "") \
    ITERATOR1(BOOL, wkeIsWebviewValid, wkeWebView webView, "") \

    ITERATOR3(const utf8*, wkeGetDocumentCompleteURL, wkeWebView webView, wkeWebFrameHandle frameId, const utf8* partialURL, "") \

    ITERATOR3(wkeMemBuf*, wkeCreateMemBuf, wkeWebView webView, void* buf, size_t length, "") \
    ITERATOR1(void, wkeFreeMemBuf, wkeMemBuf* buf, "") \

    ITERATOR1(const utf8*, wkeGetTitle, wkeWebView webView, "") \
    ITERATOR1(const wchar_t*, wkeGetTitleW, wkeWebView webView, "") \

    ITERATOR3(void, wkeResize, wkeWebView webView, int w, int h, "") \
    ITERATOR1(int, wkeGetWidth, wkeWebView webView, "") \
    ITERATOR1(int, wkeGetHeight, wkeWebView webView, "") \
    ITERATOR1(int, wkeGetContentWidth, wkeWebView webView, "") \
    ITERATOR1(int, wkeGetContentHeight, wkeWebView webView, "") \

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

    ITERATOR1(BOOL, wkeCanGoBack, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeGoBack, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeCanGoForward, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeGoForward, wkeWebView webView, "") \
    ITERATOR2(BOOL, wkeNavigateAtIndex, wkeWebView webView, int index, "") \
    ITERATOR1(int, wkeGetNavigateIndex, wkeWebView webView, "") \

    ITERATOR1(void, wkeEditorSelectAll, wkeWebView webView, "") \
    ITERATOR1(void, wkeEditorUnSelect, wkeWebView webView, "") \
    ITERATOR1(void, wkeEditorCopy, wkeWebView webView, "") \
    ITERATOR1(void, wkeEditorCut, wkeWebView webView, "") \
    ITERATOR1(void, wkeEditorPaste, wkeWebView webView, "") \
    ITERATOR1(void, wkeEditorDelete, wkeWebView webView, "") \
    ITERATOR1(void, wkeEditorUndo, wkeWebView webView, "") \
    ITERATOR1(void, wkeEditorRedo, wkeWebView webView, "") \

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

    ITERATOR2(void, wkeSetMediaVolume, wkeWebView webView, float volume, "") \
    ITERATOR1(float, wkeGetMediaVolume, wkeWebView webView, "") \

    ITERATOR5(BOOL, wkeFireMouseEvent, wkeWebView webView, unsigned int message, int x, int y, unsigned int flags, "") \
    ITERATOR4(BOOL, wkeFireContextMenuEvent, wkeWebView webView, int x, int y, unsigned int flags, "") \
    ITERATOR5(BOOL, wkeFireMouseWheelEvent, wkeWebView webView, int x, int y, int delta, unsigned int flags, "") \
    ITERATOR4(BOOL, wkeFireKeyUpEvent, wkeWebView webView, unsigned int virtualKeyCode, unsigned int flags, bool systemKey, "") \
    ITERATOR4(BOOL, wkeFireKeyDownEvent, wkeWebView webView, unsigned int virtualKeyCode, unsigned int flags, bool systemKey, "") \
    ITERATOR4(BOOL, wkeFireKeyPressEvent, wkeWebView webView, unsigned int charCode, unsigned int flags, bool systemKey, "") \
    ITERATOR6(BOOL, wkeFireWindowsMessage, wkeWebView webView, HWND hWnd, UINT message, WPARAM wParam, LPARAM lParam, LRESULT* result, "") \

    ITERATOR1(void, wkeSetFocus, wkeWebView webView, "") \
    ITERATOR1(void, wkeKillFocus, wkeWebView webView, "") \

    ITERATOR1(wkeRect, wkeGetCaretRect, wkeWebView webView, "") \
    ITERATOR1(wkeRect*, wkeGetCaretRect2, wkeWebView webView, "给一些不方便获取返回结构体的语言调用") \

    ITERATOR2(jsValue, wkeRunJS, wkeWebView webView, const utf8* script, "") \
    ITERATOR2(jsValue, wkeRunJSW, wkeWebView webView, const wchar_t* script, "") \

    ITERATOR1(jsExecState, wkeGlobalExec, wkeWebView webView, "") \
    ITERATOR2(jsExecState, wkeGetGlobalExecByFrame, wkeWebView webView, wkeWebFrameHandle frameId, "") \

    ITERATOR1(void, wkeSleep, wkeWebView webView, "") \
    ITERATOR1(void, wkeWake, wkeWebView webView, "") \
    ITERATOR1(BOOL, wkeIsAwake, wkeWebView webView, "") \

    ITERATOR2(void, wkeSetZoomFactor, wkeWebView webView, float factor, "") \
    ITERATOR1(float, wkeGetZoomFactor, wkeWebView webView, "") \
    ITERATOR0(void, wkeEnableHighDPISupport, "") \

    ITERATOR2(void, wkeSetEditable, wkeWebView webView, bool editable, "") \

    ITERATOR1(const utf8*, wkeGetString, const wkeString string, "") \
    ITERATOR1(const wchar_t*, wkeGetStringW, const wkeString string, "") \

    ITERATOR3(void, wkeSetString, wkeString string, const utf8* str, size_t len, "") \
    ITERATOR3(void, wkeSetStringWithoutNullTermination, wkeString string, const utf8* str, size_t len, "") \
    ITERATOR3(void, wkeSetStringW, wkeString string, const wchar_t* str, size_t len, "") \

    ITERATOR2(wkeString, wkeCreateString, const utf8* str, size_t len, "") \
    ITERATOR2(wkeString, wkeCreateStringW, const wchar_t* str, size_t len, "") \
    ITERATOR2(wkeString, wkeCreateStringWithoutNullTermination, const utf8* str, size_t len, "") \
    ITERATOR1(size_t, wkeGetStringLen, wkeString str, "") \
    ITERATOR1(void, wkeDeleteString, wkeString str, "") \

    ITERATOR0(wkeWebView, wkeGetWebViewForCurrentContext, "") \
    ITERATOR3(void, wkeSetUserKeyValue, wkeWebView webView, const char* key, void* value, "") \
    ITERATOR2(void*, wkeGetUserKeyValue, wkeWebView webView, const char* key, "") \

    ITERATOR1(int, wkeGetCursorInfoType, wkeWebView webView, "") \
    ITERATOR2(void, wkeSetCursorInfoType, wkeWebView webView, int type, "") \
    ITERATOR5(void, wkeSetDragFiles, wkeWebView webView, const POINT* clintPos, const POINT* screenPos, wkeString* files, int filesCount, "") \

    ITERATOR5(void, wkeSetDeviceParameter, wkeWebView webView, const char* device, const char* paramStr, int paramInt, float paramFloat, "") \
    ITERATOR1(wkeTempCallbackInfo*, wkeGetTempCallbackInfo, wkeWebView webView, "") \

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

    ITERATOR3(void, wkeOnOtherLoad, wkeWebView webView, wkeOnOtherLoadCallback callback, void* param, "") \
    ITERATOR3(void, wkeOnContextMenuItemClick, wkeWebView webView, wkeOnContextMenuItemClickCallback callback, void* param, "") \

    ITERATOR1(BOOL, wkeIsProcessingUserGesture, wkeWebView webView, "") \

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

    ITERATOR1(void, wkeNetContinueJob, wkeNetJob jobPtr, "")\
    ITERATOR1(const char*, wkeNetGetUrlByJob, wkeNetJob jobPtr, "")\
    ITERATOR1(const wkeSlist*, wkeNetGetRawHttpHead, wkeNetJob jobPtr, "")\
    ITERATOR1(const wkeSlist*, wkeNetGetRawResponseHead, wkeNetJob jobPtr, "")\

    ITERATOR1(void, wkeNetCancelRequest, wkeNetJob jobPtr, "")\
    ITERATOR1(BOOL, wkeNetHoldJobToAsynCommit, wkeNetJob jobPtr, "")\
    ITERATOR2(void, wkeNetChangeRequestUrl, wkeNetJob jobPtr, const char* url, "")\

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

    ITERATOR1(wkePostBodyElements*, wkeNetGetPostBody, wkeNetJob jobPtr, "") \
    ITERATOR2(wkePostBodyElements*, wkeNetCreatePostBodyElements, wkeWebView webView, size_t length, "") \
    ITERATOR1(void, wkeNetFreePostBodyElements, wkePostBodyElements* elements, "") \
    ITERATOR1(wkePostBodyElement*, wkeNetCreatePostBodyElement, wkeWebView webView, "") \
    ITERATOR1(void, wkeNetFreePostBodyElement, wkePostBodyElement* element, "") \

    ITERATOR9(wkeDownloadOpt, wkePopupDialogAndDownload, wkeWebView webviewHandle, const wkeDialogOptions* dialogOpt, \
        size_t expectedContentLength, const char* url, const char* mime, const char* disposition, wkeNetJob job,wkeNetJobDataBind* dataBind, wkeDownloadBind* callbackBind, "") \
    ITERATOR10(wkeDownloadOpt, wkeDownloadByPath, wkeWebView webviewHandle, const wkeDialogOptions* dialogOpt, const WCHAR* path, size_t expectedContentLength,const char* url, \
        const char* mime, const char* disposition, wkeNetJob job, wkeNetJobDataBind* dataBind, wkeDownloadBind* callbackBind, "") \

    ITERATOR2(BOOL, wkeIsMainFrame, wkeWebView webView, wkeWebFrameHandle frameId, "") \
    ITERATOR2(BOOL, wkeIsWebRemoteFrame, wkeWebView webView, wkeWebFrameHandle frameId, "") \
    ITERATOR1(wkeWebFrameHandle, wkeWebFrameGetMainFrame, wkeWebView webView, "") \
    ITERATOR4(jsValue, wkeRunJsByFrame, wkeWebView webView, wkeWebFrameHandle frameId, const utf8* script, bool isInClosure, "") \
    ITERATOR3(void, wkeInsertCSSByFrame, wkeWebView webView, wkeWebFrameHandle frameId, const utf8* cssText, "") \

    ITERATOR3(void, wkeWebFrameGetMainWorldScriptContext, wkeWebView webView, wkeWebFrameHandle webFrameId, v8ContextPtr contextOut, "") \

    ITERATOR0(v8Isolate, wkeGetBlinkMainThreadIsolate, "") \

    ITERATOR6(wkeWebView, wkeCreateWebWindow, wkeWindowType type, HWND parent, int x, int y, int width, int height, "") \
    ITERATOR1(wkeWebView, wkeCreateWebCustomWindow, const wkeWindowCreateInfo* info, "") \
    ITERATOR1(void, wkeDestroyWebWindow, wkeWebView webWindow, "") \
    ITERATOR1(HWND, wkeGetWindowHandle, wkeWebView webWindow, "") \

    ITERATOR2(void, wkeShowWindow, wkeWebView webWindow, bool show, "") \
    ITERATOR2(void, wkeEnableWindow, wkeWebView webWindow, bool enable, "") \

    ITERATOR5(void, wkeMoveWindow, wkeWebView webWindow, int x, int y, int width, int height, "") \
    ITERATOR1(void, wkeMoveToCenter, wkeWebView webWindow, "") \
    ITERATOR3(void, wkeResizeWindow, wkeWebView webWindow, int width, int height, "") \

    ITERATOR6(wkeWebDragOperation, wkeDragTargetDragEnter, wkeWebView webView, const wkeWebDragData* webDragData, const POINT* clientPoint, const POINT* screenPoint, wkeWebDragOperationsMask operationsAllowed, int modifiers, "") \
    ITERATOR5(wkeWebDragOperation, wkeDragTargetDragOver, wkeWebView webView, const POINT* clientPoint, const POINT* screenPoint, wkeWebDragOperationsMask operationsAllowed, int modifiers, "") \
    ITERATOR1(void, wkeDragTargetDragLeave, wkeWebView webView, "") \
    ITERATOR4(void, wkeDragTargetDrop, wkeWebView webView, const POINT* clientPoint, const POINT* screenPoint, int modifiers, "") \
    ITERATOR4(void, wkeDragTargetEnd, wkeWebView webView, const POINT* clientPoint, const POINT* screenPoint, wkeWebDragOperation operation, "") \

    ITERATOR1(void, wkeUtilSetUiCallback, wkeUiThreadPostTaskCallback callback, "") \
    ITERATOR1(const utf8*, wkeUtilSerializeToMHTML, wkeWebView webView, "") \
    ITERATOR3(const wkePdfDatas*, wkeUtilPrintToPdf, wkeWebView webView, wkeWebFrameHandle frameId, const wkePrintSettings* settings,"") \
    ITERATOR3(const wkeMemBuf*, wkePrintToBitmap, wkeWebView webView, wkeWebFrameHandle frameId, const wkeScreenshotSettings* settings,"") \
    ITERATOR1(void, wkeUtilRelasePrintPdfDatas, const wkePdfDatas* datas,"") \

    ITERATOR2(void, wkeSetWindowTitle, wkeWebView webWindow, const utf8* title, "") \
    ITERATOR2(void, wkeSetWindowTitleW, wkeWebView webWindow, const wchar_t* title, "") \

    ITERATOR3(void, wkeNodeOnCreateProcess, wkeWebView webView, wkeNodeOnCreateProcessCallback callback, void* param, "") \

    ITERATOR4(void, wkeOnPluginFind, wkeWebView webView, const char* mime, wkeOnPluginFindCallback callback, void* param, "") \
    ITERATOR4(void, wkeAddNpapiPlugin, wkeWebView webView, void* initializeFunc, void* getEntryPointsFunc, void* shutdownFunc, "") \

    ITERATOR4(void, wkePluginListBuilderAddPlugin, void* builder, const utf8* name, const utf8* description, const utf8* fileName, "") \
    ITERATOR3(void, wkePluginListBuilderAddMediaTypeToLastPlugin, void* builder, const utf8* name, const utf8* description, "") \
    ITERATOR2(void, wkePluginListBuilderAddFileExtensionToLastMediaType, void* builder, const utf8* fileExtension, "") \

    ITERATOR1(wkeWebView, wkeGetWebViewByNData, void* ndata, "") \

    ITERATOR5(BOOL, wkeRegisterEmbedderCustomElement, wkeWebView webView, wkeWebFrameHandle frameId, const char* name, void* options, void* outResult, "") \

    ITERATOR3(void, wkeSetMediaPlayerFactory, wkeWebView webView, wkeMediaPlayerFactory factory, wkeOnIsMediaPlayerSupportsMIMEType callback, "") \

    ITERATOR3(const utf8* , wkeGetContentAsMarkup, wkeWebView webView, wkeWebFrameHandle frame, size_t* size, "") \

    ITERATOR1(const utf8*, wkeUtilDecodeURLEscape, const utf8* url, "") \
    ITERATOR1(const utf8*, wkeUtilEncodeURLEscape, const utf8* url, "") \
    ITERATOR1(const utf8*, wkeUtilBase64Encode, const utf8* str, "") \
    ITERATOR1(const utf8*, wkeUtilBase64Decode, const utf8* str, "") \
    ITERATOR1(const wkeMemBuf*, wkeUtilCreateV8Snapshot, const utf8* str, "") \

    ITERATOR0(void, wkeRunMessageLoop, "") \

    ITERATOR1(void, wkeSaveMemoryCache, wkeWebView webView, "") \

    ITERATOR3(void, jsBindFunction, const char* name, jsNativeFunction fn, unsigned int argCount, "") \
    ITERATOR2(void, jsBindGetter, const char* name, jsNativeFunction fn, "") \
    ITERATOR2(void, jsBindSetter, const char* name, jsNativeFunction fn, "") \

    ITERATOR4(void, wkeJsBindFunction, const char* name, wkeJsNativeFunction fn, void* param, unsigned int argCount, "") \
    ITERATOR3(void, wkeJsBindGetter, const char* name, wkeJsNativeFunction fn, void* param, "") \
    ITERATOR3(void, wkeJsBindSetter, const char* name, wkeJsNativeFunction fn, void* param, "") \

    ITERATOR1(int, jsArgCount, jsExecState es, "") \
    ITERATOR2(jsType, jsArgType, jsExecState es, int argIdx, "") \
    ITERATOR2(jsValue, jsArg, jsExecState es, int argIdx, "") \

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

    ITERATOR1(jsValue, jsInt, int n, "") \
    ITERATOR1(jsValue, jsFloat, float f, "") \
    ITERATOR1(jsValue, jsDouble, double d, "") \
    ITERATOR1(jsValue, jsDoubleString, const char* str, "") \
    ITERATOR1(jsValue, jsBoolean, bool b, "") \

    ITERATOR0(jsValue, jsUndefined, "") \
    ITERATOR0(jsValue, jsNull, "") \
    ITERATOR0(jsValue, jsTrue, "") \
    ITERATOR0(jsValue, jsFalse, "") \

    ITERATOR2(jsValue, jsString, jsExecState es, const utf8* str, "") \
    ITERATOR2(jsValue, jsStringW, jsExecState es, const wchar_t* str, "") \
    ITERATOR1(jsValue, jsEmptyObject, jsExecState es, "") \
    ITERATOR1(jsValue, jsEmptyArray, jsExecState es, "") \

    ITERATOR2(jsValue, jsObject, jsExecState es, jsData* obj, "") \
    ITERATOR2(jsValue, jsFunction, jsExecState es, jsData* obj, "") \
    ITERATOR2(jsData*, jsGetData, jsExecState es, jsValue object, "") \

    ITERATOR3(jsValue, jsGet, jsExecState es, jsValue object, const char* prop, "") \
    ITERATOR4(void, jsSet, jsExecState es, jsValue object, const char* prop, jsValue v, "") \

    ITERATOR3(jsValue, jsGetAt, jsExecState es, jsValue object, int index, "") \
    ITERATOR4(void, jsSetAt, jsExecState es, jsValue object, int index, jsValue v, "") \
    ITERATOR2(jsKeys*, jsGetKeys, jsExecState es, jsValue object, "") \
    ITERATOR2(BOOL, jsIsJsValueValid, jsExecState es, jsValue object, "") \
    ITERATOR1(BOOL, jsIsValidExecState, jsExecState es, "") \
    ITERATOR3(void, jsDeleteObjectProp, jsExecState es, jsValue object, const char* prop, "") \

    ITERATOR2(int, jsGetLength, jsExecState es, jsValue object, "") \
    ITERATOR3(void, jsSetLength, jsExecState es, jsValue object, int length, "") \

    ITERATOR1(jsValue, jsGlobalObject, jsExecState es, "") \
    ITERATOR1(wkeWebView, jsGetWebView, jsExecState es, "") \

    ITERATOR2(jsValue, jsEval, jsExecState es, const utf8* str, "") \
    ITERATOR2(jsValue, jsEvalW, jsExecState es, const wchar_t* str, "") \
    ITERATOR3(jsValue, jsEvalExW, jsExecState es, const wchar_t* str, bool isInClosure, "") \

    ITERATOR5(jsValue, jsCall, jsExecState es, jsValue func, jsValue thisObject, jsValue* args, int argCount, "") \
    ITERATOR4(jsValue, jsCallGlobal, jsExecState es, jsValue func, jsValue* args, int argCount, "") \

    ITERATOR2(jsValue, jsGetGlobal, jsExecState es, const char* prop, "") \
    ITERATOR3(void, jsSetGlobal, jsExecState es, const char* prop, jsValue v, "") \

    ITERATOR0(void, jsGC, "") \
    ITERATOR2(BOOL, jsAddRef, jsExecState es, jsValue val, "") \
    ITERATOR2(BOOL, jsReleaseRef, jsExecState es, jsValue val, "") \
    ITERATOR1(jsExceptionInfo*, jsGetLastErrorIfException, jsExecState es, "") \
    ITERATOR2(jsValue, jsThrowException, jsExecState es, const utf8* exception, "") \
    ITERATOR1(const utf8*, jsGetCallstack, jsExecState es, "")
