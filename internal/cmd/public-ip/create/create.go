package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	associatedResourceIdFlag = "associated-resource-id"
	labelFlag                = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	AssociatedResourceId *string
	Labels               *map[string]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Public IP",
		Long:  "Creates a Public IP.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a public IP`,
				`$ stackit public-ip create`,
			),
			examples.NewExample(
				`Create a public IP with associated resource ID "xxx"`,
				`$ stackit public-ip create --associated-resource-id xxx`,
			),
			examples.NewExample(
				`Create a public IP with associated resource ID "xxx" and labels`,
				`$ stackit public-ip create --associated-resource-id xxx --labels key=value,foo=bar`,
			),
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
			} else if projectLabel == "" {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a public IP for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create public IP: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, *resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), associatedResourceIdFlag, "Associates the public IP with a network interface or virtual IP (ID)")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a public IP. E.g. '--labels key1=value1,key2=value2,...'")
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:      globalFlags,
		AssociatedResourceId: flags.FlagToStringPointer(p, cmd, associatedResourceIdFlag),
		Labels:               flags.FlagToStringToStringPointer(p, cmd, labelFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreatePublicIPRequest {
	req := apiClient.CreatePublicIP(ctx, model.ProjectId, model.Region)

	payload := iaas.CreatePublicIPPayload{
		NetworkInterface: iaas.NewNullableString(model.AssociatedResourceId),
		Labels:           utils.ConvertStringMapToInterfaceMap(model.Labels),
	}

	return req.CreatePublicIPPayload(payload)
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, publicIp iaas.PublicIp) error {
	return p.OutputResult(outputFormat, publicIp, func() error {
		p.Outputf("Created public IP for project %q.\nPublic IP ID: %s\n", projectLabel, utils.PtrString(publicIp.Id))
		return nil
	})
}
