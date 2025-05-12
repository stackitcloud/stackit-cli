package describe

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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

func NewCmd(params *params.CmdParams) *cobra.Command {
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

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal Secrets Manager instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(output, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal Secrets Manager instance: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		table := tables.NewTable()
		table.AddRow("ID", utils.PtrString(instance.Id))
		table.AddSeparator()
		table.AddRow("NAME", utils.PtrString(instance.Name))
		table.AddSeparator()
		table.AddRow("STATE", utils.PtrString(instance.State))
		table.AddSeparator()
		table.AddRow("SECRETS", utils.PtrString(instance.SecretCount))
		table.AddSeparator()
		table.AddRow("ENGINE", utils.PtrString(instance.SecretsEngine))
		table.AddSeparator()
		table.AddRow("CREATION DATE", utils.PtrString(instance.CreationStartDate))
		table.AddSeparator()
		// Only show ACL if it's present and not empty
		if aclList.Acls != nil && len(*aclList.Acls) > 0 {
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
