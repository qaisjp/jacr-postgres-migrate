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

	history := make(chan History, count)
	errors := make(chan error, count)

	for w := 1; w <= 100; w++ {
		go worker(db, history, errors)
	}

	res.Listen(history)

	progress := 0
	for err := range errors {
		progress++
		if err == nil {
			bar.Increment()
		} else {
			log.Println(err)
		}

		if progress == count {
			close(errors)
		}
	}

	return nil
}

func worker(db *pg.DB, hchan chan History, errs chan error) {
	for history := range hchan {
		dubID := history.DubID
		if dubID == "" {
			dubID = history.RethinkID
		}

		_, err := db.Exec(
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

		// wrap returns nil if err is nil, so this is ok
		errs <- errors.Wrapf(err, "could not insert history %s", history.RethinkID)
	}
}
