package update

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/goccy/go-yaml"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/albwaf/client"

	"github.com/spf13/cobra"
	albwaf "github.com/stackitcloud/stackit-sdk-go/services/albwaf/v1api"
)

const (
	nameArg           = "NAME"
	configurationFlag = "configuration"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name          string
	Configuration *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", nameArg),
		Short: "Updates a managed rule set of the ALB WAF",
		Long:  "Updates the rules of a managed rule set (MRS) of the Web Application Firewall (WAF) for application loadbalancers. Only the rules provided in the configuration file are updated, all other rules remain unchanged.",
		Args:  args.SingleArg(nameArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Update the rules of a managed rule set with name "my-managed-rule-set" from a configuration file`,
				"$ stackit beta alb-waf managed-rule-set update my-managed-rule-set --configuration my-rules.json"),
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

			prompt := fmt.Sprintf("Are you sure you want to update the managed rule set %q for project %q?", model.Name, projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
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
				return fmt.Errorf("update managed rule set: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resp)
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

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	name := inputArgs[0]
	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            name,
		Configuration:   flags.FlagToStringPointer(p, cmd, configurationFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *albwaf.APIClient) (req albwaf.ApiPatchManagedRuleSetRequest, err error) {
	payload, err := readPayload(model)
	if err != nil {
		return req, err
	}
	req = apiClient.DefaultAPI.PatchManagedRuleSet(ctx, model.ProjectId, model.Region, model.Name)
	return req.PatchManagedRuleSetPayload(payload), nil
}

func readPayload(model *inputModel) (payload albwaf.PatchManagedRuleSetPayload, err error) {
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

func outputResult(p *print.Printer, outputFormat, projectLabel string, resp *albwaf.GetManagedRuleSetResponse) error {
	if resp == nil {
		return fmt.Errorf("update managed rule set response is empty")
	}
	return p.OutputResult(outputFormat, resp, func() error {
		p.Outputf("Updated managed rule set %q for project %q.\n", resp.Name, projectLabel)
		return nil
	})
}
