package describe

import (
	"context"
	"encoding/json"
	"fmt"

	"strings"

	"github.com/ghodss/yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/secrets-manager/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/secretsmanager"
)

const (
	instanceIdArg = "INSTANCE_ID"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	InstanceId string
}

func NewCmd(p *print.Printer) *cobra.Command {
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
			model, err := parseInput(p, cmd, args)
			if err != nil {
				return err
			}
			// Configure API client
			apiClient, err := client.ConfigureClient(p)
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

			return outputResult(p, model.OutputFormat, instance, aclList)
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

func buildGetInstanceRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiGetInstanceRequest {
	req := apiClient.GetInstance(ctx, model.ProjectId, model.InstanceId)
	return req
}

func buildListACLsRequest(ctx context.Context, model *inputModel, apiClient *secretsmanager.APIClient) secretsmanager.ApiListACLsRequest {
	req := apiClient.ListACLs(ctx, model.ProjectId, model.InstanceId)
	return req
}

func outputResult(p *print.Printer, outputFormat string, instance *secretsmanager.Instance, aclList *secretsmanager.AclList) error {
	output := struct {
		*secretsmanager.Instance
		*secretsmanager.AclList
	}{instance, aclList}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Secrets Manager instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.Marshal(output)
		if err != nil {
			return fmt.Errorf("marshal Secrets Manager instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("ID", *instance.Id)
		table.AddSeparator()
		table.AddRow("NAME", *instance.Name)
		table.AddSeparator()
		table.AddRow("STATE", *instance.State)
		table.AddSeparator()
		table.AddRow("SECRETS", *instance.SecretCount)
		table.AddSeparator()
		table.AddRow("ENGINE", *instance.SecretsEngine)
		table.AddSeparator()
		table.AddRow("CREATION DATE", *instance.CreationStartDate)
		table.AddSeparator()
		// Only show ACL if it's present and not empty
		if aclList != nil && aclList.Acls != nil && len(*aclList.Acls) > 0 {
			var cidrs []string

			for _, acl := range *aclList.Acls {
				cidrs = append(cidrs, *acl.Cidr)
			}

			table.AddRow("ACL", strings.Join(cidrs, ","))
		}
		err := table.Display(p)
		if err != nil {
			return fmt.Errorf("render table: %w", err)
		}

		return nil
	}
}
