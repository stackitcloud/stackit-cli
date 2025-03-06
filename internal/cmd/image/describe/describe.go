package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ImageId string
}

const imageIdArg = "IMAGE_ID"

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", imageIdArg),
		Short: "Describes image",
		Long:  "Describes an image by its internal ID.",
		Args:  args.SingleArg(imageIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(`Describe image "xxx"`, `$ stackit image describe xxx`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			// Call API
			request := buildRequest(ctx, model, apiClient)

			image, err := request.Execute()
			if err != nil {
				return fmt.Errorf("get image: %w", err)
			}

			if err := outputResult(p, model.OutputFormat, image); err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiGetImageRequest {
	request := apiClient.GetImage(ctx, model.ProjectId, model.ImageId)
	return request
}

func parseInput(p *print.Printer, cmd *cobra.Command, cliArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ImageId:         cliArgs[0],
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

func outputResult(p *print.Printer, outputFormat string, resp *iaas.Image) error {
	if resp == nil {
		return fmt.Errorf("image not found")
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal image: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal image: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		if id := resp.Id; id != nil {
			table.AddRow("ID", *id)
		}
		table.AddSeparator()

		if name := resp.Name; name != nil {
			table.AddRow("NAME", *name)
			table.AddSeparator()
		}
		if format := resp.DiskFormat; format != nil {
			table.AddRow("FORMAT", *format)
			table.AddSeparator()
		}
		if diskSize := resp.MinDiskSize; diskSize != nil {
			table.AddRow("DISK SIZE", *diskSize)
			table.AddSeparator()
		}
		if ramSize := resp.MinRam; ramSize != nil {
			table.AddRow("RAM SIZE", *ramSize)
			table.AddSeparator()
		}
		if config := resp.Config; config != nil {
			if os := config.OperatingSystem; os != nil {
				table.AddRow("OPERATING SYSTEM", *os)
				table.AddSeparator()
			}
			if distro := config.OperatingSystemDistro; distro != nil && distro.IsSet() {
				table.AddRow("OPERATING SYSTEM DISTRIBUTION", *distro.Get())
				table.AddSeparator()
			}
			if version := config.OperatingSystemVersion; version != nil && version.IsSet() {
				table.AddRow("OPERATING SYSTEM VERSION", *version.Get())
				table.AddSeparator()
			}
			if uefi := config.Uefi; uefi != nil {
				table.AddRow("UEFI BOOT", *uefi)
				table.AddSeparator()
			}
		}

		if resp.Labels != nil && len(*resp.Labels) > 0 {
			labels := []string{}
			for key, value := range *resp.Labels {
				labels = append(labels, fmt.Sprintf("%s: %s", key, value))
			}
			table.AddRow("LABELS", strings.Join(labels, "\n"))
			table.AddSeparator()
		}

		if err := table.Display(p); err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
