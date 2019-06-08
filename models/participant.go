package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"github.com/gofrs/uuid"
)

type Participant struct {
	ID          uuid.UUID        `json:"id" db:"id"`
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
	Name        string           `json:"name" db:"name"`
	Phone       string           `json:"phone" db:"phone"`
	DateOfBirth nulls.Time       `json:"date_of_birth" db:"date_of_birth"`
	Classes     Classes          `many_to_many:"class_memberships" db:"-"`
	Memberships ClassMemberships `has_many:"class_memberships" db:"-"`
}

func (p Participant) SelectLabel() string {
	return p.Name
}

func (p Participant) SelectValue() interface{} {
	return p.ID
}

// String is not required by pop and may be deleted
func (p Participant) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Participants is not required by pop and may be deleted
type Participants []Participant

// String is not required by pop and may be deleted
func (p Participants) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *Participant) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: p.Name, Name: "Name"},
		&validators.StringIsPresent{Field: p.Phone, Name: "Phone"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *Participant) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *Participant) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
