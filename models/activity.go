package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Activity is an event, such as an evening meeting or a weekend trip
type Activity struct {
	ID           uuid.UUID    `json:"id" db:"id"`
	UserID       uuid.UUID    `db:"user_id"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
	Title        string       `json:"title" db:"title"`
	Duration     float64      `json:"duration" db:"duration"`
	Date         time.Time    `json:"date" db:"date"`
	ClassID      uuid.UUID    `json:"class_id" db:"class_id"`
	Class        Class        `belongs_to:"classes" db:"-"`
	Participants Participants `many_to_many:"activity_participants" db:"-"`
}

// String is not required by pop and may be deleted
func (a Activity) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Activities is not required by pop and may be deleted
type Activities []Activity

// String is not required by pop and may be deleted
func (a Activities) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Activity) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.Title, Name: "Title"},
		&validators.TimeIsPresent{Field: a.Date, Name: "Date"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Activity) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Activity) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
