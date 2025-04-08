package update

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
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
}

func NewCmd(p *print.Printer) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
			if err != nil {
				p.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update an application loadbalancer for project %q?", projectLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
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
				s := spinner.New(p)
				s.Start("updating loadbalancer")
				_, err = wait.CreateOrUpdateLoadbalancerWaitHandler(ctx, apiClient, model.ProjectId, model.Region, *resp.Name).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for loadbalancer creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(p, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(configurationFlag, "c", "", "filename of the input configuration file")
	err := flags.MarkFlagsRequired(cmd, configurationFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Configuration:   flags.FlagToStringPointer(p, cmd, configurationFlag),
	}

	if p.IsVerbosityDebug() {
		modelStr, err := print.BuildDebugStrFromInputModel(model)
		if err != nil {
			p.Debug(print.ErrorLevel, "convert model to string for debugging: %v", err)
		} else {
			p.Debug(print.DebugLevel, "parsed input values: %s", modelStr)
		}
	}

	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *alb.APIClient) (req alb.ApiUpdateLoadBalancerRequest, err error) {
	payload, err := readPayload(ctx, model)
	if err != nil {
		return req, err
	}
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
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal loadbalancer: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal loadbalancer: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		operationState := "Updated"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s application loadbalancer for %q. Name: %s\n", operationState, projectLabel, utils.PtrString(resp.Name))
		return nil
	}
}
