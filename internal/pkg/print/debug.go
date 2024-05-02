package print

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"sort"
	"strings"

	"github.com/stackitcloud/stackit-sdk-go/core/config"
)

var defaultDebugIncludeHeaders = []string{"Accept", "Content-Type", "Content-Length", "User-Agent", "Date", "Referrer-Policy"}

// BuildDebugStrFromInputModel converts an input model to a user-friendly string representation.
// This function converts the input model to a map, removes empty values, and generates a string representation of the map.
// The purpose of this function is to provide a more readable output than the default JSON representation.
// It is particularly useful when outputting to the slog logger, as the JSON format with escaped quotes does not look good.
func BuildDebugStrFromInputModel(model any) (string, error) {
	// Marshaling and Unmarshaling is the best way to convert the struct to a map
	modelBytes, err := json.Marshal(model)
	if err != nil {
		return "", fmt.Errorf("marshal model to JSON: %w", err)
	}

	var inputModelMap map[string]any
	if err := json.Unmarshal(modelBytes, &inputModelMap); err != nil {
		return "", fmt.Errorf("unmarshaling JSON to map: %w", err)
	}

	return BuildDebugStrFromMap(inputModelMap), nil
}

// BuildDebugStrFromMap converts a map to a user-friendly string representation.
// This function removes empty values and generates a string representation of the map.
// The string representation is in the format: [key1: value1, key2: value2, ...]
// The keys are ordered alphabetically to make the output deterministic.
func BuildDebugStrFromMap(inputMap map[string]any) string {
	// Sort the keys to make the output deterministic
	keys := make([]string, 0, len(inputMap))
	for key := range inputMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var keyValues []string
	for _, key := range keys {
		value := inputMap[key]
		if isEmpty(value) {
			continue
		}

		if valueMap, ok := value.(map[string]any); ok {
			value = BuildDebugStrFromMap(valueMap)
		}

		// If the value is a slice, convert it to a string representation
		if valueSlice, ok := value.([]any); ok {
			sliceStr := make([]string, len(valueSlice))
			for i, item := range valueSlice {
				if itemMap, ok := item.(map[string]any); ok {
					sliceStr[i] = BuildDebugStrFromMap(itemMap)
				} else {
					sliceStr[i] = fmt.Sprintf("%v", item)
				}
			}
			value = BuildDebugStrFromSlice(sliceStr)
		}
		keyValues = append(keyValues, fmt.Sprintf("%s: %v", key, value))
	}

	result := strings.Join(keyValues, ", ")
	return fmt.Sprintf("[%s]", result)
}

// BuildDebugStrFromSlice converts a slice to a user-friendly string representation.
func BuildDebugStrFromSlice(inputSlice []string) string {
	sliceStr := strings.Join(inputSlice, ", ")
	return fmt.Sprintf("[%s]", sliceStr)
}

// buildHeaderMap converts a map to a user-friendly string representation.
// This function also filters the headers based on the includeHeaders parameter.
// If includeHeaders is empty, the default header filters are used.
func buildHeaderMap(headers http.Header, includeHeaders []string) map[string]any {
	headersMap := make(map[string]any)
	for key, values := range headers {
		headersMap[key] = strings.Join(values, ", ")
	}

	var headersToInclude []string

	if len(includeHeaders) == 0 {
		headersToInclude = defaultDebugIncludeHeaders
	} else {
		headersToInclude = includeHeaders
	}

	for key := range headersMap {
		if slices.Contains(headersToInclude, key) {
			continue
		}
		delete(headersMap, key)
	}

	return headersMap
}

// BuildDebugStrFromHTTPRequest converts an HTTP request to a user-friendly string representation.
// This function also receives a list of headers to include in the output, if empty, the default headers are used.
// The return value is a list of strings that should be printed separately.
func BuildDebugStrFromHTTPRequest(req *http.Request, includeHeaders []string) ([]string, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}
	if req.URL == nil || req.Proto == "" || req.Method == "" {
		return nil, fmt.Errorf("request is invalid")
	}

	status := fmt.Sprintf("request to %s: %s %s", req.URL, req.Method, req.Proto)

	headersMap := buildHeaderMap(req.Header, includeHeaders)
	headers := fmt.Sprintf("request headers: %v", BuildDebugStrFromMap(headersMap))

	if req.Body == nil {
		return []string{status, headers}, nil
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return []string{status, headers}, fmt.Errorf("read request body: %w", err)
	}
	var bodyMap map[string]any
	if len(body) != 0 {
		if err := json.Unmarshal(body, &bodyMap); err != nil {
			return []string{status, headers}, fmt.Errorf("unmarshal request body: %w", err)
		}
	}

	// restore body
	req.Body = io.NopCloser(strings.NewReader(string(body)))
	payload := fmt.Sprintf("request body: %v", BuildDebugStrFromMap(bodyMap))

	return []string{status, headers, payload}, nil
}

// BuildDebugStrFromHTTPResponse converts an HTTP response to a user-friendly string representation.
// This function also receives a list of headers to include in the output, if empty, the default headers are used.
// The return value is a list of strings that should be printed separately.
func BuildDebugStrFromHTTPResponse(resp *http.Response, includeHeaders []string) ([]string, error) {
	if resp == nil {
		return nil, fmt.Errorf("response is nil")
	}

	if resp.Request == nil || resp.Proto == "" || resp.Status == "" {
		return nil, fmt.Errorf("response is invalid")
	}

	status := fmt.Sprintf("response from %s: %s %s", resp.Request.URL, resp.Proto, resp.Status)

	headersMap := buildHeaderMap(resp.Header, includeHeaders)
	headers := fmt.Sprintf("response headers: %v", BuildDebugStrFromMap(headersMap))

	if resp.Body == nil {
		return []string{status, headers}, nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []string{status, headers}, fmt.Errorf("read response body: %w", err)
	}
	var bodyMap map[string]any
	if len(body) != 0 {
		if err := json.Unmarshal(body, &bodyMap); err != nil {
			return []string{status, headers}, fmt.Errorf("unmarshal response body: %w", err)
		}
	}

	// restore body
	resp.Body = io.NopCloser(strings.NewReader(string(body)))
	payload := fmt.Sprintf("response body: %v", BuildDebugStrFromMap(bodyMap))

	return []string{status, headers, payload}, nil
}

// RequestResponseCapturer is a middleware that captures the request and response of an HTTP request.
// Receives a printer and a list of headers to include in the output
// If the list of headers is empty, the default headers are used.
// The printer is used to print the captured data.
func RequestResponseCapturer(p *Printer, includeHeaders []string) config.Middleware {
	return func(rt http.RoundTripper) http.RoundTripper {
		return &roundTripperWithCapture{rt, p, includeHeaders}
	}
}

type roundTripperWithCapture struct {
	transport    http.RoundTripper
	p            *Printer
	debugHeaders []string
}

func (rt roundTripperWithCapture) RoundTrip(req *http.Request) (*http.Response, error) {
	reqStr, err := BuildDebugStrFromHTTPRequest(req, rt.debugHeaders)
	if err != nil {
		rt.p.Debug(ErrorLevel, "printing request to debug logs: %v", err)
	}
	for _, line := range reqStr {
		rt.p.Debug(DebugLevel, line)
	}
	resp, err := rt.transport.RoundTrip(req)
	defer func() {
		if err == nil {
			respStrSlice, err := BuildDebugStrFromHTTPResponse(resp, rt.debugHeaders)
			if err != nil {
				rt.p.Debug(ErrorLevel, "printing response to debug logs: %v", err)
			}
			for _, line := range respStrSlice {
				rt.p.Debug(DebugLevel, line)
			}
		}
	}()
	return resp, err
}

// isEmpty checks if a value is empty (nil, empty string, zero value for other types)
func isEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	switch v := value.(type) {
	case string:
		return v == ""
	case []interface{}:
		return len(v) == 0
	case []string:
		return len(v) == 0
	case []int:
		return len(v) == 0
	case []bool:
		return len(v) == 0
	case []float64:
		return len(v) == 0
	case map[string]any:
		return len(v) == 0
	default:
		return false
	}
}
