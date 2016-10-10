package entities

type Project struct {
	ID            ProjectID
	Title         string
	Link          string
	Description   string
	IssueTypes    []TypeID
	ActivityTypes []TypeID
}

type Tracker struct {
	ID          int64
	URL         string
	Type        string
	Credentials Credentials
}

type Credentials struct {
	Login    string
	Password string
}

type TypeID struct {
	ID   int64
	Name string
}

type Issue struct {
	ID          IssueID
	ProjectID   ProjectID `json:"-"`
	Type        TypeID
	Title       string
	Description string
	Estimate    int64
	DueDate     int64
	Done        Progress
	Spent       int64
	URL         string
}

type NewIssue struct {
	Issue
	// Be carefule with NewIssue type
	// Type field inside Issue will be empty
	// We have only id of Type with NewIssue
	Type int64
}

type User struct {
	ID   string
	Name string
	Mail string
}

type Report struct {
	IssueID    IssueID
	ActivityID int64
	Comments   string
	Duration   int64 //In seconds
	Started    int64
}

type Pagination struct {
	Offset int
	Limit  int
}

type ProjectID string
type IssueID string
type Progress int
type IssueURL string
