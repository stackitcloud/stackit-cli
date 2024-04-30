package print

import (
	"testing"

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
