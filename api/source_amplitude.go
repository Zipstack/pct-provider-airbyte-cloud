package api

import (
	"encoding/json"
	"fmt"
)

type SourceAmplitudeID struct {
	SourceId string `json:"sourceId"`
}

type SourceAmplitude struct {
	Name                    string                    `json:"name"`
	SourceId                string                    `json:"sourceId,omitempty"`
	WorkspaceId             string                    `json:"workspaceId"`
	ConnectionConfiguration SourceAmplitudeConnConfig `json:"configuration"`
}

type SourceAmplitudeConnConfig struct {
	SourceType       string `json:"sourceType"`
	StartDate        string `json:"start_date"`
	DataRegion       string `json:"data_region"`
	RequestTimeRange int    `json:"request_time_range,omitempty"`
	ApiKey           string `json:"api_key"`
	SecretKey        string `json:"secret_key"`
}

func (c *Client) CreateAmplitudeSource(payload SourceAmplitude) (SourceAmplitude, error) {
	// logger := fwhelpers.GetLogger()
	method := "POST"
	url := c.Host + "/v1/sources"
	body, err := json.Marshal(payload)
	if err != nil {
		return SourceAmplitude{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourceAmplitude{}, err
	}
	source := SourceAmplitude{}

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

func (c *Client) ReadAmplitudeSource(sourceId string) (SourceAmplitude, error) {
	// logger := fwhelpers.GetLogger()

	method := "GET"
	url := c.Host + "/v1/sources/" + sourceId

	b, statusCode, _, _, err := c.doRequest(method, url, []byte{}, nil)
	if err != nil {
		return SourceAmplitude{}, err
	}

	source := SourceAmplitude{}
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

func (c *Client) UpdateAmplitudeSource(payload SourceAmplitude) (SourceAmplitude, error) {
	// logger := fwhelpers.GetLogger()

	method := "PUT"
	url := c.Host + "/v1/sources/" + payload.SourceId
	body, err := json.Marshal(payload)
	if err != nil {
		return SourceAmplitude{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourceAmplitude{}, err
	}

	source := SourceAmplitude{}
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

func (c *Client) DeleteAmplitudeSource(sourceId string) error {
	// logger := fwhelpers.GetLogger()

	method := "DELETE"
	url := c.Host + "/v1/sources/" + sourceId
	sId := SourceAmplitudeID{sourceId}
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
