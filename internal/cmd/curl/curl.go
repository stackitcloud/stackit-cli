package curl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"

	"github.com/spf13/cobra"
)

const (
	requestMethodFlag          = "request"
	headerFlag                 = "header"
	dataFlag                   = "data"
	includeResponseHeadersFlag = "include"
	failOnHTTPErrorFlag        = "fail"
	outputFileFlag             = "output"
)

const (
	urlArg = "URL"
)

type inputModel struct {
	URL                    string
	RequestMethod          string
	Headers                []string
	Data                   *string
	IncludeResponseHeaders bool
	FailOnHTTPError        bool
	OutputFile             *string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("curl %s", urlArg),
		Short: "Execute an authenticated HTTP request to an endpoint",
		Long:  "Execute an HTTP request to an endpoint, using the authentication provided by the CLI",
		Example: examples.Build(
			examples.NewExample(
				"Make a GET request to http://locahost:8000",
				"$ stackit curl http://locahost:8000",
			),
			examples.NewExample(
				`Make a GET request to http://locahost:8000, write complete response (headers and body) to file "./output.txt"`,
				"$ stackit curl http://locahost:8000 -include --output ./output.txt",
			),
			examples.NewExample(
				`Make a POST request to http://locahost:8000 with payload from file "./payload.json"`,
				`$ stackit curl http://locahost:8000 -X POST --data @./payload.json`,
			),
			examples.NewExample(
				`Make a POST request to http://locahost:8000 with header "Foo: Bar", fail if server returns error (such as 403 Forbidden)`,
				`$ stackit curl http://locahost:8000 -X POST -H "Foo: Bar" --fail`,
			),
		),
		Args: args.SingleArg(urlArg, validateURL),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			bearerToken, err := getBearerToken(cmd)
			if err != nil {
				return err
			}

			req, err := buildRequest(model, bearerToken)
			if err != nil {
				return err
			}

			client := http.Client{
				Timeout: 30 * time.Second,
			}
			resp, err := client.Do(req)
			if err != nil {
				return fmt.Errorf("do request: %w", err)
			}
			defer func() {
				closeErr := resp.Body.Close()
				if closeErr != nil {
					err = fmt.Errorf("close response body: %w", closeErr)
				}
			}()

			err = outputResponse(cmd, model, resp)
			if err != nil {
				return err
			}

			if model.FailOnHTTPError && resp.StatusCode >= 400 {
				os.Exit(22)
			}
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func validateURL(value string) error {
	urlStruct, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("parse URL: %w", err)
	}
	urlHost := urlStruct.Hostname()
	if urlHost == "" {
		return fmt.Errorf("bad url")
	}
	if !strings.HasSuffix(urlHost, "stackit.cloud") {
		return fmt.Errorf("only urls belonging to STACKIT are permitted")
	}
	return nil
}

func configureFlags(cmd *cobra.Command) {
	requestMethodOptions := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}
	headerFlagUsage := `Custom headers to include in the request, can be specified multiple times. If the "Authorization" header is set, it will override the authentication provided by the CLI`

	cmd.Flags().VarP(flags.EnumFlag(true, "", requestMethodOptions...), requestMethodFlag, "X", "HTTP method, defaults to GET")
	cmd.Flags().StringSliceP(headerFlag, "H", []string{}, headerFlagUsage)
	cmd.Flags().Var(flags.ReadFromFileFlag(), dataFlag, `Content to include in the request body. Can be a string or a file path prefixed with "@"`)
	cmd.Flags().Bool(includeResponseHeadersFlag, false, "If set, response headers are added to the output")
	cmd.Flags().Bool(failOnHTTPErrorFlag, false, "If set, exits with error 22 if response code is 4XX or 5XX")
	cmd.Flags().String(outputFileFlag, "", "Writes output to provided file instead of printing to console")
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	urlString := inputArgs[0]
	requestMethod := flags.FlagToStringValue(cmd, requestMethodFlag)
	if requestMethod == "" {
		requestMethod = http.MethodGet
	}

	return &inputModel{
		URL:                    urlString,
		RequestMethod:          strings.ToUpper(requestMethod),
		Headers:                flags.FlagToStringSliceValue(cmd, headerFlag),
		Data:                   flags.FlagToStringPointer(cmd, dataFlag),
		IncludeResponseHeaders: flags.FlagToBoolValue(cmd, includeResponseHeadersFlag),
		FailOnHTTPError:        flags.FlagToBoolValue(cmd, failOnHTTPErrorFlag),
		OutputFile:             flags.FlagToStringPointer(cmd, outputFileFlag),
	}, nil
}

func getBearerToken(cmd *cobra.Command) (string, error) {
	_, err := auth.AuthenticationConfig(cmd, auth.AuthorizeUser)
	if err != nil {
		return "", &errors.AuthError{}
	}
	token, err := auth.GetAuthField(auth.ACCESS_TOKEN)
	if err != nil {
		return "", fmt.Errorf("get access token: %w", err)
	}
	return token, nil
}

func buildRequest(model *inputModel, bearerToken string) (*http.Request, error) {
	var body io.Reader = http.NoBody
	if model.Data != nil {
		body = bytes.NewBufferString(*model.Data)
	}
	req, err := http.NewRequest(model.RequestMethod, model.URL, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	for _, header := range model.Headers {
		headerSplit := strings.SplitN(header, ": ", 2)
		if len(headerSplit) != 2 {
			return nil, fmt.Errorf("badly formatted header %q", header)
		}
		req.Header.Set(headerSplit[0], headerSplit[1])
	}
	return req, nil
}

func outputResponse(cmd *cobra.Command, model *inputModel, resp *http.Response) error {
	output := make([]byte, 0)
	if model.IncludeResponseHeaders {
		respHeader, err := httputil.DumpResponse(resp, false)
		if err != nil {
			return fmt.Errorf("print response headers: %w", err)
		}
		output = append(output, respHeader...)
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read respose body: %w", err)
	}
	output = append(output, respBody...)

	if model.OutputFile == nil {
		cmd.Println(string(output))
	} else {
		err = os.WriteFile(*model.OutputFile, output, 0o600)
		if err != nil {
			return fmt.Errorf("write output to file: %w", err)
		}
	}

	return nil
}
