package translaters

import (
	"github.com/go-pg/pg"
	r "gopkg.in/gorethink/gorethink.v3"
)

var translaters = make([]Translater, 0)

// Translater is an interface for every translater that
// will translate a particular object in the database
type Translater interface {
	Name() string // Name of the translater

	// We don't have these because the translater could
	// pull/push information from/to multiple tables
	// SourceTable() string      // Rethink
	// DestinationTable() string // Postgres

	Translate(*r.Session, *pg.DB) error
}

// RegisterTranslater is called by each translater
func RegisterTranslater(t Translater) {
	translaters = append(translaters, t)
}

func List() []Translater {
	return translaters
}
