package gacacher

type ICacher interface {
	Generate(int) (interface{}, error) // generate result of turn X
	Load(int) (interface{}, error)     // load result of turn X
	Save(int, value interface{})       // save result of turn Y
}
