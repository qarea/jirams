package rpcsvc

import (
	"../../entities"
	"gitlab.qarea.org/tgms/ctxtg"
)

// ProjectsReq input parameter to GetProjects
type ProjectsReq struct {
	Context ctxtg.Context
	Tracker entities.Tracker
	entities.Pagination
}

// ProjectsResp output parameter from GetProjects
type ProjectsResp struct {
	Projects []entities.Project
	Amount   int64
}

// ProjectDetailsReq input parameter to GetProjectDetails
type ProjectDetailsReq struct {
	Context   ctxtg.Context
	Tracker   entities.Tracker
	ProjectID entities.ProjectID
}

// ProjectDetailsResp output parameter from GetProjectDetails
type ProjectDetailsResp struct {
	Project entities.Project
}

// CurrentUserReq input parameter to GetCurrentUser
type CurrentUserReq struct {
	Context ctxtg.Context
	Tracker entities.Tracker
}

// CurrentUserResp output parameter from GetCurrentUser
type CurrentUserResp struct {
	User entities.User
}

// ProjectIssuesReq input parameter to GetProjectIssues
type ProjectIssuesReq struct {
	Context   ctxtg.Context
	Tracker   entities.Tracker
	ProjectID entities.ProjectID
	entities.Pagination
}

// ProjectIssuesResp output parameter from GetProjectIssues
type ProjectIssuesResp struct {
	Issues []entities.Issue
	Amount int64
}

// CreateIssueReq input parameter to CreateIssue
type CreateIssueReq struct {
	Context   ctxtg.Context
	Tracker   entities.Tracker
	Issue     entities.NewIssue
	ProjectID entities.ProjectID
}

// CreateIssueResp output parameter from CreateIssue
type CreateIssueResp struct {
	Issue entities.Issue
}

// GetIssueReq input parameter to GetIssue
type GetIssueReq struct {
	Context   ctxtg.Context
	Tracker   entities.Tracker
	IssueID   entities.IssueID
	ProjectID entities.ProjectID
}

// GetIssueResp output parameter from GetIssue
type GetIssueResp struct {
	Issue entities.Issue
}

// CreateReportReq input parameter from CreateReport
type CreateReportReq struct {
	Context   ctxtg.Context
	Tracker   entities.Tracker
	ProjectID entities.ProjectID
	Report    entities.Report
}

// GetIssueByURLReq input parameter to GetIssueByURL
type GetIssueByURLReq struct {
	Context  ctxtg.Context
	Tracker  entities.Tracker
	IssueURL entities.IssueURL
}

// GetIssueByURLResp output parameter from GetIssueByURL
type GetIssueByURLResp struct {
	Issue     entities.Issue
	ProjectID entities.ProjectID
}

// UpdateIssueProgressReq input parameter to UpdateIssueReq
type UpdateIssueProgressReq struct {
	Context   ctxtg.Context
	Tracker   entities.Tracker
	IssueID   entities.IssueID
	ProjectID entities.ProjectID
	Progress  entities.Progress
}

// GetReportsReq input parameter to GetTotalReports
type GetReportsReq struct {
	Context ctxtg.Context
	Tracker entities.Tracker
	Date    int64
}

// GetReportsResp output parameter from GetTotalReports
type GetReportsResp struct {
	Total int64
}
