package users

import (
	"time"
)

type DubtrackUser struct {
	// ID int
	Karma    int    `gorethink:"karma"`
	DubID    string `gorethink:"uid"`
	Username string `gorethink:"username"`

	// Rethink Only
	Seen struct {
		Message string    `gorethink:"message"`
		Time    time.Time `gorethink:"time"`
		Type    string    `gorethink:"type"`
	} `sql:"-" gorethink:"seen"`

	// Postgres Only
	SeenMessage string
	SeenTime    time.Time
	SeenType    string

	RethinkID string `gorethink:"id"`
}
