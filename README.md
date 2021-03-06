# rescacher

The idea is when you have an API and the value of it changes over time, rescacher will cache the next value, next next value,... of it and just throw it out on time with no delay.

Contents:
* [Usecase](#usecase)
* [Installation](#installation)
* [How to use](#how-to-use)
  * [Sample](#sample)
  * [IGen](#igen)
  * [ICacher](#icacher)
  * [Load value from cacher](#load-value-from-cacher)
* [Functions](#functions)
* [Status](#status-pre-release)

## Use cases

- The next values being calculated by server can be had errors or failed to calculate, rescacher will retry the calculation after an configured interval and cache many predicted results so server can use them gradually.
- Your users dont have to wait or the event will happen at the right time (no delay).

## Installation

`go get github.com/func25/rescacher`

## How to use 

### Sample
You can get an sample from `cacher_test.go`, basically it will cache next 50 results into database and retrieve each result in 1 second over time.

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
