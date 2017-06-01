package motd

import (
	"time"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"

	"github.com/qaisjp/jacr-postgres-migrate/translaters"

	"log"

	"encoding/json"

	r "gopkg.in/gorethink/gorethink.v3"
)

type translater struct{}

func init() {
	translaters.RegisterTranslater(&translater{})
}

func (t *translater) Name() string {
	return "motd"
}

func (t *translater) Translate(rs *r.Session, db *pg.DB) (err error) {
	log.Println("= Clearing existing motd data...")
	_, err = db.Exec("TRUNCATE TABLE public.notices RESTART IDENTITY RESTRICT; DELETE FROM settings WHERE name = 'motd';")
	if err != nil {
		return errors.Wrap(err, "could not run clear existing data in postgres")
	}

	log.Println("= Reading motd settings from RethinkDB...")
	res, err := r.Table("settings").Get("motd").Run(rs)
	if err != nil {
		return errors.Wrap(err, "could not run rethink query 1")
	}
	defer res.Close()

	var motd MOTDSetting
	err = res.One(&motd)
	if err != nil {
		return errors.Wrapf(err, "could not read motd settings")
	}

	res, err = r.Table("settings").Get("motd.dubtrack").Run(rs)
	if err != nil {
		return errors.Wrap(err, "could not run rethink query 2")
	}
	defer res.Close()

	var dub MOTDDubtrackSetting
	err = res.One(&dub)
	if err != nil {
		return errors.Wrapf(err, "could not read motd dubtrack setting")
	}

	// First insert new motd messages
	messages := make([]Notice, len(motd.Messages))
	for i, msg := range motd.Messages {
		messages[i] = Notice{
			Title:   msg[:1],
			Message: msg,
		}
	}

	err = db.Insert(&messages)
	if err != nil {
		return errors.Wrapf(err, "could not insert notices")
	}

	json, err := json.Marshal(struct {
		Enabled          bool
		Interval         int
		NextMessage      int
		LastAnnounceTime time.Time
	}{
		Enabled:          motd.Enabled,
		Interval:         motd.Interval,
		NextMessage:      dub.NextMessage,
		LastAnnounceTime: dub.LastAnnounceTime,
	})

	if err != nil {
		return errors.Wrapf(err, "could not build json")
	}

	_, err = db.Exec(`INSERT INTO settings(name, value) VALUES('motd', ?)`, string(json))
	if err != nil {
		return errors.Wrapf(err, "could not insert settings")
	}

	return nil
}
