package grifts

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/dimchansky/utfbom"
	"github.com/gobuffalo/nulls"
	"golang.org/x/crypto/bcrypt"

	"github.com/gocarina/gocsv"
	"github.com/markbates/grift/grift"
	"github.com/mogensen/fdfapp/models"
	"github.com/pkg/errors"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("truncate", "truncates a database")
	grift.Add("truncate", func(c *grift.Context) error {
		s := read("This will truncate the database!!! Are you sure? (y/N): ")

		if s == "y" || s == "yes" {
			// Add DB seeding stuff here
			if err := models.DB.TruncateAll(); err != nil {
				return errors.WithStack(err)
			}
		}
		return nil
	})

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {

		username := read("Username")
		ph, err := bcrypt.GenerateFromPassword([]byte(read("Password")), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user := models.User{
			Username:     username,
			PasswordHash: string(ph),
		}
		if err := models.DB.Create(&user); err != nil {
			return errors.WithStack(err)
		}

		classes := models.Classes{
			{UserID: user.ID, Name: "Puslinge"},
			{UserID: user.ID, Name: "Tumling"},
			{UserID: user.ID, Name: "Pilt"},
			{UserID: user.ID, Name: "Væbner"},
			{UserID: user.ID, Name: "Seniorvæbner"},
			{UserID: user.ID, Name: "Senior"},
		}

		for _, m := range classes {
			if err := models.DB.Create(&m); err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	})

	grift.Desc("fdfaa4", "fdfaa4 a database")
	grift.Add("fdfaa4", func(c *grift.Context) error {

		username := "fdfaa4"
		ph, err := bcrypt.GenerateFromPassword([]byte(read("Password")), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user := models.User{
			Username:     username,
			PasswordHash: string(ph),
		}
		if err := models.DB.Create(&user); err != nil {
			return errors.WithStack(err)
		}

		classes := models.Classes{
			{UserID: user.ID, Name: "Puslinge", Calendar: nulls.NewString("https://calendar.google.com/calendar/ical/fdf.dk_l5f2o8l8cbitivjsdclle5bm3k%40group.calendar.google.com/public/basic.ics")},
			{UserID: user.ID, Name: "Tumling", Calendar: nulls.NewString("https://calendar.google.com/calendar/ical/fdf.dk_l5f2o8l8cbitivjsdclle5bm3k%40group.calendar.google.com/public/basic.ics")},
			{UserID: user.ID, Name: "Pilt", Calendar: nulls.NewString("https://calendar.google.com/calendar/ical/fdf.dk_omj54ddubl3k4olhvq0395a99g%40group.calendar.google.com/public/basic.ics")},
			{UserID: user.ID, Name: "Væbner", Calendar: nulls.NewString("https://calendar.google.com/calendar/ical/fdf.dk_4qqn55n1v17thq0r08lb4g3i3o%40group.calendar.google.com/public/basic.ics")},
			{UserID: user.ID, Name: "Seniorvæbner", Calendar: nulls.NewString("https://calendar.google.com/calendar/ical/fdf.dk_76iisk34btll39hhkn1q0vu38s%40group.calendar.google.com/public/basic.ics")},
			{UserID: user.ID, Name: "Senior", Calendar: nulls.NewString("https://calendar.google.com/calendar/ical/fdf.dk_vmqn2cdcp1pcaa1g9lqqeu75d8%40group.calendar.google.com/public/basic.ics")},
		}

		for _, m := range classes {
			if err := models.DB.Create(&m); err != nil {
				return errors.WithStack(err)
			}
		}

		return nil
	})

	grift.Desc("load-participants", "load-participants from Carla csv into the database")
	grift.Add("load-participants", func(c *grift.Context) error {

		username := read("Username")
		user := &models.User{}
		models.DB.Where("username = ?", username).First(user)

		classes := &models.Classes{}
		models.DB.Where("user_id = ?", user.ID).All(classes)
		fmt.Println(classes)

		lines, err := ReadCsv("members.csv")
		if err != nil {
			panic(err)
		}

		// Loop through lines & turn into object
		for _, p := range lines {

			fmt.Printf("------------------\n")
			existingParticipant := &models.Participant{}
			q := models.DB.Where("member_number = ?", p.MemberNumber)
			err = q.First(existingParticipant)

			fmt.Printf("Participant: %s, %s\n", p.FirstName, p.MemberNumber)
			if err == nil {
				fmt.Printf("Participant exists: %s, %s\n", existingParticipant.FirstName, existingParticipant.MemberNumber)
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
					if err := models.DB.Create(&data); err != nil {
						return errors.WithStack(err)
					}

					// Add Participant to class
					classsMemberShip := models.ClassMembership{
						Class:       class,
						Participant: data,
					}
					if err := models.DB.Create(&classsMemberShip); err != nil {
						return errors.WithStack(err)
					}
				}
			}
		}

		return nil
	})

})

func read(info string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(info + ": ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)
	text = strings.TrimSuffix(text, "\n")
	return text
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
func ReadCsv(filename string) ([]*csvParticipant, error) {
	participants := []*csvParticipant{}

	// Open CSV file
	f, err := os.Open(filename)
	if err != nil {
		return participants, err
	}
	defer f.Close()

	o := utfbom.SkipOnly(f)

	gocsv.SetCSVReader(func(in io.Reader) gocsv.CSVReader {
		r := csv.NewReader(in)
		r.TrailingComma = true
		r.Comma = ';'
		return r
	})

	err = gocsv.Unmarshal(o, &participants)
	if err != nil {
		return nil, err
	}

	return participants, nil
}
