package exceptions

import "errors"

var (
	ErrCommitedAlready = errors.New("commited")
	ErrRollbackAlready = errors.New("rollback")
	ErrModelNotFound   = errors.New("model not found")
	ErrUnauthorized    = errors.New("unauthorized")
	ErrBadRequest      = errors.New("bad request")
)
