package util

import (
	"bytes"
	"fmt"
	"reflect"
	"unsafe"
)

func Byte(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func ClearAllZeroBytes(b []byte) []byte {
	i := bytes.IndexByte(b, 0)
	if i < 0 {
		i = len(b)
	}
	return b[:i]
}

func IsArray(val interface{}) bool {
	if _, ok := val.([]interface{}); ok {
		return true
	}
	return false
}

func ParseStringFromInterface(val interface{}) string {
	if ok, v := isStringDirectly(val); ok {
		return v
	}
	if ok, v := isStringThroughPointer(val); ok {
		return v
	}
	panic("not a string")
}

func IsString(val interface{}) (bool, string) {
	if ok, v := isStringDirectly(val); ok {
		return true, v
	}
	if ok, v := isStringThroughPointer(val); ok {
		return true, v
	}
	return false, ""
}

func isStringDirectly(val interface{}) (bool, string) {
	v, ok := val.(string)
	return ok, v
}

func GetTypeOfTheValue(val any) string {
	return fmt.Sprint(reflect.TypeOf(val))
}

func isStringThroughPointer(val interface{}) (bool, string) {
	if ptr, ok := val.(*interface{}); ok {
		if v, ok := (*ptr).(string); ok {
			return true, v
		}
	}
	return false, ""
}

func IsInt(val interface{}) bool {
	if _, ok := val.(int); ok {
		return true
	}
	return false
}

func ConvertInterfaceToStringArr(arr interface{}) []string {
	array := arr.([]interface{})
	stringSlice := make([]string, len(array))
	for i, v := range array {
		// Type assertion for each element
		if str, ok := v.(string); ok {
			stringSlice[i] = str
		} else {
			// Handle the case when the type assertion fails
			fmt.Printf("Element at index %d is not a string\n", i)
		}
	}
	return stringSlice
}