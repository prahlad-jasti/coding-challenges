package resp

import (
	"fmt"
	"redis-lite/util"
	"reflect"
	"strconv"
	"strings"
)

func ParseResp(splitCommand [][]byte) (interface{}, int) {
	firstElement := splitCommand[0][0]
	switch firstElement {
	case util.Colon:
		{
			return parseNumber(splitCommand)
		}
	case util.Plus, util.Dash:
		{
			return parseSimple(splitCommand)
		}
	case util.DollarSign:
		{
			return parseBulkString(splitCommand)
		}
	case util.Star:
		{
			return parseArray(splitCommand)
		}
	}
	panic("unknown command.")
}

func parseSimple(command [][]byte) (string, int) {
	return string(command[0][1:]), 0
}

func parseNumber(command [][]byte) (string, int) {
	return string(command[0][1:]), 0
}

func parseBulkString(command [][]byte) (string, int) {
	lengthOfTheString, _ := strconv.Atoi(string(command[0][1]))
	if lengthOfTheString == -1 {
		return "", 0
	}
	return string(command[1]), 1
}

func parseArray(command [][]byte) ([]interface{}, int) {
	sizeOfTheArray, _ := strconv.Atoi(string(command[0][1]))
	var arr = make([]interface{}, sizeOfTheArray)
	index := 0
	for i := 1; i <= sizeOfTheArray; i++ {
		val, idx := ParseResp(command[i+index:])
		arr[i-1] = val
		index += idx
	}
	return arr, sizeOfTheArray + index
}

func Serialize(val interface{}) string {
	if util.IsInt(val) {
		value := (val).(int)
		return fmt.Sprintf(":%d%s", value, util.CRLF)
	} else if val == nil || (reflect.ValueOf(val).IsZero() && reflect.ValueOf(val).IsNil()) {
		return fmt.Sprintf("$-1\r\n")
	} else if ok, v := util.IsString(val); ok {
		return fmt.Sprintf("$%d%s%s%s", len(v), util.CRLF, v, util.CRLF)
	} else if util.IsArray(val) {
		value := (val).([]interface{})
		return serializeArray(value)
	}
	panic("unknown serialization input")
}

func serializeArray(val []interface{}) string {
	var sb = strings.Builder{}
	sb.WriteString("*")
	sb.WriteString(strconv.Itoa(len(val)))
	sb.WriteString(util.CRLF)
	for i := 0; i < len(val); i++ {
		serializedData := Serialize(val[i])
		sb.WriteString(serializedData)
	}
	return sb.String()
}