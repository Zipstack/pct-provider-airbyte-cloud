package api

import (
	"encoding/json"
	"fmt"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
)

type DestinationMysqlID struct {
	DestinationId string `json:"destinationId"`
}

type DestinationMysql struct {
	Name                    string                     `json:"name"`
	DestinationId           string                     `json:"destinationId,omitempty"`
	WorkspaceId             string                     `json:"workspaceId"`
	ConnectionConfiguration DestinationMysqlConnConfig `json:"configuration"`
}

type DestinationMysqlConnConfig struct {
	DestinationType string `json:"destinationType"`
	Host            string `json:"host"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	Database        string `json:"database"`
	Port            int    `json:"port"`
}

func (c *Client) CreateMysqlDestination(payload DestinationMysql) (DestinationMysql, error) {
	// logger := fwhelpers.GetLogger()
	method := "POST"
	url := c.Host + "/v1/destinations"
	body, err := json.Marshal(payload)
	if err != nil {
		return DestinationMysql{}, err
	}
	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return DestinationMysql{}, err
	}

	destination := DestinationMysql{}

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

func (c *Client) ReadMysqlDestination(destinationId string) (DestinationMysql, error) {
	// logger := fwhelpers.GetLogger()

	method := "GET"
	url := c.Host + "/v1/destinations/" + destinationId

	b, statusCode, _, _, err := c.doRequest(method, url, []byte{}, nil)
	if err != nil {
		return DestinationMysql{}, err
	}

	destination := DestinationMysql{}
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

func (c *Client) UpdateMysqlDestination(payload DestinationMysql) (DestinationMysql, error) {
	logger := fwhelpers.GetLogger()

	logger.Print("[yellow]Update api is not yet exposed from Airbyte-Cloud[reset]")
	return DestinationMysql{}, nil
}

func (c *Client) DeleteMysqlDestination(destinationId string) error {
	// logger := fwhelpers.GetLogger()

	method := "DELETE"
	url := c.Host + "/v1/destinations/" + destinationId
	sId := DestinationMysqlID{destinationId}
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
