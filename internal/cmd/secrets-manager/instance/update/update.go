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

	aclFlag = "acl"

	kmsKeyIdFlag               = "kms-key-id"
	kmsKeyringIdFlag           = "kms-keyring-id"
	kmsKeyVersionFlag          = "kms-key-version"
	kmsServiceAccountEmailFlag = "kms-service-account-email"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string

	Acls *[]string

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
				`Update the range of IPs allowed to access a Secrets Manager instance with ID "xxx"`,
				"$ stackit secrets-manager instance update xxx --acl 1.2.3.0/24"),
			examples.NewExample(
				`Update the KMS key settings of a Secrets Manager instance with ID "xxx"`,
				"$ stackit secrets-manager instance update xxx --kms-key-id key-id --kms-keyring-id keyring-id --kms-key-version 1 --kms-service-account-email my-service-account-1234567@sa.stackit.cloud"),
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

			instanceLabel, err := secretsManagerUtils.GetInstanceName(ctx, apiClient, model.ProjectId, model.InstanceId)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get instance name: %v", err)
				instanceLabel = model.InstanceId
				if model.KmsKeyId != nil {
					return fmt.Errorf("get instance name: %w", err)
				}
			}

			prompt := fmt.Sprintf("Are you sure you want to update instance %q?", instanceLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, instanceLabel, apiClient)
			switch request := req.(type) {
			case secretsmanager.ApiUpdateInstanceRequest:
				err = request.Execute()
			case secretsmanager.ApiUpdateACLsRequest:
				err = request.Execute()
			default:
				err = fmt.Errorf("unknown request type")
			}
			if err != nil {
				return fmt.Errorf("update Secrets Manager instance: %w", err)
			}

			params.Printer.Info("Updated instance %q\n", instanceLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Var(flags.CIDRSliceFlag(), aclFlag, "List of IP networks in CIDR notation which are allowed to access this instance")

	cmd.Flags().String(kmsKeyIdFlag, "", "ID of the KMS key to use for encryption")
	cmd.Flags().String(kmsKeyringIdFlag, "", "ID of the KMS key ring")
	cmd.Flags().Int64(kmsKeyVersionFlag, 0, "Version of the KMS key")
	cmd.Flags().String(kmsServiceAccountEmailFlag, "", "Service account email for KMS access")

	cmd.MarkFlagsRequiredTogether(kmsKeyIdFlag, kmsKeyringIdFlag, kmsKeyVersionFlag, kmsServiceAccountEmailFlag)
	cmd.MarkFlagsMutuallyExclusive(aclFlag, kmsKeyIdFlag)
	cmd.MarkFlagsOneRequired(aclFlag, kmsKeyIdFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	acls := flags.FlagToStringSlicePointer(p, cmd, aclFlag)

	model := inputModel{
		GlobalFlagModel:        globalFlags,
		InstanceId:             instanceId,
		Acls:                   acls,
		KmsKeyId:               flags.FlagToStringPointer(p, cmd, kmsKeyIdFlag),
		KmsKeyringId:           flags.FlagToStringPointer(p, cmd, kmsKeyringIdFlag),
		KmsKeyVersion:          flags.FlagToInt64Pointer(p, cmd, kmsKeyVersionFlag),
		KmsServiceAccountEmail: flags.FlagToStringPointer(p, cmd, kmsServiceAccountEmailFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, instanceName string, apiClient *secretsmanager.APIClient) interface{ Execute() error } {
	if model.KmsKeyId != nil {
		return buildUpdateInstanceRequest(ctx, model, instanceName, apiClient)
	}

	return buildUpdateACLsRequest(ctx, model, apiClient)
}

func buildUpdateInstanceRequest(ctx context.Context, model *inputModel, instanceName string, apiClient *secretsmanager.APIClient) secretsmanager.ApiUpdateInstanceRequest {
	req := apiClient.UpdateInstance(ctx, model.ProjectId, model.InstanceId)

	payload := secretsmanager.UpdateInstancePayload{
		Name: &instanceName,
		KmsKey: &secretsmanager.KmsKeyPayload{
			KeyId:               model.KmsKeyId,
			KeyRingId:           model.KmsKeyringId,
			KeyVersion:          model.KmsKeyVersion,
			ServiceAccountEmail: model.KmsServiceAccountEmail,
		},
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
