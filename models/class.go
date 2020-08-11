package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Class is a logical collection of participants that have activities together
type Class struct {
	ID           uuid.UUID        `json:"id" db:"id"`
	UserID       uuid.UUID        `db:"user_id"`
	CreatedAt    time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at" db:"updated_at"`
	Name         string           `json:"name" db:"name"`
	Calendar     nulls.String     `json:"calendar" db:"calendar"`
	Participants Participants     `many_to_many:"class_memberships" db:"-" order_by:"first_name asc"`
	Memberships  ClassMemberships `has_many:"class_memberships" db:"-"`
}

// SelectLabel is used for creating dropdown boxes in plush
func (c Class) SelectLabel() string {
	return c.Name
}

// SelectValue is used for creating dropdown boxes in plush
func (c Class) SelectValue() interface{} {
	return c.ID
}

// String is not required by pop and may be deleted
func (c Class) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Classes is not required by pop and may be deleted
type Classes []Class

// String is not required by pop and may be deleted
func (c Classes) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *Class) Validate(tx *pop.Connection) (*validate.Errors, error) {
	// If calendar is set, we validate it.
	if c.Calendar.Valid && c.Calendar.String != "" {
		return validate.Validate(
			&validators.StringIsPresent{Field: c.Name, Name: "Name"},
			&validators.URLIsPresent{Field: c.Calendar.String, Name: "Calendar"},
		), nil
	}
	return validate.Validate(
		&validators.StringIsPresent{Field: c.Name, Name: "Name"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *Class) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *Class) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
