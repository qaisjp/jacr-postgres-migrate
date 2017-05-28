package responses

import (
	"github.com/go-pg/pg"
	"github.com/pkg/errors"

	"github.com/qaisjp/jacr-postgres-migrate/translaters"

	"fmt"

	"log"

	r "gopkg.in/gorethink/gorethink.v3"
)

type responsesTranslater struct{}

func init() {
	translaters.RegisterTranslater(&responsesTranslater{})
}

func (t *responsesTranslater) Name() string {
	return "responses"
}

func (t *responsesTranslater) Translate(rs *r.Session, db *pg.DB) (err error) {
	// First loading all commands and their aliases into memory...
	//
	res, err := r.Table("responses").Run(rs)

	if err != nil {
		return errors.Wrap(err, "could not run rethink query")
	}

	defer res.Close()

	var responses []RethinkResponse
	err = res.All(&responses)
	if err != nil {
		return errors.Wrap(err, "could not retrieve rethink query")
	}

	commands := buildCommands(responses)
	if !reportDuplicateCommands(commands) {
		return errors.New("contains duplicate commands")
	}

	contents := buildContents(responses)
	// contentCommands := buildContentCommands(contents, commands)

	log.Println("Commands and contents generated... trying to insert.")

	err = db.Insert(&commands)
	if err != nil {
		return errors.Wrap(err, "could not insert commands")
	}

	err = db.Insert(&contents)
	if err != nil {
		return errors.Wrap(err, "could not insert contents")
	}

	fmt.Println(commands)
	fmt.Println(contents)
	return nil
}

func buildCommands(responses []RethinkResponse) []ResponseCommand {
	// First count how many commands we want to perform this on
	count := 0
	for _, r := range responses {
		count += 1 + len(r.Aliases)
	}

	// Allocate slice
	commands := make([]ResponseCommand, count)

	// Create a command for ecah response, and their aliases
	i := 0
	for _, r := range responses {
		commands[i] = ResponseCommand{Name: r.Name}
		i++

		for _, n := range r.Aliases {
			commands[i] = ResponseCommand{Name: n}
			i++
		}
	}

	return commands
}

func reportDuplicateCommands(cmds []ResponseCommand) (unique bool) {
	// Find duplicate commands, and report these
	exists := make(map[string]struct{})

	unique = true

	for _, c := range cmds {
		if _, ok := exists[c.Name]; ok {
			fmt.Printf("Warning: command %s appeared again\n", c.Name)
			unique = false
		}

		exists[c.Name] = struct{}{}
	}

	return unique
}

func buildContents(responses []RethinkResponse) []ResponseContent {
	contents := make([]ResponseContent, len(responses))

	for i, r := range responses {
		contents[i] = ResponseContent{
			Messages: r.Responses,
		}
	}

	return contents
}

func buildContentCommands(contents []ResponseContent, commands []ResponseCommand) []ResponseContentCommands {
	contentCommands := make([]ResponseContentCommands, len(commands))

	var contentID int
	contentCount := 0

	for _, cmd := range commands {
		contentFound := false

		// Try find the content
		for _, c := range contents {
			if c.ID == cmd.ID {
				contentFound = true
				contentID = c.ID
			}
		}

		if !contentFound {
			continue // Skip over if we could not find the content
		}

		// Now we've found the content
		contentCommands[contentCount] = ResponseContentCommands{
			Command: cmd.ID,
			Content: contentID,
		}

		contentCount++

	}

	return contentCommands[:contentCount]
}
