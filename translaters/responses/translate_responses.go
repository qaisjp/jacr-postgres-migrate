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
	log.Println("Reading responses from RethinkDB...")

	var responses []RethinkResponse
	{
		res, err := r.Table("responses").Run(rs)

		if err != nil {
			return errors.Wrap(err, "could not run rethink query")
		}
		defer res.Close()

		err = res.All(&responses)
		if err != nil {
			return errors.Wrap(err, "could not retrieve rethink query")
		}
	}

	{

		log.Println("= Handling contents...")
		contents := buildContents(responses)

		err = db.Insert(&contents)
		if err != nil {
			return errors.Wrap(err, "could not insert contents")
		}

		log.Println("= Contents handled!")

		log.Println("= Handling commands...")
		commands := buildCommands(responses, contents)
		if !reportDuplicateCommands(commands) {
			return errors.New("contains duplicate commands")
		}

		err = db.Insert(&commands)
		if err != nil {
			return errors.Wrap(err, "could not insert commands")
		}

		log.Println("= Commands handled!")
	}
	return nil
}

func buildCommands(responses []RethinkResponse, contents []ResponseContent) []ResponseCommand {
	// Get all commands, irrespective of their group
	var commands []ResponseCommand

	{
		// First count how many commands we want to perform this on
		count := 0
		for _, r := range responses {
			count += 1 + len(r.Aliases)
		}

		// Allocate slice
		commands = make([]ResponseCommand, count)

		// Create a command for each response, and their aliases
		i := 0
		for _, r := range responses {
			commands[i] = ResponseCommand{Name: r.Name, RethinkID: r.RethinkID}
			i++

			for _, n := range r.Aliases {
				commands[i] = ResponseCommand{Name: n, RethinkID: r.RethinkID}
				i++
			}
		}
	}

	// Now loop over the commands we have
	// and update their group
	filteredCommands := make([]ResponseCommand, len(commands))
	contentCount := 0

	for _, cmd := range commands {
		contentID := -1

		// Try find the content
		for _, c := range contents {
			if c.RethinkID == cmd.RethinkID {
				contentID = c.ID
			}
		}

		if contentID == -1 {
			log.Printf("warning: content missing for cmd '%s', skipping. (rethink: %s)", cmd.Name, cmd.RethinkID)
			continue // Skip over if we could not find the content
		}

		// Now we've found the content
		cmd.Content = contentID
		filteredCommands[contentCount] = cmd

		contentCount++

	}

	return filteredCommands[:contentCount]
}

func reportDuplicateCommands(cmds []ResponseCommand) (unique bool) {
	// Find duplicate commands, and report these
	exists := make(map[ResponseCommand]struct{})

	unique = true

	for _, c := range cmds {
		flatCmd := c
		flatCmd.RethinkID = ""

		if _, ok := exists[flatCmd]; ok {
			fmt.Printf("Warning: command pair ('%s','%s') appeared again\n", c.Name, c.Content)
			unique = false
		}

		exists[flatCmd] = struct{}{}
	}

	return unique
}

func buildContents(responses []RethinkResponse) []ResponseContent {
	contents := make([]ResponseContent, len(responses))

	for i, r := range responses {
		contents[i] = ResponseContent{
			Messages:  r.Responses,
			RethinkID: r.RethinkID,
		}
	}

	return contents
}
