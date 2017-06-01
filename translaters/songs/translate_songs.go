package songs

import (
	"github.com/go-pg/pg"
	"github.com/pkg/errors"

	"github.com/qaisjp/jacr-postgres-migrate/translaters"

	"log"

	"gopkg.in/cheggaaa/pb.v1"
	r "gopkg.in/gorethink/gorethink.v3"
)

type responsesTranslater struct{}

func init() {
	translaters.RegisterTranslater(&responsesTranslater{})
}

func (t *responsesTranslater) Name() string {
	return "songs"
}

func (t *responsesTranslater) Translate(rs *r.Session, db *pg.DB) (err error) {
	// log.Println("= Clearing existing song data...")
	// _, err = db.Exec("TRUNCATE TABLE public.songs RESTART IDENTITY RESTRICT;")
	// if err != nil {
	// 	return errors.Wrap(err, "could not run clear existing data in postgres")
	// }

	log.Println("= Reading songs from RethinkDB...")

	res, err := r.Table("songs").Count().Run(rs)
	if err != nil {
		return errors.Wrapf(err, "could not query songs count")
	}

	var count int
	if err := res.One(&count); err != nil {
		return errors.Wrapf(err, "could not get songs count")
	}

	res, err = r.Table("songs").Run(rs)
	if err != nil {
		return errors.Wrap(err, "could not query songs")
	}
	defer res.Close()

	bar := pb.StartNew(count)
	defer bar.Finish()

	var song Song
	for res.Next(&song) {
		err = db.Insert(&song)
		if err != nil {
			return errors.Wrapf(err, "could not insert song %s", song.RethinkID)
		}
		bar.Increment()
	}

	if res.Err() != nil {
		return errors.Wrapf(err, "could not read cursor")
	}

	return nil
}
