package api

import (
	"encoding/json"
	"fmt"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
)

type SourceGoogleSheetsID struct {
	SourceId string `json:"sourceId"`
}

type SourceGoogleSheets struct {
	Name                    string                       `json:"name"`
	SourceId                string                       `json:"sourceId,omitempty"`
	WorkspaceId             string                       `json:"workspaceId"`
	ConnectionConfiguration SourceGoogleSheetsConnConfig `json:"configuration"`
}

type SourceGoogleSheetsConnConfig struct {
	SourceType    string                      `json:"sourceType"`
	RowBatchSize  int                         `json:"row_batch_size,omitempty"`
	SpreadsheetId string                      `json:"spreadsheet_id"`
	Credentials   GoogleSheetsCredConfigModel `json:"credentials"`
}
type GoogleSheetsCredConfigModel struct {
	AuthType           string `json:"auth_type"`
	ServiceAccountInfo string `json:"service_account_info"`
}

func (c *Client) CreateGoogleSheetsSource(payload SourceGoogleSheets) (SourceGoogleSheets, error) {
	// logger := fwhelpers.GetLogger()
	method := "POST"
	url := c.Host + "/v1/sources"
	body, err := json.Marshal(payload)
	if err != nil {
		return SourceGoogleSheets{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)
	if err != nil {
		return SourceGoogleSheets{}, err
	}
	source := SourceGoogleSheets{}

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

func (c *Client) ReadGoogleSheetsSource(sourceId string) (SourceGoogleSheets, error) {
	// logger := fwhelpers.GetLogger()

	method := "GET"
	url := c.Host + "/v1/sources/" + sourceId
	// sId := SourceGoogleSheetsID{sourceId}
	// body, err := json.Marshal(sId)
	// if err != nil {
	// 	return SourceGoogleSheets{}, err
	// }

	b, statusCode, _, _, err := c.doRequest(method, url, []byte{}, nil)
	if err != nil {
		return SourceGoogleSheets{}, err
	}

	source := SourceGoogleSheets{}
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

func (c *Client) UpdateGoogleSheetsSource(payload SourceGoogleSheets) (SourceGoogleSheets, error) {
	logger := fwhelpers.GetLogger()

	logger.Print("[yellow]Update api is not yet exposed from Airbyte-Cloud[reset]")
	return SourceGoogleSheets{}, nil
}

func (c *Client) DeleteGoogleSheetsSource(sourceId string) error {
	// logger := fwhelpers.GetLogger()

	method := "DELETE"
	url := c.Host + "/v1/sources/" + sourceId
	sId := SourceGoogleSheetsID{sourceId}
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
