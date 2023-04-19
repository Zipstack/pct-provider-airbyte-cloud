package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

type APIError struct {
	Message            string                    `json:"message"`
	ExceptionClassName string                    `json:"exceptionClassName"`
	ExceptionStack     []string                  `json:"exceptionStack"`
	ValidationErrors   []APIErrorValidationError `json:"validationErrors"`
}

type APIErrorValidationError struct {
	PropertyPath string `json:"propertyPath"`
	InvalidValue string `json:"invalidValue"`
	Message      string `json:"message"`
}

func (c *Client) getAPIError(body []byte) (string, error) {
	apiErr := APIError{}
	err := json.Unmarshal(body, &apiErr)
	if err != nil {
		return "", fmt.Errorf("content type mismatch or invalid provider api host or path")
	} else {
		slices := strings.Split(apiErr.Message, "at [Source:")
		return strings.TrimSpace(slices[0]), nil
	}
}
