package api

import (
	"github.com/qarea/ctxtg"
	"github.com/qarea/jirams/entities"
)

// GetProjectsRequest request arguments
type GetProjectsRequest struct {
	Context ctxtg.Context
	Tracker entities.TrackerConfig
}

// GetProjectsResponse response structure
type GetProjectsResponse struct {
	Projects []entities.Project
}

// GetCurrentUserRequest request arguments
type GetCurrentUserRequest struct {
	Context ctxtg.Context
	Tracker entities.TrackerConfig
}

// GetCurrentUserResponse response structure
type GetCurrentUserResponse struct {
	User entities.User
}

// GetProjectIssuesRequest request arguments
type GetProjectIssuesRequest struct {
	Context   ctxtg.Context
	Tracker   entities.TrackerConfig
	ProjectID entities.ProjectID
	UserID    entities.UserID
}

// GetProjectIssuesResponse response structure
type GetProjectIssuesResponse struct {
	Issues []entities.Issue
}

// CreateIssueRequest request arguments
type CreateIssueRequest struct {
	Context ctxtg.Context
	Tracker entities.TrackerConfig
	Issue   entities.NewIssue
}

// CreateIssueResponse response structure
type CreateIssueResponse struct {
	Issue entities.Issue
}

// GetIssueRequest request arguments
type GetIssueRequest struct {
	Context ctxtg.Context
	Tracker entities.TrackerConfig
	IssueID entities.IssueID
}

// GetIssueResponse response structure
type GetIssueResponse struct {
	Issue entities.Issue
}

// UpdateIssueProgressRequest request arguments
type UpdateIssueProgressRequest struct {
	Context  ctxtg.Context
	IssueID  entities.IssueID
	Progress uint64
}

// UpdateIssueProgressResponse response structure
type UpdateIssueProgressResponse struct{}

// CreateReportRequest request arguments
type CreateReportRequest struct {
	Context ctxtg.Context
	Tracker entities.TrackerConfig
	Report  entities.Report
}

// CreateReportResponse response structure
type CreateReportResponse struct{}

// GetTotalReportsRequest request arguments
type GetTotalReportsRequest struct {
	Context ctxtg.Context
	Tracker entities.TrackerConfig
	Date    entities.Timestamp
}

// GetTotalReportsResponse response structure
type GetTotalReportsResponse struct {
	Total entities.ReportsTotal
}

// GetIssueByURLRequest request arguments
type GetIssueByURLRequest struct {
	Context  ctxtg.Context
	Tracker  entities.TrackerConfig
	IssueURL string
}

// GetIssueByURLResponse response structure
type GetIssueByURLResponse struct {
	ProjectID entities.ProjectID
	Issue     entities.Issue
}
