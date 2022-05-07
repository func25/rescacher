package rescacher

import "context"

type ICacher interface {
	Load(ctx context.Context, turn int, pop bool) (interface{}, error) // load result of turn X (own db)
	Save(ctx context.Context, turn int, value interface{}) error       // save result of turn Y (own db)
	GetCacher(ctx context.Context) (int, error)                        // get current cached turn (own db)
	SetCacher(ctx context.Context, turn int) error                     // save current cached turn (own db)
}

type IGen interface {
	Generate(ctx context.Context, turn int) (interface{}, error) // generate result of turn X (custom)
	GetCurrentTurn(ctx context.Context) (int, error)             // get current turn (custom)
}

type Locker interface {
	Lock() error
	Unlock() error
}
