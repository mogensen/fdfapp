package actions

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/mogensen/fdfapp/models"
)

// ReportShow default implementation.
func ReportShow(c buffalo.Context) error {
	year := time.Now().Year()
	return c.Redirect(302, fmt.Sprintf("/report/%d", year))
}

// ReportYearShow default implementation.
func ReportYearShow(c buffalo.Context) error {
	c.Set("TIME_FORMAT", "02 Jan 2006")
	activities := &models.Activities{}

	year := time.Now().Year()
	var err error
	yearParam := c.Param("year")
	if yearParam != "" {
		year, err = strconv.Atoi(yearParam)
		if err != nil {
			return err
		}
		if year < 2000 || year > 2050 {
			return fmt.Errorf("Year is not supported: %d", year)
		}
	}

	then := time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
	q := scope(c).Where("date >= ?", then).Where("date <= ?", then.AddDate(1, 0, 0)).Order("date desc")

	// Retrieve all Activities from the DB
	if err := q.Eager().All(activities); err != nil {
		return err
	}

	c.Set("activities", activities)

	report := Report{
		Summary:     "test",
		Description: "much stuff",
		AgeUnknown:  0,
		Age3_5:      0,
		Age6_18:     0,
		Age19_24:    0,
		Age25_up:    0,
	}

	for _, activity := range *activities {
		for _, participant := range activity.Participants {

			// We should create a warning here
			years := 10
			if participant.DateOfBirth.Valid {
				years, _, _, _, _, _ = diff(participant.DateOfBirth.Time, time.Now())
			} else {
				report.AgeUnknown += activity.Duration
			}

			if years >= 3 && years <= 5 {
				report.Age3_5 += activity.Duration
			}

			if years >= 6 && years <= 18 {
				report.Age6_18 += activity.Duration
			}

			if years >= 19 && years <= 24 {
				report.Age19_24 += activity.Duration
			}

			if years >= 25 {
				report.Age25_up += activity.Duration
			}
		}
	}

	c.Set("report", report)
	return c.Render(200, r.HTML("report/show.html"))
}

type Report struct {
	AgeUnknown  float64
	Age3_5      float64
	Age6_18     float64
	Age19_24    float64
	Age25_up    float64
	Summary     string
	Description string
}
