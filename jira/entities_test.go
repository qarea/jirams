package jira

import (
	"testing"

	"github.com/qarea/jirams/entities"
	"github.com/stretchr/testify/assert"
)

func TestIssue(t *testing.T) {

	var (
		test = Issue{
			ID:  "1",
			Key: "TP-1",
			URL: "http://domain.com:80/rest/api/2/issue/1",
			Fields: IssueFields{
				Type:     NamedID{ID: "2", Name: "Test"},
				Title:    "Test issue",
				Estimate: 3600,
				Spent:    1800,
				DueDate:  "2016-12-25",
				Progress: IssueProgress{
					Percent: 50,
				},
			},
		}
		expected = entities.Issue{
			ID:       1,
			Type:     entities.NamedID{ID: 2, Name: "Test"},
			URL:      "http://domain.com:80/browse/TP-1",
			Title:    "Test issue",
			Estimate: 3600,
			DueDate:  1482624000,
			Spent:    1800,
			Done:     50,
		}
	)

	assert.Equal(t, expected, test.toIssue())
}
