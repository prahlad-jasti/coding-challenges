package db

import (
	"redis-lite/util"
	"reflect"
	"sync"
	"time"
)

var storage = Storage{db: make(map[string]*RedisObj), mu: sync.Mutex{}}

type Storage struct {
	mu sync.Mutex
	db map[string]*RedisObj
}

type RedisObj struct {
	Data interface{}
	Typ  string
	TTL  int64
}

func (redisObj *RedisObj) clearIfTTLExpired(key string) {
	if redisObj.TTL != 0 && time.Now().Unix() > redisObj.TTL {
		GetStorage().Delete(key)
	}
}

type SetOption struct {
	ExpiryType  string
	ExpireValue int
}

func (so *SetOption) GetDuration() time.Duration {
	if so.ExpiryType == "EX" {
		return time.Second * time.Duration(so.ExpireValue)
	} else if so.ExpiryType == "PX" {
		return time.Millisecond * time.Duration(so.ExpireValue)
	}
	panic("expire type not defined.")
}

func (storage *Storage) SetWithOptions(key string, value string, setOption SetOption) {
	storage.mu.Lock()
	redisObj := &RedisObj{Data: value, Typ: reflect.TypeOf(value).String()}
	if setOption.ExpireValue > 0 {
		newTTL := time.Now().Add(setOption.GetDuration())
		redisObj.TTL = newTTL.Unix()
	}
	storage.db[key] = redisObj
	defer storage.mu.Unlock()
}

func (storage *Storage) Set(key string, value string) {
	storage.mu.Lock()
	redisObj := &RedisObj{Data: value, Typ: reflect.TypeOf(value).String()}
	storage.db[key] = redisObj
	defer storage.mu.Unlock()
}

func (storage *Storage) SetArray(key string, value []interface{}) {
	storage.mu.Lock()
	redisObj := RedisObj{Data: value, Typ: reflect.TypeOf(value).String()}
	storage.db[key] = &redisObj
	defer storage.mu.Unlock()
}

func (storage *Storage) Get(key string) interface{} {
	if val, ok := storage.db[key]; ok {
		val.clearIfTTLExpired(key)
		if util.IsArray(val.Data) {
			return val.Data.([]interface{})
		}
		return &val.Data
	}
	return nil
}

func (storage *Storage) Delete(key string) bool {
	if storage.Exists(key) {
		storage.mu.Lock()
		delete(storage.db, key)
		defer storage.mu.Unlock()
		return true
	} else {
		return false
	}
}

func (storage *Storage) Exists(key string) bool {
	if _, ok := storage.db[key]; ok {
		return true
	}
	return false
}

func GetStorage() *Storage {
	return &storage
}