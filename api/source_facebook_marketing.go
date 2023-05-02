package api

import (
	"encoding/json"
	"fmt"

	"github.com/zipstack/pct-plugin-framework/fwhelpers"
)

type SourceFacebookMarketingID struct {
	SourceId string `json:"sourceId"`
}

type SourceFacebookMarketing struct {
	Name                    string                            `json:"name"`
	SourceId                string                            `json:"sourceId,omitempty"`
	WorkspaceId             string                            `json:"workspaceId"`
	ConnectionConfiguration SourceFacebookMarketingConnConfig `json:"configuration"`
}

type SourceFacebookMarketingConnConfig struct {
	SourceType  string `json:"sourceType"`
	AccountId   string `json:"account_id"`
	StartDate   string `json:"start_date"`
	AccessToken string `json:"access_token"`

	EndDate              string `json:"end_date,omitempty"`
	IncludeDeleted       bool   `json:"include_deleted,omitempty"`
	FetchThumbnailImages bool   `json:"fetch_thumbnail_images,omitempty"`
	//CustomInsights       any    `json:"custom_insights"`
	PageSize                   int  `json:"page_size,omitempty"`
	InsightsLookbackWindow     int  `json:"insights_lookback_window,omitempty"`
	MaxBatchSize               int  `json:"max_batch_size,omitempty"`
	ActionBreakdownsAllowEmpty bool `json:"action_breakdowns_allow_empty,omitempty"`
}

func (c *Client) CreateFacebookMarketingSource(payload SourceFacebookMarketing) (SourceFacebookMarketing, error) {
	// logger := fwhelpers.GetLogger()
	method := "POST"
	url := c.Host + "/v1/sources"
	body, err := json.Marshal(payload)
	if err != nil {
		return SourceFacebookMarketing{}, err
	}

	b, statusCode, _, _, err := c.doRequest(method, url, body, nil)

	if err != nil {
		return SourceFacebookMarketing{}, err
	}
	source := SourceFacebookMarketing{}

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

func (c *Client) ReadFacebookMarketingSource(sourceId string) (SourceFacebookMarketing, error) {
	// logger := fwhelpers.GetLogger()

	method := "GET"
	url := c.Host + "/v1/sources/" + sourceId

	b, statusCode, _, _, err := c.doRequest(method, url, []byte{}, nil)
	if err != nil {
		return SourceFacebookMarketing{}, err
	}

	source := SourceFacebookMarketing{}
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

func (c *Client) UpdateFacebookMarketingSource(payload SourceFacebookMarketing) (SourceFacebookMarketing, error) {
	logger := fwhelpers.GetLogger()

	logger.Print("[yellow]Update api is not yet exposed from Airbyte-Cloud[reset]")
	return SourceFacebookMarketing{}, nil
}

func (c *Client) DeleteFacebookMarketingSource(sourceId string) error {
	// logger := fwhelpers.GetLogger()

	method := "DELETE"
	url := c.Host + "/v1/sources/" + sourceId
	sId := SourceFacebookMarketingID{sourceId}
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
