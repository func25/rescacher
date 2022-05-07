package rcredistest

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/func25/rescacher"
	"github.com/func25/rescacher/rcredis"
	"github.com/go-redis/redis/v8"
)

var client *redis.Client
var cacher *rescacher.TurnCacher
var gen *Generator = &Generator{}

func init() {
	var err error

	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "",
		Password: "",
	})

	client.FlushAll(context.Background())

	cacher, err = rcredis.NewCacher(client, rcredis.CacherConfig{
		Name:       "example",
		Gennerator: gen,
	}, rescacher.OptResetTurnIfNotFound())

	if err != nil {
		log.Fatal(err)
	}
}

func TestCacher(t *testing.T) {
	ctx := context.Background()
	cacher.Start()

	time.Sleep(2 * time.Second)
	for ; ; gen.CurrentTurn++ {
		res, err := cacher.PopOrGen(ctx, gen.CurrentTurn)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(gen.CurrentTurn, res)
		}
		time.Sleep(time.Second)
	}
}
