package domain

import "errors"

var (
	ErrNotFound = errors.New("repository not found")

	ErrAccepted = errors.New("request accepted")

	ErrTimeout = errors.New("request timeout")

	ErrInternal = errors.New("internal error")
)
