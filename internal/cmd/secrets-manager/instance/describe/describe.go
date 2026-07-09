package describe

import (
	"context"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	secretsmanager "github.com/stackitcloud/stackit-sdk-go/services/secretsmanager/v1api"
)

const (
	instanceIdArg = "INSTANCE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("describe %s", instanceIdArg),
		Short: "Shows details of a Secrets Manager instance",
		Long:  "Shows details of a Secrets Manager instance.",
		Args:  args.SingleArg(instanceIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Get details of a Secrets Manager instance with ID "xxx"`,
				"$ stackit secrets-manager instance describe xxx"),
			examples.NewExample(
				`Get details of a Secrets Manager instance with ID "xxx" in JSON format`,
				"$ stackit secrets-manager instance describe xxx --output-format json"),
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

			// Call API to get instance details
			req := buildGetInstanceRequest(ctx, model, apiClient)
			instance, err := req.Execute()
			if err != nil {
				return fmt.Errorf("read Secrets Manager instance: %w", err)
			}

			// Call API to get instance acls
			listACLsReq := buildListACLsRequest(ctx, model, apiClient)
			aclList, err := listACLsReq.Execute()
			if err != nil {
				return fmt.Errorf("read Secrets Manager instance ACLs: %w", err)
			}

			return outputResult(params.Printer, model.OutputFormat, instance, aclList)
		},
	}
	return cmd
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	instanceId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		InstanceId:      instanceId,
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildGetInstanceRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiGetInstanceRequest {
	req := apiClient.DefaultAPI.GetInstance(ctx, model.ProjectId, model.InstanceId)
	return req
}

func buildListACLsRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiListACLsRequest {
	req := apiClient.DefaultAPI.ListACLs(ctx, model.ProjectId, model.InstanceId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, instance *secretsmanager.Instance, aclList *secretsmanager.ListACLsResponse) error {
	if instance == nil {
		return fmt.Errorf("instance is nil")
	} else if aclList == nil {
		return fmt.Errorf("aclList is nil")
	}

	output := struct {
		*secretsmanager.Instance
		*secretsmanager.ListACLsResponse
	}{instance, aclList}

	return p.OutputResult(outputFormat, output, func() error {
		table := tables.NewTable()
		table.AddRow("ID", instance.Id)
		table.AddSeparator()
		table.AddRow("NAME", instance.Name)
		table.AddSeparator()
		table.AddRow("STATE", instance.State)
		table.AddSeparator()
		table.AddRow("SECRETS", instance.SecretCount)
		table.AddSeparator()
		table.AddRow("ENGINE", instance.SecretsEngine)
		table.AddSeparator()
		table.AddRow("CREATION DATE", instance.CreationStartDate)
		table.AddSeparator()
		kmsKey := instance.KmsKey
		if kmsKey != nil {
			table.AddRow("KMS KEY ID", kmsKey.KeyId)
			table.AddSeparator()
			table.AddRow("KMS KEYRING ID", kmsKey.KeyRingId)
			table.AddSeparator()
			table.AddRow("KMS KEY VERSION", kmsKey.KeyVersion)
			table.AddSeparator()
			table.AddRow("KMS SERVICE ACCOUNT EMAIL", kmsKey.ServiceAccountEmail)
		}
		// Only show ACL if it's present and not empty
		if len(aclList.Acls) > 0 {
			var cidrs []string

			for _, acl := range aclList.Acls {
				cidrs = append(cidrs, acl.Cidr)
			}

			if kmsKey != nil {
				table.AddSeparator()
			}
			table.AddRow("ACL", strings.Join(cidrs, ","))
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	})
}
