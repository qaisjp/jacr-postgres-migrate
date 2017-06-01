package motd

import (
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
	return "motd"
}

func (t *translater) Translate(rs *r.Session, db *pg.DB) (err error) {
	log.Println("= Clearing existing motd data...")
	_, err = db.Exec("TRUNCATE TABLE public.notices RESTART IDENTITY RESTRICT;")
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

	return nil
}
