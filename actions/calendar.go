package actions

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/apognu/gocal"
	"github.com/gobuffalo/buffalo"
	"github.com/mogensen/fdfapp/models"
)

// CalendarShow default implementation.
func CalendarShow(c buffalo.Context) error {
	c.Set("TIME_FORMAT", "02 Jan 2006")
	// Allocate an empty Class
	class := &models.Class{}

	// To find the Class the parameter class_id is used.
	if err := scope(c).Eager().Find(class, c.Param("class_id")); err != nil {
		return c.Error(404, err)
	}

	filteredEvents, err := getFilteredEvents(c, *class)
	if err != nil {
		c.Logger().Error(err)
	}

	c.Set("class", class)
	c.Set("events", filteredEvents)

	return c.Render(200, r.HTML("calendar/show.html"))
}

func getFilteredEvents(c buffalo.Context, class models.Class) ([]calEvent, error) {

	events := getCalenerEvents(class)
	activities := &models.Activities{}

	// Retrieve all Activities from the DB
	if err := scope(c).Where("class_id = (?)", class.ID).All(activities); err != nil {
		return nil, err
	}
	filteredEvents := []calEvent{}
	for _, event := range events {
		hasActivity := eventHasActivity(event, activities)
		if !hasActivity {
			filteredEvents = append(filteredEvents, event)
		}
	}
	return filteredEvents, nil
}

func eventHasActivity(e calEvent, activities *models.Activities) bool {
	year, month, day := e.Start.Date()
	for _, act := range *activities {
		ay, am, ad := act.Date.Date()
		if ay == year && am == month && ad == day {
			fmt.Printf("Already created activity for: %s == %s\n ", e.Summary, act.Title)
			return true
		}
	}
	return false
}

func getCalenerEvents(class models.Class) []calEvent {
	events := []calEvent{}

	if !class.Calendar.Valid {
		return events
	}

	if class.Calendar.String == "" {
		return events
	}

	rc, err := downloadFile(class.Calendar.String)
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	// 180 days ago
	start := time.Now().Add(-time.Duration(180) * 24 * time.Hour)

	// Earlist time that we show events for. (events before this is in the old system)
	earliest, _ := time.Parse("2006-01-02", "2020-08-01")
	if earliest.After(start) {
		start = earliest
	}

	// 1 day ahead
	end := time.Now().Add(24 * time.Hour)

	c := gocal.NewParser(rc)
	c.Start, c.End = &start, &end
	c.Parse()

	sort.Slice(c.Events, func(i, j int) bool {
		return c.Events[i].Start.After(*(c.Events[j].Start))
	})

	for _, event := range c.Events {

		events = append(events, calEvent{
			Summary:     event.Summary,
			Description: event.Summary,
			Start:       event.Start,
			End:         event.End,
			Location:    event.Location,
			Duration:    event.End.Sub(*event.Start).Hours(),
		})
	}
	return events

}

func downloadFile(url string) (io.ReadCloser, error) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

type calEvent struct {
	Summary     string
	Description string
	Start       *time.Time
	End         *time.Time
	Location    string
	Duration    float64
}
