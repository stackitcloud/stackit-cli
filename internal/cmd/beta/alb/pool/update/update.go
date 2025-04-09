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
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
)

const (
	configurationFlag = "configuration"
	albNameFlag       = "name"
	poolNameFlag      = "pool"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Configuration *string
	AlbName       *string
	Poolname      *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates an application target pool",
		Long:  "Updates an application target pool.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Update an application target pool from a configuration file`,
				"$ stackit beta alb update --configuration my-target pool.json"),
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
				prompt := fmt.Sprintf("Are you sure you want to update an application target pool for project %q?", projectLabel)
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
				return fmt.Errorf("update application target pool: %w", err)
			}

			return outputResult(p, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(configurationFlag, "c", "", "filename of the input configuration file")
	cmd.Flags().StringP(albNameFlag, "n", "", "name of the target pool name to update")
	cmd.Flags().StringP(poolNameFlag, "t", "", "name of the target pool to update")
	err := flags.MarkFlagsRequired(cmd, configurationFlag, albNameFlag, poolNameFlag)
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
		Poolname:        flags.FlagToStringPointer(p, cmd, poolNameFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *alb.APIClient) (req alb.ApiUpdateTargetPoolRequest, err error) {
	payload, err := readPayload(ctx, model)
	if err != nil {
		return req, err
	}
	req = apiClient.UpdateTargetPool(ctx, model.ProjectId, model.Region, *model.AlbName, *model.Poolname)
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
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal target pool: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal target pool: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		operationState := "Updated"
		if model.Async {
			operationState = "Triggered update of"
		}
		p.Outputf("%s application target pool for %q. Name: %s\n", operationState, projectLabel, utils.PtrString(resp.Name))
		return nil
	}
}
