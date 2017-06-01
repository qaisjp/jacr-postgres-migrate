package motd

import "time"

type MOTDSetting struct {
	Enabled  bool     `gorethink:"enabled"`
	Interval int      `gorethink:"interval"`
	Messages []string `gorethink:"messages"`
}

type MOTDDubtrackSetting struct {
	LastAnnounceTime time.Time `gorethink:"lastAnnounceTime"`
	NextMessage      int       `gorethink:"nextMessage"`
}

type Notice struct {
	// ID int
	Title   string
	Message string
}
