package rescacher

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// functask: create algorithm to cache slowly from 1 to gapTurn, so it will not stress cpu

var emptyStruct = struct{}{}

// functask: slowly warm up first
type TurnCacher struct {
	// metadata
	name string

	// plug
	cacher ICacher
	gentor IGen

	// config option
	prefix          string        // prefix of key
	gapTurn         int           // default gap turn is 50, from current -> cached
	checkTime       time.Duration // default check time is 1 second
	locker          Locker        // locker using for distributed lock
	resetIfNotFound bool

	// optimized option

	// state
	stop bool
	mtx  sync.Mutex
}

func NewTurnCacher(name string, cacher ICacher, gentor IGen, opts ...CacherOption) *TurnCacher {
	turnCacher := &TurnCacher{
		name:            name,
		cacher:          cacher,
		gentor:          gentor,
		gapTurn:         50,
		checkTime:       time.Second,
		stop:            false,
		locker:          nil,
		prefix:          "resc_",
		resetIfNotFound: false,
	}

	turnCacher.applyOpts(opts...)
	return turnCacher
}

// -- progress
func (t *TurnCacher) Start() {
	go func() {
		defer func() {
			any := recover()
			if any != nil {
				fmt.Println("[rescacher] panic", any)
			}
		}()

		failed := 0

		for ; !t.stop; t.sleepAndUnlock() {
			ctx := context.Background()

			// lock
			if t.locker != nil {
				if err := t.locker.Lock(); err != nil {
					continue
				}
			}
			t.mtx.Lock()

			// get current turn
			turn, err := t.gentor.GetCurrentTurn(ctx)
			if err != nil {
				continue
			}

			// get cached turn
			cachedTurn, err := t.cacher.GetCacher(ctx)
			if err != nil || cachedTurn-turn >= t.gapTurn {
				continue
			}

			if cachedTurn < turn {
				cachedTurn = turn
			}

			// generate new turn & save
			for i := cachedTurn + 1; i <= turn+t.gapTurn; i++ {
				v, err := t.gentor.Generate(ctx, i)

				if err != nil {
					failed++
					continue
				}

				failed = 0
				cachedTurn++

				if err := t.cacher.Save(ctx, i, v); err != nil {
					break
				}
			}

			if err := t.cacher.SetCacher(ctx, cachedTurn); err != nil {
				break
			}
		}
	}()
}

func (t *TurnCacher) sleepAndUnlock() {
	// unlock
	if t.locker != nil {
		t.locker.Unlock()
	}
	t.mtx.Unlock()

	time.Sleep(t.checkTime)
}

func (t *TurnCacher) Stop() {
	t.stop = true
}

// -- get
func (t *TurnCacher) GetOrGen(ctx context.Context, turn int) (interface{}, error) {
	return t.getOrGen(ctx, turn, false)
}

func (t *TurnCacher) PopOrGen(ctx context.Context, turn int) (interface{}, error) {
	return t.getOrGen(ctx, turn, true)
}

func (t *TurnCacher) getOrGen(ctx context.Context, turn int, pop bool) (interface{}, error) {
	item, err := t.cacher.Load(ctx, turn, pop)
	if err == nil {
		if item != nil {
			return item, nil
		} else {
			// not found case
			if t.resetIfNotFound {
				t.mtx.Lock()
				defer t.mtx.Unlock()
				err := t.cacher.SetCacher(ctx, turn)
				if err != nil {
					// todo: log
				}
			}
		}
	}

	return t.gentor.Generate(ctx, turn)
}

func (t *TurnCacher) GetOnly(ctx context.Context, turn int) (interface{}, error) {
	return t.cacher.Load(ctx, turn, false)
}

// -- key
func (t *TurnCacher) GetTurnKey(turn int) string {
	return fmt.Sprintf("%v%v_%v", t.prefix, t.name, turn)
}

func (t *TurnCacher) GetCacherKey() string {
	return fmt.Sprintf("%v%v", t.prefix, t.name)
}
