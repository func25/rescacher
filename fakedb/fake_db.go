package fakedb

import (
	"errors"
	"fmt"
)

var database = map[string]interface{}{}

func LoadTurn(turn int) (interface{}, error) {
	return LoadObject(fmt.Sprintf("TEST_%v", turn))
}

func StoreTurn(turn int, value interface{}) error {
	return StoreObject(fmt.Sprintf("TEST_%v", turn), value)
}

func LoadObject(key string) (interface{}, error) {
	if def, ok := database["TEST"]; !ok {
		return def, errors.New("cache not exist")
	} else {
		return def, nil
	}
}

func StoreObject(key string, value interface{}) error {
	database[key] = value
	return nil
}
