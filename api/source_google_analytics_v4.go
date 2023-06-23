package api

import (
	"encoding/json"
	"fmt"
)

type SourceGoogleAnalyticsV4ID struct {
	SourceId string `json:"sourceId"`
}

type SourceGoogleAnalyticsV4 struct {
	Name                    string                            `json:"name"`
	SourceId                string                            `json:"sourceId,omitempty"`
	WorkspaceId             string                            `json:"workspaceId"`
	ConnectionConfiguration SourceGoogleAnalyticsV4ConnConfig `json:"configuration"`
}

type SourceGoogleAnalyticsV4ConnConfig struct {
	SourceType    string                           `json:"sourceType"`
	StartDate     string                           `json:"start_date"`
	ViewId        string                           `json:"view_id,omitempty"`
	CustomReports string                           `json:"custom_reports"`
	WindowInDays  int                              `json:"window_in_days,omitempty"`
	Credentials   GoogleAnalyticsV4CredConfigModel `json:"credentials"`
}
type GoogleAnalyticsV4CredConfigModel struct {
	AuthType        string `json:"auth_type"`
	CredentialsJson string `json:"credentials_json"`
}

func (c *Client) CreateGoogleAnalyticsV4Source(payload SourceGoogleAnalyticsV4) (SourceGoogleAnalyticsV4, error) {
	// logger := fwhelpers.GetLogger()
	method := "POST"
	url := c.Host + "/v1/sources"
	body, err := json.Marshal(payload)
	if err != nil {
		return SourceGoogleAnalyticsV4{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourceGoogleAnalyticsV4{}, err
	}
	source := SourceGoogleAnalyticsV4{}

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

func (c *Client) ReadGoogleAnalyticsV4Source(sourceId string) (SourceGoogleAnalyticsV4, error) {
	// logger := fwhelpers.GetLogger()

	method := "GET"
	url := c.Host + "/v1/sources/" + sourceId

	b, statusCode, _, _, err := c.doRequest(method, url, []byte{}, nil)
	if err != nil {
		return SourceGoogleAnalyticsV4{}, err
	}

	source := SourceGoogleAnalyticsV4{}
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

func (c *Client) UpdateGoogleAnalyticsV4Source(payload SourceGoogleAnalyticsV4) (SourceGoogleAnalyticsV4, error) {
	// logger := fwhelpers.GetLogger()
	method := "PUT"
	url := c.Host + "/v1/sources/" + payload.SourceId
	body, err := json.Marshal(payload)
	if err != nil {
		return SourceGoogleAnalyticsV4{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourceGoogleAnalyticsV4{}, err
	}

	source := SourceGoogleAnalyticsV4{}
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

func (c *Client) DeleteGoogleAnalyticsV4Source(sourceId string) error {
	// logger := fwhelpers.GetLogger()

	method := "DELETE"
	url := c.Host + "/v1/sources/" + sourceId
	sId := SourceGoogleAnalyticsV4ID{sourceId}
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
