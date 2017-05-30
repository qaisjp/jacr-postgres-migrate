package main

import (
	"flag"
	"log"

	"github.com/qaisjp/jacr-postgres-migrate/translaters"

	// Translaters to make use of
	_ "github.com/qaisjp/jacr-postgres-migrate/translaters/responses"
	_ "github.com/qaisjp/jacr-postgres-migrate/translaters/songs"
	_ "github.com/qaisjp/jacr-postgres-migrate/translaters/users"

	"github.com/go-pg/pg"
	"github.com/pkg/errors"
	r "gopkg.in/gorethink/gorethink.v3"
)

var pgPassword = flag.String("pg_pass", "", "Postgres Password")
var activeTranslaters = make([]translaters.Translater, 0)

func init() {

	translaters := translaters.List()
	enabled := make([]*bool, len(translaters))
	for i, t := range translaters {
		enabled[i] = flag.Bool(
			"t-"+t.Name(),
			false,
			"Should use translater "+t.Name(),
		)
	}

	flag.Parse()

	for i, b := range enabled {
		if *b {
			activeTranslaters = append(activeTranslaters, translaters[i])
		}
	}
}

func main() {
	rs, err := r.Connect(r.ConnectOpts{
		Address:  "localhost:28015",
		Database: "jacr_dev",
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Print("Connected to RethinkDB.")

	db := pg.Connect(&pg.Options{
		Addr:     "localhost:5432",
		User:     "postgres",
		Database: "jacr_dev",
		Password: *pgPassword,
	})

	_, err = db.Exec("SELECT 1")
	if err != nil {
		log.Print("Postgres connection error!\n")
		panic(err)
	}
	log.Print("Connected to PostgreSQL.\n")

	for _, t := range activeTranslaters {
		log.Printf("== Running translater '%s' ==", t.Name())
		err := t.Translate(rs, db)

		if err != nil {
			log.Panicln(errors.Wrapf(err, "translation failed in '%s'", t.Name()))
		}
		log.Printf("== DONE! Translater '%s' completed! ==", t.Name())
	}

	// Close RethinkDB
	err = rs.Close()
	if err != nil {
		log.Print("Could not close RethinkDB connection.")
		panic(err)
	}

	// Close PostgreSQL
	err = db.Close()
	if err != nil {
		log.Print("Could not close Postgres connection.")
		panic(err)
	}
}
