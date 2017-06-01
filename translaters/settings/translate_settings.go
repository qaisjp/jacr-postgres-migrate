package settings

import (
	"encoding/json"
	"time"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"

	"github.com/qaisjp/jacr-postgres-migrate/translaters"

	"log"

	r "gopkg.in/gorethink/gorethink.v3"
)

type translater struct{}

func init() {
	translaters.RegisterTranslater(&translater{})
}

func (t *translater) Name() string {
	return "settings"
}

func (t *translater) Translate(rs *r.Session, db *pg.DB) (err error) {
	log.Println("== Clearing existing settings data...")
	_, err = db.Exec("TRUNCATE TABLE public.settings RESTART IDENTITY RESTRICT;")
	if err != nil {
		return errors.Wrap(err, "could not clear settings table")
	}

	funcs := []struct {
		Name string
		Func func(*r.Session, *pg.DB) error
	}{
		{Name: "motd", Func: handleMOTD},
	}

	for _, fn := range funcs {
		log.Println()
		log.Printf("== Running setting handler '%s'", fn.Name)
		err = fn.Func(rs, db)
		if err != nil {
			return errors.Wrap(err, fn.Name)
		}
		log.Printf("== Handled setting handler '%s'!", fn.Name)
		log.Println()
	}

	return nil
}

func handleMOTD(rs *r.Session, db *pg.DB) (err error) {
	log.Println("=== Reading motd settings from RethinkDB...")
	res, err := r.Table("settings").Get("motd").Run(rs)
	if err != nil {
		return errors.Wrap(err, "could not run query")
	}
	defer res.Close()

	var motd MOTDSetting
	err = res.One(&motd)
	if err != nil {
		return errors.Wrapf(err, "could not read setting")
	}

	res, err = r.Table("settings").Get("motd.dubtrack").Run(rs)
	if err != nil {
		return errors.Wrap(err, "could not run dubtrack query")
	}
	defer res.Close()

	var dub MOTDDubtrackSetting
	err = res.One(&dub)
	if err != nil {
		return errors.Wrapf(err, "could not read dubtrack setting")
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
