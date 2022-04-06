package fakedb

import (
	"fmt"
	"testing"
	"time"

	"github.com/func25/rescacher"
)

func TestCacher(t *testing.T) {
	turnCacher := rescacher.NewTurnCacher(0, &DumpCacher{})
	turnCacher.Start()

	for {
		time.Sleep(time.Second)
		mutex.Lock()
		fmt.Println(database)
		mutex.Unlock()
	}
}
