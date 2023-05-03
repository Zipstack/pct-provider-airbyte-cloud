package api

import (
	"encoding/json"
	"fmt"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
)

type SourceShopifyID struct {
	SourceId string `json:"sourceId"`
}

type SourceShopify struct {
	Name                    string                  `json:"name"`
	SourceId                string                  `json:"sourceId,omitempty"`
	WorkspaceId             string                  `json:"workspaceId"`
	ConnectionConfiguration SourceShopifyConnConfig `json:"configuration"`
}

type SourceShopifyConnConfig struct {
	SourceType  string                 `json:"sourceType"`
	StartDate   string                 `json:"start_date"`
	Shop        string                 `json:"shop"`
	Credentials ShopifyCredConfigModel `json:"credentials"`
}
type ShopifyCredConfigModel struct {
	AuthMethod  string `json:"auth_method"`
	ApiPassword string `json:"api_password"`
	// ClientSecret string `json:"client_secret,omitempty"`
	// AccessToken  string `json:"access_token,omitempty"`
	// ClientId     string `json:"client_id,omitempty"`
}

func (c *Client) CreateShopifySource(payload SourceShopify) (SourceShopify, error) {
	// logger := fwhelpers.GetLogger()
	method := "POST"
	url := c.Host + "/v1/sources"
	body, err := json.Marshal(payload)
	if err != nil {
		return SourceShopify{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourceShopify{}, err
	}
	source := SourceShopify{}
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

func (c *Client) ReadShopifySource(sourceId string) (SourceShopify, error) {
	// logger := fwhelpers.GetLogger()

	method := "GET"
	url := c.Host + "/v1/sources/" + sourceId

	b, statusCode, _, _, err := c.doRequest(method, url, []byte{}, nil)
	if err != nil {
		return SourceShopify{}, err
	}

	source := SourceShopify{}
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

func (c *Client) UpdateShopifySource(payload SourceShopify) (SourceShopify, error) {
	logger := fwhelpers.GetLogger()

	logger.Print("[yellow]Update api is not yet exposed from Airbyte-Cloud[reset]")
	return SourceShopify{}, nil
}

func (c *Client) DeleteShopifySource(sourceId string) error {
	// logger := fwhelpers.GetLogger()

	method := "DELETE"
	url := c.Host + "/v1/sources/" + sourceId
	sId := SourceShopifyID{sourceId}
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
