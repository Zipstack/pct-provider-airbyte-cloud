package api

import (
	"encoding/json"
	"fmt"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
)

type ConnectionResourceID struct {
	ConnecctionID string `json:"connectionId"`
}

type ConnectionResource struct {
	Name          string `json:"name"`
	SourceID      string `json:"sourceId,omitempty"`
	DestinationID string `json:"destinationId,omitempty"`
	ConnectionID  string `json:"connectionId,omitempty"`
	// SyncCatalog                      interface{}      `json:"syncCatalog"`
	DataResidency                    string           `json:"dataResidency,omitempty"`
	NamespaceDefinition              string           `json:"namespaceDefinition,omitempty"`
	NamespaceFormat                  string           `json:"namespaceFormat,omitempty"`
	NonBreakingSchemaUpdatesBehavior string           `json:"nonBreakingSchemaUpdatesBehavior,omitempty"`
	Prefix                           string           `json:"prefix,omitempty"`
	Status                           string           `json:"status,omitempty"`
	Schedule                         ConnScheduleData `json:"schedule"`
	//OperatorConfiguration connOperatorConfig `json:"operator_configuration"`
}
type ConnScheduleData struct {
	ScheduleType   string `json:"scheduleType"`
	CronExpression string `json:"cronExpression,omitempty"`
}

type DiscoverSourceSchemaCatalog struct {
	SourceID     string `json:"sourceId"`
	DisableCache bool   `json:"disable_cache"`
}

func (c *Client) CreateConnectionResource(payload ConnectionResource) (ConnectionResource, error) {
	// logger := fwhelpers.GetLogger()

	method := "POST"
	url := c.Host + "/v1/connections"
	body, err := json.Marshal(payload)
	if err != nil {
		return ConnectionResource{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return ConnectionResource{}, err
	}

	connection := ConnectionResource{}
	if statusCode >= 200 && statusCode <= 299 {
		err = json.Unmarshal(b, &connection)
		return connection, err
	} else {
		msg, err := c.getAPIError(b)
		if err != nil {
			return connection, err
		} else {
			return connection, fmt.Errorf(msg)
		}
	}
}

func (c *Client) ReadConnectionResource(connectionId string) (ConnectionResource, error) {
	// logger := fwhelpers.GetLogger()

	method := "GET"
	url := c.Host + "/v1/connections/" + connectionId

	b, statusCode, _, _, err := c.doRequest(method, url, []byte{}, nil)
	if err != nil {
		return ConnectionResource{}, err
	}
	connection := ConnectionResource{}
	if statusCode >= 200 && statusCode <= 299 {
		err = json.Unmarshal(b, &connection)
		return connection, err
	} else {
		msg, err := c.getAPIError(b)
		if err != nil {
			return connection, err
		} else {
			return connection, fmt.Errorf(msg)
		}
	}
}

func (c *Client) UpdateConnectionResource(payload ConnectionResource) (ConnectionResource, error) {
	logger := fwhelpers.GetLogger()

	logger.Print("[yellow]Update api is not yet exposed from Airbyte-Cloud[reset]")
	return ConnectionResource{}, nil
}

func (c *Client) DeleteConnectionResource(connectionId string) error {
	// logger := fwhelpers.GetLogger()

	method := "DELETE"
	url := c.Host + "/v1/connections"

	b, statusCode, _, _, err := c.doRequest(method, url, []byte{}, nil)
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
