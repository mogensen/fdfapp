package actions

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/mogensen/fdfapp/models"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// AuthNew loads the signin page
func AuthNew(c buffalo.Context) error {

	if currentUser(c) != nil {
		return c.Redirect(302, "/classes")
	}

	c.Set("user", models.User{})
	return c.Render(200, r.HTML("auth/new.html"))
}

// AuthCreate attempts to log the user in with an existing account.
func AuthCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)

	// find a user with the username
	err := tx.Where("username = ?", strings.ToLower(strings.TrimSpace(u.Username))).First(u)

	// helper function to handle bad attempts
	bad := func() error {
		c.Set("user", u)
		verrs := validate.NewErrors()
		verrs.Add("username", "invalid username/password")
		c.Set("errors", verrs)
		return c.Render(422, r.HTML("auth/new.html"))
	}

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// couldn't find an user with the supplied username address.
			return bad()
		}
		fmt.Println("teste")
		return errors.WithStack(err)
	}

	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return bad()
	}
	c.Session().Set("current_user_id", u.ID)
	c.Flash().Add("success", T.Translate(c, "auth.login.success"))

	redirectURL := "/"
	if redir, ok := c.Session().Get("redirectURL").(string); ok {
		redirectURL = redir
	}

	return c.Redirect(302, redirectURL)
}

// AuthDestroy clears the session and logs a user out
func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", T.Translate(c, "auth.logout.success"))
	return c.Redirect(302, "/")
}
