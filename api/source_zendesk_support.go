package api

import (
	"encoding/json"
	"fmt"
)

type SourceZendeskSupportID struct {
	SourceId string `json:"sourceId"`
}

type SourceZendeskSupport struct {
	Name                    string                         `json:"name"`
	SourceId                string                         `json:"sourceId,omitempty"`
	WorkspaceId             string                         `json:"workspaceId"`
	ConnectionConfiguration SourceZendeskSupportConnConfig `json:"configuration"`
}

type SourceZendeskSupportConnConfig struct {
	SourceType      string                         `json:"sourceType"`
	StartDate       string                         `json:"start_date"`
	IgnorPagination bool                           `json:"ignore_pagination,omitempty"`
	Subdomain       string                         `json:"subdomain"`
	Credentials     SourceZendeskSupportCredConfig `json:"credentials"`
}

type SourceZendeskSupportCredConfig struct {
	Credentials string `json:"credentials"`
	Email       string `json:"email"`
	ApiToken    string `json:"api_token"`
}

func (c *Client) CreateZendeskSupportSource(payload SourceZendeskSupport) (SourceZendeskSupport, error) {
	// logger := fwhelpers.GetLogger()
	method := "POST"
	url := c.Host + "/v1/sources"
	body, err := json.Marshal(payload)
	if err != nil {
		return SourceZendeskSupport{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourceZendeskSupport{}, err
	}
	source := SourceZendeskSupport{}

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

func (c *Client) ReadZendeskSupportSource(sourceId string) (SourceZendeskSupport, error) {
	// logger := fwhelpers.GetLogger()

	method := "GET"
	url := c.Host + "/v1/sources/" + sourceId
	// sId := SourceZendeskSupportID{sourceId}
	// body, err := json.Marshal(sId)
	// if err != nil {
	// 	return SourceZendeskSupport{}, err
	// }

	b, statusCode, _, _, err := c.doRequest(method, url, []byte{}, nil)
	if err != nil {
		return SourceZendeskSupport{}, err
	}

	source := SourceZendeskSupport{}
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

func (c *Client) UpdateZendeskSupportSource(payload SourceZendeskSupport) (SourceZendeskSupport, error) {
	// logger := fwhelpers.GetLogger()

	/*method := "PUT"
	url := c.Host + "/v1/sources/" + payload.SourceId
	body, err := json.Marshal(payload)
	if err != nil {
		return SourceZendeskSupport{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourceZendeskSupport{}, err
	}

	source := SourceZendeskSupport{}
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
	// sId := SourceZendeskSupportID{sourceId}
	// body, err := json.Marshal(sId)
	// if err != nil {
	// 	return SourceZendeskSupport{}, err
	// }
	b, statusCode, _, _, err := c.doRequest(method, url, []byte{}, nil)
	if err != nil {
		return SourceZendeskSupport{}, err
	}

	source := SourceZendeskSupport{}
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

func (c *Client) DeleteZendeskSupportSource(sourceId string) error {
	// logger := fwhelpers.GetLogger()

	method := "DELETE"
	url := c.Host + "/v1/sources/" + sourceId
	sId := SourceZendeskSupportID{sourceId}
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
