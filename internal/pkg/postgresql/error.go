package postgresql

import "errors"

var (
	ErrNotFound = errors.New("not found from database")
)
