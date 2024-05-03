package print

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"slices"
	"sort"
	"strings"

	"github.com/stackitcloud/stackit-sdk-go/core/config"
)

var defaultHTTPHeaders = []string{"Accept", "Content-Type", "Content-Length", "User-Agent", "Date", "Referrer-Policy"}

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
	if inputMap == nil {
		return "[]"
	}
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

		valueStr := fmt.Sprintf("%v", value)

		switch value := value.(type) {
		case map[string]any:
			valueStr = BuildDebugStrFromMap(value)
		case []any:
			sliceStr := make([]string, len(value))
			for i, item := range value {
				if itemMap, ok := item.(map[string]any); ok {
					sliceStr[i] = BuildDebugStrFromMap(itemMap)
				} else {
					sliceStr[i] = fmt.Sprintf("%v", item)
				}
			}
			valueStr = BuildDebugStrFromSlice(sliceStr)
		}

		keyValues = append(keyValues, fmt.Sprintf("%s: %v", key, valueStr))
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

	headersToInclude := defaultHTTPHeaders
	if len(includeHeaders) != 0 {
		headersToInclude = includeHeaders
	}

	for key := range headersMap {
		if !slices.Contains(headersToInclude, key) {
			delete(headersMap, key)
		}
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

	bodyMap, err := dumpReqBody(req)
	if err != nil {
		return []string{status, headers}, fmt.Errorf("read request body: %w", err)
	}
	if bodyMap == nil {
		return []string{status, headers}, nil
	}
	body := fmt.Sprintf("request body: %s", BuildDebugStrFromMap(bodyMap))

	return []string{status, headers, body}, nil
}

// drainBody reads all of b to memory and then returns two equivalent
// ReadClosers yielding the same bytes.
//
// It returns an error if the initial slurp of all bytes fails. It does not attempt
// to make the returned ReadClosers have identical error-matching behavior.
// Taken direclty from the httputil package
// https://cs.opensource.google/go/go/+/refs/tags/go1.22.2:src/net/http/httputil/dump.go;drc=1d45a7ef560a76318ed59dfdb178cecd58caf948;l=25
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if b == nil || b == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, b, err
	}
	if err := b.Close(); err != nil {
		return nil, b, err
	}
	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
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

	bodyMap, err := dumpRespBody(resp)
	if err != nil {
		return []string{status, headers}, fmt.Errorf("read response body: %w", err)
	}
	if bodyMap == nil {
		return []string{status, headers}, nil
	}
	body := fmt.Sprintf("response body: %s", BuildDebugStrFromMap(bodyMap))

	return []string{status, headers, body}, nil
}

// dumpRespBody reads the response body and returns a string representation of it.
// Based on code from httputil package
// https://pkg.go.dev/net/http/httputil#DumpResponse
func dumpRespBody(resp *http.Response) (map[string]any, error) {
	if resp == nil {
		return nil, fmt.Errorf("response is nil")
	}
	if resp.Body == nil || resp.ContentLength == 0 {
		return nil, nil
	}
	var err error
	var buf bytes.Buffer
	var save io.ReadCloser

	savecl := resp.ContentLength

	save, resp.Body, err = drainBody(resp.Body)
	if err != nil {
		return nil, err
	}
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	if err = resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("close response body: %w", err)
	}

	resp.Body = save
	resp.ContentLength = savecl

	var bodyMap map[string]any
	if len(buf.Bytes()) != 0 {
		if err := json.Unmarshal(buf.Bytes(), &bodyMap); err != nil {
			return nil, fmt.Errorf("unmarshal response body: %w", err)
		}
	}

	return bodyMap, nil
}

// dumpReqBody reads the request body and returns a string representation of it.
// Based on code from httputil package
// https://pkg.go.dev/net/http/httputil#DumpRequest
func dumpReqBody(req *http.Request) (map[string]any, error) {
	if req == nil {
		return nil, fmt.Errorf("request is nil")
	}
	if req.Body == nil {
		return nil, nil
	}
	var err error
	var b bytes.Buffer
	var save io.ReadCloser

	save, req.Body, err = drainBody(req.Body)
	if err != nil {
		return nil, err
	}

	var dest io.Writer = &b
	chunked := len(req.TransferEncoding) > 0 && req.TransferEncoding[0] == "chunked"
	if chunked {
		dest = httputil.NewChunkedWriter(&b)
	}
	_, err = io.Copy(dest, req.Body)
	if chunked {
		if closer, ok := dest.(io.Closer); ok {
			err = closer.Close()
			if err != nil {
				return nil, fmt.Errorf("close chunked writer: %w", err)
			}
		}
		_, err = b.WriteString("\r\n")
		if err != nil {
			return nil, fmt.Errorf("write chunked trailer: %w", err)
		}
	}

	req.Body = save
	if err != nil {
		return nil, err
	}

	// marshall body to map
	var bodyMap map[string]any
	if len(b.Bytes()) != 0 {
		if err := json.Unmarshal(b.Bytes(), &bodyMap); err != nil {
			return nil, fmt.Errorf("unmarshal request body: %w", err)
		}
	}

	return bodyMap, nil
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
	transport        http.RoundTripper
	p                *Printer
	debugHttpHeaders []string
}

func (rt roundTripperWithCapture) RoundTrip(req *http.Request) (*http.Response, error) {
	reqStr, err := BuildDebugStrFromHTTPRequest(req, rt.debugHttpHeaders)
	if err != nil {
		rt.p.Debug(ErrorLevel, "printing request to debug logs: %v", err)
	}
	for _, line := range reqStr {
		rt.p.Debug(DebugLevel, line)
	}
	resp, err := rt.transport.RoundTrip(req)
	defer func() {
		if err == nil {
			respStrSlice, tempErr := BuildDebugStrFromHTTPResponse(resp, rt.debugHttpHeaders)
			if tempErr != nil {
				rt.p.Debug(ErrorLevel, "printing HTTP response to debug logs: %v", tempErr)
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
