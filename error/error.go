package error

import "errors"

var (
	ErrItemAlreadyExist = errors.New("item already exist")
)
