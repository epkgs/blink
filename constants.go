package blink

const DLL_FILE = "blink.dll" // DLL file name

var (
	// AppID is the unique id of the app.
	AppID = RandString(8)
)
