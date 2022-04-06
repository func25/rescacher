package rescacher

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
	gapTime      time.Duration       // default gap time is 1 second
	locker       Locker              // locker using for distributed lock
	fCurrentTurn func() (int, error) // get the current cached turn of cacher KEY

	// optimized option

	// state
	stop bool
}

func NewTurnCacher(turn int, cacher ICacher, opts ...CacherOption) *turnCacher {
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
	return turnCacher
}

func (t *turnCacher) Start() {
	go func() {
		failed := 0

		for ; !t.stop; time.Sleep(t.gapTime) {
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

			cachedTurn, err := t.cacher.GetCachedTurn()
			if err != nil || cachedTurn-turn >= t.gapTurn {
				continue
			}

			if cachedTurn < turn {
				cachedTurn = turn + 1
			}

			// generate new turn & save
			for i := 0; i < t.gapTurn; i++ {
				v, err := cacher.Generate(cachedTurn)

				if err != nil {
					failed++
					continue
				}

				failed = 0
				cachedTurn++

				if err := t.cacher.Save(cachedTurn, v); err != nil {
					break
				}
			}

			if err := t.cacher.SetCachedTurn(cachedTurn); err != nil {
				break
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
