package actions

import "github.com/mogensen/fdfapp/models"

func (as *ActionSuite) Test_ActivitiesResource_List() {
	u := &models.User{
		Username:             "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}
	verrs, err := u.Create(as.DB)
	as.NoError(err)
	as.False(verrs.HasAny())
	as.Session.Set("current_user_id", u.ID)

	res := as.HTML("/activities").Get()
	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_ActivitiesResource_Show() {
	// as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_ActivitiesResource_Create() {
	// as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_ActivitiesResource_Update() {
	// as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_ActivitiesResource_Destroy() {
	// as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_ActivitiesResource_New() {
	// as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_ActivitiesResource_Edit() {
	// as.Fail("Not Implemented!")
}
