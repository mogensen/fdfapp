package actions

func (as *ActionSuite) Test_ActivitiesResource_List() {
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
