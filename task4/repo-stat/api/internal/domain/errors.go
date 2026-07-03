package domain

import "errors"

var (
	ErrNotFound = errors.New("repository not found")

	ErrMovedPermanently = errors.New("repository moved permanently")

	ErrForbidden = errors.New("access forbidden")

	ErrUnauthorized = errors.New("unauthorized")

	ErrRateLimit = errors.New("rate limit exceeded")

	ErrInvalidInput = errors.New("invalid input")

	ErrTimeout = errors.New("request timeout")

	ErrInternal = errors.New("internal error")

	ErrSubscriptionAlreadyExists = errors.New("subscription is already exists")
)
