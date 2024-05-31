package blink

import (
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type BindFunctionCallback func(es JsExecState)

type JS struct {
	mb *Blink
}

func newJS(blink *Blink) *JS {
	js := &JS{
		mb: blink,
	}

	return js
}

func (js *JS) bindFunction(funcName string, funcArgCount uint32, callback BindFunctionCallback) {
	var cb WkeJsNativeFunction = func(es JsExecState, param uintptr) (voidRes uintptr) {

		callback(es)

		return 0
	}
	js.mb.CallFunc("wkeJsBindFunction", StringToPtr(funcName), CallbackToPtr(cb), 0, uintptr(funcArgCount))
}

// 获取页面主frame的jsExecState
func (js *JS) GlobalExec(viewHandle WkeHandle) (es JsExecState) {

	ptr, _, _ := js.mb.CallFunc("wkeGlobalExec", uintptr(viewHandle))

	return JsExecState(ptr)
}

func (js *JS) GetWebView(es JsExecState) WkeHandle {
	p, _, _ := js.mb.CallFunc("jsGetWebView", uintptr(es))
	return WkeHandle(p)
}

func (js *JS) TypeOf(val JsValue) JsType {
	ptr, _, _ := js.mb.CallFunc("jsTypeOf", uintptr(val))
	return JsType(ptr)
}

func (js *JS) Undefined() JsValue {
	r, _, _ := js.mb.CallFunc("jsUndefined")
	return JsValue(r)
}

func (js *JS) Int(value int32) JsValue {
	r, _, _ := js.mb.CallFunc("jsInt", uintptr(value))
	return JsValue(r)
}

func (js *JS) Double(value float64) JsValue {
	r, _, _ := js.mb.CallFunc("jsDouble", uintptr(value))
	return JsValue(r)
}

func (js *JS) Boolean(value bool) JsValue {
	r, _, _ := js.mb.CallFunc("jsBoolean", BoolToPtr(value))
	return JsValue(r)
}

func (js *JS) ArgCount(es JsExecState) uint32 {
	ptr, _, _ := js.mb.CallFunc("jsArgCount", uintptr(es))
	return uint32(ptr)
}

// 判断第argIdx个参数的参数类型。argIdx从是个0开始计数的值。如果超出jsArgCount返回的值，将发生崩溃
func (js *JS) ArgType(es JsExecState, argIdx uint32) JsType {
	ptr, _, _ := js.mb.CallFunc("jsArgType", uintptr(es), uintptr(argIdx))
	return JsType(ptr)
}

// 获取第argIdx对应的参数的jsValue值。
func (js *JS) Arg(es JsExecState, argIdx uint32) JsValue {
	ptr, _, _ := js.mb.CallFunc("jsArg", uintptr(es), uintptr(argIdx))
	return JsValue(ptr)
}

// str的代码会在mb内部自动被包裹在一个function(){}中。所以使用的变量会被隔离 注意：要获取返回值，请写return。
func (js *JS) Eval(es JsExecState, str string) JsValue {
	ptr, _, _ := js.mb.CallFunc("jsEvalW", uintptr(es), StringToWCharPtr(str))
	return JsValue(ptr)
}

// 如果object是个js的object，则获取prop指定的属性。如果object不是js object类型，则返回 nil
func (js *JS) Get(es JsExecState, object JsValue, prop string) JsValue {

	ptr, _, _ := js.mb.CallFunc("jsGet", uintptr(es), uintptr(object), StringToPtr(prop))

	return JsValue(ptr)
}

// 设置object的属性
func (js *JS) Set(es JsExecState, object JsValue, prop string, value JsValue) {

	js.mb.CallFunc("jsSet", uintptr(es), uintptr(object), StringToPtr(prop), uintptr(value))
}

// 获取window上的属性
func (js *JS) GetGlobal(es JsExecState, prop string) JsValue {

	ptr, _, _ := js.mb.CallFunc("jsGetGlobal", uintptr(es), StringToPtr(prop))

	return JsValue(ptr)
}

// 设置window上的属性
func (js *JS) SetGlobal(es JsExecState, prop string, value JsValue) {

	js.mb.CallFunc("jsSetGlobal", uintptr(es), StringToPtr(prop), uintptr(value))
}

// 设置js arrary的第index个成员的值，object必须是js array才有用，否则会返回nil
func (js *JS) GetAt(es JsExecState, object JsValue, index uint32) JsValue {
	p, _, _ := js.mb.CallFunc("jsGetAt", uintptr(es), uintptr(object), uintptr(index))
	return JsValue(p)
}

// 设置js arrary的第index个成员的值，object必须是js array才有用。
func (js *JS) SetAt(es JsExecState, object JsValue, index uint32, value JsValue) {

	js.mb.CallFunc("jsSetAt", uintptr(es), uintptr(object), uintptr(index), uintptr(value))
}

// 获取object有哪些key
func (js *JS) GetKeys(es JsExecState, object JsValue) []string {

	p, _, _ := js.mb.CallFunc("jsGetKeys", uintptr(es), uintptr(object))

	keys := *((*JsKeys)(unsafe.Pointer(p)))

	items := make([]string, keys.Length)
	for i := 0; i < len(items); i++ {
		items[i] = PtrToString(*((*uintptr)(unsafe.Pointer(keys.First))))
		keys.First += unsafe.Sizeof(uintptr(0))
	}
	return items
}

func (js *JS) String(es JsExecState, value string) JsValue {
	r, _, _ := js.mb.CallFunc("jsString", uintptr(es), StringToPtr(value))
	return JsValue(r)
}
func (js *JS) EmptyArray(es JsExecState) JsValue {
	r, _, _ := js.mb.CallFunc("jsEmptyArray", uintptr(es))
	return JsValue(r)
}
func (js *JS) EmptyObject(es JsExecState) JsValue {
	r, _, _ := js.mb.CallFunc("jsEmptyObject", uintptr(es))
	return JsValue(r)
}

// 获取js arrary的长度，object必须是js array才有用。
func (js *JS) GetLength(es JsExecState, object JsValue) int {
	p, _, _ := js.mb.CallFunc("jsGetLength", uintptr(es), uintptr(object))
	return int(p)
}

func (js *JS) SetLength(es JsExecState, object JsValue, length uint32) {
	js.mb.CallFunc("jsSetLength", uintptr(es), uintptr(object), uintptr(length))
}

func (js *JS) ToDouble(es JsExecState, value JsValue) float64 {
	p, _, _ := js.mb.CallFunc("jsToDouble", uintptr(es), uintptr(value))
	return float64(p)
}

func (js *JS) ToBoolean(es JsExecState, value JsValue) bool {
	p, _, _ := js.mb.CallFunc("jsToBoolean", uintptr(es), uintptr(value))
	return p != 0
}

func (js *JS) ToTempString(es JsExecState, value JsValue) string {
	p, _, _ := js.mb.CallFunc("jsToTempString", uintptr(es), uintptr(value))
	return PtrToString(p)
}

func (js *JS) ToString(es JsExecState, value JsValue) string {
	p, _, _ := js.mb.CallFunc("jsToString", uintptr(es), uintptr(value))
	return PtrToString(p)
}

func (js *JS) Call(es JsExecState, fn, thisValue JsValue, args []JsValue) JsValue {
	var ptr = uintptr(0)
	l := len(args)
	if l > 0 {
		ptr = uintptr(unsafe.Pointer(&args[0]))
	}

	r, _, _ := js.mb.CallFunc("jsCall", uintptr(es), uintptr(fn), uintptr(thisValue), ptr, uintptr(l))
	return JsValue(r)
}

func (js *JS) ToJsValue(es JsExecState, value interface{}) JsValue {
	if value == nil {
		return js.Undefined()
	}
	switch val := value.(type) {
	case int:
		return js.Int(int32(val))
	case int8:
		return js.Int(int32(val))
	case int16:
		return js.Int(int32(val))
	case int32:
		return js.Int(val)
	case int64:
		return js.Double(float64(val))
	case uint:
		return js.Int(int32(val))
	case uint8:
		return js.Int(int32(val))
	case uint16:
		return js.Int(int32(val))
	case uint32:
		return js.Int(int32(val))
	case uint64:
		return js.Double(float64(val))
	case float32:
		return js.Double(float64(val))
	case float64:
		return js.Double(val)
	case bool:
		return js.Boolean(val)
	case string:
		return js.String(es, val)
	case time.Time:
		return js.Double(float64(val.Unix()))
	default:
		break
	}
	rt := reflect.TypeOf(value)
	rv := reflect.ValueOf(value)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		length := rv.Len()
		arr := js.EmptyArray(es)
		js.SetLength(es, arr, uint32(length))
		for i := 0; i < length; i++ {
			v := js.ToJsValue(es, rv.Index(i).Interface())
			js.SetAt(es, arr, uint32(i), v)
		}
		return arr
	case reflect.Map:
		obj := js.EmptyObject(es)
		kv := rv.MapRange()
		for kv.Next() && kv.Key().Kind() == reflect.String {
			k := kv.Key().Interface().(string)
			v := js.ToJsValue(es, kv.Value().Interface())
			js.Set(es, obj, k, v)
		}
		return obj
	case reflect.Struct:
		obj := js.EmptyObject(es)
		for i := 0; i < rv.NumField(); i++ {
			f := rt.Field(i)
			if strings.ToUpper(f.Name)[0] == f.Name[0] {
				fname := rt.Field(i).Name
				fvalue := rv.Field(i).Interface()
				v := js.ToJsValue(es, fvalue)
				js.Set(es, obj, fname, v)
			}
		}
		return obj
	}
	// TODO: 移除 panic，应该使用返回 error
	panic("不支持的go类型：" + rv.Kind().String() + "(" + rv.Type().String() + ")")
}

func (js *JS) ToGoValue(es JsExecState, value JsValue) interface{} {
	switch js.TypeOf(value) {
	case JsType_NULL, JsType_UNDEFINED:
		return nil
	case JsType_NUMBER:
		return js.ToDouble(es, value)
	case JsType_BOOLEAN:
		return js.ToBoolean(es, value)
	case JsType_STRING:
		return js.ToTempString(es, value)
	case JsType_ARRAY:
		length := js.GetLength(es, value)
		ps := make([]interface{}, length)
		for i := 0; i < length; i++ {
			v := js.GetAt(es, value, uint32(i))
			ps[i] = js.ToGoValue(es, v)
		}
		return ps
	case JsType_OBJECT:
		ps := make(map[string]interface{})
		keys := js.GetKeys(es, value)
		for _, k := range keys {
			v := js.Get(es, value, k)
			ps[k] = js.ToGoValue(es, v)
		}
		return ps
	default:
		// TODO: 移除 panic，应该使用返回 error
		panic("不支持的js类型：" + strconv.Itoa(int(value)))
	}
}
