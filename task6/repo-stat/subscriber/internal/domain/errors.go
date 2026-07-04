package domain

import "errors"

var (
	ErrSubscriptionAlreadyExists = errors.New("subscription already exists")
	ErrRepositoryNotFound        = errors.New("repository not found on github")
)
