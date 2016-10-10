package jira

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/qarea/jirams/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testRequest struct {
	method  string
	url     string
	body    io.Reader
	headers map[string]string
	error   error
	result  interface{}
	data    interface{}
}

var (
	testTracker = entities.TrackerConfig{
		URL: "https://tracker.com",
		Credentials: entities.TrackerCredentials{
			Login:    "tester",
			Password: "test",
		},
	}
	testEstimate = entities.Duration(3600)
	testDueDate  = entities.Timestamp(1482624000)
	testIssue    = entities.Issue{
		ID:       10000,
		Type:     entities.NamedID{ID: 10001, Name: "Task"},
		URL:      "https://tracker.com/browse/10000",
		Title:    "Test Issue",
		Estimate: testEstimate,
		DueDate:  testDueDate,
		Spent:    1800,
		Done:     50,
	}
	testJiraIssue = Issue{
		ID:  "10000",
		URL: "https://tracker.com/browse/10000",
		Fields: IssueFields{
			Type:     NamedID{ID: "10001", Name: "Task"},
			Title:    "Test Issue",
			Estimate: 3600,
			DueDate:  "2016-12-25",
			Spent:    1800,
			Progress: IssueProgress{
				Percent: 50,
			},
			ProjectID: EntityID{ID: "10000"},
		},
	}
)

type MockJiraRequester struct {
	mock.Mock
}

func (t *MockJiraRequester) Request(tracker entities.TrackerConfig, request *http.Request, res interface{}) error {
	args := t.Called(tracker, request)
	if r, ok := res.(*[]Project); ok {
		*r = args.Get(0).([]Project)
	}
	if r, ok := res.(*Project); ok {
		*r = args.Get(0).(Project)
	}
	if r, ok := res.(*User); ok {
		*r = args.Get(0).(User)
	}
	if r, ok := res.(*Issues); ok {
		*r = args.Get(0).(Issues)
	}
	if r, ok := res.(*Issue); ok {
		*r = args.Get(0).(Issue)
	}
	if r, ok := res.(*EntityID); ok {
		*r = args.Get(0).(EntityID)
	}
	if r, ok := res.(*WorklogPage); ok {
		*r = args.Get(0).(WorklogPage)
	}
	return args.Error(1)
}

func (t *MockJiraRequester) IterateRequest(tracker entities.TrackerConfig, url string, dataContainer interface{}, callback func(interface{}) (int, int, error)) error {
	args := t.Called(tracker, url)
	if _, ok := dataContainer.(*IssueIDPage); ok {
		data := args.Get(0).(IssueIDPage)
		_, _, _ = callback(&data)
	}

	if _, ok := dataContainer.(*WorklogPage); ok {
		data := args.Get(0).(WorklogPage)
		_, _, _ = callback(&data)
	}
	return args.Error(1)
}

type MockStore struct {
	mock.Mock
}

func (t *MockStore) Init() {}
func (t *MockStore) GetID(trackerID entities.TrackerID, key entities.UserKey) (res entities.UserID, err error) {
	return entities.UserID(1), nil
}
func (t *MockStore) GetKey(trackerID entities.TrackerID, userID entities.UserID) (res string, err error) {
	return "KEY", nil
}

func TestGetProjects(t *testing.T) {
	var (
		requests = map[string]testRequest{
			"OK": {
				method: "GET",
				url:    "https://tracker.com/rest/api/2/project",
				data: []Project{
					{
						ID:    "10000",
						Title: "Test Project",
					},
				},
			},
			"OK Subrequest": {
				method: "GET",
				url:    "https://tracker.com/rest/api/2/project/10000",
				data: Project{
					ID:          "10000",
					Title:       "Test Project",
					Description: "Test Project Description",
					IssueTypes: []NamedID{
						{ID: "10000", Name: "Task"},
						{ID: "10001", Name: "Sub-task"},
					},
				},
			},
		}
	)
	testRequester := new(MockJiraRequester)
	for _, r := range requests {
		request, _ := http.NewRequest(r.method, r.url, r.body)
		testRequester.On("Request", testTracker, request).
			Return(r.data, r.error)
	}

	client := Client{&MockStore{}, testRequester}

	// TEST VALID REQUEST
	var (
		expected = []entities.Project{
			{
				ID:          10000,
				Title:       "Test Project",
				Description: "Test Project Description",
				IssueTypes: []entities.NamedID{
					{ID: 10000, Name: "Task"},
					{ID: 10001, Name: "Sub-task"},
				},
				ActivityTypes: make([]entities.NamedID, 0),
			},
		}
		result []entities.Project
	)

	err := client.GetProjects(testTracker, &result)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	testRequester.AssertExpectations(t)
}

func TestGetProjectsEmpty(t *testing.T) {
	var (
		emptyProjects []Project
		requests      = map[string]testRequest{
			"Empty": {
				method: "GET",
				url:    "https://tracker.com/rest/api/2/project",
				data:   emptyProjects,
			},
		}
	)

	testRequester := new(MockJiraRequester)
	for _, r := range requests {
		request, _ := http.NewRequest(r.method, r.url, r.body)
		testRequester.On("Request", testTracker, request).
			Return(r.data, r.error)
	}

	client := Client{&MockStore{}, testRequester}

	// TEST EMPTY PROJECTS LIST
	var (
		expected = make([]entities.Project, 0)
		result   []entities.Project
	)
	err := client.GetProjects(testTracker, &result)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	testRequester.AssertExpectations(t)
}

func TestGetProjectsInvalidTracker(t *testing.T) {
	type test struct {
		Cfg entities.TrackerConfig
		Err error
	}
	var (
		badTrackers = map[string]test{
			"Empty": {
				Cfg: entities.TrackerConfig{},
				Err: entities.ErrInvalidTrackerURL,
			},
			"No URL": {
				Cfg: entities.TrackerConfig{
					Credentials: entities.TrackerCredentials{
						Login:    "test",
						Password: "test",
					},
				},
				Err: entities.ErrInvalidTrackerURL,
			},
			"No login": {
				Cfg: entities.TrackerConfig{
					URL: "http://tracker.com",
					Credentials: entities.TrackerCredentials{
						Password: "test",
					},
				},
				Err: entities.ErrInvalidCredentials,
			},
			"No password": {
				Cfg: entities.TrackerConfig{
					URL: "http://tracker.com",
					Credentials: entities.TrackerCredentials{
						Login: "test",
					},
				},
				Err: entities.ErrInvalidCredentials,
			},
		}
	)
	testRequester := new(MockJiraRequester)
	client := Client{&MockStore{}, testRequester}

	var (
		result []entities.Project
		err    error
	)
	for _, test := range badTrackers {
		err = client.GetProjects(test.Cfg, &result)
		if assert.Error(t, err, "Error expected") {
			assert.Equal(t, test.Err, err)
		}
	}
}

func TestGetCurrentUser(t *testing.T) {
	requests := map[string]testRequest{
		"OK": {
			method: "GET",
			url:    "https://tracker.com/rest/api/2/myself",
			data: User{
				Key:  "user",
				Name: "John Smith",
				Mail: "john@smith.com",
			},
		},
	}

	testRequester := new(MockJiraRequester)
	for _, r := range requests {
		request, _ := http.NewRequest(r.method, r.url, r.body)
		testRequester.On("Request", testTracker, request).
			Return(r.data, r.error)
	}

	client := Client{&MockStore{}, testRequester}

	var (
		expected = entities.User{
			ID:   1,
			Name: "John Smith",
			Mail: "john@smith.com",
		}
		result entities.User
	)
	err := client.GetCurrentUser(testTracker, &result)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	testRequester.AssertExpectations(t)
}

func TestGetProjectIssues(t *testing.T) {
	requests := map[string]testRequest{
		"OK": {
			method: "GET",
			url:    "https://tracker.com/rest/api/2/search?jql=project%3D1+AND+assignee%3DcurrentUser%28%29+AND+resolution%3DUnresolved&startAt=0",
			data: Issues{
				StartAt:    0,
				MaxResults: 50,
				Total:      1,
				Issues:     []Issue{testJiraIssue},
			},
		},
	}

	testRequester := new(MockJiraRequester)
	for _, r := range requests {
		request, _ := http.NewRequest(r.method, r.url, r.body)
		testRequester.On("Request", testTracker, request).
			Return(r.data, r.error)
	}

	client := Client{&MockStore{}, testRequester}

	var (
		expected = []entities.Issue{testIssue}
		result   []entities.Issue
	)

	err := client.GetProjectIssues(testTracker, entities.ProjectID(1), entities.UserID(2), &result)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	testRequester.AssertExpectations(t)
}

func TestGetIssue(t *testing.T) {
	requests := map[string]testRequest{
		"OK": {
			method: "GET",
			url:    "https://tracker.com/rest/api/2/issue/10000",
			data:   testJiraIssue,
		},
	}

	testRequester := new(MockJiraRequester)
	for _, r := range requests {
		request, _ := http.NewRequest(r.method, r.url, r.body)
		testRequester.On("Request", testTracker, request).
			Return(r.data, r.error)
	}

	client := Client{&MockStore{}, testRequester}

	var (
		expected = testIssue
		result   entities.Issue
	)

	err := client.GetIssue(testTracker, entities.IssueID(10000), &result)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	testRequester.AssertExpectations(t)
}

func TestGetIssueByURL(t *testing.T) {
	requests := map[string]testRequest{
		"OK": {
			method: "GET",
			url:    "https://tracker.com/rest/api/2/issue/10000",
			data:   testJiraIssue,
		},
	}

	testRequester := new(MockJiraRequester)
	for _, r := range requests {
		request, _ := http.NewRequest(r.method, r.url, r.body)
		testRequester.On("Request", testTracker, request).
			Return(r.data, r.error)
	}

	client := Client{&MockStore{}, testRequester}

	var (
		expected  = testIssue
		expected2 = entities.ProjectID(10000)
		result    entities.Issue
		result2   entities.ProjectID
	)

	err := client.GetIssueByURL(testTracker, "https://tracker.com/browse/10000", &result, &result2)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	assert.Equal(t, expected2, result2)
	testRequester.AssertExpectations(t)
}

func TestCreateIssue(t *testing.T) {
	var (
		testIssue = NewIssue{
			Fields: NewIssueFields{
				Title:    "Test Issue",
				Project:  ID{ID: "10000"},
				Type:     ID{ID: "10001"},
				Assignee: UserKey{Key: "KEY"},
				TimeTracking: TimeTracking{
					OriginalEstimate:  60,
					RemainingEstimate: 60,
				},
			},
		}
	)

	payloadBytes, _ := json.Marshal(testIssue)
	var requests = map[string]testRequest{
		"OK": {
			method: "POST",
			url:    "https://tracker.com/rest/api/2/issue",
			body:   bytes.NewBuffer(payloadBytes),
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			data: EntityID{ID: "1"},
		},
		"Sub-request": {
			method: "GET",
			url:    "https://tracker.com/rest/api/2/issue/1",
			data: Issue{
				ID: "1",
				Fields: IssueFields{
					Type:     NamedID{ID: "10001", Name: "Task"},
					Title:    "Test Issue",
					Estimate: 3600,
				},
			},
		},
	}

	testRequester := new(MockJiraRequester)
	for _, r := range requests {
		request, _ := http.NewRequest(r.method, r.url, r.body)
		for header, value := range r.headers {
			request.Header.Set(header, value)
		}
		testRequester.On("Request", testTracker, request).
			Return(r.data, r.error)
	}

	client := Client{&MockStore{}, testRequester}

	var (
		estimate = entities.Duration(3600)
		expected = entities.Issue{
			ID:       1,
			Type:     entities.NamedID{ID: 10001, Name: "Task"},
			Title:    "Test Issue",
			Estimate: estimate,
		}
		result   entities.Issue
		newIssue = entities.NewIssue{
			ProjectID: 10000,
			Assignee:  1,
			Type:      entities.EntityID(10001),
			Title:     "Test Issue",
			Estimate:  3600,
		}
	)
	err := client.CreateIssue(testTracker, newIssue, &result)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	testRequester.AssertExpectations(t)
}

func TestCreateReport(t *testing.T) {
	reportTime := entities.Timestamp(1482667200)
	report := Worklog{
		Started: time.Unix(int64(reportTime), 0).Format(jiraTimestampLayout),
		Spent:   3600,
		Comment: "Test Report",
	}
	payloadBytes, _ := json.Marshal(report)
	requests := map[string]testRequest{
		"OK": {
			method: "POST",
			url:    "https://tracker.com/rest/api/2/issue/1/worklog",
			body:   bytes.NewBuffer(payloadBytes),
			headers: map[string]string{
				"Content-Type": "application/json",
			},
		},
	}

	testRequester := new(MockJiraRequester)
	for _, r := range requests {
		request, _ := http.NewRequest(r.method, r.url, r.body)
		for header, value := range r.headers {
			request.Header.Set(header, value)
		}
		testRequester.On("Request", testTracker, request).
			Return(r.data, r.error)
	}

	client := Client{&MockStore{}, testRequester}

	err := client.CreateReport(testTracker, entities.Report{
		IssueID:  1,
		Started:  reportTime,
		Duration: 3600,
		Comments: "Test Report",
	})

	assert.Nil(t, err)
	testRequester.AssertExpectations(t)
}

func TestGetTotalReports(t *testing.T) {
	var (
		issues = IssueIDPage{
			StartAt:    0,
			MaxResults: 50,
			Total:      1,
			IssueIDs: []ID{
				{ID: "10000"},
			},
		}
		worklogs = WorklogPage{
			StartAt:    0,
			MaxResults: 50,
			Total:      1,
			Worklogs: []Worklog{
				{
					Started: "2016-12-25T14:00:00.000+0200",
					Spent:   3600,
					Comment: "Test Report",
				},
			},
		}
	)

	testRequester := new(MockJiraRequester)
	testRequester.On("IterateRequest", testTracker, "https://tracker.com/rest/api/2/search?jql=worklogAuthor%3DcurrentUser%28%29+AND+worklogDate%3D%222016%2F12%2F25%22&fields=id").
		Return(issues, nil)
	testRequester.On("IterateRequest", testTracker, "https://tracker.com/rest/api/2/issue/10000/worklog").
		Return(worklogs, nil)

	client := Client{&MockStore{}, testRequester}

	var (
		expected = entities.ReportsTotal(3600)
		result   entities.ReportsTotal
	)
	err := client.GetTotalReports(testTracker, 1482624000, &result)

	assert.Nil(t, err)
	assert.Equal(t, expected, result)
	testRequester.AssertExpectations(t)
}
