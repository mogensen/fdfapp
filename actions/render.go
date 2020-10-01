package actions

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/helpers/hctx"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush/v4"
	"github.com/gobuffalo/tags"
	"github.com/gofrs/uuid"
	"github.com/mogensen/fdfapp/models"
)

var r *render.Engine
var assetsBox = packr.New("app:assets", "../public")

func init() {
	plush.DefaultTimeFormat = "02 January 2006"

	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",

		// Box containing all of the templates:
		TemplatesBox: packr.New("app:templates", "../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			"format": func(time nulls.Time, format string) string {
				return time.Time.Format(format)
			},
			"partipantHoursA": func(activity *models.Activity) string {
				hours := float64(len(activity.Participants)) * activity.Duration
				return fmt.Sprintf("%.02f", hours)
			},
			"partipantHours": func(activity models.Activity) string {
				hours := float64(len(activity.Participants)) * activity.Duration
				return fmt.Sprintf("%.02f", hours)
			},
			"checkboxChecked": func(id uuid.UUID, slice models.Participants) string {
				for _, c := range slice {
					if id == c.ID {
						return "checked"
					}
				}
				return ""
			},

			"isActive": func(name string, help hctx.HelperContext) string {
				if cp, ok := help.Value("current_route").(buffalo.RouteInfo); ok {
					if strings.HasPrefix(cp.PathName, name) {
						return "active"
					}
				}
				return "inactive"
			},
			"getParticipant": func(cms models.ClassMembership, participants *models.Participants) models.Participant {
				for _, p := range *participants {
					if p.ID == cms.ParticipantID {
						return p
					}
				}
				return models.Participant{}
			},
			"getClass": func(cms models.ClassMembership, classes *models.Classes) models.Class {
				for _, c := range *classes {
					if c.ID == cms.ClassID {
						return c
					}
				}
				return models.Class{}
			},

			"buttonGroupButton": func(text, icon, link string, help hctx.HelperContext) (template.HTML, error) {
				a := tags.New("a", tags.Options{"class": "btn btn-light btn-sm", "href": link})
				i := tags.New("i", tags.Options{"class": fmt.Sprintf("fas fa-%s", icon)})
				a.Append(i)
				a.Append(tags.New("br", tags.Options{}))
				a.Append(text)
				return a.HTML(), nil
			},

			"buttonGroup": func(floatRight bool, help hctx.HelperContext) (template.HTML, error) {
				float := "float-right"
				if !floatRight {
					float = "btn-block"
				}
				group := tags.New("div", tags.Options{"class": "btn-group " + float, "role": "group"})
				if help.HasBlock() {
					bc, err := help.Block()
					if err != nil {
						return "", err
					}
					group.Append(bc)
				}
				return group.HTML(), nil
			},
			"image": func(img *models.ParticipantsImage) bool {
				if img != nil && img.ID != uuid.Nil {
					return true
				}
				return false
			},
			"uuid": func(id uuid.UUID) bool {
				if id == uuid.Nil {
					return false
				}
				return true
			},

			"age": func(t nulls.Time, classes *models.Classes) string {
				year, _, _, _, _, _ := diff(t.Time, time.Now())
				return fmt.Sprintf("%d år", year)
			},

			// for non-bootstrap form helpers uncomment the lines
			// below and import "github.com/gobuffalo/helpers/forms"
			// forms.FormKey:     forms.Form,
			// forms.FormForKey:  forms.FormFor,
		},
	})
}

func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}
