package users

import (
	"github.com/go-pg/pg"
	"github.com/pkg/errors"

	"github.com/qaisjp/jacr-postgres-migrate/translaters"

	"log"

	r "gopkg.in/gorethink/gorethink.v3"
)

type responsesTranslater struct{}

func init() {
	translaters.RegisterTranslater(&responsesTranslater{})
}

func (t *responsesTranslater) Name() string {
	return "users"
}

func (t *responsesTranslater) Translate(rs *r.Session, db *pg.DB) (err error) {
	log.Println("= Reading users from RethinkDB...")

	res, err := r.Table("users").Run(rs)

	if err != nil {
		return errors.Wrap(err, "could not run rethink query")
	}
	defer res.Close()

	var user DubtrackUser
	for res.Next(&user) {
		user.SeenMessage = user.Seen.Message
		user.SeenTime = user.Seen.Time
		user.SeenType = user.Seen.Type

		err = db.Insert(&user)
		if err != nil {
			return errors.Wrapf(err, "could not insert user %s", user.RethinkID)
		}
	}
	if res.Err() != nil {
		// error
	}

	return nil
}
