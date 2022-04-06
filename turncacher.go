package gacacher

import (
	"time"
)

var CACHER_LOCK = "cacherlock"
var emptyStruct = struct{}{}

type turnCacher struct {
	// plug
	cacher      ICacher
	currentTurn int // current turn

	// config option
	gapTurn      int                 // default gap turn is 50
	gapTime      time.Duration       // default gap time is 1
	locker       Locker              // locker using for distributed lock
	fCurrentTurn func() (int, error) // get the current cached turn of cacher KEY

	// optimized option

	// state
	stop bool
}

func NewTurnCacher[T any](turn int, cacher ICacher, opts ...CacherOption) (*turnCacher, error) {
	turnCacher := &turnCacher{
		currentTurn:  turn,
		cacher:       cacher,
		gapTurn:      50,
		gapTime:      time.Second,
		stop:         false,
		locker:       nil,
		fCurrentTurn: func() (int, error) { return turn, nil },
	}

	turnCacher.applyOpts(opts...)
	return turnCacher, nil
}

func (t *turnCacher) Start(amount int, gapTime time.Duration) {
	go func() {
		failed := 0

		for ; !t.stop; time.Sleep(gapTime) {
			cacher := t.cacher

			// locking
			if t.locker != nil {
				if err := t.locker.Lock(); err != nil {
					continue
				}
			}

			// get current turn
			turn, err := t.fCurrentTurn()
			if err != nil {
				continue
			}

			// generate new turn & save
			for i := 0; i < amount; i++ {
				v, err := cacher.Generate(turn)

				if err != nil {
					failed++
					continue
				}

				failed = 0
				t.cacher.Save(turn, v)
			}

			// unlocking
			if t.locker != nil {
				t.locker.Unlock()
			}
		}
	}()
}

func (t *turnCacher) GetOrGen(turn int) (interface{}, error) {
	item, err := t.cacher.Load(turn)
	if err == nil && item != nil {
		return item, nil
	}

	return t.cacher.Generate(turn)
}

func (t *turnCacher) GetOnly(turn int) (interface{}, error) {
	return t.cacher.Load(turn)
}

func (t *turnCacher) Stop() {
	t.stop = true
}
