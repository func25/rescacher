package rescacher

type ICacher interface {
	Generate(turn int) (interface{}, error) // generate result of turn X
	Load(turn int) (interface{}, error)     // load result of turn X
	Save(turn int, value interface{}) error // save result of turn Y
	GetCachedTurn() (int, error)            // get current cached turn
	SetCachedTurn(turn int) error           // save current cached turn
}
