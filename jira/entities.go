package jira

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/qarea/jirams/entities"
	"github.com/qarea/jirams/store"
)

const (
	dueDateLayout = "2006-01-02"
)

// Project - JIRA project structure
type Project struct {
	ID          string    `json:"id"`
	Title       string    `json:"name"`
	Description string    `json:"description,omitempty"`
	IssueTypes  []NamedID `json:"issueTypes,omitempty"`
}

func (project *Project) toProject() (res entities.Project) {
	id, _ := strconv.ParseUint(project.ID, 10, 64)
	issueTypes := make([]entities.NamedID, len(project.IssueTypes))
	for i, issueType := range project.IssueTypes {
		issueTypes[i] = issueType.toNamedID()
	}
	res = entities.Project{
		ID:            entities.ProjectID(id),
		Title:         project.Title,
		Description:   project.Description,
		IssueTypes:    issueTypes,
		ActivityTypes: make([]entities.NamedID, 0),
	}
	return
}

// NamedID - generic entity consisting of ID and name (issue type etc)
type NamedID struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (t *NamedID) toNamedID() (res entities.NamedID) {
	id, _ := strconv.ParseUint(t.ID, 10, 64)
	res = entities.NamedID{ID: id, Name: t.Name}
	return
}

// User - JIRA user structure
type User struct {
	Key  entities.UserKey `json:"key"`
	Name string           `json:"displayName"`
	Mail string           `json:"emailAddress"`
}

func (user *User) toUser(trackerID entities.TrackerID, store store.UserKeyMapper) (res entities.User, err error) {
	id, err := store.GetID(trackerID, user.Key)
	if err != nil {
		return *new(entities.User), err
	}
	return entities.User{
		ID:   id,
		Name: user.Name,
		Mail: user.Mail}, nil
}

// IssueFields - JIRA issue properties structure
type IssueFields struct {
	Type      NamedID       `json:"issuetype"`
	Title     string        `json:"summary"`
	Estimate  uint64        `json:"timeoriginalestimate,omitempty"`
	DueDate   string        `json:"duedate,omitempty"`
	Spent     uint64        `json:"timespent,omitempty"`
	Progress  IssueProgress `json:"progress"`
	ProjectID EntityID      `json:"project"`
}

// IssueProgress - JIRA time management data structure
type IssueProgress struct {
	Percent uint8 `json:"percent"`
}

// Issue - JIRA issue structure
type Issue struct {
	ID     string      `json:"id"`
	URL    string      `json:"self"`
	Key    string      `json:"key"`
	Fields IssueFields `json:"fields"`
}

func (issue *Issue) toIssue() entities.Issue {
	var (
		id      uint64
		dueDate uint64
	)
	id, _ = strconv.ParseUint(issue.ID, 10, 64)

	if issue.Fields.DueDate != "" {
		dueDateTime, err := time.Parse(dueDateLayout, issue.Fields.DueDate)
		if err == nil {
			dueDateTS := uint64(dueDateTime.Unix())
			dueDate = dueDateTS
		}
	}

	return entities.Issue{
		ID:       entities.IssueID(id),
		Type:     (&NamedID{issue.Fields.Type.ID, issue.Fields.Type.Name}).toNamedID(),
		URL:      strings.Replace(issue.URL, fmt.Sprintf("/rest/api/2/issue/%d", id), fmt.Sprintf("/browse/%v", issue.Key), 1),
		Title:    issue.Fields.Title,
		Estimate: entities.Duration(issue.Fields.Estimate),
		DueDate:  entities.Timestamp(dueDate),
		Spent:    entities.Duration(issue.Fields.Spent),
		Done:     issue.Fields.Progress.Percent,
	}
}

// Issues - JIRA issues collection
type Issues struct {
	StartAt    int     `json:"startAt"`
	MaxResults int     `json:"maxResults"`
	Total      int     `json:"total"`
	Issues     []Issue `json:"issues"`
}

// NewIssue - JIRA issue creation parameters
type NewIssue struct {
	Fields NewIssueFields `json:"fields"`
}

// NewIssueFields - New JIRA issue parameters
type NewIssueFields struct {
	Title        string       `json:"summary"`
	Description  string       `json:"description,omitempty"`
	Project      ID           `json:"project"`
	Type         ID           `json:"issuetype"`
	Assignee     UserKey      `json:"assignee"`
	TimeTracking TimeTracking `json:"timetracking,omitempty"`
}

// ID - Generic entity with ID string
type ID struct {
	ID string `json:"id"`
}

func (s *ID) toUint64() uint64 {
	id, _ := strconv.ParseUint(s.ID, 10, 64)
	return id
}

// UserKey - Key only JIRA user structure
type UserKey struct {
	Key string `json:"name"`
}

// TimeTracking - JIRA time tracking data structure
type TimeTracking struct {
	OriginalEstimate  float64 `json:"originalEstimate"`
	RemainingEstimate float64 `json:"remainingEstimate"`
}

// EntityID - Generic string entity ID
type EntityID struct {
	ID string `json:"id"`
}

func (s *EntityID) toIssueID() entities.IssueID {
	id, _ := strconv.ParseUint(s.ID, 10, 64)
	return entities.IssueID(id)
}

func (s *EntityID) toProjectID() entities.ProjectID {
	id, _ := strconv.ParseUint(s.ID, 10, 64)
	return entities.ProjectID(id)
}

// Worklog - JIRA worklog structure
type Worklog struct {
	Started string `json:"started"`
	Spent   uint64 `json:"timeSpentSeconds"`
	Comment string `json:"comment,omitempty"`
}

// IssueIDPage - JIRA issue ID collection page
type IssueIDPage struct {
	StartAt    int  `json:"startAt"`
	MaxResults int  `json:"maxResults"`
	Total      int  `json:"total"`
	IssueIDs   []ID `json:"issues"`
}

// WorklogPage - JIRA worklog collection page
type WorklogPage struct {
	StartAt    int       `json:"startAt"`
	MaxResults int       `json:"maxResults"`
	Total      int       `json:"total"`
	Worklogs   []Worklog `json:"worklogs"`
}
