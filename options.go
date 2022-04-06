package gacacher

import "time"

type CacherOption func(*turnCacher)

func (c *turnCacher) applyOpts(opts ...CacherOption) *turnCacher {
	for i := range opts {
		opts[i](c)
	}

	return c
}

func OptGapTurn(gapTurn int) CacherOption {
	return func(q *turnCacher) {
		q.gapTurn = gapTurn
	}
}

func OptGapTime(gapTime time.Duration) CacherOption {
	return func(q *turnCacher) {
		q.gapTime = gapTime
	}
}

func OptLocker(locker Locker) CacherOption {
	return func(q *turnCacher) {
		q.locker = locker
	}
}

func OptRefreshCurrentTurn(f func() (int, error)) CacherOption {
	return func(tc *turnCacher) {
		tc.fCurrentTurn = f
	}
}
