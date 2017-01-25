package jira

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"github.com/powerman/narada-go/narada"
	"github.com/qarea/jirams/entities"
	"github.com/qarea/jirams/store"
)

var l = narada.NewLog("jira client: ")

const (
	basePath            = "/rest/api/2/"
	projectResource     = "project"
	currentUserResource = "myself"
	searchResource      = "search?jql="
	issueResource       = "issue"
	jiraTimestampLayout = "2006-01-02T15:04:05.000-0700"
	jiraDateLayout      = "2006/01/02"
)

// Client implements TrackerClient interface for JIRA tracker
type Client struct {
	Store store.UserKeyMapper
	Jira  TrackerRequester
}

// NewClient creates new instance of Client
func NewClient(store store.UserKeyMapper) *Client {
	return &Client{Store: store, Jira: &Requester{}}
}

// GetProjects fetches and returns a list of projects for current user
func (client *Client) GetProjects(tracker entities.TrackerConfig, res *[]entities.Project) (err error) {
	var (
		projects []Project
		baseURL  = tracker.URL + basePath
	)

	if err = validateTrackerConfig(tracker); err != nil {
		return
	}

	request, _ := http.NewRequest("GET", baseURL+projectResource, nil)
	if err = client.Jira.Request(tracker, request, &projects); err != nil {
		return
	}

	for i, project := range projects {
		request, _ := http.NewRequest("GET", baseURL+projectResource+"/"+project.ID, nil)
		err = client.Jira.Request(tracker, request, &project)
		if err != nil {
			continue
		}
		projects[i] = project
	}

	*res = make([]entities.Project, len(projects))
	for i, project := range projects {
		(*res)[i] = project.toProject()
	}

	return
}

// GetCurrentUser retrieves current user information from tracker
func (client *Client) GetCurrentUser(tracker entities.TrackerConfig, res *entities.User) (err error) {
	baseURL := tracker.URL + basePath
	var user User
	request, _ := http.NewRequest("GET", baseURL+currentUserResource, nil)
	err = client.Jira.Request(tracker, request, &user)
	if err != nil {
		return
	}
	*res, err = user.toUser(tracker.ID, client.Store)
	return
}

// GetProjectIssues retrieves list of issues from specified project assigned to current user
func (client *Client) GetProjectIssues(tracker entities.TrackerConfig, projectID entities.ProjectID, userID entities.UserID, res *[]entities.Issue) error {
	baseURL := tracker.URL + basePath

	var (
		startAt  = 0
		finished = false
		query    = fmt.Sprintf("project=%d AND assignee=currentUser() AND resolution=Unresolved", projectID)
		issues   []Issue
	)

	for !finished {
		limit := fmt.Sprintf("&startAt=%d", startAt)
		request, err := http.NewRequest("GET", baseURL+searchResource+url.QueryEscape(query)+limit, nil)
		if err != nil {
			return err
		}
		var chunk Issues
		err = client.Jira.Request(tracker, request, &chunk)
		if err != nil {
			return err
		}
		startAt += len(chunk.Issues)
		finished = startAt >= chunk.Total
		issues = append(issues, chunk.Issues...)
	}

	*res = make([]entities.Issue, len(issues))
	for i, issue := range issues {
		(*res)[i] = issue.toIssue()
	}
	return nil
}

// GetIssue retrieves issue information by issue ID
func (client *Client) GetIssue(tracker entities.TrackerConfig, issueID entities.IssueID, res *entities.Issue) error {
	baseURL := tracker.URL + basePath
	request, _ := http.NewRequest("GET", baseURL+issueResource+"/"+fmt.Sprintf("%d", issueID), nil)
	var issue Issue
	err := client.Jira.Request(tracker, request, &issue)
	if err != nil {
		if err == entities.ErrNotFound {
			err = entities.ErrIssueNotFound
		}
		return err
	}
	*res = issue.toIssue()
	return nil
}

var re = regexp.MustCompile("(issues|browse)\\/([0-9A-Z-]+)")

// GetIssueByURL attempts to parse provided URL and retrieve corresponding issue
func (client *Client) GetIssueByURL(tracker entities.TrackerConfig, issueURL string, res *entities.Issue, res2 *entities.ProjectID) error {
	matches := re.FindStringSubmatch(issueURL)
	if matches == nil {
		return errors.New("Failed to parse Issue URL")
	}

	baseURL := tracker.URL + basePath
	request, _ := http.NewRequest("GET", baseURL+issueResource+"/"+url.QueryEscape(matches[2]), nil)
	var issue Issue
	err := client.Jira.Request(tracker, request, &issue)
	if err != nil {
		if err == entities.ErrNotFound {
			err = entities.ErrIssueNotFound
		}
		return err
	}
	*res = issue.toIssue()
	*res2 = issue.Fields.ProjectID.toProjectID()

	return nil
}

// CreateIssue creates new issue with provided parameters
func (client *Client) CreateIssue(tracker entities.TrackerConfig, newIssue entities.NewIssue, res *entities.Issue) error {
	baseURL := tracker.URL + basePath
	userKey, err := client.Store.GetKey(tracker.ID, newIssue.Assignee)
	if err != nil || userKey == "" {
		return errors.New("Unknown User ID")
	}
	estimateHours := float64(int((newIssue.Estimate/60)*100) / 100)
	payload := NewIssue{
		Fields: NewIssueFields{
			Title:    newIssue.Title,
			Project:  ID{ID: fmt.Sprintf("%d", uint64(newIssue.ProjectID))},
			Type:     ID{ID: fmt.Sprintf("%d", uint64(newIssue.Type))},
			Assignee: UserKey{Key: userKey},
			TimeTracking: TimeTracking{
				OriginalEstimate:  estimateHours,
				RemainingEstimate: estimateHours}}}
	payloadBytes, _ := json.Marshal(payload)
	request, _ := http.NewRequest("POST", baseURL+issueResource, bytes.NewBuffer(payloadBytes))
	request.Header.Set("Content-Type", "application/json")
	var newIssueID EntityID
	err = client.Jira.Request(tracker, request, &newIssueID)
	if err != nil {
		return err
	}
	issueID := newIssueID.toIssueID()

	return client.GetIssue(tracker, issueID, res)
}

// CreateReport creates a work time report for specified issue
func (client *Client) CreateReport(tracker entities.TrackerConfig, report entities.Report) error {
	baseURL := tracker.URL + basePath
	started := time.Unix(int64(report.Started), 0).Format(jiraTimestampLayout)
	payload := Worklog{
		Started: started,
		Spent:   uint64(report.Duration),
		Comment: report.Comments}
	payloadBytes, _ := json.Marshal(payload)
	request, _ := http.NewRequest("POST", baseURL+issueResource+"/"+fmt.Sprintf("%d", report.IssueID)+"/worklog", bytes.NewBuffer(payloadBytes))
	request.Header.Set("Content-Type", "application/json")

	return client.Jira.Request(tracker, request, nil)
}

// GetTotalReports returns total time worked on specified date
func (client *Client) GetTotalReports(tracker entities.TrackerConfig, date entities.Timestamp, res *entities.ReportsTotal) error {
	baseURL := tracker.URL + basePath
	jiraDate := time.Unix(int64(date), 0).Format(jiraDateLayout)
	query := fmt.Sprintf("worklogAuthor=currentUser() AND worklogDate=\"%s\"", jiraDate)
	url := baseURL + searchResource + url.QueryEscape(query) + "&fields=id"
	// fetch list of issues that user reported to on given date
	var (
		issues   IssueIDPage
		issueIDs []entities.IssueID
	)
	err := client.Jira.IterateRequest(tracker, url, &issues, func(data interface{}) (loaded int, total int, err error) {
		if data, ok := data.(*IssueIDPage); ok {
			loaded = len(data.IssueIDs)
			total = data.Total
			issueIDsChunk := make([]entities.IssueID, loaded)
			for i, issueID := range data.IssueIDs {
				issueIDsChunk[i] = entities.IssueID(issueID.toUint64())
			}
			issueIDs = append(issueIDs, issueIDsChunk...)
		} else {
			err = errors.New("Expected data to be of type *IssueIDPage")
		}
		return
	})

	if err != nil {
		return err
	}

	// fetch worklog entries for each issue and calculate total report
	var worklogs WorklogPage
	for _, id := range issueIDs {
		url := baseURL + issueResource + "/" + fmt.Sprintf("%d", id) + "/worklog"
		err := client.Jira.IterateRequest(tracker, url, &worklogs, func(data interface{}) (loaded int, total int, err error) {
			if data, ok := data.(*WorklogPage); ok {
				loaded = len(data.Worklogs)
				total = data.Total
				for _, worklog := range data.Worklogs {
					worklogDate, err := time.Parse(jiraTimestampLayout, worklog.Started)
					worklogDateTruncated := worklogDate.Truncate(time.Hour * 24)
					if err == nil && worklogDateTruncated.Unix() == int64(date) {
						*res += entities.ReportsTotal(worklog.Spent)
					}
				}
			} else {
				err = errors.New("Expected data to be of type *WorklogPage")
			}
			return
		})
		if err != nil {
			return err
		}
	}
	return nil
}
