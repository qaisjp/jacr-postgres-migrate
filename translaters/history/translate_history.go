package history

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
	return "history"
}

func (t *responsesTranslater) Translate(rs *r.Session, db *pg.DB) (err error) {
	log.Println("= Reading history from RethinkDB...")

	res, err := r.Table("history").Count().Run(rs)
	if err != nil {
		return errors.Wrapf(err, "could not query history count")
	}

	var count int
	if err := res.One(&count); err != nil {
		return errors.Wrapf(err, "could not get history count")
	}

	res, err = r.Table("history").Run(rs)
	if err != nil {
		return errors.Wrap(err, "could not query history")
	}
	defer res.Close()

	bar := pb.StartNew(count)
	defer bar.Finish()

	var history History
	for res.Next(&history) {
		dubID := history.DubID
		if dubID == "" {
			dubID = history.RethinkID
		}

		_, err = db.Exec(
			`
			WITH
				u as (SELECT id FROM dubtrack_users WHERE rethink_id = ?),
				s as (SELECT id FROM songs WHERE rethink_id = ?)
			INSERT INTO
			history (dub_id, score_down, score_grab, score_up, song, "user", time)
			VALUES (?, ?, ?, ?, (SELECT id FROM s), (SELECT id from u), ?)
			`,
			history.User,
			history.Song,
			dubID,
			history.Score.Down,
			history.Score.Grab,
			history.Score.Up,
			history.Time,
		)

		if err != nil {
			return errors.Wrapf(err, "could not insert history %s", history.RethinkID)
		}
		bar.Increment()
	}

	if res.Err() != nil {
		return errors.Wrapf(err, "could not read cursor")
	}

	return nil
}
