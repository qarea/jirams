package api

import (
	"testing"

	"github.com/qarea/ctxtg"
	"github.com/qarea/jirams/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockParser struct{ mock.Mock }

// Parse mock
func (t *mockParser) Parse(token ctxtg.Token) (*ctxtg.Claims, error) {
	args := t.Called(token)
	claims := args.Get(0).(ctxtg.Claims)
	return &claims, args.Error(1)
}

// ParseWithClaims mock
func (t *mockParser) ParseWithClaims(token ctxtg.Token, f ctxtg.ClaimsFunc) error {
	args := t.Called(token)
	claims := args.Get(0).(ctxtg.Claims)
	_ = f(claims)
	return args.Error(1)
}

// ParseCtxWithClaims mock
func (t *mockParser) ParseCtxWithClaims(c ctxtg.Context, f ctxtg.CtxClaimsFunc) error {
	args := t.Called(c)
	return args.Error(1)
}

// Each request should be authorized
func TestAuthorization(t *testing.T) {
	ctx := ctxtg.Context{}
	parser := mockParser{}
	parser.On("ParseCtxWithClaims", ctx).Return(ctxtg.Claims{UserID: 1}, nil).Times(8)

	tracker := entities.TrackerConfig{}
	api := NewRPCAPI(nil, &parser)

	err := api.GetProjects(GetProjectsRequest{Context: ctx, Tracker: tracker}, nil)
	assert.Nil(t, err)

	err = api.GetCurrentUser(GetCurrentUserRequest{Context: ctx, Tracker: tracker}, nil)
	assert.Nil(t, err)

	err = api.GetProjectIssues(GetProjectIssuesRequest{Context: ctx, Tracker: tracker, ProjectID: 1, UserID: 1}, nil)
	assert.Nil(t, err)

	err = api.GetIssue(GetIssueRequest{Context: ctx, Tracker: tracker, IssueID: 1}, nil)
	assert.Nil(t, err)

	err = api.GetIssueByURL(GetIssueByURLRequest{Context: ctx, Tracker: tracker, IssueURL: "http://tracker.com/browse/1"}, nil)
	assert.Nil(t, err)

	err = api.CreateIssue(CreateIssueRequest{Context: ctx, Tracker: tracker, Issue: entities.NewIssue{}}, nil)
	assert.Nil(t, err)

	err = api.CreateReport(CreateReportRequest{Context: ctx, Tracker: tracker, Report: entities.Report{}}, nil)
	assert.Nil(t, err)

	err = api.GetTotalReports(GetTotalReportsRequest{Context: ctx, Tracker: tracker, Date: 15000000}, nil)
	assert.Nil(t, err)

	parser.AssertExpectations(t)
}
