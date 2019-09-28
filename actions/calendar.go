package actions

import (
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
	// Allocate an empty Class
	class := &models.Class{}

	// To find the Class the parameter class_id is used.
	if err := scope(c).Eager().Find(class, c.Param("class_id")); err != nil {
		return c.Error(404, err)
	}

	events := getCalenerEvents(class)

	c.Set("class", class)
	c.Set("events", events)

	return c.Render(200, r.HTML("calendar/show.html"))
}

func getCalenerEvents(class *models.Class) []gocal.Event {

	if !class.Calendar.Valid {
		return []gocal.Event{}
	}

	if class.Calendar.String == "" {
		return []gocal.Event{}
	}

	rc, err := downloadFile(class.Calendar.String)
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	start := time.Now().Add(-time.Duration(180) * 24 * time.Hour)
	end := time.Now()

	c := gocal.NewParser(rc)
	c.Start, c.End = &start, &end
	c.Parse()

	sort.Slice(c.Events, func(i, j int) bool {
		return c.Events[i].Start.Before(*(c.Events[j].Start))
	})

	return c.Events

}

func downloadFile(url string) (io.ReadCloser, error) {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
