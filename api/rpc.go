package api

import (
	"context"

	"github.com/powerman/narada-go/narada"
	"github.com/qarea/ctxtg"
	"github.com/qarea/jirams/entities"
)

var l = narada.NewLog("api: ")

// NewRPCAPI creates an instance of API, implementing RPC interface
func NewRPCAPI(client TrackerClient, parser ctxtg.TokenParser) *API {
	return &API{Client: client, Parser: parser}
}

// TrackerClient defines interface for tracker client business logic implementation
type TrackerClient interface {
	GetProjects(tracker entities.TrackerConfig, res *[]entities.Project) error
	GetCurrentUser(tracker entities.TrackerConfig, res *entities.User) error
	GetProjectIssues(tracker entities.TrackerConfig, projectID entities.ProjectID, userID entities.UserID, res *[]entities.Issue) error
	GetIssue(tracker entities.TrackerConfig, issueID entities.IssueID, res *entities.Issue) error
	GetIssueByURL(tracker entities.TrackerConfig, issueURL string, res *entities.Issue, res2 *entities.ProjectID) error
	CreateIssue(tracker entities.TrackerConfig, issue entities.NewIssue, res *entities.Issue) error
	CreateReport(tracker entities.TrackerConfig, report entities.Report) error
	GetTotalReports(tracker entities.TrackerConfig, date entities.Timestamp, res *entities.ReportsTotal) error
}

// API implements service RPC interface
type API struct {
	Client TrackerClient
	Parser ctxtg.TokenParser
}

// GetProjects provides corresponding API method
func (api *API) GetProjects(req GetProjectsRequest, res *GetProjectsResponse) (err error) {
	err = api.Parser.ParseCtxWithClaims(req.Context, func(ctx context.Context, claims ctxtg.Claims) error {
		err = api.Client.GetProjects(req.Tracker, &res.Projects)
		if err != nil {
			err = entities.NewLoggedError(l, ctx, err, "Failed to retrieve projects")
		}
		return err
	})
	return
}

// GetCurrentUser provides corresponding API method
func (api *API) GetCurrentUser(req GetCurrentUserRequest, res *GetCurrentUserResponse) (err error) {
	err = api.Parser.ParseCtxWithClaims(req.Context, func(ctx context.Context, claims ctxtg.Claims) error {
		err = api.Client.GetCurrentUser(req.Tracker, &res.User)
		if err != nil {
			err = entities.NewLoggedError(l, ctx, err, "Failed to retrieve current user")
		}
		return err
	})
	return
}

// GetProjectIssues provides corresponding API method
func (api *API) GetProjectIssues(req GetProjectIssuesRequest, res *GetProjectIssuesResponse) (err error) {
	err = api.Parser.ParseCtxWithClaims(req.Context, func(ctx context.Context, claims ctxtg.Claims) error {
		err = api.Client.GetProjectIssues(req.Tracker, req.ProjectID, req.UserID, &res.Issues)
		if err != nil {
			err = entities.NewLoggedError(l, ctx, err, "Failed to retrieve project issues")
		}
		return err
	})
	return
}

// CreateIssue provides corresponding API method
func (api *API) CreateIssue(req CreateIssueRequest, res *CreateIssueResponse) (err error) {
	err = api.Parser.ParseCtxWithClaims(req.Context, func(ctx context.Context, claims ctxtg.Claims) error {
		req.Issue.Assignee = entities.UserID(claims.UserID)
		err = api.Client.CreateIssue(req.Tracker, req.Issue, &res.Issue)
		if err != nil {
			err = entities.NewLoggedError(l, ctx, err, "Failed to create issue")
		}
		return err
	})
	return
}

// GetIssue provides corresponding API method
func (api *API) GetIssue(req GetIssueRequest, res *GetIssueResponse) (err error) {
	err = api.Parser.ParseCtxWithClaims(req.Context, func(ctx context.Context, claims ctxtg.Claims) error {
		err = api.Client.GetIssue(req.Tracker, req.IssueID, &res.Issue)
		if err != nil {
			err = entities.NewLoggedError(l, ctx, err, "Failed to retrieve issue")
		}
		return err
	})
	return
}

// CreateReport provides corresponding API method
func (api *API) CreateReport(req CreateReportRequest, res *CreateReportResponse) (err error) {
	err = api.Parser.ParseCtxWithClaims(req.Context, func(ctx context.Context, claims ctxtg.Claims) error {
		err = api.Client.CreateReport(req.Tracker, req.Report)
		if err != nil {
			err = entities.NewLoggedError(l, ctx, err, "Failed to create report")
		}
		return err
	})
	return
}

// GetTotalReports provides corresponding API method
func (api *API) GetTotalReports(req GetTotalReportsRequest, res *GetTotalReportsResponse) (err error) {
	err = api.Parser.ParseCtxWithClaims(req.Context, func(ctx context.Context, claims ctxtg.Claims) error {
		err = api.Client.GetTotalReports(req.Tracker, req.Date, &res.Total)
		if err != nil {
			err = entities.NewLoggedError(l, ctx, err, "Failed to retrieve reports")
		}
		return err
	})
	return
}

// GetIssueByURL provides corresponding API method
func (api *API) GetIssueByURL(req GetIssueByURLRequest, res *GetIssueByURLResponse) (err error) {
	err = api.Parser.ParseCtxWithClaims(req.Context, func(ctx context.Context, claims ctxtg.Claims) error {
		err = api.Client.GetIssueByURL(req.Tracker, req.IssueURL, &res.Issue, &res.ProjectID)
		if err != nil {
			err = entities.NewLoggedError(l, ctx, err, "Failed to retrieve issue")
		}
		return err
	})
	return
}

// UpdateIssueProgress provides dummy implementation of corrcponding method
// as the feature is not supported by JIRA
func (api *API) UpdateIssueProgress(req UpdateIssueProgressRequest, res *UpdateIssueProgressResponse) (err error) {
	return
}
