package template

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-sdk-go/services/alb"
)

const (
	formatFlag = "format"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Format *string
}

//go:embed template.json
var template string

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "template",
		Short: "create an alb template",
		Long:  "creates a json or yaml template file for creating/updating an application loadbalancer.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Creat a yaml template`,
				`$ stackit beta alb template --format=yaml`,
			),
			examples.NewExample(
				`Creat a json template`,
				`$ stackit beta alb template --format=json`,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			model, err := parseInput(p, cmd)
			if err != nil {
				return err
			}

			if model.Format == nil || *model.Format == "json" {
				p.Outputln(template)
			} else if *model.Format == "yaml" {
				var target alb.CreateLoadBalancerPayload
				if err := json.Unmarshal([]byte(template), &target); err != nil {
					return fmt.Errorf("cannot unmarshal template: %w", err)
				}
				encoder := yaml.NewEncoder(os.Stdout, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
				if encoder.Encode(target); err != nil {
					return fmt.Errorf("cannot marshal template to yaml: %w", err)
				}
			} else {
				return fmt.Errorf("invalid format %q defined. Must be 'json' or 'yaml'", *model.Format)
			}

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().VarP(flags.EnumFlag(true, "json", "json", "yaml"), formatFlag, "f", "Defines the output format (yaml or json), default is json")
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Format:          flags.FlagToStringPointer(p, cmd, formatFlag),
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
