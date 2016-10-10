package entities

// EntityID - generic entity ID type
type EntityID uint64

// Project - TG project entity
type Project struct {
	ID            ProjectID
	Title         string
	Description   string
	IssueTypes    []NamedID
	ActivityTypes []NamedID
}

// ProjectID - project ID
type ProjectID uint64

// NamedID - id + name entity structure
type NamedID struct {
	ID   uint64
	Name string
}

// TrackerConfig - TG tracker configuration
type TrackerConfig struct {
	ID          TrackerID
	URL         string
	Credentials TrackerCredentials
}

// TrackerID - tracker id
type TrackerID uint64

// TrackerCredentials - TG tracker credentials
type TrackerCredentials struct {
	Login    string
	Password string
}

// User - TG user
type User struct {
	ID   UserID
	Name string
	Mail string
}

// UserID - user ID
type UserID uint64

// UserKey - JIRA user indetifier string
type UserKey string

// Issue - TG issue
type Issue struct {
	ID       IssueID
	Type     NamedID
	URL      string
	Title    string
	Estimate Duration
	DueDate  Timestamp
	Spent    Duration
	Done     uint8
}

// IssueID - issue ID
type IssueID uint64

// NewIssue - set of parameters for new issue
type NewIssue struct {
	ProjectID ProjectID
	Assignee  UserID
	Type      EntityID
	Title     string
	Estimate  uint64
}

// Report - set of parameters for new work time report
type Report struct {
	IssueID  IssueID
	Started  Timestamp
	Duration Duration
	Comments string
}

// Timestamp - unix timestamp
type Timestamp uint64

// Duration - time duration in seconds
type Duration uint64

// ReportsTotal - total amount of reported time in seconds
type ReportsTotal uint64
