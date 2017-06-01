package songs

import (
	"time"
)

type Song struct {
	Fkid string
	Type string

	LastPlay   time.Time
	Name       string
	SkipReason string

	RecentPlays int
	TotalPlays  int

	Retagged     bool
	Autoretagged bool

	RethinkID string `gorethink:"id"`
}
