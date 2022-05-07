# rescacher

The idea is when you have an API and the value of it changes over time (you **should** know what result of api call in specific future time), and you want to cache the results that it will "throw" next.

Contents:
* [Usecase](#usecase)
* [Installation](#installation)
* [How to use](#how-to-use)
  * [IGen](#igen)
  * [ICacher](#icacher)
  * [Load value from cacher](#load-value-from-cacher)
* [Functions](#functions)
* [Status](#status-pre-release)

## Usecase

- The next value being calculated by server may be had errors or failed to calculate, rescacher will retry after an interval.
- Your users dont have to wait.

## Installation

`go get github.com/func25/rescacher`

## How to use 

### IGen
First, we want to know how to generate the result and the current "result id" of the api right now by using an object implement IGen interface:

```go
type IGen interface {
	Generate(ctx context.Context, turn int) (interface{}, error) // generate result of turn X
	GetCurrentTurn(ctx context.Context) (int, error)             // get current turn
}
```

### ICacher

Next, you need a ICacher that can store/ load the "predicted value" and cacher itself into the database, we already made an cacher based on redis:

```go
cacher, err = rcredis.NewCacher(client, rcredis.CacherConfig{
		Name:       "example",
		Gennerator: gen, // object implement IGen
})
```

### Load value from cacher

```go
// return the result of turn "id"
res, err := cacher.PopOrGen(ctx, id) 
fmt.Println(res)
```

## Functions
- GetOrGen: get the result of turn X (if it is not exist then generate), but the result still remain in database

- PopOrGen: like GetOrGen but we remove the result out of database when successfully retrieve


## Status: pre-release
This lib is under developing, please notice when using it
