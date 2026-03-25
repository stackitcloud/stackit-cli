package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"

	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/client"
	secretsManagerUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/secretsmanager"
)

const (
	instanceIdArg = "INSTANCE_ID"

	instanceNameFlag = "name"
	aclFlag          = "acl"

	kmsKeyIdFlag               = "kms-key-id"
	kmsKeyringIdFlag           = "kms-keyring-id"
	kmsKeyVersionFlag          = "kms-key-version"
	kmsServiceAccountEmailFlag = "kms-service-account-email"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string

	InstanceName *string
	Acls         *[]string

	KmsKeyId               *string
	KmsKeyringId           *string
	KmsKeyVersion          *int64
	KmsServiceAccountEmail *string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", instanceIdArg),
		Short: "Updates a Secrets Manager instance",
		Long:  "Updates a Secrets Manager instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the name of a Secrets Manager instance with ID "xxx"`,
				"$ stackit secrets-manager instance update xxx --name my-new-name"),
			examples.NewExample(
				`Update the range of IPs allowed to access a Secrets Manager instance with ID "xxx"`,
				"$ stackit secrets-manager instance update xxx --acl 1.2.3.0/24"),
			examples.NewExample(
				`Update the name and ACLs of a Secrets Manager instance with ID "xxx"`,
				"$ stackit secrets-manager instance update xxx --name my-new-name --acl 1.2.3.0/24"),
			examples.NewExample(
				`Update the KMS key settings of a Secrets Manager instance with ID "xxx"`,
				"$ stackit secrets-manager instance update xxx --name my-instance --kms-key-id key-id --kms-keyring-id keyring-id --kms-key-version 1 --kms-service-account-email my-service-account-1234567@sa.stackit.cloud"),
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

			existingInstanceName, err := secretsManagerUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				existingInstanceName = model.InstanceId
			}

			prompt := fmt.Sprintf("Are you sure you want to update instance %q?", existingInstanceName)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API - execute UpdateInstance and/or UpdateACLs based on flags
			if model.InstanceName != nil {
				req := buildUpdateInstanceRequest(ctx, model, apiClient)
				err = req.Execute()
				if err != nil {
					return fmt.Errorf("update Secrets Manager instance: %w", err)
				}
			}

			if model.Acls != nil {
				req := buildUpdateACLsRequest(ctx, model, apiClient)
				err = req.Execute()
				if err != nil {
					if model.InstanceName != nil {
						return fmt.Errorf(`the Secrets Manager instance was successfully updated, but the configuration of the ACLs failed.

If you want to retry configuring the ACLs, you can do it via: 
  $ stackit secrets-manager instance update %s --acl %s`, model.InstanceId, *model.Acls)
					}
					return fmt.Errorf("update Secrets Manager instance ACLs: %w", err)
				}
			}

			params.Printer.Info("Updated instance %q\n", existingInstanceName)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(instanceNameFlag, "n", "", "Instance name")
	cmd.Flags().Var(flags.CIDRSliceFlag(), aclFlag, "List of IP networks in CIDR notation which are allowed to access this instance")

	cmd.Flags().String(kmsKeyIdFlag, "", "ID of the KMS key to use for encryption")
	cmd.Flags().String(kmsKeyringIdFlag, "", "ID of the KMS key ring")
	cmd.Flags().Int64(kmsKeyVersionFlag, 0, "Version of the KMS key")
	cmd.Flags().String(kmsServiceAccountEmailFlag, "", "Service account email for KMS access")

	cmd.MarkFlagsRequiredTogether(kmsKeyIdFlag, kmsKeyringIdFlag, kmsKeyVersionFlag, kmsServiceAccountEmailFlag)
	cmd.MarkFlagsOneRequired(aclFlag, instanceNameFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:        globalFlags,
		InstanceId:             instanceId,
		InstanceName:           flags.FlagToStringPointer(p, cmd, instanceNameFlag),
		Acls:                   flags.FlagToStringSlicePointer(p, cmd, aclFlag),
		KmsKeyId:               flags.FlagToStringPointer(p, cmd, kmsKeyIdFlag),
		KmsKeyringId:           flags.FlagToStringPointer(p, cmd, kmsKeyringIdFlag),
		KmsKeyVersion:          flags.FlagToInt64Pointer(p, cmd, kmsKeyVersionFlag),
		KmsServiceAccountEmail: flags.FlagToStringPointer(p, cmd, kmsServiceAccountEmailFlag),
	}

	if model.KmsKeyId != nil && model.InstanceName == nil {
		return nil, fmt.Errorf("--name is required when using KMS flags")
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildUpdateInstanceRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiUpdateInstanceRequest {
	req := apiClient.UpdateInstance(ctx, model.ProjectId, model.InstanceId)

	payload := secretsmanager.UpdateInstancePayload{
		Name: model.InstanceName,
	}

	if model.KmsKeyId != nil {
		payload.KmsKey = &secretsmanager.KmsKeyPayload{
			KeyId:               model.KmsKeyId,
			KeyRingId:           model.KmsKeyringId,
			KeyVersion:          model.KmsKeyVersion,
			ServiceAccountEmail: model.KmsServiceAccountEmail,
		}
	}

	req = req.UpdateInstancePayload(payload)

	return req
}

func buildUpdateACLsRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiUpdateACLsRequest {
	req := apiClient.UpdateACLs(ctx, model.ProjectId, model.InstanceId)

	cidrs := []secretsmanager.UpdateACLPayload{}

	for _, acl := range *model.Acls {
		cidrs = append(cidrs, secretsmanager.UpdateACLPayload{Cidr: utils.Ptr(acl)})
	}

	req = req.UpdateACLsPayload(secretsmanager.UpdateACLsPayload{Cidrs: &cidrs})

	return req
}
