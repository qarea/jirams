package jira

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/qarea/jirams/entities"
)

// TrackerRequester interface defines capability to make one or iteration of requests to the tracker
type TrackerRequester interface {
	Request(entities.TrackerConfig, *http.Request, interface{}) error
	IterateRequest(entities.TrackerConfig, string, interface{}, func(interface{}) (int, int, error)) error
}

// Requester implements TrackerRequester for JIRA tracker
type Requester struct{}

// Request performs a request to specified JIRA API URL and unmarshals the response to data structure
func (requester *Requester) Request(tracker entities.TrackerConfig, request *http.Request, res interface{}) error {
	httpClient := &http.Client{}
	request.SetBasicAuth(tracker.Credentials.Login, tracker.Credentials.Password)

	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	if response.StatusCode >= 400 {
		switch response.StatusCode {
		case 401:
			return entities.ErrInvalidCredentials
		case 404:
			return entities.ErrNotFound
		default:
			if response.StatusCode >= 500 {
				return entities.ErrServerUnavailable
			}
			l.ERR(response.Status)
			return entities.ErrInvalidRequest
		}
	}

	if res != nil {
		if err = json.Unmarshal(body, res); err != nil {
			l.ERR(err.Error())
			return err
		}
	}

	return nil
}

// IterateRequest performs requests to specified URL until all items are retrieved
// Each chunk of entities in passed to the callback function
func (requester *Requester) IterateRequest(tracker entities.TrackerConfig, url string, dataContainer interface{}, callback func(interface{}) (int, int, error)) error {
	var (
		startAt  = 0
		finished = false
		glue     = "?"
	)
	for !finished {
		paginationParams := fmt.Sprintf("startAt=%d", startAt)
		if strings.Contains(url, "?") {
			glue = "&"
		}
		request, _ := http.NewRequest("GET", url+glue+paginationParams, nil)
		err := requester.Request(tracker, request, dataContainer)
		if err != nil {
			return err
		}
		loaded, total, err := callback(dataContainer)
		if err != nil {
			return err
		}
		startAt += loaded
		finished = startAt >= total
	}

	return nil
}

func validateTrackerConfig(cfg entities.TrackerConfig) (err error) {
	switch {
	case cfg.URL == "":
		err = entities.ErrInvalidTrackerURL
	case cfg.Credentials.Login == "":
		err = entities.ErrInvalidCredentials
	case cfg.Credentials.Password == "":
		err = entities.ErrInvalidCredentials
	}
	return
}
