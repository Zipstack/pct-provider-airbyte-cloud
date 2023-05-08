package api

import (
	"encoding/json"
	"fmt"
)

type SourceFreshdeskID struct {
	SourceId string `json:"sourceId"`
}

type SourceFreshdesk struct {
	Name                    string                    `json:"name"`
	SourceId                string                    `json:"sourceId,omitempty"`
	WorkspaceId             string                    `json:"workspaceId"`
	ConnectionConfiguration SourceFreshdeskConnConfig `json:"configuration"`
}

type SourceFreshdeskConnConfig struct {
	SourceType        string `json:"sourceType"`
	StartDate         string `json:"start_date"`
	Domain            string `json:"domain"`
	ApiKey            string `json:"api_key"`
	RequestsPerMinute int    `json:"requests_per_minute,omitempty"`
}

func (c *Client) CreateFreshdeskSource(payload SourceFreshdesk) (SourceFreshdesk, error) {
	// logger := fwhelpers.GetLogger()
	method := "POST"
	url := c.Host + "/v1/sources"
	body, err := json.Marshal(payload)
	if err != nil {
		return SourceFreshdesk{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourceFreshdesk{}, err
	}
	source := SourceFreshdesk{}

	if statusCode >= 200 && statusCode <= 299 {
		err = json.Unmarshal(b, &source)
		return source, err
	} else {
		msg, err := c.getAPIError(b)
		if err != nil {
			return source, err
		} else {
			return source, fmt.Errorf(msg)
		}
	}
}

func (c *Client) ReadFreshdeskSource(sourceId string) (SourceFreshdesk, error) {
	// logger := fwhelpers.GetLogger()

	method := "GET"
	url := c.Host + "/v1/sources/" + sourceId

	b, statusCode, _, _, err := c.doRequest(method, url, []byte{}, nil)
	if err != nil {
		return SourceFreshdesk{}, err
	}

	source := SourceFreshdesk{}
	if statusCode >= 200 && statusCode <= 299 {
		err = json.Unmarshal(b, &source)
		return source, err
	} else {
		msg, err := c.getAPIError(b)
		if err != nil {
			return source, err
		} else {
			return source, fmt.Errorf(msg)
		}
	}
}

func (c *Client) UpdateFreshdeskSource(payload SourceFreshdesk) (SourceFreshdesk, error) {
	// logger := fwhelpers.GetLogger()
	return SourceFreshdesk{}, fmt.Errorf("update resource is not supported")
}

func (c *Client) DeleteFreshdeskSource(sourceId string) error {
	// logger := fwhelpers.GetLogger()

	method := "DELETE"
	url := c.Host + "/v1/sources/" + sourceId
	sId := SourceFreshdeskID{sourceId}
	body, err := json.Marshal(sId)
	if err != nil {
		return err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return err
	}

	if statusCode >= 200 && statusCode <= 299 {
		return nil
	} else {
		msg, err := c.getAPIError(b)
		if err != nil {
			return err
		} else {
			return fmt.Errorf(msg)
		}
	}
}
