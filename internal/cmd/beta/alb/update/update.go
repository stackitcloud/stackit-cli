package update

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/alb/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
	"github.com/stackitcloud/stackit-sdk-go/services/alb/wait"
)

const (
	configurationFlag = "configuration"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Configuration *string
	Version       *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates an application loadbalancer",
		Long:  "Updates an application loadbalancer.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Update an application loadbalancer from a configuration file`,
				"$ stackit beta alb update --configuration my-loadbalancer.json"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update an application loadbalancer for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// for updates of an existing ALB the current version must be passed to the request
			model.Version, err = getCurrentAlbVersion(ctx, apiClient, model)
			if err != nil {
				return err
			}
			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return err
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update application loadbalancer: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("updating loadbalancer")
				_, err = wait.CreateOrUpdateLoadbalancerWaitHandler(ctx, apiClient, model.ProjectId, model.Region, *resp.Name).
					WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for loadbalancer update: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(configurationFlag, "c", "", "Filename of the input configuration file")
	err := flags.MarkFlagsRequired(cmd, configurationFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Configuration:   flags.FlagToStringPointer(p, cmd, configurationFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func getCurrentAlbVersion(ctx context.Context, apiClient *alb.APIClient, model *inputModel) (*string, error) {
	// use the configuration file to find the name of the loadbalancer
	updatePayload, err := readPayload(ctx, model)
	if err != nil {
		return nil, err
	}
	if updatePayload.Name == nil {
		return nil, fmt.Errorf("no name found in configuration")
	}
	if err != nil {
		return nil, err
	}
	resp, err := apiClient.GetLoadBalancer(ctx, model.ProjectId, model.Region, *updatePayload.Name).Execute()
	if err != nil {
		return nil, err
	}
	return resp.Version, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *alb.APIClient) (req alb.ApiUpdateLoadBalancerRequest, err error) {
	payload, err := readPayload(ctx, model)
	if err != nil {
		return req, err
	}
	if payload.Name == nil {
		return req, fmt.Errorf("no name found in loadbalancer configuration")
	}
	payload.Version = model.Version
	req = apiClient.UpdateLoadBalancer(ctx, model.ProjectId, model.Region, *payload.Name)
	return req.UpdateLoadBalancerPayload(payload), nil
}

func readPayload(_ context.Context, model *inputModel) (payload alb.UpdateLoadBalancerPayload, err error) {
	if model.Configuration == nil {
		return payload, fmt.Errorf("no configuration file defined")
	}
	file, err := os.Open(*model.Configuration)
	if err != nil {
		return payload, fmt.Errorf("cannot open configuration file %q: %w", *model.Configuration, err)
	}
	defer file.Close() // nolint:errcheck // at this point close errors are not relevant anymore

	if strings.HasSuffix(*model.Configuration, ".yaml") {
		decoder := yaml.NewDecoder(bufio.NewReader(file), yaml.UseJSONUnmarshaler())
		if err := decoder.Decode(&payload); err != nil {
			return payload, fmt.Errorf("cannot deserialize yaml configuration from %q: %w", *model.Configuration, err)
		}
	} else if strings.HasSuffix(*model.Configuration, ".json") {
		decoder := json.NewDecoder(bufio.NewReader(file))
		if err := decoder.Decode(&payload); err != nil {
			return payload, fmt.Errorf("cannot deserialize json configuration from %q: %w", *model.Configuration, err)
		}
	} else {
		return payload, fmt.Errorf("cannot determine configuration fileformat of %q by extension. Must be '.json' or '.yaml'", *model.Configuration)
	}

	return payload, nil
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, resp *alb.LoadBalancer) error {
	if resp == nil {
		return fmt.Errorf("update loadbalancer response is empty")
	}
	return p.OutputResult(model.OutputFormat, resp, func() error {
		operationState := "Updated"
		if model.Async {
			operationState = "Triggered update of"
		}
		p.Outputf("%s application loadbalancer for %q. Name: %s\n", operationState, projectLabel, utils.PtrString(resp.Name))
		return nil
	})
}
