package api

import (
	"encoding/json"
	"fmt"
)

type DestinationPostgresID struct {
	DestinationId string `json:"destinationId"`
}

type DestinationPostgres struct {
	Name                    string                        `json:"name"`
	DestinationId           string                        `json:"destinationId,omitempty"`
	WorkspaceId             string                        `json:"workspaceId"`
	ConnectionConfiguration DestinationPostgresConnConfig `json:"configuration"`
}

type DestinationPostgresConnConfig struct {
	DestinationType    string             `json:"destinationType"`
	Host               string             `json:"host"`
	Username           string             `json:"username"`
	Password           string             `json:"password"`
	Database           string             `json:"database"`
	Port               int                `json:"port"`
	Schema             string             `json:"schema"`
	SslModeConfig      SslModeConfig      `json:"ssl_mode"`
	TunnelMethodConfig TunnelMethodConfig `json:"tunnel_method"`
}
type SslModeConfig struct {
	Mode string `json:"mode"`
}
type TunnelMethodConfig struct {
	TunnelMethod string `json:"tunnel_method"`
}

func (c *Client) CreatePostgresDestination(payload DestinationPostgres) (DestinationPostgres, error) {
	// logger := fwhelpers.GetLogger()
	method := "POST"
	url := c.Host + "/v1/destinations"
	body, err := json.Marshal(payload)
	if err != nil {
		return DestinationPostgres{}, err
	}
	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return DestinationPostgres{}, err
	}

	destination := DestinationPostgres{}

	if statusCode >= 200 && statusCode <= 299 {
		err = json.Unmarshal(b, &destination)
		return destination, err
	} else {
		msg, err := c.getAPIError(b)
		if err != nil {
			return destination, err
		} else {
			return destination, fmt.Errorf(msg)
		}
	}
}

func (c *Client) ReadPostgresDestination(destinationId string) (DestinationPostgres, error) {
	// logger := fwhelpers.GetLogger()

	method := "GET"
	url := c.Host + "/v1/destinations/" + destinationId

	b, statusCode, _, _, err := c.doRequest(method, url, []byte{}, nil)
	if err != nil {
		return DestinationPostgres{}, err
	}

	destination := DestinationPostgres{}
	if statusCode >= 200 && statusCode <= 299 {
		err = json.Unmarshal(b, &destination)
		return destination, err
	} else {
		msg, err := c.getAPIError(b)
		if err != nil {
			return destination, err
		} else {
			return destination, fmt.Errorf(msg)
		}
	}
}

func (c *Client) UpdatePostgresDestination(payload DestinationPostgres) (DestinationPostgres, error) {
	// logger := fwhelpers.GetLogger()

	return DestinationPostgres{}, fmt.Errorf("update resource is not supported")
}

func (c *Client) DeletePostgresDestination(destinationId string) error {
	// logger := fwhelpers.GetLogger()

	method := "DELETE"
	url := c.Host + "/v1/destinations/" + destinationId
	sId := DestinationPostgresID{destinationId}
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
