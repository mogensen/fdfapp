package grifts

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/brianvoe/gofakeit"
	"github.com/gobuffalo/nulls"
	"github.com/ipsn/go-adorable"
	"github.com/mogensen/fdfapp/models"
)

func createParticipant(user models.User, year int, class models.Class) {
	dobStr := randomdata.FullDateInRange(fmt.Sprintf("%d-01-01", year-1), fmt.Sprintf("%d-12-31", year))
	fmt.Printf("dob %s\n", dobStr)
	dob, err := time.Parse("Monday 2 Jan 2006", dobStr)
	if err != nil {
		log.Fatalf("Error parsing date: %v", err)
	}
	participant := models.Participant{
		UserID:       user.ID,
		FirstName:    randomdata.FirstName(randomdata.RandomGender),
		LastName:     randomdata.LastName(),
		MemberNumber: randomdata.StringNumber(4, ""),
		DateOfBirth:  nulls.NewTime(dob),
		Phone:        randomdata.Digits(8),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	models.DB.Create(&participant)
	models.DB.Create(&models.ClassMembership{
		ClassID:       class.ID,
		ParticipantID: participant.ID,
	})

	avatar := adorable.Random()

	pImage := models.ParticipantsImage{
		ImageData:     []byte(base64.StdEncoding.EncodeToString(avatar)),
		ParticipantID: participant.ID,
	}
	models.DB.Create(&pImage)
}

// StartYear and EndYear is relative to now
//  eg: randomTimestamp(4, 2) gives a random time between two and four years afo
func randomTimestamp(startYear, endYear int) time.Time {
	start := time.Now().AddDate(-startYear, 0, 0)
	end := time.Now().AddDate(-endYear, 0, 0)
	return gofakeit.DateRange(start, end)
}
