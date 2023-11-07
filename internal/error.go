package internal

import (
	"errors"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidInputFile = errors.New("invalid input file")
)
