package domain

import (
	"time"
)

type Subscription struct {
	ID        int
	Owner     string
	Repo      string
	CreatedAt time.Time
}
