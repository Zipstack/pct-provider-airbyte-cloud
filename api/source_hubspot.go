package api

import (
	"encoding/json"
	"fmt"
)

type SourceHubspotID struct {
	SourceId string `json:"sourceId"`
}

type SourceHubspot struct {
	Name                    string                  `json:"name"`
	SourceId                string                  `json:"sourceId,omitempty"`
	WorkspaceId             string                  `json:"workspaceId"`
	ConnectionConfiguration SourceHubspotConnConfig `json:"configuration"`
}

type SourceHubspotConnConfig struct {
	SourceType  string                 `json:"sourceType"`
	StartDate   string                 `json:"start_date"`
	Credentials HubspotCredConfigModel `json:"credentials"`
}
type HubspotCredConfigModel struct {
	CredentialsTitle string `json:"credentials_title"`
	AccessToken      string `json:"access_token"`
}

func (c *Client) CreateHubspotSource(payload SourceHubspot) (SourceHubspot, error) {
	// logger := fwhelpers.GetLogger()
	method := "POST"
	url := c.Host + "/v1/sources"
	body, err := json.Marshal(payload)
	if err != nil {
		return SourceHubspot{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourceHubspot{}, err
	}
	source := SourceHubspot{}

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

func (c *Client) ReadHubspotSource(sourceId string) (SourceHubspot, error) {
	// logger := fwhelpers.GetLogger()

	method := "GET"
	url := c.Host + "/v1/sources/" + sourceId
	// sId := SourceHubspotID{sourceId}
	// body, err := json.Marshal(sId)
	// if err != nil {
	// 	return SourceHubspot{}, err
	// }

	b, statusCode, _, _, err := c.doRequest(method, url, []byte{}, nil)
	if err != nil {
		return SourceHubspot{}, err
	}

	source := SourceHubspot{}
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

func (c *Client) UpdateHubspotSource(payload SourceHubspot) (SourceHubspot, error) {
	// logger := fwhelpers.GetLogger()

	/*method := "PUT"
	url := c.Host + "/v1/sources/" + payload.SourceId
	body, err := json.Marshal(payload)
	if err != nil {
		return SourceHubspot{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourceHubspot{}, err
	}

	source := SourceHubspot{}
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
	}*/
	method := "GET"
	url := c.Host + "/v1/sources/" + payload.SourceId
	// sId := SourceHubspotID{sourceId}
	// body, err := json.Marshal(sId)
	// if err != nil {
	// 	return SourceHubspot{}, err
	// }
	b, statusCode, _, _, err := c.doRequest(method, url, []byte{}, nil)
	if err != nil {
		return SourceHubspot{}, err
	}

	source := SourceHubspot{}
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

func (c *Client) DeleteHubspotSource(sourceId string) error {
	// logger := fwhelpers.GetLogger()

	method := "DELETE"
	url := c.Host + "/v1/sources/" + sourceId
	sId := SourceHubspotID{sourceId}
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
