package rpcsvc

import (
	"context"
	"reflect"
	"testing"

	"../../entities"
	"github.com/pkg/errors"
	"github.com/powerman/narada-go/narada"
	"gitlab.qarea.org/tgms/ctxtg"
	"gitlab.qarea.org/tgms/ctxtg/ctxtgtest"
)

func TestVERSION(t *testing.T) {
	api := &API{}
	var args struct{}
	var res string
	err := api.Version(&args, &res)
	if err != nil {
		t.Errorf("Version(), err = %v", err)
	}
	if ver, _ := narada.Version(); res != ver {
		t.Errorf("Version() = %v, want %v", res, ver)
	}
}

func TestProjectDetails(t *testing.T) {
	type test struct {
		project   entities.Project
		projectID entities.ProjectID
		tracker   entities.Tracker
		token     ctxtg.Token
		tokenErr  error
		err       error
	}
	tests := map[string]test{
		"Return project": {
			token:     "123412",
			projectID: "55",
			tracker:   testTracker,
			project: entities.Project{
				ID:          "1",
				Title:       "t",
				Description: "d",
			},
		},
		"Token parse error": {
			token:    "invalid token",
			tokenErr: ctxtg.ErrInvalidToken,
		},
		"Return error": {
			err: errors.New("hi"),
		},
	}

	for label, test := range tests {
		rc := TestTrackerClient{
			project: func(ctx context.Context, tr entities.Tracker, pid entities.ProjectID) (*entities.Project, error) {
				if test.tokenErr != nil {
					t.Error("Should not be called", label)
				}
				checkTracker(t, label, tr)
				checkCtx(t, label, ctx)
				return &test.project, test.err
			},
		}
		p := &ctxtgtest.Parser{
			Err:           test.tokenErr,
			TokenExpected: test.token,
		}

		r := newAPI(rc, p)
		var resp ProjectDetailsResp
		err := r.GetProjectDetails(&ProjectDetailsReq{
			Context: testContext(test.token),
			Tracker: testTracker,
		}, &resp)

		if err := p.Error(); err != nil {
			t.Errorf("Parser error in test %v: %v", label, err)
		}
		if (test.err != nil || test.tokenErr != nil) && err == nil {
			t.Errorf("Test %s should return err", label)
		}
		if !reflect.DeepEqual(test.project, resp.Project) {
			t.Errorf("Test %s unexpected projects resp", label)
		}
	}
}

func TestGetProjects(t *testing.T) {
	type test struct {
		pagination entities.Pagination
		projects   []entities.Project
		amount     int64
		token      ctxtg.Token
		tokenErr   error
		err        error
	}
	tests := map[string]test{
		"Return projects": {
			token: "123412",
			pagination: entities.Pagination{
				Limit:  30,
				Offset: 0,
			},
			amount: 20,
			projects: []entities.Project{
				{
					ID:          "1",
					Title:       "t",
					Description: "d",
				},
				{
					ID:          "1",
					Title:       "t2",
					Description: "d2",
				},
			},
		},
		"Token parse error": {
			token:    "invalid token",
			tokenErr: ctxtg.ErrInvalidToken,
		},
		"Return error": {
			projects: []entities.Project{},
			pagination: entities.Pagination{
				Limit:  30,
				Offset: 0,
			},
			err: errors.New("hi"),
		},
		"Empty project list": {
			projects: []entities.Project{},
		},
	}

	for label, test := range tests {
		rc := TestTrackerClient{
			projects: func(ctx context.Context, tr entities.Tracker, p entities.Pagination) ([]entities.Project, int64, error) {
				if test.tokenErr != nil {
					t.Error("Should not be called", label)
				}
				checkCtx(t, label, ctx)
				checkTracker(t, label, tr)
				if p != test.pagination {
					t.Errorf("Test %s invalid pagination", label)
				}
				return test.projects, test.amount, test.err
			},
		}
		p := &ctxtgtest.Parser{
			Err:           test.tokenErr,
			TokenExpected: test.token,
		}

		r := newAPI(rc, p)
		var resp ProjectsResp
		err := r.GetProjects(&ProjectsReq{
			Context:    testContext(test.token),
			Tracker:    testTracker,
			Pagination: test.pagination,
		}, &resp)

		if err := p.Error(); err != nil {
			t.Errorf("Parser error in test %v: %v", label, err)
		}
		if (test.err != nil || test.tokenErr != nil) && err == nil {
			t.Errorf("Test %s should return err", label)
		}
		if !reflect.DeepEqual(test.projects, resp.Projects) {
			t.Errorf("Test %s unexpected projects resp label", label)
		}
		if !reflect.DeepEqual(test.amount, resp.Amount) {
			t.Errorf("Test %s unexpected amount resp label %v != %v", label, resp.Amount, test.amount)
		}
	}
}

func TestGetCurrentUser(t *testing.T) {
	type test struct {
		user     *entities.User
		err      error
		token    ctxtg.Token
		tokenErr error
	}
	tests := map[string]test{
		"Return user": {
			user: &entities.User{
				ID:   "1",
				Name: "tolya",
			},
		},
		"Token parse error": {
			token:    "invalid token",
			tokenErr: ctxtg.ErrInvalidToken,
		},
		"Error": {
			err: errors.New("user err"),
		},
	}

	for label, test := range tests {
		rc := TestTrackerClient{
			userInfo: func(ctx context.Context, tr entities.Tracker) (*entities.User, error) {
				if test.tokenErr != nil {
					t.Errorf("Should not be called %v", label)
				}
				checkCtx(t, label, ctx)
				checkTracker(t, label, tr)
				return test.user, test.err
			},
		}

		p := &ctxtgtest.Parser{
			Err:           test.tokenErr,
			TokenExpected: test.token,
		}

		r := newAPI(rc, p)

		var resp CurrentUserResp
		err := r.GetCurrentUser(&CurrentUserReq{
			Context: testContext(test.token),
			Tracker: testTracker,
		}, &resp)

		if err := p.Error(); err != nil {
			t.Errorf("Parser error in test %v: %v", label, err)
		}
		if (test.err != nil || test.tokenErr != nil) && err == nil {
			t.Errorf("Test %s should return err", label)
		}
		if test.user != nil && *test.user != resp.User {
			t.Errorf("Test %s invalid user returned", label)
		}
	}
}

func TestGetProjectIssues(t *testing.T) {
	type test struct {
		issues     []entities.Issue
		pagination entities.Pagination
		amount     int64
		err        error
		token      ctxtg.Token
		tokenErr   error
	}
	tests := map[string]test{
		"Issues": {
			issues: []entities.Issue{
				{
					ID:    "1",
					Title: "t",
				},
				{
					ID:    "2",
					Title: "t2",
				},
			},
			amount: 2,
		},
		"Token parse error": {
			token:    "invalid token",
			tokenErr: ctxtg.ErrInvalidToken,
		},
		"Error": {
			err: errors.New("issues err"),
		},
	}

	for label, test := range tests {
		rc := TestTrackerClient{
			projectIssues: func(ctx context.Context, tr entities.Tracker, id entities.ProjectID, p entities.Pagination) ([]entities.Issue, int64, error) {
				if test.tokenErr != nil {
					t.Errorf("Should not be called %v", label)
				}
				checkCtx(t, label, ctx)
				checkTracker(t, label, tr)
				if p != test.pagination {
					t.Errorf("Test %s invalid pagination", label)
				}
				return test.issues, test.amount, test.err
			},
		}
		p := &ctxtgtest.Parser{
			Err:           test.tokenErr,
			TokenExpected: test.token,
		}

		r := newAPI(rc, p)

		var resp ProjectIssuesResp
		err := r.GetProjectIssues(&ProjectIssuesReq{
			Context: testContext(test.token),
			Tracker: testTracker,
		}, &resp)

		if err := p.Error(); err != nil {
			t.Errorf("Parser error in test %v: %v", label, err)
		}
		if (test.err != nil || test.tokenErr != nil) && err == nil {
			t.Errorf("Test %s should return err", label)
		}
		if !reflect.DeepEqual(test.issues, resp.Issues) {
			t.Errorf("Test %s unexpected issues resp label", label)
		}
		if !reflect.DeepEqual(test.amount, resp.Amount) {
			t.Errorf("Test %s unexpected amount resp label %v != %v", label, resp.Amount, test.amount)
		}
	}
}

func TestCreateIssue(t *testing.T) {
	type test struct {
		issue         entities.NewIssue
		issueToReturn *entities.Issue
		projectID     entities.ProjectID
		err           error
		token         ctxtg.Token
		tokenErr      error
	}
	tests := map[string]test{
		"Create issue": {
			issue: entities.NewIssue{
				Issue: entities.Issue{
					Title: "issue1",
				},
			},
			issueToReturn: &entities.Issue{
				ID:    "1",
				Title: "issue1",
			},
			projectID: "1",
		},
		"Token parse error": {
			token:    "invalid token",
			tokenErr: ctxtg.ErrInvalidToken,
		},
		"Error response": {
			issue: entities.NewIssue{
				Issue: entities.Issue{
					Title: "issue1",
				},
			},
			err: entities.ErrForbidden,
		},
	}

	for label, test := range tests {
		rc := TestTrackerClient{
			createIssue: func(ctx context.Context, tr entities.Tracker, is entities.NewIssue, id entities.ProjectID) (*entities.Issue, error) {
				if test.tokenErr != nil {
					t.Errorf("Should not be called %v", label)
				}
				checkCtx(t, label, ctx)
				checkTracker(t, label, tr)
				if is != test.issue {
					t.Errorf("Test %s invalid issue passed", label)
				}
				if id != test.projectID {
					t.Errorf("Test %s invalid projectID passed", label)
				}
				return test.issueToReturn, test.err
			},
		}
		p := &ctxtgtest.Parser{
			Err:           test.tokenErr,
			TokenExpected: test.token,
		}

		r := newAPI(rc, p)

		var resp CreateIssueResp
		err := r.CreateIssue(&CreateIssueReq{
			Context:   testContext(test.token),
			Tracker:   testTracker,
			Issue:     test.issue,
			ProjectID: test.projectID,
		}, &resp)

		if err := p.Error(); err != nil {
			t.Errorf("Parser error in test %v: %v", label, err)
		}
		if (test.err != nil || test.tokenErr != nil) && err == nil {
			t.Errorf("Test %s should return err", label)
		}
		if test.issueToReturn != nil && *test.issueToReturn != resp.Issue {
			t.Errorf("Test %s invalid issue returned", label)
		}
	}
}

func TestGetIssue(t *testing.T) {
	type test struct {
		issue     *entities.Issue
		issueID   entities.IssueID
		projectID entities.ProjectID
		err       error
		token     ctxtg.Token
		tokenErr  error
	}
	tests := map[string]test{
		"Token parse error": {
			token:    "invalid token",
			tokenErr: ctxtg.ErrInvalidToken,
		},
		"Create issue": {
			issue: &entities.Issue{
				Title: "issue1",
			},
			issueID: "1",
		},
		"Error response": {
			issueID: "2",
			err:     entities.ErrForbidden,
		},
	}

	for label, test := range tests {
		rc := TestTrackerClient{
			issue: func(ctx context.Context, tr entities.Tracker, pid entities.ProjectID, id entities.IssueID) (*entities.Issue, error) {
				if test.tokenErr != nil {
					t.Error("Should not be called", label)
				}
				if test.projectID != pid {
					t.Error("Invalid project ID")
				}
				checkCtx(t, label, ctx)
				checkTracker(t, label, tr)
				if id != test.issueID {
					t.Errorf("Test %s invalid issueID passed", label)
				}
				return test.issue, test.err
			},
		}
		p := &ctxtgtest.Parser{
			Err:           test.tokenErr,
			TokenExpected: test.token,
		}

		r := newAPI(rc, p)

		var resp GetIssueResp
		err := r.GetIssue(&GetIssueReq{
			Context:   testContext(test.token),
			Tracker:   testTracker,
			IssueID:   test.issueID,
			ProjectID: test.projectID,
		}, &resp)

		if err := p.Error(); err != nil {
			t.Errorf("Parser error in test %v: %v", label, err)
		}
		if (test.err != nil || test.tokenErr != nil) && err == nil {
			t.Errorf("Test %s should return err", label)
		}
		if test.issue != nil && *test.issue != resp.Issue {
			t.Errorf("Test %s invalid issue returned", label)
		}
	}
}

func TestGetIssueByURL(t *testing.T) {
	type test struct {
		issue    *entities.Issue
		issueURL entities.IssueURL
		err      error
		token    ctxtg.Token
		tokenErr error
	}
	tests := map[string]test{
		"Ok issue": {
			issue: &entities.Issue{
				Title:     "issue1",
				ProjectID: "2",
			},
			issueURL: "/issue/2",
		},
		"Token parse error": {
			token:    "invalid token",
			tokenErr: ctxtg.ErrInvalidToken,
		},
		"Error response": {
			issueURL: "/issue/1",
			err:      entities.ErrForbidden,
		},
	}

	for label, test := range tests {
		rc := TestTrackerClient{
			issueByURL: func(ctx context.Context, tr entities.Tracker, url entities.IssueURL) (*entities.Issue, error) {
				if test.tokenErr != nil {
					t.Error("Should not be called", label)
				}
				checkCtx(t, label, ctx)
				checkTracker(t, label, tr)
				if url != test.issueURL {
					t.Errorf("Test %s invalid issue URL passed", label)
				}
				return test.issue, test.err
			},
		}

		p := &ctxtgtest.Parser{
			Err:           test.tokenErr,
			TokenExpected: test.token,
		}

		r := newAPI(rc, p)

		var resp GetIssueByURLResp
		err := r.GetIssueByURL(&GetIssueByURLReq{
			Context:  testContext(test.token),
			Tracker:  testTracker,
			IssueURL: test.issueURL,
		}, &resp)

		if err := p.Error(); err != nil {
			t.Errorf("Parser error in test %v: %v", label, err)
		}
		if (test.err != nil || test.tokenErr != nil) && err == nil {
			t.Errorf("Test %s should return err", label)
		}
		if test.issue != nil && *test.issue != resp.Issue && test.issue.ProjectID != resp.ProjectID {
			t.Errorf("Test %s invalid issue returned", label)
		}
	}
}

func TestUpdateIssue(t *testing.T) {
	type test struct {
		progress      entities.Progress
		issueToReturn *entities.Issue
		projectID     entities.ProjectID
		err           error
		token         ctxtg.Token
		tokenErr      error
	}
	tests := map[string]test{
		"Update issue": {
			progress: 50,
			issueToReturn: &entities.Issue{
				ID:    "1",
				Title: "issue1",
			},
			projectID: "1",
		},
		"Token parse error": {
			token:    "invalid token",
			tokenErr: ctxtg.ErrInvalidToken,
		},
		"Error response": {
			progress: 20,
			err:      entities.ErrForbidden,
		},
	}

	for label, test := range tests {
		rc := TestTrackerClient{
			updateIssueProgress: func(ctx context.Context, tr entities.Tracker, pid entities.ProjectID, iid entities.IssueID, pr entities.Progress) error {
				if test.tokenErr != nil {
					t.Error("Should not be called", label)
				}
				if pid != test.projectID {
					t.Error("Invalid projectID")
				}
				checkCtx(t, label, ctx)
				checkTracker(t, label, tr)
				if pr != test.progress {
					t.Errorf("Test %s invalid issue passed", label)
				}
				return test.err
			},
		}
		p := &ctxtgtest.Parser{
			Err:           test.tokenErr,
			TokenExpected: test.token,
		}

		r := newAPI(rc, p)

		err := r.UpdateIssueProgress(&UpdateIssueProgressReq{
			Context:   testContext(test.token),
			Tracker:   testTracker,
			Progress:  test.progress,
			ProjectID: test.projectID,
		}, nil)

		if err := p.Error(); err != nil {
			t.Errorf("Parser error in test %v: %v", label, err)
		}
		if (test.err != nil || test.tokenErr != nil) && err == nil {
			t.Errorf("Test %s should return err", label)
		}
	}
}

func TestCreateReport(t *testing.T) {
	type test struct {
		report    entities.Report
		projectID entities.ProjectID
		err       error
		token     ctxtg.Token
		tokenErr  error
	}
	tests := map[string]test{
		"Create report": {
			projectID: "2",
			report: entities.Report{
				IssueID:    "1",
				ActivityID: 3,
				Comments:   "comme",
				Duration:   4,
				Started:    5,
			},
		},
		"Token parse error": {
			token:    "invalid token",
			tokenErr: ctxtg.ErrInvalidToken,
		},
		"Error response": {
			err: entities.ErrForbidden,
		},
	}

	for label, test := range tests {
		rc := TestTrackerClient{
			createReport: func(ctx context.Context, tr entities.Tracker, pid entities.ProjectID, rep entities.Report) error {
				if test.tokenErr != nil {
					t.Error("Should not be called", label)
				}
				if test.projectID != pid {
					t.Error("Invalid projectID")
				}
				checkCtx(t, label, ctx)
				checkTracker(t, label, tr)
				if rep != test.report {
					t.Errorf("Test %s invalid report passed", label)
				}
				return test.err
			},
		}
		p := &ctxtgtest.Parser{
			Err:           test.tokenErr,
			TokenExpected: test.token,
		}

		r := newAPI(rc, p)

		err := r.CreateReport(&CreateReportReq{
			Context:   testContext(test.token),
			Tracker:   testTracker,
			ProjectID: test.projectID,
			Report:    test.report,
		}, &struct{}{})

		if err := p.Error(); err != nil {
			t.Errorf("Parser error in test %v: %v", label, err)
		}
		if (test.err != nil || test.tokenErr != nil) && err == nil {
			t.Errorf("Test %s should return err", label)
		}
	}
}

func TestGetTotalReports(t *testing.T) {
	type test struct {
		date     int64
		result   int64
		err      error
		token    ctxtg.Token
		tokenErr error
	}
	tests := map[string]test{
		"Successful": {
			date:   1,
			result: 2,
		},
		"Token parse error": {
			token:    "invalid token",
			tokenErr: ctxtg.ErrInvalidToken,
		},
		"Error": {
			err: entities.ErrRemoteServer,
		},
	}

	for label, test := range tests {
		rc := TestTrackerClient{
			totalReports: func(ctx context.Context, tr entities.Tracker, date int64) (int64, error) {
				if test.tokenErr != nil {
					t.Error("Should not be called", label)
				}
				checkCtx(t, label, ctx)
				checkTracker(t, label, tr)
				if date != test.date {
					t.Errorf("Test %s invalid date passed", label)
				}
				return test.result, test.err
			},
		}
		p := &ctxtgtest.Parser{
			Err:           test.tokenErr,
			TokenExpected: test.token,
		}

		r := newAPI(rc, p)
		var resp GetReportsResp
		err := r.GetTotalReports(&GetReportsReq{
			Context: testContext(test.token),
			Tracker: testTracker,
			Date:    test.date,
		}, &resp)

		if err := p.Error(); err != nil {
			t.Errorf("Parser error in test %v: %v", label, err)
		}
		if (test.err != nil || test.tokenErr != nil) && err == nil {
			t.Errorf("Test %s should return err", label)
		}
		if test.result != resp.Total {
			t.Errorf("Test %s invalid response", label)
		}
	}
}

func checkCtx(t *testing.T, label string, ctx context.Context) {
	if ctx == nil {
		t.Errorf("Test %s passed nil context", label)
	}
}

func checkTracker(t *testing.T, label string, tr entities.Tracker) {
	if tr != testTracker {
		t.Errorf("Test %s Invalid tracker passed", label)
	}
}

var testTracker = entities.Tracker{
	ID:   30,
	URL:  "testurl",
	Type: "RED",
	Credentials: entities.Credentials{
		Login:    "str",
		Password: "pass",
	},
}

func testContext(token ctxtg.Token) ctxtg.Context {
	return ctxtg.Context{
		Token: token,
	}
}

type TestTrackerClient struct {
	projects            func(context.Context, entities.Tracker, entities.Pagination) ([]entities.Project, int64, error)
	project             func(context.Context, entities.Tracker, entities.ProjectID) (*entities.Project, error)
	projectIssues       func(context.Context, entities.Tracker, entities.ProjectID, entities.Pagination) ([]entities.Issue, int64, error)
	userInfo            func(context.Context, entities.Tracker) (*entities.User, error)
	issue               func(context.Context, entities.Tracker, entities.ProjectID, entities.IssueID) (*entities.Issue, error)
	issueByURL          func(context.Context, entities.Tracker, entities.IssueURL) (*entities.Issue, error)
	createIssue         func(context.Context, entities.Tracker, entities.NewIssue, entities.ProjectID) (*entities.Issue, error)
	updateIssueProgress func(context.Context, entities.Tracker, entities.ProjectID, entities.IssueID, entities.Progress) error
	totalReports        func(ctx context.Context, t entities.Tracker, date int64) (int64, error)
	createReport        func(context.Context, entities.Tracker, entities.ProjectID, entities.Report) error
}

func (r TestTrackerClient) Projects(ctx context.Context, t entities.Tracker, p entities.Pagination) ([]entities.Project, int64, error) {
	return r.projects(ctx, t, p)
}

func (r TestTrackerClient) Project(ctx context.Context, t entities.Tracker, pid entities.ProjectID) (*entities.Project, error) {
	return r.project(ctx, t, pid)
}

func (r TestTrackerClient) ProjectIssues(ctx context.Context, t entities.Tracker, id entities.ProjectID, p entities.Pagination) ([]entities.Issue, int64, error) {
	return r.projectIssues(ctx, t, id, p)
}

func (r TestTrackerClient) UserInfo(ctx context.Context, t entities.Tracker) (*entities.User, error) {
	return r.userInfo(ctx, t)
}

func (r TestTrackerClient) Issue(ctx context.Context, t entities.Tracker, pid entities.ProjectID, id entities.IssueID) (*entities.Issue, error) {
	return r.issue(ctx, t, pid, id)
}

func (r TestTrackerClient) IssueByURL(ctx context.Context, t entities.Tracker, url entities.IssueURL) (*entities.Issue, error) {
	return r.issueByURL(ctx, t, url)
}

func (r TestTrackerClient) CreateIssue(ctx context.Context, t entities.Tracker, i entities.NewIssue, projectID entities.ProjectID) (*entities.Issue, error) {
	return r.createIssue(ctx, t, i, projectID)
}

func (r TestTrackerClient) UpdateIssueProgress(ctx context.Context, t entities.Tracker, pid entities.ProjectID, iid entities.IssueID, pr entities.Progress) error {
	return r.updateIssueProgress(ctx, t, pid, iid, pr)
}

func (r TestTrackerClient) TotalReports(ctx context.Context, t entities.Tracker, date int64) (int64, error) {
	return r.totalReports(ctx, t, date)
}

func (r TestTrackerClient) CreateReport(ctx context.Context, t entities.Tracker, pid entities.ProjectID, rep entities.Report) error {
	return r.createReport(ctx, t, pid, rep)
}
