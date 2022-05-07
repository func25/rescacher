package rcredistest

import (
	"context"
	"time"
)

type Generator struct {
	CurrentTurn int
}

func (*Generator) Generate(ctx context.Context, turn int) (interface{}, error) {
	return time.Now().Unix(), nil
}

func (g *Generator) GetCurrentTurn(ctx context.Context) (int, error) {
	return g.CurrentTurn, nil
}
