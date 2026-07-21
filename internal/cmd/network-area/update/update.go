package update

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	rmClient "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	rmUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	areaIdArg = "AREA_ID"

	nameFlag           = "name"
	organizationIdFlag = "organization-id"
	areaIdFlag         = "area-id"
	labelFlag          = "labels"
)

// NetworkAreaResponses is a workaround, to keep the two responses of the iaas v2 api together for the json and yaml output
// Should be removed when the deprecated flags are removed
type NetworkAreaResponses struct {
	NetworkArea  iaas.NetworkArea   `json:"network_area"`
	RegionalArea *iaas.RegionalArea `json:"regional_area"`
}

type inputModel struct {
	*globalflags.GlobalFlagModel
	AreaId         string
	Name           *string
	OrganizationId *string
	Labels         map[string]any
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", areaIdArg),
		Short: "Updates a STACKIT Network Area (SNA)",
		Long:  "Updates a STACKIT Network Area (SNA) in an organization.",
		Args:  args.SingleArg(areaIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update network area with ID "xxx" in organization with ID "yyy" with new name "network-area-1-new"`,
				"$ stackit network-area update xxx --organization-id yyy --name network-area-1-new",
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

			var orgLabel string
			rmApiClient, err := rmClient.ConfigureClient(params.Printer, params.CliVersion)
			if err == nil {
				orgLabel, err = rmUtils.GetOrganizationName(ctx, rmApiClient.DefaultAPI, *model.OrganizationId)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get organization name: %v", err)
					orgLabel = *model.OrganizationId
				} else if orgLabel == "" {
					orgLabel = *model.OrganizationId
				}
			} else {
				params.Printer.Debug(print.ErrorLevel, "configure resource manager client: %v", err)
			}

			prompt := fmt.Sprintf("Are you sure you want to update a network area for organization %q?", orgLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update network area: %w", err)
			}

			if resp == nil || resp.Id == nil {
				return fmt.Errorf("update network area: empty response")
			}

			responses := NetworkAreaResponses{
				NetworkArea: *resp,
			}

			return outputResult(params.Printer, model.OutputFormat, orgLabel, responses)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "Network area name")
	cmd.Flags().Var(flags.UUIDFlag(), organizationIdFlag, "Organization ID")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a network-area. E.g. '--labels key1=value1,key2=value2,...'")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	areaId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Name:            flags.FlagToStringPointer(p, cmd, nameFlag),
		OrganizationId:  flags.FlagToStringPointer(p, cmd, organizationIdFlag),
		AreaId:          areaId,
		Labels:          flags.FlagToStringToAny(p, cmd, labelFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiPartialUpdateNetworkAreaRequest {
	req := apiClient.DefaultAPI.PartialUpdateNetworkArea(ctx, *model.OrganizationId, model.AreaId)

	payload := iaas.PartialUpdateNetworkAreaPayload{
		Name:   model.Name,
		Labels: model.Labels,
	}

	return req.PartialUpdateNetworkAreaPayload(payload)
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, responses NetworkAreaResponses) error {
	prettyOutputFunc := func() error {
		p.Outputf("Updated STACKIT Network Area for project %q.\n", projectLabel)
		return nil
	}

	// If RegionalArea is NOT set in the reponses, then no deprecated Flags were set.
	// In this case, only the response of NetworkArea should be printed in JSON and yaml output, to avoid breaking changes after the deprecated fields are removed
	if responses.RegionalArea == nil {
		return p.OutputResult(outputFormat, responses.NetworkArea, prettyOutputFunc)
	}

	return p.OutputResult(outputFormat, responses, func() error {
		p.Outputf("Updated STACKIT Network Area for project %q.\n", projectLabel)
		return nil
	})
}
