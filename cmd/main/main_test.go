//+build integration

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	. "gopkg.in/ory-am/dockertest.v2"

	"github.com/boltdb/bolt"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/qarea/ctxtg"
	"github.com/qarea/jirams/api"
	"github.com/qarea/jirams/cfg"
	"github.com/qarea/jirams/entities"
	"github.com/stretchr/testify/assert"
)

var (
	db        *bolt.DB
	listener  net.Listener
	jwt       = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1Nzc5MzQyNDUsInN1YiI6IjEifQ.Cg0CbkaV9YpX9kig1zP0c9_vwCOAVXACgQvkD1lSxr8znMkG8cfhpfDIVcxXc5zRtsdO-SoyngV7Y1zRmBoDmpz4H24QdiQKCLkYefKArg5SV67KKGoU2_e2wjDoIfFotS43wYLrRYyAr9Dgx22wrBYLvmwh0XcYRvwJesPPd4Q3fWN-gO5Iz2P8ytj3ds5x_04mT7pCj4GUAPZ7oqszfI75kv2PUaQJzxv24ptZ9Rdr0hFiEqLwCXNLH4Vb6cIvCLOunyOZ2sqGBwEhV-q8mPJXND_h6-Y1ijKID5d83Te-5JSxy70utJFSPdRSVdBg8y9ZJ8W5ahZeor8RkjT9Ng"
	publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApmHKN/PQAscFDPMKeuEa
ex9BFQ3G7gnDfM3Q51TWMHlfyTX3teiZeCLn+fObZ3FQfDG4JkYODq43b8ltAjsr
dToseETIZukjOGHoulTZpS5B64yy3QDT6hqsEwwo514P6QI3IUP5ZJUtlIcF7HyK
rBBkB3QgoQ6KuzGQmELcaWjp3Ai1k8f68l+gXa4fnjgwqt4FeMbfNvrNK3ohWPVK
eTwfeelfgkehBua2BazuPyx38ry5DjXtXKfRCvhbz+KF4zsaU9CWgmmY/esjw0B5
SAjakNxq/e9YwLn9K2r/eWkzRwIuXcnLqqU1GqmE5YAzADSuEeznspREVr1vo2sF
swIDAQAB
-----END PUBLIC KEY-----`
	context = ctxtg.Context{
		Token: ctxtg.Token(jwt),
	}
	tracker = entities.TrackerConfig{
		URL: "http://test-jira:8080",
		Credentials: entities.TrackerCredentials{
			Login:    "admin",
			Password: "adminpass",
		},
	}
	estimate   = entities.Duration(3600)
	dueDate    = entities.Timestamp(1482624000)
	reportDate = entities.Timestamp(1462104000)
	testIssue  = entities.Issue{
		ID:       10000,
		Type:     entities.NamedID{ID: 10000, Name: "Task"},
		URL:      "http://localhost:8080/browse/TP-1",
		Title:    "Test task",
		Estimate: estimate,
		DueDate:  dueDate,
		Spent:    1800,
		Done:     50,
	}

	jiraContainer ContainerID
	jiraIP        string
	jiraPort      int
	jiraURL       string
)

type rpcRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"string"`
	ID      int                    `json:"id"`
}

func TestMain(m *testing.M) {
	cfg.HTTP.Listen = "127.0.0.1:8000"

	var err error

	// START JIRA CONTAINER
	BindDockerToLocalhost = "1"
	jiraContainer, jiraIP, jiraPort, err = SetupCustomContainer("tgms/test-jira", 8080, 120*time.Second)
	if err != nil {
		log.Fatalf("Could not setup container: %s", err)
	}
	tracker.URL = fmt.Sprintf("http://%v:%v", jiraIP, jiraPort)
	testIssue.URL = fmt.Sprintf("%v/browse/TP-1", tracker.URL)

	err = ConnectToCustomContainer(fmt.Sprintf("http://%v:%v/rest/api/2/mypermissions", jiraIP, jiraPort), 120, time.Second*3, func(url string) bool {
		request, err := http.NewRequest("GET", url, nil)
		if err != nil {
			panic(err)
		}
		client := &http.Client{}

		response, err := client.Do(request)
		if err != nil {
			return false
		}
		defer response.Body.Close()

		return response.StatusCode == 200
	})

	if err != nil {
		l.Fatal(err)
	}

	// START SERVER
	db, err = bolt.Open("var/bolt/store.db", 0666, nil)
	if err != nil {
		l.Fatal(err)
	}

	go start(appParams{
		BoltDB:    db,
		PublicKey: []byte(publicKey),
	})

	time.Sleep(time.Second)

	code := m.Run()

	db.Close()
	jiraContainer.KillRemove()
	os.Exit(code)
}

func TestUnauthorized(t *testing.T) {
	client := jsonrpc2.NewHTTPClient("http://127.0.0.1:8000/rpc")
	defer client.Close()

	var (
		result interface{}
		err    error
	)

	err = client.Call("API.GetProjects", map[string]interface{}{"Context": ctxtg.Context{}, "Tracker": tracker}, &result)
	assert.Equal(t, *entities.ErrUnauthorized, newRPCError(err))
	assert.Nil(t, result)
}

func TestGetProjects(t *testing.T) {
	client := jsonrpc2.NewHTTPClient("http://127.0.0.1:8000/rpc")
	defer client.Close()

	type ProjectsResult struct {
		Projects []entities.Project
	}

	var (
		result   ProjectsResult
		err      error
		expected = ProjectsResult{
			Projects: []entities.Project{
				{
					IssueTypes: []entities.NamedID{
						{ID: 10000, Name: "Task"},
						{ID: 10001, Name: "Sub-task"}},
					ActivityTypes: []entities.NamedID{},
					ID:            10000,
					Title:         "Test Project"}}}
	)

	err = client.Call("API.GetProjects", map[string]interface{}{"Context": context, "Tracker": tracker}, &result)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestGetCurrentUser(t *testing.T) {
	client := jsonrpc2.NewHTTPClient("http://127.0.0.1:8000/rpc")
	defer client.Close()

	type UserResult struct {
		User entities.User
	}

	var (
		result   UserResult
		err      error
		expected = UserResult{
			User: entities.User{
				ID:   1,
				Name: "Admin",
				Mail: "admin@nowhere.com",
			},
		}
	)

	err = client.Call("API.GetCurrentUser", map[string]interface{}{"Context": context, "Tracker": tracker}, &result)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestCreateReport(t *testing.T) {
	client := jsonrpc2.NewHTTPClient("http://127.0.0.1:8000/rpc")
	defer client.Close()

	var (
		result    interface{}
		err       error
		newReport = entities.Report{
			IssueID:  10000,
			Started:  entities.Timestamp(reportDate),
			Duration: 1800,
			Comments: "Test report",
		}
	)

	err = client.Call("API.CreateReport", map[string]interface{}{"Context": context, "Tracker": tracker, "Report": newReport}, &result)
	assert.Nil(t, err)
}

func TestGetProjectIssues(t *testing.T) {
	client := jsonrpc2.NewHTTPClient("http://127.0.0.1:8000/rpc")
	defer client.Close()

	type IssuesResponse struct {
		Issues []entities.Issue
	}

	var (
		result   IssuesResponse
		err      error
		expected = IssuesResponse{
			Issues: []entities.Issue{testIssue},
		}
	)

	err = client.Call("API.GetProjectIssues", map[string]interface{}{
		"Context":   context,
		"Tracker":   tracker,
		"ProjectID": 10000,
		"UserID":    1,
	}, &result)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func TestGetIssue(t *testing.T) {
	client := jsonrpc2.NewHTTPClient("http://127.0.0.1:8000/rpc")
	defer client.Close()

	type IssueResponse struct {
		Issue entities.Issue
	}

	var (
		result   IssueResponse
		err      error
		expected = IssueResponse{
			Issue: testIssue,
		}
	)

	err = client.Call("API.GetIssue", map[string]interface{}{"Context": context, "Tracker": tracker, "IssueID": 10000}, &result)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)

	err = client.Call("API.GetIssue", map[string]interface{}{"Context": context, "Tracker": tracker, "IssueID": 1000}, &result)
	if assert.Error(t, err) {
		assert.Equal(t, *entities.ErrIssueNotFound, newRPCError(err))
	}
}

func TestGetIssueByURL(t *testing.T) {
	client := jsonrpc2.NewHTTPClient("http://127.0.0.1:8000/rpc")
	defer client.Close()

	var (
		result   api.GetIssueByURLResponse
		url      = tracker.URL + "/browse/10000"
		err      error
		expected = api.GetIssueByURLResponse{
			ProjectID: 10000,
			Issue:     testIssue,
		}
	)

	err = client.Call("API.GetIssueByURL", map[string]interface{}{"Context": context, "Tracker": tracker, "IssueURL": url}, &result)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)

	err = client.Call("API.GetIssueByURL", map[string]interface{}{"Context": context, "Tracker": tracker, "IssueURL": url[:len(url)-1]}, &result)
	if assert.Error(t, err) {
		assert.Equal(t, *entities.ErrIssueNotFound, newRPCError(err))
	}
}

func TestCreateIssue(t *testing.T) {
	client := jsonrpc2.NewHTTPClient("http://127.0.0.1:8000/rpc")
	defer client.Close()

	type IssueResponse struct {
		Issue entities.Issue
	}

	var (
		result   IssueResponse
		err      error
		newIssue = entities.NewIssue{
			ProjectID: 10000,
			Assignee:  1,
			Type:      entities.EntityID(10000),
			Title:     "Test task 2",
			Estimate:  1800,
		}
		expected = IssueResponse{
			Issue: entities.Issue{
				ID:       1,
				URL:      fmt.Sprintf("%v/browse/TP-2", tracker.URL),
				Type:     entities.NamedID{ID: 10000, Name: "Task"},
				Title:    "Test task 2",
				Estimate: 1800,
			},
		}
	)

	err = client.Call("API.CreateIssue", map[string]interface{}{"Context": context, "Tracker": tracker, "Issue": newIssue}, &result)
	assert.Nil(t, err)
	assert.True(t, result.Issue.ID > 10000)
	result.Issue.ID = 1
	assert.Equal(t, expected, result)
}

func TestGetTotalReports(t *testing.T) {
	client := jsonrpc2.NewHTTPClient("http://127.0.0.1:8000/rpc")
	defer client.Close()

	type TotalResponse struct {
		Total int
	}

	var (
		result   TotalResponse
		err      error
		expected = TotalResponse{Total: 1800}
	)

	err = client.Call("API.GetTotalReports", map[string]interface{}{"Context": context, "Tracker": tracker, "Date": 1462060800}, &result)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}

func newRPCError(err error) jsonrpc2.Error {
	var result jsonrpc2.Error
	_ = json.Unmarshal([]byte(err.Error()), &result)
	return result
}
