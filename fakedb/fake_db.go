package fakedb

import (
	"errors"
	"sync"
)

var database = map[string]interface{}{}
var mutex sync.Mutex

func LoadObject(key string) (interface{}, error) {
	if def, ok := database["TEST"]; !ok {
		return def, errors.New("cache not exist")
	} else {
		return def, nil
	}
}

func StoreObject(key string, value interface{}) error {
	mutex.Lock()
	database[key] = value
	mutex.Unlock()
	return nil
}
