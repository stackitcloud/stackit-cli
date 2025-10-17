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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
)

const (
	configurationFlag = "configuration"
	albNameFlag       = "name"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Configuration *string
	AlbName       *string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates an application target pool",
		Long:  "Updates an application target pool.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Update an application target pool from a configuration file (the name of the pool is read from the file)`,
				"$ stackit beta alb update --configuration my-target-pool.json --name my-load-balancer"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
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
				prompt := fmt.Sprintf("Are you sure you want to update an application target pool for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
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
				return fmt.Errorf("update application target pool: %w", err)
			}

			return outputResult(params.Printer, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(configurationFlag, "c", "", "Filename of the input configuration file")
	cmd.Flags().StringP(albNameFlag, "n", "", "Name of the target pool name to update")
	err := flags.MarkFlagsRequired(cmd, configurationFlag, albNameFlag)
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
		AlbName:         flags.FlagToStringPointer(p, cmd, albNameFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *alb.APIClient) (req alb.ApiUpdateTargetPoolRequest, err error) {
	payload, err := readPayload(ctx, model)
	if err != nil {
		return req, err
	}
	if payload.Name == nil {
		return req, fmt.Errorf("update target pool: no poolname provided")
	}
	req = apiClient.UpdateTargetPool(ctx, model.ProjectId, model.Region, *model.AlbName, *payload.Name)
	return req.UpdateTargetPoolPayload(payload), nil
}

func readPayload(_ context.Context, model *inputModel) (payload alb.UpdateTargetPoolPayload, err error) {
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

func outputResult(p *print.Printer, model *inputModel, projectLabel string, resp *alb.TargetPool) error {
	if resp == nil {
		return fmt.Errorf("update target pool response is empty")
	}
	return p.OutputResult(model.OutputFormat, resp, func() error {
		operationState := "Updated"
		if model.Async {
			operationState = "Triggered update of"
		}
		p.Outputf("%s application target pool for %q. Name: %s\n", operationState, projectLabel, utils.PtrString(resp.Name))
		return nil
	})
}
