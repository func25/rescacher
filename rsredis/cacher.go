package rsredis

import (
	"context"

	"github.com/func25/rescacher"
	"github.com/go-redis/redis/v8"
)

type cachedHandler struct {
	client *redis.Client
	cacher *rescacher.TurnCacher
	config CacherConfig
}

var _ rescacher.ICacher = (*cachedHandler)(nil)

func NewCacher(client *redis.Client, config CacherConfig, opts ...rescacher.CacherOption) (*rescacher.TurnCacher, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	handler := &cachedHandler{
		client: client,
		config: config,
	}
	handler.cacher = rescacher.NewTurnCacher(handler.config.Name, handler, handler.config.Gennerator, opts...)

	return handler.cacher, nil
}

func (c *cachedHandler) Load(ctx context.Context, turn int, pop bool) (interface{}, error) {
	key := c.cacher.GetTurnKey(turn)
	res, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}

	if pop {
		c.client.Del(ctx, key)
	}

	return res, err
}

func (c *cachedHandler) Save(ctx context.Context, turn int, value interface{}) error {
	return c.client.Set(ctx, c.cacher.GetTurnKey(turn), value, c.config.TurnExpired).Err()
}

func (c *cachedHandler) SetCacher(ctx context.Context, cachedTurn int) error {
	return c.client.Set(ctx, c.cacher.GetCacherKey(), cachedTurn, c.config.CacherExpired).Err()
}

func (c *cachedHandler) GetCacher(ctx context.Context) (int, error) {
	res, err := c.client.Get(ctx, c.cacher.GetCacherKey()).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return res, err
}
