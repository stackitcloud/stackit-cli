package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/secretsmanager"

	"github.com/spf13/cobra"
)

const (
	instanceNameFlag = "name"
	aclFlag          = "acl"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	InstanceName *string
	Acls         *[]string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a Secrets Manager instance",
		Long:  "Creates a Secrets Manager instance.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a Secrets Manager instance with name "my-instance"`,
				`$ stackit secrets-manager instance create --name my-instance`),
			examples.NewExample(
				`Create a Secrets Manager instance with name "my-instance" and specify IP range which is allowed to access it`,
				`$ stackit secrets-manager instance create --name my-instance --acl 1.2.3.0/24`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			model, err := parseInput(cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, cmd, p)
			if err != nil {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a Secrets Manager instance for project %q?", projectLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API to create instance
			req := buildCreateInstanceRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Secrets Manager instance: %w", err)
			}
			instanceId := *resp.Id

			// Call API to create ACLs for instance, if ACLs are provided
			if model.Acls != nil {
				updateReq := buildUpdateACLsRequest(ctx, model, instanceId, apiClient)
				err = updateReq.Execute()
				if err != nil {
					return fmt.Errorf(`the Secrets Manager instance was successfully created, but the configuration of the ACLs failed. The default behavior is to have no ACL.

If you want to retry configuring the ACLs, you can do it via: 
  $ stackit secrets-manager instance update %s --acl %s`, instanceId, *model.Acls)
				}
			}

			p.Outputf("Created instance for project %q. Instance ID: %s\n", projectLabel, instanceId)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(instanceNameFlag, "n", "", "Instance name")
	cmd.Flags().Var(flags.CIDRSliceFlag(), aclFlag, "List of IP networks in CIDR notation which are allowed to access this instance")

	err := flags.MarkFlagsRequired(cmd, instanceNameFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		InstanceName:    flags.FlagToStringPointer(cmd, instanceNameFlag),
		Acls:            flags.FlagToStringSlicePointer(cmd, aclFlag),
	}, nil
}

func buildCreateInstanceRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiCreateInstanceRequest {
	req := apiClient.CreateInstance(ctx, model.ProjectId)

	req = req.CreateInstancePayload(secretsmanager.CreateInstancePayload{
		Name: model.InstanceName,
	})

	return req
}

func buildUpdateACLsRequest(ctx context.Context, model *inputModel, instanceId string, apiClient *secretsmanager.APIClient) secretsmanager.ApiUpdateACLsRequest {
	req := apiClient.UpdateACLs(ctx, model.ProjectId, instanceId)

	cidrs := make([]secretsmanager.AclUpdate, len(*model.Acls))

	for i, acl := range *model.Acls {
		cidrs[i] = secretsmanager.AclUpdate{Cidr: utils.Ptr(acl)}
	}

	req = req.UpdateACLsPayload(secretsmanager.UpdateACLsPayload{Cidrs: &cidrs})

	return req
}
