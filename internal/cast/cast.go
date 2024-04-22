package cast

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"unsafe"
)

func StrToInt64(value string) int64 {
	v, _ := strconv.ParseInt(value, 10, 64)
	return v
}
func StrToInt32(value string) int32 {
	v, _ := strconv.ParseInt(value, 10, 32)
	return int32(v)
}

func StrToFloat64(value string) float64 {
	v, _ := strconv.ParseFloat(value, 64)
	return v
}

func StrToFloat32(value string) float32 {
	v, _ := strconv.ParseFloat(value, 32)
	return float32(v)
}

// GetParamOf 获取参数指针
func GetParamOf(index int, ptr uintptr) uintptr {
	return *(*uintptr)(unsafe.Pointer(ptr + uintptr(index)*unsafe.Sizeof(ptr)))
}

// GetParamPtr 根据指定指针位置开始 偏移获取指针
func GetParamPtr(ptr uintptr, offset int) unsafe.Pointer {
	return unsafe.Pointer(ptr + uintptr(offset))
}

// ToString 接口转 string
func ToString(value interface{}) string {
	return fmt.Sprintf("%v", value)
}

// ToBool bool
func ToBool(v interface{}) bool {
	switch v := v.(type) {
	case []byte:
		if len(v) == 1 {
			return ByteToInt8(v[0]) > 0
		} else if len(v) == 2 {
			return BytesToInt16(v) > 0
		} else if len(v) == 4 {
			return BytesToInt32(v) > 0
		} else if len(v) == 8 {
			return BytesToInt64(v) > 0
		}
		return len(v) > 0
	case string:
		return len(v) > 0
	case float32:
		return v > 0
	case float64:
		return v > 0
	case bool:
		return v
	case int:
		return v > 0
	case int8:
		return v > 0
	case int16:
		return v > 0
	case int32:
		return v > 0
	case int64:
		return v > 0
	case uintptr:
		return v > 0
	default:
		return false
	}
}

func ToFloat64(v interface{}) float64 {
	switch v := v.(type) {
	case []byte:
		if len(v) == 4 {
			return float64(BytesToFloat32(v))
		} else if len(v) == 8 {
			return BytesToFloat64(v)
		}
		return 0.0
	case string:
		return StrToFloat64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	case bool:
		if v {
			return 1
		} else {
			return 0
		}
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uintptr:
		return float64(v)
	default:
		return 0
	}
}

type Number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 | uintptr |
		float32 | float64
}

func ToNumber[T Number](v interface{}) (result T, ok bool) {
	ok = true
	switch v := v.(type) {
	case []byte:
		if len(v) == 1 {
			result = T(ByteToInt8(v[0]))
		} else if len(v) == 2 {
			result = T(BytesToInt16(v))
		} else if len(v) == 4 {
			result = T(BytesToInt32(v))
		} else if len(v) == 8 {
			result = T(BytesToInt64(v))
		}
	case string:
		result = T(StrToInt64(v))
	case float32:
		result = T(v)
	case float64:
		result = T(v)
	case bool:
		if v {
			result = 1
		} else {
			result = 0
		}
	case int:
		result = T(v)
	case int8:
		result = T(v)
	case int16:
		result = T(v)
	case int32:
		result = T(v)
	case int64:
		result = T(v)
	case uint:
		result = T(v)
	case uint8:
		result = T(v)
	case uint16:
		result = T(v)
	case uint32:
		result = T(v)
	case uint64:
		result = T(v)
	case uintptr:
		result = T(v)
	default:
		ok = false
	}

	return
}

func IntToBytes(i int) []byte {
	buf := bytes.NewBuffer([]byte{})
	if strconv.IntSize == 32 {
		if err := binary.Write(buf, binary.BigEndian, int32(i)); err == nil {
			return buf.Bytes()
		}
	} else {
		if err := binary.Write(buf, binary.BigEndian, int64(i)); err == nil {
			return buf.Bytes()
		}
	}
	return nil
}

func UIntToBytes(i uint) []byte {
	buf := bytes.NewBuffer([]byte{})
	if strconv.IntSize == 32 {
		if err := binary.Write(buf, binary.BigEndian, uint32(i)); err == nil {
			return buf.Bytes()
		}
	} else {
		if err := binary.Write(buf, binary.BigEndian, uint64(i)); err == nil {
			return buf.Bytes()
		}
	}
	return nil
}

func Int8ToBytes(i int8) []byte {
	return []byte{byte(i)}
}

func UInt8ToBytes(i uint8) []byte {
	return []byte{byte(i)}
}

func Int16ToBytes(i int16) []byte {
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.BigEndian, i); err == nil {
		return buf.Bytes()
	}
	return nil
}

func UInt16ToBytes(i uint16) []byte {
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.BigEndian, i); err == nil {
		return buf.Bytes()
	}
	return nil
}

func Int32ToBytes(i int32) []byte {
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.BigEndian, i); err == nil {
		return buf.Bytes()
	}
	return nil
}

func UInt32ToBytes(i uint32) []byte {
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.BigEndian, i); err == nil {
		return buf.Bytes()
	}
	return nil
}

func Int64ToBytes(i int64) []byte {
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.BigEndian, i); err == nil {
		return buf.Bytes()
	}
	return nil
}

func UInt64ToBytes(i uint64) []byte {
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.BigEndian, i); err == nil {
		return buf.Bytes()
	}
	return nil
}

func BytesToInt(b []byte) int {
	var i int64
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &i)
	if err != nil {
		return 0
	}
	return int(i)
}

func BytesToUInt(b []byte) uint {
	var i uint64
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &i)
	if err != nil {
		return 0
	}
	return uint(i)
}

func ByteToInt8(b byte) int8 {
	return int8(b)
}

func ByteToUInt8(b byte) uint8 {
	return uint8(b)
}

func BytesToInt16(b []byte) int16 {
	var i int16
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &i)
	if err != nil {
		return 0
	}
	return i
}

func BytesToUInt16(b []byte) uint16 {
	var i uint16
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &i)
	if err != nil {
		return 0
	}
	return i
}

func BytesToInt32(b []byte) int32 {
	var i int32
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &i)
	if err != nil {
		return 0
	}
	return i
}

func BytesToUInt32(b []byte) uint32 {
	var i uint32
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &i)
	if err != nil {
		return 0
	}
	return i
}

func BytesToInt64(b []byte) int64 {
	var i int64
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &i)
	if err != nil {
		return 0
	}
	return i
}

func BytesToUInt64(b []byte) uint64 {
	var i uint64
	err := binary.Read(bytes.NewReader(b), binary.BigEndian, &i)
	if err != nil {
		return 0
	}
	return i
}

func BytesToString(data []byte) string {
	return *(*string)(unsafe.Pointer(&data))
}

//func StringToBytes(data string) []byte {
//	return *(*[]byte)(unsafe.Pointer(&data))
//}

// String转换Bytes数组，isDStr转换DString 默认GoString
func StringToBytes(s string, isDStr ...bool) []byte {
	if len(isDStr) > 0 && isDStr[0] {
		temp := []byte(s)
		utf8StrArr := make([]byte, len(temp)+1)
		copy(utf8StrArr, temp)
		return utf8StrArr
	} else {
		return []byte(s)
	}
}

// Float64ToBytes Float64转byte
func Float64ToBytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

// BytesToFloat64 byte转Float64
func BytesToFloat64(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	return math.Float64frombits(bits)
}

// Float32ToBytes Float64转byte
func Float32ToBytes(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}

// BytesToFloat32 byte转Float64
func BytesToFloat32(bytes []byte) float32 {
	bits := binary.LittleEndian.Uint32(bytes)
	return math.Float32frombits(bits)
}

func BoolToByte(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func ByteToBool(b byte) bool {
	if b == 1 {
		return true
	}
	return false
}

// 将map[string]interface{}转换为结构体
func MapToStruct(m map[string]interface{}, s interface{}) error {
	structValue, ok := s.(reflect.Value)
	if !ok {
		sValue := reflect.ValueOf(s)
		if sValue.Kind() != reflect.Ptr || sValue.Elem().Kind() != reflect.Struct {
			return fmt.Errorf("s must be a pointer to a struct")
		}
		structValue = sValue.Elem()
	}

	structType := structValue.Type()

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		// 处理嵌套结构体
		if fieldType.Type.Kind() == reflect.Struct {
			// 如果字段是结构体，则递归调用mapToStruct进行填充
			nestedStructValue := reflect.New(fieldType.Type).Elem()
			if err := MapToStruct(m, nestedStructValue); err != nil {
				return err
			}
			field.Set(nestedStructValue)
		} else if fieldType.Type.Kind() == reflect.Ptr && fieldType.Type.Elem().Kind() == reflect.Struct {
			// 处理指针类型的嵌套结构体
			mapValue, ok := m[fieldType.Name]
			if !ok {
				continue // 字段不存在于map中，跳过
			}

			// 如果传入的参数是指针类型的map，则解引用取得实际的map[string]interface{}值
			if reflect.TypeOf(mapValue).Kind() == reflect.Ptr {
				mapValue = reflect.ValueOf(mapValue).Elem().Interface()
			}

			// 创建新的结构体实例并填充数据
			nestedStructValue := reflect.New(fieldType.Type.Elem()).Elem()
			if err := MapToStruct(mapValue.(map[string]interface{}), nestedStructValue); err != nil {
				return err
			}
			// 将新创建的结构体指针赋值给字段
			field.Set(reflect.New(fieldType.Type.Elem()))
			field.Elem().Set(nestedStructValue)
		} else {
			// 处理普通字段
			mapValue, ok := m[fieldType.Name]
			if !ok {
				continue // 字段不存在于map中，跳过
			}

			// 检查字段类型是否匹配
			if reflect.TypeOf(mapValue).AssignableTo(fieldType.Type) {
				field.Set(reflect.ValueOf(mapValue))
			} else {
				return fmt.Errorf("Type mismatch for field '%s'", fieldType.Name)
			}
		}
	}

	return nil
}

// 将结构体转换为 map
func StructToMap(s interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	// 获取结构体的类型信息
	structValue := reflect.ValueOf(s)
	structType := structValue.Type()

	// 遍历结构体的字段
	for i := 0; i < structValue.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		// 如果字段是结构体类型，则递归地将其转换为 map
		if fieldValue.Kind() == reflect.Struct {
			nestedMap := StructToMap(fieldValue.Interface())
			result[field.Name] = nestedMap
		} else {
			// 否则直接添加到 map 中
			result[field.Name] = fieldValue.Interface()
		}
	}

	return result
}

// 调用 callback
func Param(param reflect.Type, input interface{}) (reflectVal reflect.Value, err error) {

	pKind := param.Kind()

	switch pKind {
	case reflect.String:
		reflectVal = reflect.ValueOf(ToString(input))
	case reflect.Int:
		val, _ := ToNumber[int](input)
		reflectVal = reflect.ValueOf(val)
	case reflect.Int8:
		val, _ := ToNumber[int8](input)
		reflectVal = reflect.ValueOf(val)
	case reflect.Int16:
		val, _ := ToNumber[int16](input)
		reflectVal = reflect.ValueOf(val)
	case reflect.Int32:
		val, _ := ToNumber[int32](input)
		reflectVal = reflect.ValueOf(val)
	case reflect.Int64:
		val, _ := ToNumber[int64](input)
		reflectVal = reflect.ValueOf(val)
	case reflect.Uint:
		val, _ := ToNumber[uint](input)
		reflectVal = reflect.ValueOf(val)
	case reflect.Uint8:
		val, _ := ToNumber[uint8](input)
		reflectVal = reflect.ValueOf(val)
	case reflect.Uint16:
		val, _ := ToNumber[uint16](input)
		reflectVal = reflect.ValueOf(val)
	case reflect.Uint32:
		val, _ := ToNumber[uint32](input)
		reflectVal = reflect.ValueOf(val)
	case reflect.Uint64:
		val, _ := ToNumber[uint64](input)
		reflectVal = reflect.ValueOf(val)
	case reflect.Uintptr:
		val, _ := ToNumber[uintptr](input)
		reflectVal = reflect.ValueOf(val)
	case reflect.Float32:
		val, _ := ToNumber[float32](input)
		reflectVal = reflect.ValueOf(val)
	case reflect.Float64:
		val, _ := ToNumber[float64](input)
		reflectVal = reflect.ValueOf(val)
	case reflect.Bool:
		reflectVal = reflect.ValueOf(ToBool(input))
	case reflect.Ptr:
		if param.Elem().Kind() == reflect.Struct {
			// 如果参数是指针类型的结构体，并且传入的是*map[string]interface{}，则进行转换
			inputKind := reflect.TypeOf(input).Kind()
			var m map[string]interface{}
			if inputKind == reflect.Ptr {
				m = *input.(*map[string]interface{})
			} else {
				m = input.(map[string]interface{})
			}
			reflectVal = reflect.New(param.Elem()).Elem()
			if err = MapToStruct(m, reflectVal); err != nil {
				return
			}
		} else {
			// 仅支持 struct 指针
			err = fmt.Errorf("failed to convert input to pointer")
			return
		}
	case reflect.Struct:
		inputKind := reflect.TypeOf(input).Kind()
		if inputKind == reflect.Map {
			// 如果参数是结构体且传入的是map[string]interface{}，则进行转换
			reflectVal = reflect.New(param).Elem()
			if err = MapToStruct(input.(map[string]interface{}), reflectVal); err != nil {
				return
			}
		} else if inputKind == reflect.Ptr {
			// 如果参数是指针类型的结构体，并且传入的是*map[string]interface{}，则进行转换
			reflectVal = reflect.New(param).Elem()
			if err = MapToStruct(*input.(*map[string]interface{}), reflectVal); err != nil {
				return
			}
		} else if inputKind == reflect.Struct {
			reflectVal = reflect.ValueOf(input)
		} else {
			err = fmt.Errorf("failed to convert input to struct")
			return
		}

	case reflect.Map:
		inputKind := reflect.TypeOf(input).Kind()
		switch inputKind {
		case reflect.Map:
			// 如果参数是map，则直接使用
			reflectVal = reflect.ValueOf(input)
		case reflect.Ptr:
			// 如果参数是指针类型，则解引用
			reflectVal = reflect.ValueOf(input).Elem()
		default:
			err = fmt.Errorf("failed to convert input to map")
			return
		}
	case reflect.Slice:
		inputKind := reflect.TypeOf(input).Kind()
		switch inputKind {
		case reflect.Slice:
			// 如果输入是slice，则直接使用
			reflectVal = reflect.ValueOf(input)
		case reflect.Ptr:
			// 如果输入是指针类型，则解引用
			reflectVal = reflect.ValueOf(input).Elem()
		default:
			err = fmt.Errorf("type mismatch, expect %v input %v(%v)", pKind, input, inputKind)
			return
		}

	default:
		reflectVal = reflect.ValueOf(input)
	}

	return
}
