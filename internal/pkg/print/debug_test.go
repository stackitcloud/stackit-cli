package print

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"
)

type testGlobalFlags struct {
	Async        bool
	AssumeYes    bool
	OutputFormat string
	ProjectId    string
	Verbosity    string
}

type testInputModel struct {
	*testGlobalFlags
	InstanceId   string
	HidePassword bool
	JobName      string
	MaxCount     int
}

const (
	testJobName = "test-job"
)

var (
	testInstanceId = uuid.NewString()
	testProjectId  = uuid.NewString()
)

func fixtureInputModel(mods ...func(model *testInputModel)) *testInputModel {
	model := &testInputModel{
		testGlobalFlags: &testGlobalFlags{
			Async:        false,
			AssumeYes:    false,
			OutputFormat: "pretty",
			ProjectId:    testProjectId,
			Verbosity:    "info",
		},
		InstanceId:   testInstanceId,
		JobName:      testJobName,
		MaxCount:     10,
		HidePassword: true,
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureHTTPRequest(mods ...func(req *http.Request)) *http.Request {
	testBody, err := json.Marshal(map[string]string{"key": "value"})
	if err != nil {
		return nil
	}
	request, err := http.NewRequest("GET", "http://example.com", bytes.NewReader(testBody))
	if err != nil {
		return nil
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("Content-Length", "15")

	for _, mod := range mods {
		mod(request)
	}

	return request
}

func fixtureHTTPResponse(mods ...func(resp *http.Response)) *http.Response {
	testBody, err := json.Marshal(map[string]string{"key": "value"})
	if err != nil {
		return nil
	}
	response := &http.Response{
		Body:          io.NopCloser(bytes.NewReader(testBody)),
		StatusCode:    http.StatusOK,
		Proto:         "HTTP/1.1",
		Status:        "200 OK",
		ContentLength: int64(len(testBody)),
		Request:       &http.Request{Method: "GET", URL: &url.URL{Host: "example.com", Scheme: "http"}},
		Header: http.Header{
			"Content-Type":   []string{"application/json"},
			"Accept":         []string{"application/json"},
			"Content-Length": []string{"15"},
		},
	}

	for _, mod := range mods {
		mod(response)
	}

	return response
}

func TestBuildDebugStrFromInputModel(t *testing.T) {
	tests := []struct {
		description string
		model       any
		expected    string
		isValid     bool
	}{
		{
			description: "base",
			model:       fixtureInputModel(),
			expected:    `[AssumeYes: false, Async: false, HidePassword: true, InstanceId: ` + testInstanceId + `, JobName: ` + testJobName + `, MaxCount: 10, OutputFormat: pretty, ProjectId: ` + testProjectId + `, Verbosity: info]`,
			isValid:     true,
		},
		{
			description: "empty string, zero values for int and bool fields",
			model: fixtureInputModel(func(model *testInputModel) {
				model.JobName = ""
				model.MaxCount = 0
				model.HidePassword = false
			}),
			expected: `[AssumeYes: false, Async: false, HidePassword: false, InstanceId: ` + testInstanceId + `, MaxCount: 0, OutputFormat: pretty, ProjectId: ` + testProjectId + `, Verbosity: info]`,
			isValid:  true,
		},
		{
			description: "empty input model",
			model:       &testInputModel{},
			expected:    `[HidePassword: false, MaxCount: 0]`,
			isValid:     true,
		},
		{
			description: "nil input model",
			model:       nil,
			expected:    `[]`,
			isValid:     true,
		},
		{
			description: "invalid input model",
			model:       []int{},
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			model := tt.model
			actual, err := BuildDebugStrFromInputModel(model)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.isValid {
				t.Fatalf("expected error, got nil")
			}
			if actual != tt.expected {
				t.Fatalf("expected: %s, actual: %s", tt.expected, actual)
			}
		})
	}
}

func TestBuildDebugStrFromMap(t *testing.T) {
	tests := []struct {
		description string
		inputMap    map[string]any
		expected    string
	}{
		{
			description: "base",
			inputMap: map[string]any{
				"key1": "value1",
				"key2": "value2",
				"key3": 123,
				"key4": false,
			},
			expected: "[key1: value1, key2: value2, key3: 123, key4: false]",
		},
		{
			description: "nested map",
			inputMap: map[string]any{
				"key1": "value1",
				"key2": map[string]any{
					"nestedKey1": "nestedValue1",
					"nestedKey2": "nestedValue2",
				},
			},
			expected: "[key1: value1, key2: [nestedKey1: nestedValue1, nestedKey2: nestedValue2]]",
		},
		{
			description: "nested slice of string",
			inputMap: map[string]any{
				"key1": "value1",
				"key2": []any{"value1", "value2"},
			},
			expected: "[key1: value1, key2: [value1, value2]]",
		},
		{
			description: "nested slice of int",
			inputMap: map[string]any{
				"key1": "value1",
				"key2": []any{1, 2},
			},
			expected: "[key1: value1, key2: [1, 2]]",
		},
		{
			description: "nested slice of map",
			inputMap: map[string]any{
				"key1": "value1",
				"key2": []any{
					map[string]any{
						"nestedKey1": "nestedValue1",
					},
					map[string]any{
						"nestedKey2": "nestedValue2",
					},
				},
			},
			expected: "[key1: value1, key2: [[nestedKey1: nestedValue1], [nestedKey2: nestedValue2]]]",
		},
		{
			description: "empty values",
			inputMap: map[string]any{
				"key1": "",
				"key2": nil,
				"key3": 0,
				"key4": false,
			},
			expected: "[key3: 0, key4: false]",
		},
		{
			description: "empty map",
			inputMap:    map[string]any{},
			expected:    "[]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			actual := BuildDebugStrFromMap(tt.inputMap)
			if actual != tt.expected {
				t.Fatalf("expected: %s, actual: %s", tt.expected, actual)
			}
		})
	}
}

func TestBuildDebugStrFromSlice(t *testing.T) {
	tests := []struct {
		description string
		inputSlice  []string
		expected    string
	}{
		{
			description: "base",
			inputSlice:  []string{"value1", "value2", "value3"},
			expected:    "[value1, value2, value3]",
		},
		{
			description: "empty slice",
			inputSlice:  []string{},
			expected:    "[]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			actual := BuildDebugStrFromSlice(tt.inputSlice)
			if actual != tt.expected {
				t.Fatalf("expected: %s, actual: %s", tt.expected, actual)
			}
		})
	}
}

func TestBuildHeaderMap(t *testing.T) {
	tests := []struct {
		description         string
		inputHeader         http.Header
		inputIncludeHeaders []string
		expected            map[string]any
	}{
		{
			description: "base",
			inputHeader: http.Header{
				"key1": []string{"value1"},
				"key2": []string{"value2"},
				"key3": []string{"value3"},
			},
			inputIncludeHeaders: []string{"key1", "key2"},
			expected: map[string]any{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			description: "no include headers",
			inputHeader: http.Header{
				"Accept": []string{"value1"},
				"key2":   []string{"value2"},
				"Date":   []string{"value3"},
			},
			inputIncludeHeaders: []string{},
			expected: map[string]any{
				"Accept": "value1",
				"Date":   "value3",
			},
		},
		{
			description:         "empty header",
			inputHeader:         http.Header{},
			inputIncludeHeaders: []string{},
			expected:            map[string]any{},
		},
		{
			description: "empty header, some include headers",
			inputHeader: http.Header{},
			inputIncludeHeaders: []string{
				"key1",
				"key2",
			},
			expected: map[string]any{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			actual := buildHeaderMap(tt.inputHeader, tt.inputIncludeHeaders)
			diff := cmp.Diff(actual, tt.expected)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildDebugStrFromHTTPRequest(t *testing.T) {
	tests := []struct {
		description         string
		inputReq            *http.Request
		inputIncludeHeaders []string
		expected            []string
		isValid             bool
	}{
		{
			description: "base",
			inputReq:    fixtureHTTPRequest(),
			expected: []string{
				"request to http://example.com: GET HTTP/1.1",
				"request headers: [Accept: application/json, Content-Length: 15, Content-Type: application/json]",
				"request body: [key: value]",
			},
			isValid: true,
		},
		{
			description:         "include headers",
			inputReq:            fixtureHTTPRequest(),
			inputIncludeHeaders: []string{"Content-Type", "Accept"},
			expected: []string{
				"request to http://example.com: GET HTTP/1.1",
				"request headers: [Accept: application/json, Content-Type: application/json]",
				"request body: [key: value]",
			},
			isValid: true,
		},
		{
			description: "empty request",
			inputReq:    &http.Request{},
			isValid:     false,
		},
		{
			description: "nil request",
			inputReq:    nil,
			isValid:     false,
		},
		{
			description: "empty headers",
			inputReq: fixtureHTTPRequest(func(req *http.Request) {
				req.Header = http.Header{}
			}),
			expected: []string{
				"request to http://example.com: GET HTTP/1.1",
				"request headers: []",
				"request body: [key: value]",
			},
			isValid: true,
		},
		{
			description: "empty body",
			inputReq: fixtureHTTPRequest(func(req *http.Request) {
				req.Body = nil
			}),
			expected: []string{
				"request to http://example.com: GET HTTP/1.1",
				"request headers: [Accept: application/json, Content-Length: 15, Content-Type: application/json]",
			},
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			actual, err := BuildDebugStrFromHTTPRequest(tt.inputReq, tt.inputIncludeHeaders)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.isValid {
				t.Fatalf("expected error, got nil")
			}
			diff := cmp.Diff(actual, tt.expected)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildDebugStrFromHTTPResponse(t *testing.T) {
	tests := []struct {
		description         string
		inputResp           *http.Response
		inputIncludeHeaders []string
		expected            []string
		isValid             bool
	}{
		{
			description: "base",
			inputResp:   fixtureHTTPResponse(), // nolint:bodyclose // false positive, body is closed in the test
			expected: []string{
				"response from http://example.com: HTTP/1.1 200 OK",
				"response headers: [Accept: application/json, Content-Length: 15, Content-Type: application/json]",
				"response body: [key: value]",
			},
			isValid: true,
		},
		{
			description: "empty response",
			inputResp:   &http.Response{},
			isValid:     false,
		},
		{
			description: "nil response",
			inputResp:   nil,
			isValid:     false,
		},
		{
			description: "empty headers",
			inputResp: fixtureHTTPResponse(func(resp *http.Response) { // nolint:bodyclose // false positive, body is closed in the test
				resp.Header = http.Header{}
			}),
			expected: []string{
				"response from http://example.com: HTTP/1.1 200 OK",
				"response headers: []",
				"response body: [key: value]",
			},
			isValid: true,
		},
		{
			description: "empty body",
			inputResp: fixtureHTTPResponse(func(resp *http.Response) { // nolint:bodyclose // false positive, body is closed in the test
				resp.Body = nil
			}),
			expected: []string{
				"response from http://example.com: HTTP/1.1 200 OK",
				"response headers: [Accept: application/json, Content-Length: 15, Content-Type: application/json]",
			},
			isValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			var err error
			if tt.inputResp != nil && tt.inputResp.Body != nil {
				defer func() {
					err = tt.inputResp.Body.Close()
				}()
			}
			actual, err := BuildDebugStrFromHTTPResponse(tt.inputResp, tt.inputIncludeHeaders)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.isValid {
				t.Fatalf("expected error, got nil")
			}
			diff := cmp.Diff(actual, tt.expected)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		description string
		value       interface{}
		expected    bool
	}{
		{
			description: "nil value",
			value:       nil,
			expected:    true,
		},
		{
			description: "empty string",
			value:       "",
			expected:    true,
		},
		{
			description: "zero value",
			value:       0,
			expected:    false,
		},
		{
			description: "non-empty string",
			value:       "test",
			expected:    false,
		},
		{
			description: "non-empty str slice",
			value:       []string{"test"},
			expected:    false,
		},
		{
			description: "empty str slice",
			value:       []string{},
			expected:    true,
		},
		{
			description: "non-empty int slice",
			value:       []int{1},
			expected:    false,
		},
		{
			description: "empty int slice",
			value:       []int{},
			expected:    true,
		},
		{
			description: "non-empty map",
			value:       map[string]any{"key": "value"},
			expected:    false,
		},
		{
			description: "empty map",
			value:       map[string]any{},
			expected:    true,
		},
		{
			description: "non-empty bool slice",
			value:       []bool{true},
			expected:    false,
		},
		{
			description: "empty bool slice",
			value:       []bool{},
			expected:    true,
		},
		{
			description: "non-empty float64 slice",
			value:       []float64{1.1},
			expected:    false,
		},
		{
			description: "empty float64 slice",
			value:       []float64{},
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			actual := isEmpty(tt.value)
			if actual != tt.expected {
				t.Fatalf("expected: %t, actual: %t", tt.expected, actual)
			}
		})
	}
}

func TestDumpRespBody(t *testing.T) {
	tests := []struct {
		description string
		resp        *http.Response
		expected    map[string]any
		isValid     bool
	}{
		{
			description: "base",
			resp:        fixtureHTTPResponse(), // nolint:bodyclose // false positive, body is closed in the test
			expected: map[string]any{
				"key": "value",
			},
			isValid: true,
		},
		{
			description: "empty response",
			resp:        &http.Response{},
			isValid:     true,
			expected:    nil,
		},
		{
			description: "nil response",
			resp:        nil,
			isValid:     false,
		},
		{
			description: "empty body",
			resp: fixtureHTTPResponse(func(resp *http.Response) { // nolint:bodyclose // false positive, body is closed in the test
				resp.Body = nil
			}),
			isValid: true,
		},
		{
			description: "invalid body",
			resp: fixtureHTTPResponse(func(resp *http.Response) { // nolint:bodyclose // false positive, body is closed in the test
				resp.Body = io.NopCloser(bytes.NewReader([]byte("invalid")))
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if tt.resp != nil && tt.resp.Body != nil {
				defer func() {
					_ = tt.resp.Body.Close()
				}()
			}
			actual, err := dumpRespBody(tt.resp)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.isValid {
				t.Fatalf("expected error, got nil")
			}
			diff := cmp.Diff(actual, tt.expected)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestDumpReqBody(t *testing.T) {
	tests := []struct {
		description string
		req         *http.Request
		expected    map[string]any
		isValid     bool
	}{
		{
			description: "base",
			req:         fixtureHTTPRequest(),
			expected: map[string]any{
				"key": "value",
			},
			isValid: true,
		},
		{
			description: "empty request",
			req:         &http.Request{},
			isValid:     true,
			expected:    nil,
		},
		{
			description: "nil request",
			req:         nil,
			isValid:     false,
		},
		{
			description: "empty body",
			req: fixtureHTTPRequest(func(req *http.Request) {
				req.Body = nil
			}),
			isValid: true,
		},
		{
			description: "invalid body",
			req: fixtureHTTPRequest(func(req *http.Request) {
				req.Body = io.NopCloser(bytes.NewReader([]byte("invalid")))
			}),
			isValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			if tt.req != nil && tt.req.Body != nil {
				defer func() {
					_ = tt.req.Body.Close()
				}()
			}
			actual, err := dumpReqBody(tt.req)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("unexpected error: %v", err)
			}
			if !tt.isValid {
				t.Fatalf("expected error, got nil")
			}
			diff := cmp.Diff(actual, tt.expected)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}
