package actions

import (
	"errors"
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/mogensen/fdfapp/models"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (ClassMembership)
// DB Table: Plural (class_memberships)
// Resource: Plural (ClassMemberships)
// Path: Plural (/class_memberships)
// View Template Folder: Plural (/templates/class_memberships/)

// ClassMembershipsResource is the resource for the ClassMembership model
type ClassMembershipsResource struct {
	buffalo.Resource
}

// List gets all ClassMemberships. This function is mapped to the path
// GET /class_memberships
func (v ClassMembershipsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	classMemberships := &models.ClassMemberships{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all ClassMemberships from the DB
	if err := q.All(classMemberships); err != nil {
		return err
	}

	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)

	return c.Render(200, r.Auto(c, classMemberships))
}

// Show gets the data for one ClassMembership. This function is mapped to
// the path GET /class_memberships/{class_membership_id}
func (v ClassMembershipsResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	// Allocate an empty ClassMembership
	classMembership := &models.ClassMembership{}

	// To find the ClassMembership the parameter class_membership_id is used.
	if err := tx.Eager().Find(classMembership, c.Param("class_membership_id")); err != nil {
		return c.Error(404, err)
	}

	return c.Render(200, r.Auto(c, classMembership))
}

// New renders the form for creating a new ClassMembership.
// This function is mapped to the path GET /class_memberships/new
func (v ClassMembershipsResource) New(c buffalo.Context) error {
	if err := bindClasses(c); err != nil {
		return errors.New("No Classes found")
	}

	if err := bindParticipants(c); err != nil {
		return errors.New("No Participants found")
	}

	return c.Render(200, r.Auto(c, &models.ClassMembership{}))
}

// Create adds a ClassMembership to the DB. This function is mapped to the
// path POST /class_memberships
func (v ClassMembershipsResource) Create(c buffalo.Context) error {
	// Allocate an empty ClassMembership
	classMembership := &models.ClassMembership{}

	// Bind classMembership to the html form elements
	if err := c.Bind(classMembership); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(classMembership)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the new.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, classMembership))
	}

	participant := &models.Participant{}

	// To find the Participant the parameter participant_id is used.
	if err := tx.Eager().Find(participant, classMembership.ParticipantID); err != nil {
		return c.Error(404, err)
	}

	// If there are no errors set a success message
	c.Flash().Add("success", T.Translate(c, "classMembership.created.success"))
	// and redirect to the class_memberships index page
	x := fmt.Sprintf("/participants/%s", classMembership.ParticipantID)
	return c.Redirect(302, x)
}

// Edit renders a edit form for a ClassMembership. This function is
// mapped to the path GET /class_memberships/{class_membership_id}/edit
func (v ClassMembershipsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	// Allocate an empty ClassMembership
	classMembership := &models.ClassMembership{}

	if err := tx.Find(classMembership, c.Param("class_membership_id")); err != nil {
		return c.Error(404, err)
	}

	classes := &models.Classes{}

	// Retrieve all Classes from the DB
	if err := tx.All(classes); err != nil {
		return err
	}

	fmt.Println("Finding classes")
	// Add the paginator to the context so it can be used in the template.
	c.Set("classes", classes)

	return c.Render(200, r.Auto(c, classMembership))
}

// Update changes a ClassMembership in the DB. This function is mapped to
// the path PUT /class_memberships/{class_membership_id}
func (v ClassMembershipsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	// Allocate an empty ClassMembership
	classMembership := &models.ClassMembership{}

	if err := tx.Find(classMembership, c.Param("class_membership_id")); err != nil {
		return c.Error(404, err)
	}

	// Bind ClassMembership to the html form elements
	if err := c.Bind(classMembership); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(classMembership)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		// Make the errors available inside the html template
		c.Set("errors", verrs)

		// Render again the edit.html template that the user can
		// correct the input.
		return c.Render(422, r.Auto(c, classMembership))
	}

	// If there are no errors set a success message
	c.Flash().Add("success", T.Translate(c, "classMembership.updated.success"))
	// and redirect to the class_memberships index page
	return c.Render(200, r.Auto(c, classMembership))
}

// Destroy deletes a ClassMembership from the DB. This function is mapped
// to the path DELETE /class_memberships/{class_membership_id}
func (v ClassMembershipsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return errors.New("no transaction found")
	}

	// Allocate an empty ClassMembership
	classMembership := &models.ClassMembership{}

	// To find the ClassMembership the parameter class_membership_id is used.
	if err := tx.Find(classMembership, c.Param("class_membership_id")); err != nil {
		return c.Error(404, err)
	}

	if err := tx.Destroy(classMembership); err != nil {
		return err
	}

	// If there are no errors set a flash message
	c.Flash().Add("success", T.Translate(c, "classMembership.destroyed.success"))
	// Redirect to the class_memberships index page
	x := fmt.Sprintf("/participants/%s", classMembership.ParticipantID)
	return c.Redirect(302, x)
}
