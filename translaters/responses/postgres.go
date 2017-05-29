package responses

// Postgres
type ResponseCommand struct {
	ID        int
	Name      string
	Content   int
	RethinkID string `sql:"-"`
}

// Postgres
type ResponseContent struct {
	ID        int
	Messages  []string `pg:",array"`
	RethinkID string   `sql:"-"`
}
