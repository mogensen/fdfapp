package actions

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/dimchansky/utfbom"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/binding"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v5"
	"github.com/gocarina/gocsv"
	"github.com/mogensen/fdfapp/models"
	"github.com/pkg/errors"
)

// ParticipantsUploadsResource is the resource for the ParticipantsUpload model
type ParticipantsUploadsResource struct {
	buffalo.Resource
}

// ParticipantsUploads is used for holding the form data
type ParticipantsUploads struct {
	MyFile binding.File `db:"-" form:"carlaCsvAlleMedlemmer"`
}

// New renders the form for creating a new ParticipantsUpload.
// This function is mapped to the path GET /participants_uploads/new
func (v ParticipantsUploadsResource) New(c buffalo.Context) error {
	return c.Render(200, r.Auto(c, &ParticipantsUploads{}))
}

// Create adds a ParticipantsUpload to the DB. This function is mapped to the
// path POST /participants_uploads
func (v ParticipantsUploadsResource) Create(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	// Allocate an empty ParticipantsUpload
	participantsUpload := &ParticipantsUploads{}

	f, err := c.File("carlaCsvAlleMedlemmer")
	if err != nil {
		return errors.WithStack(err)
	}

	if !f.Valid() {
		return nil
	}

	classes := &models.Classes{}
	scope(c).All(classes)

	lines, err := readCsv(f)
	if err != nil {
		return c.Error(500, err)
	}
	user := currentUser(c)

	created := 0
	updatedClass := 0
	missingClass := 0

	// Loop through lines & turn into object
	for _, p := range lines {

		fmt.Printf("------------------\n")
		existingParticipant := &models.Participant{}
		q := scope(c).Where("member_number = ?", p.MemberNumber).Eager()
		err = q.First(existingParticipant)

		fmt.Printf("Participant: %s, %s\n", p.FirstName, p.MemberNumber)
		if err == nil {
			fmt.Printf("Participant exists: %s, %s\n", existingParticipant.FirstName, existingParticipant.MemberNumber)
			updatedClass++

			for _, v := range existingParticipant.Memberships {

				classMembership := &models.ClassMembership{}

				// To find the ClassMembership the parameter class_membership_id is used.
				if err := tx.Find(classMembership, v.ID); err != nil {
					return c.Error(404, err)
				}

				if err := tx.Destroy(classMembership); err != nil {
					return err
				}
			}
			existingParticipant.Memberships = models.ClassMemberships{}

			err := tx.Update(existingParticipant)
			if err != nil {
				return err
			}

			for _, class := range *classes {
				// Create participant only if his class exist
				if class.Name == p.Class {
					// Add Participant to class
					classsMemberShip := models.ClassMembership{
						Class:       class,
						Participant: *existingParticipant,
					}
					verrs, err := tx.ValidateAndCreate(&classsMemberShip)
					if err != nil {
						return err
					}
					if verrs.HasAny() {
						// Make the errors available inside the html template
						c.Set("errors", verrs)
						c.Logger().Warn(verrs)

						// Render again the new.html template that the user can
						// correct the input.
						return c.Render(422, r.Auto(c, participantsUpload))
					}
				}
			}

			continue
		}

		data := models.Participant{
			UserID:       user.ID,
			FirstName:    strings.TrimSpace(p.FirstName + " " + p.MiddelName),
			LastName:     p.LastName,
			MemberNumber: p.MemberNumber,
			DateOfBirth:  nulls.NewTime(p.DateOfBirth.Time),
			Phone:        p.Phone,
			Classes:      models.Classes{},
			Memberships:  models.ClassMemberships{},
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		fmt.Printf("%v\n", data)

		for _, class := range *classes {
			// Create participant only if his class exist
			if class.Name == p.Class {
				verrs, err := tx.ValidateAndCreate(&data)
				created++
				if err != nil {
					return err
				}

				if verrs.HasAny() {
					// Make the errors available inside the html template
					c.Set("errors", verrs)
					c.Logger().Warn(verrs)

					// Render again the new.html template that the user can
					// correct the input.
					return c.Render(422, r.Auto(c, participantsUpload))
				}

				// Add Participant to class
				classsMemberShip := models.ClassMembership{
					Class:       class,
					Participant: data,
				}
				verrs, err = tx.ValidateAndCreate(&classsMemberShip)
				if err != nil {
					return err
				}

				if verrs.HasAny() {
					// Make the errors available inside the html template
					c.Set("errors", verrs)
					c.Logger().Warn(verrs)

					// Render again the new.html template that the user can
					// correct the input.
					return c.Render(422, r.Auto(c, participantsUpload))
				}
				continue
			}
		}
		missingClass++
	}

	// If there are no errors set a success message
	c.Flash().Add("success", fmt.Sprintf("Success. %d medlemmer oprettet, %d medlemmer opdateret klasse, %d medlemmer hvor klassen ikke findes.", created, updatedClass, missingClass))
	// and redirect to the participants_uploads index page
	return c.Redirect(302, "/participants")
}

type dateTime struct {
	time.Time
}

type csvParticipant struct {
	MemberNumber string   `csv:"Medlemsnummer"`
	FirstName    string   `csv:"Fornavn"`
	MiddelName   string   `csv:"Mellemnavn"`
	LastName     string   `csv:"Efternavn"`
	Phone        string   `csv:"Telefon"`
	DateOfBirth  dateTime `csv:"Fødselsdag"`
	Class        string   `csv:"Klasse"`
}

// UnmarshalCSV converts the CSV string as internal date
func (date *dateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("02-01-2006", csv)
	return err
}

// ReadCsv accepts a file and returns its content as csvParticipants
func readCsv(reader io.Reader) ([]*csvParticipant, error) {
	participants := []*csvParticipant{}

	o := utfbom.SkipOnly(reader)

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.FieldsPerRecord = -1
		r.Comma = ';'
		return r
	})

	err := gocsv.Unmarshal(o, &participants)
	if err != nil {
		return nil, err
	}

	return participants, nil
}
