package gacacher

import "sync"

type Locker interface {
	Lock() error
	Unlock() error
}

type AbcLocker struct {
	mutex sync.Mutex
}

func (a AbcLocker) Lock() error {
	a.mutex.Lock()
	return nil
}

func (a AbcLocker) Unlock() error {
	a.mutex.Unlock()
	return nil
}
