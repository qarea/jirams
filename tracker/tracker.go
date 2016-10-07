package tracker

import (
	"context"
	"fmt"
	"net/http"

	"github.com/powerman/narada-go/narada"

	"../cfg"
	"../entities"
)

var log = narada.NewLog("tracker client: ")

// NewClient returns new instance of tracker client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout:       cfg.HTTP.Timeout,
			CheckRedirect: copyHeadersOnRedirect,
		},
	}
}

// Client to tracker
type Client struct {
	httpClient *http.Client
}

// Project return project by id or err if not foind project
func (c *Client) Project(ctx context.Context, tr entities.Tracker, pid entities.ProjectID) (*entities.Project, error) {
	panic("not implemented")
}

// Projects return project list and total amount of projects
func (c *Client) Projects(ctx context.Context, tr entities.Tracker, p entities.Pagination) ([]entities.Project, int64, error) {
	panic("not implemented")
}

// ProjectIssues return issues assigned to user and total amount
func (c *Client) ProjectIssues(ctx context.Context, tr entities.Tracker, pid entities.ProjectID, p entities.Pagination) ([]entities.Issue, int64, error) {
	panic("not implemented")
}

// UserInfo returns user info from tracker t
func (c *Client) UserInfo(ctx context.Context, tr entities.Tracker) (*entities.User, error) {
	panic("not implemented")
}

// Issue query issue by ID from tracker t or err if not found
func (c *Client) Issue(ctx context.Context, tr entities.Tracker, pid entities.ProjectID, iid entities.IssueID) (*entities.Issue, error) {
	panic("not implemented")
}

// IssueByURL query issue by URL ar err if not found
func (c *Client) IssueByURL(ctx context.Context, tr entities.Tracker, url entities.IssueURL) (*entities.Issue, error) {
	panic("not implemented")
}

// CreateIssue for projectID and assign it to user
func (c *Client) CreateIssue(ctx context.Context, tr entities.Tracker, issue entities.NewIssue, pid entities.ProjectID) (*entities.Issue, error) {
	panic("not implemented")
}

// UpdateIssueProgress updates issue progress by issue id
func (c *Client) UpdateIssueProgress(ctx context.Context, tr entities.Tracker, pid entities.ProjectID, is entities.IssueID, pr entities.Progress) error {
	panic("not implemented")
}

// TotalReports receive date as UNIX timestamp (seconds) and return total reported time at this day in seconds
func (c *Client) TotalReports(ctx context.Context, t entities.Tracker, date int64) (int64, error) {
	panic("not implemented")
}

// CreateReport for user
func (c *Client) CreateReport(ctx context.Context, tr entities.Tracker, pid entities.ProjectID, rep entities.Report) error {
	panic("not implemented")
}

func copyHeadersOnRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return fmt.Errorf("too many redirects")
	}
	if len(via) == 0 {
		return nil
	}
	for attr, val := range via[0].Header {
		if _, ok := req.Header[attr]; !ok {
			req.Header[attr] = val
		}
	}
	return nil
}
