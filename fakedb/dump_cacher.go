package fakedb

import (
	"fmt"
	"time"

	"github.com/func25/rescacher"
)

type DumpCacher struct{}

var _ rescacher.ICacher = (*DumpCacher)(nil)

func (*DumpCacher) Load(turn int) (interface{}, error) {
	return LoadObject(fmt.Sprintf("TEST_%v", turn))
}

func (*DumpCacher) Save(turn int, value interface{}) error {

	err := StoreObject(fmt.Sprintf("TEST_%v", turn), value)
	if err == nil {
		StoreObject("TEST", turn)
	}

	return err
}

func (*DumpCacher) Generate(turn int) (interface{}, error) {
	return time.Now().Unix(), nil
}

func (*DumpCacher) GetCachedTurn() (int, error) {
	if v, err := LoadObject("TEST"); err == nil {
		return 0, err
	} else {
		res, _ := v.(int)
		return res, nil
	}
}

func (*DumpCacher) SetCachedTurn(turn int) error {
	return StoreObject("TEST", turn)
}
