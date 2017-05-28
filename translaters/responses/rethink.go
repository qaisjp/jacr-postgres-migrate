package responses

type RethinkResponse struct {
	RethinkID string   `gorethink:"id"`
	Name      string   `gorethink:"name"`
	Aliases   []string `gorethink:"aliases"`
	Responses []string `gorethink:"responses"`
}
