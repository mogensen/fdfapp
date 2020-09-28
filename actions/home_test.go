package actions

import "github.com/mogensen/fdfapp/models"

func (as *ActionSuite) Test_HomeHandler() {
	res := as.HTML("/").Get()
	as.Equal(302, res.Code)
	as.Contains(res.Body.String(), "/signin")
}
func (as *ActionSuite) Test_SigninHandler() {
	res := as.HTML("/signin").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Sign In")
}

func (as *ActionSuite) Test_HomeHandler_LoggedIn() {
	u := &models.User{
		Username:             "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}
	verrs, err := u.Create(as.DB)
	as.NoError(err)
	as.False(verrs.HasAny())
	as.Session.Set("current_user_id", u.ID)

	// Redirect to classes if logged in
	res := as.HTML("/signin").Get()
	as.Equal(302, res.Code)
	as.Contains(res.Body.String(), "/classes")

	as.Session.Clear()

	// Redirect to / if not loggedin
	res = as.HTML("/classes").Get()
	as.Equal(302, res.Code)
	as.Contains(res.Body.String(), "/")

	// Redirect to / if not loggedin
	res = as.HTML("/").Get()
	as.Equal(302, res.Code)
	as.Contains(res.Body.String(), "/signin")
}
