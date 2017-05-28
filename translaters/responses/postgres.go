package responses

// Postgres
type ResponseCommand struct {
	ID        int
	Name      string
	RethinkID string `sql:"-"`
}

// Postgres
type ResponseContentCommands struct {
	Command int
	Content int
}

// Postgres
type ResponseContent struct {
	ID        int
	Messages  []string `pg:",array"`
	RethinkID string   `sql:"-"`
}
