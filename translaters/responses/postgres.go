package responses

// Postgres
type ResponseCommand struct {
	ID        int
	Name      string
	Group     int
	RethinkID string `sql:"-"`
}

// Postgres
type ResponseGroup struct {
	ID        int
	Messages  []string `pg:",array"`
	RethinkID string   `sql:"-"`
}
