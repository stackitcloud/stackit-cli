package curl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

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

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("curl %s", urlArg),
		Short: "Executes an authenticated HTTP request to an endpoint",
		Long:  "Executes an HTTP request to an endpoint, using the authentication provided by the CLI.",
		Example: examples.Build(
			examples.NewExample(
				"Get all the DNS zones for project with ID xxx via GET request to https://dns.api.stackit.cloud/v1/projects/xxx/zones",
				"$ stackit curl https://dns.api.stackit.cloud/v1/projects/xxx/zones",
			),
			examples.NewExample(
				`Get all the DNS zones for project with ID xxx via GET request to https://dns.api.stackit.cloud/v1/projects/xxx/zones, write complete response (headers and body) to file "./output.txt"`,
				"$ stackit curl https://dns.api.stackit.cloud/v1/projects/xxx/zones --include --output ./output.txt",
			),
			examples.NewExample(
				`Create a new DNS zone for project with ID xxx via POST request to https://dns.api.stackit.cloud/v1/projects/xxx/zones with payload from file "./payload.json"`,
				`$ stackit curl https://dns.api.stackit.cloud/v1/projects/xxx/zones -X POST --data @./payload.json`,
			),
			examples.NewExample(
				`Get all the DNS zones for project with ID xxx via GET request to https://dns.api.stackit.cloud/v1/projects/xxx/zones, with header "Authorization: Bearer yyy", fail if server returns error (such as 403 Forbidden)`,
				`$ stackit curl https://dns.api.stackit.cloud/v1/projects/xxx/zones -X POST -H "Authorization: Bearer yyy" --fail`,
			),
		),
		Args: args.SingleArg(urlArg, utils.ValidateURLDomain),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			bearerToken, err := getBearerToken(params.Printer)
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

			err = outputResponse(params.Printer, model, resp)
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

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	urlString := inputArgs[0]
	requestMethod := flags.FlagToStringValue(p, cmd, requestMethodFlag)
	if requestMethod == "" {
		requestMethod = http.MethodGet
	}

	model := inputModel{
		URL:                    urlString,
		RequestMethod:          strings.ToUpper(requestMethod),
		Headers:                flags.FlagToStringSliceValue(p, cmd, headerFlag),
		Data:                   flags.FlagToStringPointer(p, cmd, dataFlag),
		IncludeResponseHeaders: flags.FlagToBoolValue(p, cmd, includeResponseHeadersFlag),
		FailOnHTTPError:        flags.FlagToBoolValue(p, cmd, failOnHTTPErrorFlag),
		OutputFile:             flags.FlagToStringPointer(p, cmd, outputFileFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func getBearerToken(p *print.Printer) (string, error) {
	_, err := auth.AuthenticationConfig(p, auth.AuthorizeUser)
	if err != nil {
		p.Debug(print.ErrorLevel, "configure authentication: %v", err)
		return "", &errors.AuthError{}
	}

	userSessionExpired, err := auth.UserSessionExpired()
	if err != nil {
		return "", err
	}
	if userSessionExpired {
		return "", &errors.SessionExpiredError{}
	}

	accessToken, err := auth.GetValidAccessToken(p)
	if err != nil {
		p.Debug(print.ErrorLevel, "get valid access token: %v", err)
		return "", &errors.SessionExpiredError{}
	}

	return accessToken, nil
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

func outputResponse(p *print.Printer, model *inputModel, resp *http.Response) error {
	if resp == nil {
		return fmt.Errorf("http response is empty")
	}
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
		return fmt.Errorf("read response body: %w", err)
	}

	if strings.Contains(strings.ToLower(string(respBody)), "jwt is expired") {
		return &errors.SessionExpiredError{}
	}

	if strings.Contains(strings.ToLower(string(respBody)), "jwt is missing") {
		return &errors.AuthError{}
	}

	var prettyJSON bytes.Buffer
	if json.Valid(respBody) {
		if err := json.Indent(&prettyJSON, respBody, "", "  "); err == nil {
			respBody = prettyJSON.Bytes()
		} // if indenting fails, fall back to original body
	}

	output = append(output, respBody...)

	if model.OutputFile == nil {
		p.Outputln(string(output))
	} else {
		err = os.WriteFile(utils.PtrString(model.OutputFile), output, 0o600)
		if err != nil {
			return fmt.Errorf("write output to file: %w", err)
		}
	}

	return nil
}
