package motd

type MOTDSetting struct {
	Enabled  bool     `gorethink:"enabled"`
	Interval int      `gorethink:"interval"`
	Messages []string `gorethink:"messages"`
}

type Notice struct {
	// ID int
	Title   string
	Message string
}
