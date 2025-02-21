package curl

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

var testURL = "https://some-service.api.stackit.cloud/v1/foo?bar=baz"
var testToken = "auth-token"

func fixtureArgValues(mods ...func(argValues []string)) []string {
	argValues := []string{
		testURL,
	}
	for _, mod := range mods {
		mod(argValues)
	}
	return argValues
}

func fixtureFlagValues(mods ...func(flagValues map[string]string)) map[string]string {
	flagValues := map[string]string{
		requestMethodFlag:          "post",
		headerFlag:                 "Test-header-1: Test value 1",
		dataFlag:                   "data",
		includeResponseHeadersFlag: "true",
		failOnHTTPErrorFlag:        "true",
		outputFileFlag:             "./output.txt",
	}
	for _, mod := range mods {
		mod(flagValues)
	}
	return flagValues
}

func fixtureInputModel(mods ...func(model *inputModel)) *inputModel {
	model := &inputModel{
		URL:                    testURL,
		RequestMethod:          "POST",
		Headers:                []string{"Test-header-1: Test value 1"},
		Data:                   utils.Ptr("data"),
		IncludeResponseHeaders: true,
		FailOnHTTPError:        true,
		OutputFile:             utils.Ptr("./output.txt"),
	}
	for _, mod := range mods {
		mod(model)
	}
	return model
}

func fixtureRequest(mods ...func(request *http.Request)) *http.Request {
	req, err := http.NewRequest("POST", testURL, bytes.NewBufferString("data"))
	req.Header.Set("Test-header-1", "Test value 1")
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testToken))
	for _, mod := range mods {
		mod(req)
	}
	return req
}

func TestParseInput(t *testing.T) {
	tests := []struct {
		description      string
		argValues        []string
		flagValues       map[string]string
		headerFlagValues []string
		allowedURLDomain string
		isValid          bool
		expectedModel    *inputModel
	}{
		{
			description:   "base",
			argValues:     fixtureArgValues(),
			flagValues:    fixtureFlagValues(),
			isValid:       true,
			expectedModel: fixtureInputModel(),
		},
		{
			description: "no arg values",
			argValues:   []string{},
			flagValues:  fixtureFlagValues(),
			isValid:     false,
		},
		{
			description: "no flag values",
			argValues:   fixtureArgValues(),
			flagValues:  map[string]string{},
			isValid:     true,
			expectedModel: &inputModel{
				URL:           testURL,
				RequestMethod: "GET",
			},
		},
		{
			description: "invalid URL 1",
			argValues: []string{
				"",
			},
			flagValues: fixtureFlagValues(),
			isValid:    false,
		},
		{
			description: "invalid URL 2",
			argValues: []string{
				"foo",
			},
			flagValues: fixtureFlagValues(),
			isValid:    false,
		},
		{
			description: "URL outside STACKIT",
			argValues: []string{
				"https://www.example.website.com/",
			},
			flagValues:       fixtureFlagValues(),
			allowedURLDomain: "",
			isValid:          true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.URL = "https://www.example.website.com/"
			}),
		},
		{
			description: "invalid method 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[requestMethodFlag] = ""
			}),
			isValid: false,
		},
		{
			description: "invalid method 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[requestMethodFlag] = "foo"
			}),
			isValid: false,
		},
		{
			description: "invalid method 3",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[requestMethodFlag] = " GET"
			}),
			isValid: false,
		},
		{
			description: "valid method 1",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[requestMethodFlag] = "put"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.RequestMethod = "PUT"
			}),
		},
		{
			description: "valid method 2",
			argValues:   fixtureArgValues(),
			flagValues: fixtureFlagValues(func(flagValues map[string]string) {
				flagValues[requestMethodFlag] = "pAtCh"
			}),
			isValid: true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.RequestMethod = "PATCH"
			}),
		},
		{
			description:      "repeated header flags",
			argValues:        fixtureArgValues(),
			flagValues:       fixtureFlagValues(),
			headerFlagValues: []string{"Test-header-2: Test value 2", "Test-header-3: Test value 3"},
			isValid:          true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Headers = append(
					model.Headers,
					"Test-header-2: Test value 2",
					"Test-header-3: Test value 3",
				)
			}),
		},
		{
			description:      "repeated header flags with list value",
			argValues:        fixtureArgValues(),
			flagValues:       fixtureFlagValues(),
			headerFlagValues: []string{"Test-header-2: Test value 2,Test-header-3: Test value 3"},
			isValid:          true,
			expectedModel: fixtureInputModel(func(model *inputModel) {
				model.Headers = append(
					model.Headers,
					"Test-header-2: Test value 2",
					"Test-header-3: Test value 3",
				)
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			p := print.NewPrinter()
			cmd := NewCmd(p)
			err := globalflags.Configure(cmd.Flags())
			if err != nil {
				t.Fatalf("configure global flags: %v", err)
			}

			viper.Reset()
			viper.Set(config.AllowedUrlDomainKey, tt.allowedURLDomain)

			for flag, value := range tt.flagValues {
				err := cmd.Flags().Set(flag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", flag, value, err)
				}
			}

			for _, value := range tt.headerFlagValues {
				err := cmd.Flags().Set(headerFlag, value)
				if err != nil {
					if !tt.isValid {
						return
					}
					t.Fatalf("setting flag --%s=%s: %v", headerFlag, value, err)
				}
			}

			err = cmd.ValidateArgs(tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating args: %v", err)
			}

			err = cmd.ValidateRequiredFlags()
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error validating flags: %v", err)
			}

			model, err := parseInput(p, cmd, tt.argValues)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error parsing input: %v", err)
			}

			if !tt.isValid {
				t.Fatalf("did not fail on invalid input")
			}
			diff := cmp.Diff(model, tt.expectedModel)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestBuildRequest(t *testing.T) {
	defaultReq, err := http.NewRequest("GET", testURL, http.NoBody)
	if err != nil {
		t.Fatalf("failed to create new request: %v", err)
	}
	defaultReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testToken))

	tests := []struct {
		description     string
		model           *inputModel
		isValid         bool
		expectedRequest *http.Request
	}{
		{
			description:     "base",
			model:           fixtureInputModel(),
			isValid:         true,
			expectedRequest: fixtureRequest(),
		},
		{
			description: "default values",
			model: &inputModel{
				URL:           testURL,
				RequestMethod: "GET",
			},
			isValid:         true,
			expectedRequest: defaultReq,
		},
		{
			description: "invalid header 1",
			model: fixtureInputModel(func(model *inputModel) {
				model.Headers = append(model.Headers, "foo")
			}),
			isValid: false,
		},
		{
			description: "invalid header 2",
			model: fixtureInputModel(func(model *inputModel) {
				model.Headers = append(model.Headers, "foo bar")
			}),
			isValid: false,
		},
		{
			description: "invalid header 3",
			model: fixtureInputModel(func(model *inputModel) {
				model.Headers = append(model.Headers, "foo:")
			}),
			isValid: false,
		},
		{
			description: "extra headers 1",
			model: fixtureInputModel(func(model *inputModel) {
				model.Headers = append(
					model.Headers,
					"Test-header-2: Test value 2",
					"Test-header-3: Test value 3",
				)
			}),
			isValid: true,
			expectedRequest: fixtureRequest(func(request *http.Request) {
				request.Header.Set("Test-header-2", "Test value 2")
				request.Header.Set("Test-header-3", "Test value 3")
			}),
		},
		{
			description: "extra headers 2",
			model: fixtureInputModel(func(model *inputModel) {
				model.Headers = append(
					model.Headers,
					"Test-header-2: Test value 2",
					"Test-header-3: Test value 3",
					"Test-header-2: Test value 4",
				)
			}),
			isValid: true,
			expectedRequest: fixtureRequest(func(request *http.Request) {
				request.Header.Set("Test-header-2", "Test value 4")
				request.Header.Set("Test-header-3", "Test value 3")
			}),
		},
		{
			description: "extra headers 3",
			model: fixtureInputModel(func(model *inputModel) {
				model.Headers = append(
					model.Headers,
					"Test-header-2: Test value 2",
					"Test-header-3: Test value 3",
					"Authorization: Test value 4",
				)
			}),
			isValid: true,
			expectedRequest: fixtureRequest(func(request *http.Request) {
				request.Header.Set("Test-header-2", "Test value 2")
				request.Header.Set("Test-header-3", "Test value 3")
				request.Header.Set("Authorization", "Test value 4")
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			request, err := buildRequest(tt.model, testToken)
			if err != nil {
				if !tt.isValid {
					return
				}
				t.Fatalf("error building request: %v", err)
			}

			diff := cmp.Diff(request, tt.expectedRequest,
				cmp.AllowUnexported(http.Request{}),
				cmpopts.IgnoreFields(http.Request{}, "GetBody"), // Function, not relevant for the test
				cmp.Comparer(func(x, y *bytes.Buffer) bool { // Used to compare request bodies
					xBytes := x.Bytes()
					yBytes := y.Bytes()
					return bytes.Equal(xBytes, yBytes)
				}),
				cmpopts.EquateComparable(context.Background()),
			)
			if diff != "" {
				t.Fatalf("Data does not match: %s", diff)
			}
		})
	}
}

func TestOutputResponse(t *testing.T) {
	type args struct {
		model *inputModel
		resp  *http.Response
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "empty",
			args:    args{},
			wantErr: true,
		},
		{
			name: "http response as argument",
			args: args{
				model: fixtureInputModel(),
				resp:  &http.Response{Body: http.NoBody},
			},
			wantErr: false,
		},
	}
	p := print.NewPrinter()
	p.Cmd = NewCmd(p)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := outputResponse(p, tt.args.model, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("outputResponse() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.args.model != nil {
				if _, err := os.Stat(*tt.args.model.OutputFile); err == nil {
					if err := os.Remove(*tt.args.model.OutputFile); err != nil {
						t.Errorf("remove output file error = %v", err)
					}
				}
			}
		})
	}
}
