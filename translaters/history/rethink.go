package history

import (
	"time"
)

type History struct {
	DubID string `gorethink:"platformID"`

	Score struct {
		Down int
		Up   int
		Grab int
	}

	Song string
	Time time.Time
	User string

	RethinkID string `gorethink:"id"`
}
