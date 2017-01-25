package api

import (
	"errors"
	"testing"

	"github.com/qarea/ctxtg"
	"github.com/qarea/ctxtg/ctxtgtest"
	"github.com/qarea/jirams/entities"
	"github.com/stretchr/testify/assert"
)

type TestTrackerClient struct {
	getProjects      func(tracker entities.TrackerConfig, res *[]entities.Project) error
	getCurrentUser   func(tracker entities.TrackerConfig, res *entities.User) error
	getProjectIssues func(tracker entities.TrackerConfig, projectID entities.ProjectID, userID entities.UserID, res *[]entities.Issue) error
	createIssue      func(tracker entities.TrackerConfig, issue entities.NewIssue, res *entities.Issue) error
	getIssue         func(tracker entities.TrackerConfig, issueID entities.IssueID, res *entities.Issue) error
	createReport     func(tracker entities.TrackerConfig, report entities.Report) error
	getTotalReports  func(tracker entities.TrackerConfig, date entities.Timestamp, res *entities.ReportsTotal) error
	getIssueByUrl    func(tracker entities.TrackerConfig, issueURL string, res *entities.Issue, res2 *entities.ProjectID) error
}

func (t *TestTrackerClient) GetProjects(tracker entities.TrackerConfig, res *[]entities.Project) error {
	return t.getProjects(tracker, res)
}

func (t *TestTrackerClient) GetCurrentUser(tracker entities.TrackerConfig, res *entities.User) error {
	return t.getCurrentUser(tracker, res)
}

func (t *TestTrackerClient) GetProjectIssues(tracker entities.TrackerConfig, projectID entities.ProjectID, userID entities.UserID, res *[]entities.Issue) error {
	return t.getProjectIssues(tracker, projectID, userID, res)
}

func (t *TestTrackerClient) CreateIssue(tracker entities.TrackerConfig, issue entities.NewIssue, res *entities.Issue) error {
	return t.createIssue(tracker, issue, res)
}

func (t *TestTrackerClient) GetIssue(tracker entities.TrackerConfig, issueID entities.IssueID, res *entities.Issue) error {
	return t.getIssue(tracker, issueID, res)
}

func (t *TestTrackerClient) CreateReport(tracker entities.TrackerConfig, report entities.Report) error {
	return t.createReport(tracker, report)
}

func (t *TestTrackerClient) GetTotalReports(tracker entities.TrackerConfig, date entities.Timestamp, res *entities.ReportsTotal) error {
	return t.getTotalReports(tracker, date, res)
}

func (t *TestTrackerClient) GetIssueByURL(tracker entities.TrackerConfig, issueURL string, res *entities.Issue, res2 *entities.ProjectID) error {
	return t.getIssueByUrl(tracker, issueURL, res, res2)
}

func TestGetProjects(t *testing.T) {
	a := assert.New(t)

	var uid ctxtg.UserID = 1
	var token ctxtg.Token = "dsa"

	p := &ctxtgtest.Parser{
		TokenExpected: token,
		Claims: ctxtg.Claims{
			UserID: uid,
		},
		Err: nil,
	}

	var prs = []entities.Project{
		{
			ID:            1,
			Title:         "test title",
			Description:   "desc",
			IssueTypes:    []entities.NamedID{{1, "fdsdadad"}},
			ActivityTypes: []entities.NamedID{{2, "fdsdadadwwwer"}},
		},
	}

	req := GetProjectsRequest{
		Context: ctxtg.Context{Token: token},
		Tracker: entities.TrackerConfig{
			ID:  1,
			URL: "http://tracker.com",
			Credentials: entities.TrackerCredentials{
				Login:    "login",
				Password: "password",
			},
		},
	}

	st := &TestTrackerClient{
		getProjects: func(tracker entities.TrackerConfig, res *[]entities.Project) error {
			a.Equal(req.Tracker, tracker)
			*res = prs
			return nil
		},
	}

	api := API{st, p}
	var res GetProjectsResponse
	err := api.GetProjects(req, &res)
	a.NoError(err)
	a.Equal(prs, res.Projects)

	a.NoError(p.Error())

}

func TestGetProjectsWithClientError(t *testing.T) {
	a := assert.New(t)

	var uid ctxtg.UserID = 1
	var token ctxtg.Token = "dsa"

	p := &ctxtgtest.Parser{
		TokenExpected: token,
		Claims: ctxtg.Claims{
			UserID: uid,
		},
		Err: nil,
	}

	req := GetProjectsRequest{
		Context: ctxtg.Context{Token: token},
		Tracker: entities.TrackerConfig{
			ID:  1,
			URL: "http://tracker.com",
			Credentials: entities.TrackerCredentials{
				Login:    "login",
				Password: "password",
			},
		},
	}
	//trClientError := errors.New("My error")
	st := &TestTrackerClient{
		getProjects: func(tracker entities.TrackerConfig, res *[]entities.Project) error {
			a.Equal(req.Tracker, tracker)
			return errors.New("My error")
		},
	}

	api := API{st, p}
	var res GetProjectsResponse
	err := api.GetProjects(req, &res)
	//a.Empty() look here
	a.Len(res.Projects, 0)
	a.Error(err)
	//a.Equal(trClientError, err)
	a.NoError(p.Error())
}

func TestGetCurrentUser(t *testing.T) {
	a := assert.New(t)
	var uid ctxtg.UserID = 1
	var token ctxtg.Token = "dsa"
	p := &ctxtgtest.Parser{
		TokenExpected: token,
		Claims: ctxtg.Claims{
			UserID: uid,
		},
		Err: nil,
	}
	req := GetCurrentUserRequest{
		Context: ctxtg.Context{Token: token},
		Tracker: entities.TrackerConfig{
			ID:  1,
			URL: "http://tracker.com",
			Credentials: entities.TrackerCredentials{
				Login:    "login",
				Password: "password",
			},
		},
	}

	usr := entities.User{1, "user1", "dada@company.com"}
	st := &TestTrackerClient{
		getCurrentUser: func(tracker entities.TrackerConfig, res *entities.User) error {
			a.Equal(req.Tracker, tracker)
			*res = usr
			return nil
		},
	}

	api := API{st, p}
	var res GetCurrentUserResponse
	err := api.GetCurrentUser(req, &res)
	a.NoError(err)
	a.Equal(usr, res.User, "Should be equal")
}

func TestGetCurrentUserWithError(t *testing.T) {
	a := assert.New(t)
	var uid ctxtg.UserID = 1
	var token ctxtg.Token = "dsa"
	p := &ctxtgtest.Parser{
		TokenExpected: token,
		Claims: ctxtg.Claims{
			UserID: uid,
		},
		Err: nil,
	}
	req := GetCurrentUserRequest{
		Context: ctxtg.Context{Token: token},
		Tracker: entities.TrackerConfig{
			ID:  1,
			URL: "http://tracker.com",
			Credentials: entities.TrackerCredentials{
				Login:    "login",
				Password: "password",
			},
		},
	}

	st := &TestTrackerClient{
		getCurrentUser: func(tracker entities.TrackerConfig, res *entities.User) error {
			a.Equal(req.Tracker, tracker)
			return errors.New("GetCurrentUsers error")
		},
	}

	api := API{st, p}
	var res GetCurrentUserResponse
	err := api.GetCurrentUser(req, &res)
	a.NotNil(err)

}

func TestGetProjectIssues(t *testing.T) {
	a := assert.New(t)
	var uid ctxtg.UserID = 1
	// var pid entities.ProjectID = 123
	// var usrID entities.UserID = 2
	var token ctxtg.Token = "dsa"
	p := &ctxtgtest.Parser{
		TokenExpected: token,
		Claims: ctxtg.Claims{
			UserID: uid,
		},
		Err: nil,
	}

	req := GetProjectIssuesRequest{
		Context: ctxtg.Context{Token: token},
		Tracker: entities.TrackerConfig{
			ID:  1,
			URL: "http://tracker.com",
			Credentials: entities.TrackerCredentials{
				Login:    "login",
				Password: "password",
			},
		},
	}

	entity := entities.Issue{Title: "title", URL: "http://tracker.com"}
	var issues []entities.Issue
	issues = append(issues, entity)

	st := &TestTrackerClient{
		getProjectIssues: func(tracker entities.TrackerConfig, projectID entities.ProjectID, userID entities.UserID, res *[]entities.Issue) error {
			a.Equal(req.Tracker, tracker)
			*res = issues
			return nil
		},
	}

	api := NewRPCAPI(st, p)
	var res GetProjectIssuesResponse
	err := api.GetProjectIssues(req, &res)
	a.NoError(err)
	a.Equal(issues, res.Issues)

}

func TestGetProjectIssuesWithError(t *testing.T) {
	a := assert.New(t)
	var uid ctxtg.UserID = 1
	var token ctxtg.Token = "dsa"
	p := &ctxtgtest.Parser{
		TokenExpected: token,
		Claims: ctxtg.Claims{
			UserID: uid,
		},
		Err: nil,
	}

	req := GetProjectIssuesRequest{
		Context: ctxtg.Context{Token: token},
		Tracker: entities.TrackerConfig{
			ID:  1,
			URL: "http://tracker.com",
			Credentials: entities.TrackerCredentials{
				Login:    "login",
				Password: "password",
			},
		},
	}

	entity := entities.Issue{Title: "title", URL: "http://tracker.com"}
	var issues []entities.Issue
	issues = append(issues, entity)

	st := &TestTrackerClient{
		getProjectIssues: func(tracker entities.TrackerConfig, projectID entities.ProjectID, userID entities.UserID, res *[]entities.Issue) error {
			a.Equal(req.Tracker, tracker)
			return errors.New("Error with getting project issues")
		},
	}

	api := &API{st, p}
	var res GetProjectIssuesResponse
	err := api.GetProjectIssues(req, &res)
	a.Error(err)

}

func TestCreateIssue(t *testing.T) {
	a := assert.New(t)
	var uid ctxtg.UserID = 1
	var token ctxtg.Token = "dsa"
	p := &ctxtgtest.Parser{
		TokenExpected: token,
		Claims: ctxtg.Claims{
			UserID: uid,
		},
		Err: nil,
	}

	req := CreateIssueRequest{
		Context: ctxtg.Context{Token: token},
		Tracker: entities.TrackerConfig{
			ID:  1,
			URL: "http://tracker.com",
			Credentials: entities.TrackerCredentials{
				Login:    "login",
				Password: "password",
			},
		},
	}

	iss := entities.Issue{Title: "title", URL: "site.com"}
	st := &TestTrackerClient{
		createIssue: func(tracker entities.TrackerConfig, issue entities.NewIssue, res *entities.Issue) error {
			a.Equal(req.Tracker, tracker)
			*res = iss
			return nil
		},
	}

	api := &API{st, p}
	var res CreateIssueResponse
	err := api.CreateIssue(req, &res)
	a.NoError(err)
	a.Equal(iss, res.Issue, "Should be equal")

}

func TestCreateIssueWithError(t *testing.T) {
	a := assert.New(t)
	var uid ctxtg.UserID = 1
	var token ctxtg.Token = "dsa"
	p := &ctxtgtest.Parser{
		TokenExpected: token,
		Claims: ctxtg.Claims{
			UserID: uid,
		},
		Err: nil,
	}

	req := CreateIssueRequest{
		Context: ctxtg.Context{Token: token},
		Tracker: entities.TrackerConfig{
			ID:  1,
			URL: "http://tracker.com",
			Credentials: entities.TrackerCredentials{
				Login:    "login",
				Password: "password",
			},
		},
	}

	st := &TestTrackerClient{
		createIssue: func(tracker entities.TrackerConfig, issue entities.NewIssue, res *entities.Issue) error {
			a.Equal(req.Tracker, tracker)
			return errors.New("Error")
		},
	}

	api := &API{st, p}
	var res CreateIssueResponse
	err := api.CreateIssue(req, &res)
	a.NotNil(err)

}

func TestGetIssue(t *testing.T) {

	a := assert.New(t)
	var uid ctxtg.UserID = 1
	var token ctxtg.Token = "dsa"
	p := &ctxtgtest.Parser{
		TokenExpected: token,
		Claims: ctxtg.Claims{
			UserID: uid,
		},
		Err: nil,
	}

	req := GetIssueRequest{
		Context: ctxtg.Context{Token: token},
		Tracker: entities.TrackerConfig{
			ID:  1,
			URL: "http://tracker.com",
			Credentials: entities.TrackerCredentials{
				Login:    "login",
				Password: "password",
			},
		},
	}
	iss := entities.Issue{URL: "somesite.com", Title: "title"}
	st := &TestTrackerClient{
		getIssue: func(tracker entities.TrackerConfig, issueID entities.IssueID, res *entities.Issue) error {
			a.Equal(req.Tracker, tracker)
			*res = iss
			return nil
		},
	}

	api := &API{st, p}
	var res GetIssueResponse
	err := api.GetIssue(req, &res)
	a.Equal(iss, res.Issue)
	a.NoError(err)
}

func TestGetIssueWithError(t *testing.T) {

	a := assert.New(t)
	var uid ctxtg.UserID = 1
	var token ctxtg.Token = "dsa"
	p := &ctxtgtest.Parser{
		TokenExpected: token,
		Claims: ctxtg.Claims{
			UserID: uid,
		},
		Err: nil,
	}

	req := GetIssueRequest{
		Context: ctxtg.Context{Token: token},
		Tracker: entities.TrackerConfig{
			ID:  1,
			URL: "http://tracker.com",
			Credentials: entities.TrackerCredentials{
				Login:    "login",
				Password: "password",
			},
		},
	}

	st := &TestTrackerClient{
		getIssue: func(tracker entities.TrackerConfig, issueID entities.IssueID, res *entities.Issue) error {
			a.Equal(req.Tracker, tracker)
			return errors.New("Error while getting by id")
		},
	}

	api := &API{st, p}
	var res GetIssueResponse
	err := api.GetIssue(req, &res)
	a.Error(err)
}

func TestCreateReport(t *testing.T) {

	a := assert.New(t)
	var uid ctxtg.UserID = 1
	var token ctxtg.Token = "dsa"
	p := &ctxtgtest.Parser{
		TokenExpected: token,
		Claims: ctxtg.Claims{
			UserID: uid,
		},
		Err: nil,
	}

	rep := entities.Report{Comments: "Comment"}
	req := CreateReportRequest{
		Context: ctxtg.Context{Token: token},
		Tracker: entities.TrackerConfig{
			ID:  1,
			URL: "http://tracker.com",
			Credentials: entities.TrackerCredentials{
				Login:    "login",
				Password: "password",
			},
		},
		Report: rep,
	}

	st := &TestTrackerClient{
		createReport: func(tracker entities.TrackerConfig, report entities.Report) error {
			a.Equal(req.Tracker, tracker)
			a.Equal(rep, report)
			return nil
		},
	}

	api := &API{st, p}
	var res CreateReportResponse
	err := api.CreateReport(req, &res)
	a.NoError(err)
}

func TestCreateReportWithError(t *testing.T) {

	a := assert.New(t)
	var uid ctxtg.UserID = 1
	var token ctxtg.Token = "dsa"
	p := &ctxtgtest.Parser{
		TokenExpected: token,
		Claims: ctxtg.Claims{
			UserID: uid,
		},
		Err: nil,
	}

	rep := entities.Report{Comments: "Comment"}
	req := CreateReportRequest{
		Context: ctxtg.Context{Token: token},
		Tracker: entities.TrackerConfig{
			ID:  1,
			URL: "http://tracker.com",
			Credentials: entities.TrackerCredentials{
				Login:    "login",
				Password: "password",
			},
		},
		Report: rep,
	}

	st := &TestTrackerClient{
		createReport: func(tracker entities.TrackerConfig, report entities.Report) error {
			a.Equal(req.Tracker, tracker)
			a.Equal(rep, report)
			return errors.New("Error in reports creating")
		},
	}

	api := &API{st, p}
	var res CreateReportResponse
	err := api.CreateReport(req, &res)
	a.NotNil(err)
}

func TestGetTotalReports(t *testing.T) {
	a := assert.New(t)
	var uid ctxtg.UserID = 1
	var token ctxtg.Token = "dsa"
	p := &ctxtgtest.Parser{
		TokenExpected: token,
		Claims: ctxtg.Claims{
			UserID: uid,
		},
		Err: nil,
	}

	req := GetTotalReportsRequest{
		Context: ctxtg.Context{Token: token},
		Tracker: entities.TrackerConfig{
			ID:  1,
			URL: "http://tracker.com",
			Credentials: entities.TrackerCredentials{
				Login:    "login",
				Password: "password",
			},
		},
	}

	r := entities.ReportsTotal(1)
	st := &TestTrackerClient{
		getTotalReports: func(tracker entities.TrackerConfig, date entities.Timestamp, res *entities.ReportsTotal) error {
			a.Equal(req.Tracker, tracker)
			*res = r
			return nil
		},
	}

	api := &API{st, p}
	var res GetTotalReportsResponse
	err := api.GetTotalReports(req, &res)
	a.Equal(r, res.Total, "Should be equal")
	a.NoError(err)
}

func TestGetTotalReportsWithError(t *testing.T) {
	a := assert.New(t)
	var uid ctxtg.UserID = 1
	var token ctxtg.Token = "dsa"
	p := &ctxtgtest.Parser{
		TokenExpected: token,
		Claims: ctxtg.Claims{
			UserID: uid,
		},
		Err: nil,
	}

	req := GetTotalReportsRequest{
		Context: ctxtg.Context{Token: token},
		Tracker: entities.TrackerConfig{
			ID:  1,
			URL: "http://tracker.com",
			Credentials: entities.TrackerCredentials{
				Login:    "login",
				Password: "password",
			},
		},
	}

	st := &TestTrackerClient{
		getTotalReports: func(tracker entities.TrackerConfig, date entities.Timestamp, res *entities.ReportsTotal) error {
			a.Equal(req.Tracker, tracker)
			return errors.New("Error while getting total reports")
		},
	}

	api := &API{st, p}
	var res GetTotalReportsResponse
	err := api.GetTotalReports(req, &res)
	a.NotNil(err)
}

func TestGetIssueByURL(t *testing.T) {
	a := assert.New(t)
	var uid ctxtg.UserID = 1
	var token ctxtg.Token = "dsa"
	p := &ctxtgtest.Parser{
		TokenExpected: token,
		Claims: ctxtg.Claims{
			UserID: uid,
		},
		Err: nil,
	}

	req := GetIssueByURLRequest{
		Context: ctxtg.Context{Token: token},
		Tracker: entities.TrackerConfig{
			ID:  1,
			URL: "http://tracker.com",
			Credentials: entities.TrackerCredentials{
				Login:    "login",
				Password: "password",
			},
		},
	}

	i := entities.Issue{URL: "tracker.com"}
	pid := entities.ProjectID(1)
	st := &TestTrackerClient{
		getIssueByUrl: func(tracker entities.TrackerConfig, issueURL string, res *entities.Issue, res2 *entities.ProjectID) error {
			a.Equal(req.Tracker, tracker)
			*res = i
			*res2 = pid
			return nil
		},
	}

	api := &API{st, p}
	var res GetIssueByURLResponse
	err := api.GetIssueByURL(req, &res)
	a.NoError(err)
	a.Equal(i, res.Issue)
	a.Equal(pid, res.ProjectID)
}

func TestGetIssueByURLWithError(t *testing.T) {
	a := assert.New(t)
	var uid ctxtg.UserID = 1
	var token ctxtg.Token = "dsa"
	p := &ctxtgtest.Parser{
		TokenExpected: token,
		Claims: ctxtg.Claims{
			UserID: uid,
		},
		Err: nil,
	}

	req := GetIssueByURLRequest{
		Context: ctxtg.Context{Token: token},
		Tracker: entities.TrackerConfig{
			ID:  1,
			URL: "http://tracker.com",
			Credentials: entities.TrackerCredentials{
				Login:    "login",
				Password: "password",
			},
		},
	}

	st := &TestTrackerClient{
		getIssueByUrl: func(tracker entities.TrackerConfig, issueURL string, res *entities.Issue, res2 *entities.ProjectID) error {
			a.Equal(req.Tracker, tracker)
			return errors.New("Error while getting by URL")
		},
	}

	api := &API{st, p}
	var res GetIssueByURLResponse
	err := api.GetIssueByURL(req, &res)
	a.NotNil(err)
}
