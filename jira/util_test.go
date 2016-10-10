package jira

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/qarea/jirams/entities"
	"github.com/stretchr/testify/assert"
)

type testUtilRequest struct {
	method       string
	url          string
	body         string
	responseCode int
	response     string
	expectError  bool
	error        error
	expected     interface{}
}

type TestEntity struct {
	Foo string
}

type TestEntityPage struct {
	StartAt    int
	MaxResults int
	Total      int
	Entities   []TestEntity
}

var testTrackerConfig = entities.TrackerConfig{
	Credentials: entities.TrackerCredentials{
		Login:    "testLogin",
		Password: "testPassword",
	},
}

func TestRequest(t *testing.T) {

	requests := map[string]testUtilRequest{
		"GET": {
			method:       "GET",
			responseCode: http.StatusOK,
			response:     `{"Foo":"bar"}`,
			expected:     TestEntity{Foo: "bar"},
		},
		"POST": {
			method:       "POST",
			body:         `{"Test":"abc"}`,
			responseCode: http.StatusOK,
			response:     `{"Foo":"bar"}`,
			expected:     TestEntity{Foo: "bar"},
		},
		"400": {
			method:       "GET",
			responseCode: http.StatusBadRequest,
			expectError:  true,
			error:        entities.ErrInvalidRequest,
		},
		"401": {
			method:       "GET",
			responseCode: http.StatusUnauthorized,
			expectError:  true,
			error:        entities.ErrInvalidCredentials,
		},
		"404": {
			method:       "GET",
			responseCode: http.StatusNotFound,
			expectError:  true,
			error:        entities.ErrNotFound,
		},
		"500": {
			method:       "GET",
			responseCode: http.StatusInternalServerError,
			expectError:  true,
			error:        entities.ErrServerUnavailable,
		},
	}

	requester := Requester{}

	for _, test := range requests {
		srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			header := res.Header()
			header["Content-Type"] = []string{"application/json"}
			res.WriteHeader(test.responseCode)
			_, _ = res.Write([]byte(test.response))
		}))

		url := srv.URL

		request, _ := http.NewRequest(test.method, url, bytes.NewBuffer([]byte(test.body)))

		var result TestEntity

		err := requester.Request(testTrackerConfig, request, &result)

		if test.expectError {
			if assert.Error(t, err) {
				assert.Equal(t, test.error, err)
			}
		} else {
			assert.Nil(t, err)
			assert.Equal(t, test.expected, result)
		}

		srv.Close()
	}
}

func TestIterateRequest(t *testing.T) {
	requests := map[string]string{
		"0": `{"startAt":0,"maxResults":1,"total":2,"entities":[{"foo":"bar"}]}`,
		"1": `{"startAt":1,"maxResults":1,"total":2,"entities":[{"foo":"baz"}]}`,
	}

	requester := Requester{}

	srv := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startAt := req.URL.Query()["startAt"][0]
		header := res.Header()
		header["Content-Type"] = []string{"application/json"}
		res.WriteHeader(http.StatusOK)
		_, _ = res.Write([]byte(requests[startAt]))
	}))

	url := srv.URL

	var (
		tpl    TestEntityPage
		result []TestEntity
	)

	err := requester.IterateRequest(testTrackerConfig, url, &tpl, func(data interface{}) (loaded int, total int, err error) {
		if data, ok := data.(*TestEntityPage); ok {
			loaded = len(data.Entities)
			total = data.Total
			result = append(result, data.Entities...)
		} else {
			err = errors.New("Incorrect type")
		}
		return
	})

	assert.Nil(t, err)
	assert.Equal(t, []TestEntity{
		{Foo: "bar"},
		{Foo: "baz"},
	}, result)
}
