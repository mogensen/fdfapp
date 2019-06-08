package actions

import (
	"errors"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/mogensen/fdfapp/models"
)

func bindClasses(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}
	return bindClassesWithConnection(c, tx)
}

func bindClassesWithConnection(c buffalo.Context, tx *pop.Connection) error {
	classes := &models.Classes{}

	// Retrieve all Classes from the DB
	if err := tx.All(classes); err != nil {
		return err
	}

	c.Set("classes", classes)
	return nil
}

func bindParticipants(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	return bindClassesWithConnection(c, tx)
}

func bindParticipantsWithConnection(c buffalo.Context, tx *pop.Connection) error {
	participants := &models.Participants{}

	// Retrieve all participants from the DB
	if err := tx.All(participants); err != nil {
		return err
	}

	c.Set("participants", participants)
	return nil
}
