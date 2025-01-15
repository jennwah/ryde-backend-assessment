package postgresql

import "errors"

var (
	ErrNotFound                      = errors.New("not found from database")
	ErrUniqueConstraintViolationCode = "23505"
	ErrUserAlreadyExist              = errors.New("user already exists")
	ErrFriendshipAlreadyExist        = errors.New("friendship already exists")
)
