package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/client"
	observabilityUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/observability/utils"

	"github.com/spf13/cobra"
	observability "github.com/stackitcloud/stackit-sdk-go/services/observability/v1api"
)

const (
	instanceIdFlag = "instance-id"
	aclFlag        = "acl"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
	Acls       []string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Updates the ACLs of an Observability instance",
		Long:  "Updates the ACLs of an Observability instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Set the ACLs of Observability instance "xxx" to the IP networks "1.2.3.0/24" and "5.6.7.8/32"`,
				"$ stackit observability acl update --instance-id xxx --acl 1.2.3.0/24,5.6.7.8/32"),
			examples.NewExample(
				`Remove all ACLs of Observability instance "xxx"`,
				`$ stackit observability acl update --instance-id xxx --acl ""`),
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

			instanceLabel, err := observabilityUtils.GetInstanceName(ctx, apiClient.DefaultAPI, model.InstanceId, model.ProjectId)
			if err != nil || instanceLabel == "" {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
			}

			prompt := fmt.Sprintf("Are you sure you want to replace the ACLs of Observability instance %q?", instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("update Observability instance ACLs: %w", err)
			}

			params.Printer.Info("Updated ACLs of Observability instance %q\n", instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.UUIDFlag(), instanceIdFlag, "Instance ID")
	cmd.Flags().Var(flags.CIDRSliceFlag(), aclFlag, "List of IP networks in CIDR notation which are allowed to access the instance. Replaces all existing ACLs, an empty value removes all ACLs")

	err := flags.MarkFlagsRequired(cmd, instanceIdFlag, aclFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	acls := flags.FlagToStringSliceValue(p, cmd, aclFlag)
	if acls == nil {
		// An empty --acl value is explicitly allowed to remove all ACLs
		acls = []string{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      flags.FlagToStringValue(p, cmd, instanceIdFlag),
		Acls:            acls,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient observability.DefaultAPI) observability.ApiUpdateACLRequest {
	req := apiClient.UpdateACL(ctx, model.InstanceId, model.ProjectId)
	req = req.UpdateACLPayload(observability.UpdateACLPayload{
		Acl: model.Acls,
	})
	return req
}
