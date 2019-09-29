package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/validate"
	"github.com/gofrs/uuid"
)

// ClassMembership details that a participant is member of a given class
type ClassMembership struct {
	ID            uuid.UUID   `json:"id" db:"id"`
	CreatedAt     time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at" db:"updated_at"`
	ClassID       uuid.UUID   `json:"class_id" db:"class_id"`
	ParticipantID uuid.UUID   `json:"participant_id" db:"participant_id"`
	Class         Class       `belongs_to:"class" db:"-"`
	Participant   Participant `belongs_to:"participant" db:"-"`
}

// String is not required by pop and may be deleted
func (c ClassMembership) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// ClassMemberships is not required by pop and may be deleted
type ClassMemberships []ClassMembership

// String is not required by pop and may be deleted
func (c ClassMemberships) String() string {
	jc, _ := json.Marshal(c)
	return string(jc)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (c *ClassMembership) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (c *ClassMembership) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (c *ClassMembership) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
