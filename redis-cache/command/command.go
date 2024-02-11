package command

import (
	"fmt"
	"redis-lite/db"
	"redis-lite/resp"
	"redis-lite/util"
	"slices"
	"strconv"
	"strings"
)

var expiryCommands = []string{"EX", "PX"}

func CommandFactory(arr []string) Command {
	head := arr[0]
	switch strings.ToUpper(head) {
	case "PING":
		return Ping{}
	case "SET":
		return Set{}
	case "GET":
		return Get{}
	case "EXISTS":
		return Exists{}
	case "DEL":
		return Del{}
	case "INCR":
		return Incr{}
	case "DECR":
		return Decr{}
	case "LPUSH":
		return LPush{}
	case "LRANGE":
		return LRange{}
	default:
		panic("unrecognized command detected.")
	}
}

type Command interface {
	execute(arr []string) interface{}
}

type Ping struct {
}

func (ping Ping) execute(arr []string) interface{} {
	return resp.Serialize("PONG")
}

type Set struct {
}

func (set Set) execute(arr []string) interface{} {
	var exValue = 0
	var exType string
	for i := 3; i < len(arr); i++ {
		upperVal := strings.ToUpper(arr[i])
		if slices.Contains(expiryCommands, upperVal) {
			if i+1 >= len(arr) {
				panic("if you use EX or PX command with set, you have to provide seconds or milliseconds")
			}
			exSec, err := strconv.Atoi(arr[i+1])
			if err != nil {
				panic("expiration time must be an integer")
			}
			exValue = exSec
			exType = upperVal
		}
	}
	key := arr[1]
	value := arr[2]
	if exValue > 0 {
		db.GetStorage().SetWithOptions(key, value, db.SetOption{ExpireValue: exValue, ExpiryType: exType})
	} else {
		db.GetStorage().Set(key, value)
	}
	return resp.Serialize("OK")
}

type Get struct {
}

func (get Get) execute(arr []string) interface{} {
	value := db.GetStorage().Get(arr[1])
	return resp.Serialize(value)
}

type Exists struct {
}

func (exist Exists) execute(arr []string) interface{} {
	if db.GetStorage().Exists(arr[1]) {
		return resp.Serialize(1)
	}
	return resp.Serialize(0)
}

type Del struct {
}

func (del Del) execute(arr []string) interface{} {
	var removedCount = 0
	for _, s := range arr[1:] {
		b := db.GetStorage().Delete(s)
		if b {
			removedCount++
		}
	}
	return resp.Serialize(removedCount)
}

type Incr struct {
}

func (incr Incr) execute(arr []string) interface{} {
	storage := db.GetStorage()
	key := arr[1]
	value := storage.Get(key)
	if value == nil {
		newValue := 1
		storage.Set(key, strconv.Itoa(newValue))
		return resp.Serialize(newValue)
	}
	valueAsString := util.ParseStringFromInterface(value)
	intValue, err := strconv.Atoi(valueAsString)
	if err != nil {
		panic(fmt.Sprintf("key holds not an integer value. key : %s, value : %s ", arr[1], valueAsString))
	}
	newValue := intValue + 1
	storage.Set(key, strconv.Itoa(newValue))
	return resp.Serialize(newValue)

}

type Decr struct {
}

func (decr Decr) execute(arr []string) interface{} {
	storage := db.GetStorage()
	key := arr[1]
	value := storage.Get(key)
	if value == nil {
		newValue := -1
		storage.Set(key, strconv.Itoa(newValue))
		return resp.Serialize(newValue)
	}
	valueAsString := util.ParseStringFromInterface(value)

	intValue, err := strconv.Atoi(valueAsString)
	if err != nil {
		panic(fmt.Sprintf("key holds not an integer value. key : %s, value : %s ", arr[1], valueAsString))
	}
	newValue := intValue - 1
	storage.Set(key, strconv.Itoa(newValue))
	return resp.Serialize(newValue)

}

type LPush struct {
}

func (lPush LPush) execute(arr []string) interface{} {
	storage := db.GetStorage()
	key := arr[1]
	val := arr[2]
	value := storage.Get(key)
	if value == nil {
		newValue := []interface{}{val}
		storage.SetArray(key, newValue)
		return resp.Serialize(len(newValue))
	}
	if !util.IsArray(value) {
		panic(fmt.Sprintf("key does not hold an array. key: %s, value: %s", key, value))
	}
	array := storage.Get(key).([]interface{})
	newArray := make([]interface{}, 0)
	newArray = append(newArray, val)
	for _, v := range array {
		newArray = append(newArray, v)
	}
	storage.SetArray(key, newArray)
	return resp.Serialize(len(newArray))
}

type LRange struct {
}

func (lRange LRange) execute(arr []string) interface{} {
	storage := db.GetStorage()
	key := arr[1]
	start, err := strconv.Atoi(arr[2])
	stop, err2 := strconv.Atoi(arr[3])
	if err != nil || err2 != nil {
		panic(fmt.Sprintf("start and stop values must be integer"))
	}
	value := storage.Get(key)
	if value == nil {
		return resp.Serialize(nil)
	}
	if !util.IsArray(value) {
		panic(fmt.Sprintf("key does not hold an array. key: %s, value: %s", key, value))
	}
	array := storage.Get(key).([]interface{})

	if start > len(array)-1 || start > stop {
		return resp.Serialize(make([]interface{}, 0))
	}
	if stop > len(array)-1 {
		stop = len(array) - 1
	}
	if start < 0 && stop < 0 {
		start = len(array) + start
		stop = len(array) + stop
	} else if start < len(array)*-1 {
		start = 0
	}
	return resp.Serialize(array[start : stop+1])
}