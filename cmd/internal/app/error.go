package app

import (
	"errors"
)

// Errors.
var (
	ErrNotFound      = errors.New("not found")
	ErrTelegramExist = errors.New("telegram exists")
	ErrCardExists    = errors.New("card exists")
)
