// Package rpcsvc provides handlers for JSON-RPC 2.0.
package rpcsvc

import (
	"context"
	"net"
	"net/http"
	"net/rpc"

	"gitlab.qarea.org/tgms/ctxtg"

	"github.com/pkg/errors"
	"github.com/powerman/narada-go/narada"

	"../../cfg"
	"../../entities"

	"github.com/powerman/rpc-codec/jsonrpc2"
)

var log = narada.NewLog("rpcsvc: ")

// TrackerClient required interface for new tracker
type TrackerClient interface {
	Project(context.Context, entities.Tracker, entities.ProjectID) (*entities.Project, error)
	//Projects return project list and total amount of projects
	Projects(context.Context, entities.Tracker, entities.Pagination) ([]entities.Project, int64, error)
	//ProjectIssues return issues assigned to user and total amount
	ProjectIssues(context.Context, entities.Tracker, entities.ProjectID, entities.Pagination) ([]entities.Issue, int64, error)
	UserInfo(context.Context, entities.Tracker) (*entities.User, error)
	Issue(context.Context, entities.Tracker, entities.ProjectID, entities.IssueID) (*entities.Issue, error)
	IssueByURL(context.Context, entities.Tracker, entities.IssueURL) (*entities.Issue, error)
	CreateIssue(context.Context, entities.Tracker, entities.NewIssue, entities.ProjectID) (*entities.Issue, error)
	UpdateIssueProgress(context.Context, entities.Tracker, entities.ProjectID, entities.IssueID, entities.Progress) error
	//TotalReports receive date as UNIX timestamp (seconds) and return total reported time at this day in seconds
	TotalReports(ctx context.Context, t entities.Tracker, date int64) (int64, error)
	CreateReport(context.Context, entities.Tracker, entities.ProjectID, entities.Report) error
}

// Init registers JSON-RPC handlers
func Init(r TrackerClient, p ctxtg.TokenParser) {
	if err := rpc.Register(newAPI(r, p)); err != nil {
		log.Fatal(err)
	}
	http.Handle(cfg.HTTP.BasePath+"/rpc", jsonrpc2.HTTPHandler(nil))
}

func newAPI(r TrackerClient, p ctxtg.TokenParser) *API {
	return &API{
		tracker:     r,
		tokenParser: p,
	}
}

// API JSON-RPC 2.0 struct
type API struct {
	tracker     TrackerClient
	tokenParser ctxtg.TokenParser
}

// Version return current service narada version
func (*API) Version(args *struct{}, res *string) error {
	*res, _ = narada.Version()
	return nil
}

// GetProjects return paginated projects list for user
func (r *API) GetProjects(req *ProjectsReq, resp *ProjectsResp) error {
	err := r.tokenParser.ParseCtxWithClaims(req.Context, func(ctx context.Context, c ctxtg.Claims) error {
		projects, amount, err := r.tracker.Projects(ctx, req.Tracker, req.Pagination)
		*resp = ProjectsResp{
			Projects: projects,
			Amount:   amount,
		}
		return err
	})
	return errWithLog(req.Context, "fail to GetProjects", err)
}

// GetProjectDetails returns full information by projectID
func (r *API) GetProjectDetails(req *ProjectDetailsReq, resp *ProjectDetailsResp) error {
	err := r.tokenParser.ParseCtxWithClaims(req.Context, func(ctx context.Context, c ctxtg.Claims) error {
		project, err := r.tracker.Project(ctx, req.Tracker, req.ProjectID)
		if project != nil {
			*resp = ProjectDetailsResp{
				Project: *project,
			}
		}
		return err
	})
	return errWithLog(req.Context, "fail to GetProjectDetails", err)
}

// GetCurrentUser returns current user info
func (r *API) GetCurrentUser(req *CurrentUserReq, resp *CurrentUserResp) error {
	err := r.tokenParser.ParseCtxWithClaims(req.Context, func(ctx context.Context, c ctxtg.Claims) error {
		u, err := r.tracker.UserInfo(ctx, req.Tracker)
		if u != nil {
			*resp = CurrentUserResp{
				User: *u,
			}
		}
		return err
	})
	return errWithLog(req.Context, "current user info err", err)
}

// GetProjectIssues returns user's issues by project ID
func (r *API) GetProjectIssues(req *ProjectIssuesReq, resp *ProjectIssuesResp) error {
	err := r.tokenParser.ParseCtxWithClaims(req.Context, func(ctx context.Context, c ctxtg.Claims) error {
		is, amount, err := r.tracker.ProjectIssues(ctx, req.Tracker, req.ProjectID, req.Pagination)
		*resp = ProjectIssuesResp{
			Issues: is,
			Amount: amount,
		}
		return err
	})
	return errWithLog(req.Context, "issues err", err)
}

// CreateIssue creates issue on tracker
func (r *API) CreateIssue(req *CreateIssueReq, resp *CreateIssueResp) error {
	err := r.tokenParser.ParseCtxWithClaims(req.Context, func(ctx context.Context, c ctxtg.Claims) error {
		issue, err := r.tracker.CreateIssue(ctx, req.Tracker, req.Issue, req.ProjectID)
		if issue != nil {
			*resp = CreateIssueResp{
				Issue: *issue,
			}
		}
		return err
	})
	return errWithLog(req.Context, "create issue err", err)
}

// GetIssue returns Issue by ID
func (r *API) GetIssue(req *GetIssueReq, resp *GetIssueResp) error {
	err := r.tokenParser.ParseCtxWithClaims(req.Context, func(ctx context.Context, c ctxtg.Claims) error {
		issue, err := r.tracker.Issue(ctx, req.Tracker, req.ProjectID, req.IssueID)
		if issue != nil {
			*resp = GetIssueResp{
				Issue: *issue,
			}
		}
		return err
	})
	return errWithLog(req.Context, "issue by ID err", err)
}

// UpdateIssueProgress updates issue progress in percents
func (r *API) UpdateIssueProgress(req *UpdateIssueProgressReq, _ *struct{}) error {
	err := r.tokenParser.ParseCtxWithClaims(req.Context, func(ctx context.Context, c ctxtg.Claims) error {
		return r.tracker.UpdateIssueProgress(ctx, req.Tracker, req.ProjectID, req.IssueID, req.Progress)
	})
	return errWithLog(req.Context, "update issue err", err)
}

// CreateReport reports time on tracker for user ID
func (r *API) CreateReport(req *CreateReportReq, _ *struct{}) error {
	err := r.tokenParser.ParseCtxWithClaims(req.Context, func(ctx context.Context, c ctxtg.Claims) error {
		return r.tracker.CreateReport(ctx, req.Tracker, req.ProjectID, req.Report)
	})
	return errWithLog(req.Context, "create report err", err)
}

// GetTotalReports receive UNIX timestamp of date and aggregate reported time for user for this day
func (r *API) GetTotalReports(req *GetReportsReq, resp *GetReportsResp) error {
	err := r.tokenParser.ParseCtxWithClaims(req.Context, func(ctx context.Context, c ctxtg.Claims) error {
		time, err := r.tracker.TotalReports(ctx, req.Tracker, req.Date)
		*resp = GetReportsResp{
			Total: time,
		}
		return err
	})
	return errWithLog(req.Context, "total reports err", err)
}

// GetIssueByURL parse incoming URL and return issue and project ID
func (r *API) GetIssueByURL(req *GetIssueByURLReq, resp *GetIssueByURLResp) error {
	err := r.tokenParser.ParseCtxWithClaims(req.Context, func(ctx context.Context, c ctxtg.Claims) error {
		issue, err := r.tracker.IssueByURL(ctx, req.Tracker, req.IssueURL)
		if issue == nil {
			return entities.ErrIssueNotFound
		}
		*resp = GetIssueByURLResp{
			Issue:     *issue,
			ProjectID: issue.ProjectID,
		}
		return err
	})
	return errWithLog(req.Context, "issue by URL err", err)
}

func errWithLog(ctx ctxtg.Context, prefix string, err error) error {
	if err == nil {
		return nil
	}
	log.ERR("tracking id: %s, token: %s, %s: %+v", ctx.TracingID, ctx.Token, prefix, err)
	err = errors.Cause(err)
	if err == context.DeadlineExceeded || isHTTPTimeoutErr(err) {
		return entities.ErrTimeout
	}
	return err
}

func isHTTPTimeoutErr(err error) bool {
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return true
	}
	return false
}
