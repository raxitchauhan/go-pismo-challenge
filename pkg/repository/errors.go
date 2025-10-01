package repository

import "errors"

var (
	ErrDuplicate = errors.New("duplicate request received")
	ErrNoRows    = errors.New("no rows found")
)
