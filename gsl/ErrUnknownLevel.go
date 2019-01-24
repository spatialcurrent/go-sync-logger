package gsl

import (
	"fmt"
)

type ErrUnknownLevel struct {
	Level string
}

func (e *ErrUnknownLevel) Error() string {
	return fmt.Sprintf("unknown level %s", e.Level)
}
