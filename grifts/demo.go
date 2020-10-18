package grifts

import (
	"math/rand"
	"time"

	"github.com/markbates/grift/grift"
	"github.com/mogensen/fdfapp/models"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

var _ = grift.Namespace("demo", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {

		ph, err := bcrypt.GenerateFromPassword([]byte("demo"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user := models.User{
			Username:     "demo",
			PasswordHash: string(ph),
		}
		if err := models.DB.Create(&user); err != nil {
			return errors.WithStack(err)
		}

		classes := models.Classes{
			{UserID: user.ID, Name: "Pusling"},
			{UserID: user.ID, Name: "Tumling"},
			{UserID: user.ID, Name: "Pilt"},
			{UserID: user.ID, Name: "Væbner"},
			{UserID: user.ID, Name: "Seniorvæbner"},
			{UserID: user.ID, Name: "Senior"},
		}

		for i, m := range classes {
			if err := models.DB.Create(&m); err != nil {
				return errors.WithStack(err)
			}
			// Store new ID
			classes[i] = m
		}

		leder := models.Class{UserID: user.ID, Name: "Leder"}
		if err := models.DB.Create(&leder); err != nil {
			return errors.WithStack(err)
		}

		year := (time.Now().Year() - 6)
		grownUpYear := (time.Now().Year() - 22)

		for _, class := range classes {
			for i := 0; i < 10; i++ {
				createParticipant(user, year, class)
			}

			// Create two grown ups in each Class
			createParticipant(user, grownUpYear, class, leder)
			createParticipant(user, grownUpYear-10, class, leder)

			year -= 2
		}

		durations := []float64{0.5, 1, 1.25, 1.5, 2, 3, 5}
		for _, class := range classes {
			for i := 0; i < 50; i++ {

				c := &models.Class{}

				// Retrieve all participants from the DB
				if err := models.DB.Eager("Participants.Image").Find(c, class.ID); err != nil {
					return err
				}

				parts := models.Participants{}
				for i := 0; i < rand.Intn(len(c.Participants)); i++ {
					parts = append(parts, c.Participants[rand.Intn(len(c.Participants))])
				}

				activityTime := randomTimestamp(4, 0)
				act := &models.Activity{
					Title:        randActivityTitle(),
					Duration:     durations[rand.Intn(len(durations))],
					UserID:       user.ID,
					ClassID:      class.ID,
					Date:         activityTime,
					Participants: parts,
				}
				models.DB.Create(act)
			}
		}

		return nil
	})

})
