package rescacher

import "time"

type CacherOption func(*TurnCacher)

func (c *TurnCacher) applyOpts(opts ...CacherOption) *TurnCacher {
	for i := range opts {
		opts[i](c)
	}

	return c
}

func OptGapTurn(gapTurn int) CacherOption {
	return func(q *TurnCacher) {
		q.gapTurn = gapTurn
	}
}

func OptGapTime(gapTime time.Duration) CacherOption {
	return func(q *TurnCacher) {
		q.checkTime = gapTime
	}
}

func OptLocker(locker Locker) CacherOption {
	return func(q *TurnCacher) {
		q.locker = locker
	}
}

func OptKeyPrefix(prefix string) CacherOption {
	return func(tc *TurnCacher) {
		tc.prefix = prefix
	}
}

func OptResetTurnIfNotFound() CacherOption {
	return func(tc *TurnCacher) {
		tc.resetIfNotFound = true
	}
}
