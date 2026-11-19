package bucket

import (
	"encoding/xml"
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

type APIError struct {
	StatusCode int
	Code       string
	Message    string
	RequestID  string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("Object Storage API error %d: %s: %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("Object Storage API error %d: %s", e.StatusCode, e.Message)
}

func IsNotFound(err error) bool {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr.Code == "NoSuchKey" || apiErr.Code == "NoSuchBucket" || apiErr.StatusCode == 404
	}
	return errors.Is(err, ErrNotFound)
}

func parseAPIError(op string, statusCode int, body []byte) error {
	var apiErr struct {
		XMLName   xml.Name `xml:"Error"`
		Code      string   `xml:"Code"`
		Message   string   `xml:"Message"`
		RequestID string   `xml:"RequestId"`
	}
	if err := xml.Unmarshal(body, &apiErr); err == nil && apiErr.Code != "" {
		return &APIError{
			StatusCode: statusCode,
			Code:       apiErr.Code,
			Message:    apiErr.Message,
			RequestID:  apiErr.RequestID,
		}
	}
	return &APIError{
		StatusCode: statusCode,
		Message:    fmt.Sprintf("%s: unexpected response", op),
	}
}
