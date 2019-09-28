package grifts

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/dimchansky/utfbom"
	"github.com/gobuffalo/nulls"

	"github.com/gocarina/gocsv"
	"github.com/markbates/grift/grift"
	"github.com/mogensen/fdfapp/models"
	"github.com/pkg/errors"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		// Add DB seeding stuff here
		if err := models.DB.TruncateAll(); err != nil {
			return errors.WithStack(err)
		}
		classes := models.Classes{
			// {Name: "Numlinge"},
			// {Name: "Puslinge"},
			{Name: "Tumling", Calendar: nulls.NewString("https://calendar.google.com/calendar/ical/fdf.dk_l5f2o8l8cbitivjsdclle5bm3k%40group.calendar.google.com/public/basic.ics")},
			{Name: "Pilt", Calendar: nulls.NewString("https://calendar.google.com/calendar/ical/fdf.dk_omj54ddubl3k4olhvq0395a99g%40group.calendar.google.com/public/basic.ics")},
			{Name: "Væbner", Calendar: nulls.NewString("https://calendar.google.com/calendar/ical/fdf.dk_4qqn55n1v17thq0r08lb4g3i3o%40group.calendar.google.com/public/basic.ics")},
			{Name: "Seniorvæbner", Calendar: nulls.NewString("https://calendar.google.com/calendar/ical/fdf.dk_76iisk34btll39hhkn1q0vu38s%40group.calendar.google.com/public/basic.ics")},
			{Name: "Senior", Calendar: nulls.NewString("https://calendar.google.com/calendar/ical/fdf.dk_vmqn2cdcp1pcaa1g9lqqeu75d8%40group.calendar.google.com/public/basic.ics")},
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

		classes := &models.Classes{}
		models.DB.All(classes)
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

type DateTime struct {
	time.Time
}

type CsvParticipant struct {
	MemberNumber string   `csv:"Medlemsnummer"`
	FirstName    string   `csv:"Fornavn"`
	MiddelName   string   `csv:"Mellemnavn"`
	LastName     string   `csv:"Efternavn"`
	Phone        string   `csv:"Telefon"`
	DateOfBirth  DateTime `csv:"Fødselsdag"`
	Class        string   `csv:"Klasse"`
}

// UnmarshalCSV converts the CSV string as internal date
func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("02-01-2006", csv)
	return err
}

// ReadCsv accepts a file and returns its content as a multi-dimentional type
// with lines and each column. Only parses to string type.
func ReadCsv(filename string) ([]*CsvParticipant, error) {
	participants := []*CsvParticipant{}

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
