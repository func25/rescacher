package rcredis

import (
	"errors"
	"strings"
	"time"

	"github.com/func25/rescacher"
)

type CacherConfig struct {
	Name          string
	Gennerator    rescacher.IGen
	CacherExpired time.Duration
	TurnExpired   time.Duration
}

func (c CacherConfig) Validate() error {
	c.Name = strings.Trim(c.Name, " ")
	if len(c.Name) == 0 {
		return errors.New("cacher name cannot be empty")
	}

	if c.Gennerator == nil {
		return errors.New("generator cannot be nil")
	}

	return nil
}
