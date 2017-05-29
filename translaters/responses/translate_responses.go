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

		log.Println("= Handling groups...")
		groups := buildGroups(responses)

		err = db.Insert(&groups)
		if err != nil {
			return errors.Wrap(err, "could not insert groups")
		}

		log.Println("= Groups handled!")

		log.Println("= Handling commands...")
		commands := buildCommands(responses, groups)
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

func buildCommands(responses []RethinkResponse, groups []ResponseGroup) []ResponseCommand {
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
	groupCount := 0

	for _, cmd := range commands {
		groupID := -1

		// Try find the group
		for _, c := range groups {
			if c.RethinkID == cmd.RethinkID {
				groupID = c.ID
			}
		}

		if groupID == -1 {
			log.Printf("warning: group missing for cmd '%s', skipping. (rethink: %s)", cmd.Name, cmd.RethinkID)
			continue // Skip over if we could not find the group
		}

		// Now we've found the group
		cmd.Group = groupID
		filteredCommands[groupCount] = cmd

		groupCount++

	}

	return filteredCommands[:groupCount]
}

func reportDuplicateCommands(cmds []ResponseCommand) (unique bool) {
	// Find duplicate commands, and report these
	exists := make(map[ResponseCommand]struct{})

	unique = true

	for _, c := range cmds {
		flatCmd := c
		flatCmd.RethinkID = ""

		if _, ok := exists[flatCmd]; ok {
			fmt.Printf("Warning: command pair ('%s','%s') appeared again\n", c.Name, c.Group)
			unique = false
		}

		exists[flatCmd] = struct{}{}
	}

	return unique
}

func buildGroups(responses []RethinkResponse) []ResponseGroup {
	groups := make([]ResponseGroup, len(responses))

	for i, r := range responses {
		groups[i] = ResponseGroup{
			Messages:  r.Responses,
			RethinkID: r.RethinkID,
		}
	}

	return groups
}
