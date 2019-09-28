package actions

import (
	"github.com/gobuffalo/buffalo"
	"github.com/mogensen/fdfapp/models"
)

func bindClasses(c buffalo.Context) error {
	classes := &models.Classes{}

	// Retrieve all Classes from the DB
	if err := scope(c).All(classes); err != nil {
		return err
	}

	c.Set("classes", classes)
	return nil
}

func bindParticipants(c buffalo.Context) error {
	participants := &models.Participants{}

	// Retrieve all participants from the DB
	if err := scope(c).All(participants); err != nil {
		return err
	}

	c.Set("participants", participants)
	return nil
}
