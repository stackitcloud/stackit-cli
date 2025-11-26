package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"
	"github.com/stackitcloud/stackit-sdk-go/services/intake/wait"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

const (
	intakeIdArg = "INTAKE_ID"

	// Top-level flags
	displayNameFlag = "display-name"
	runnerIdFlag    = "runner-id"
	descriptionFlag = "description"
	labelsFlag      = "labels"

	// Catalog flags
	catalogURIFlag       = "catalog-uri"
	catalogWarehouseFlag = "catalog-warehouse"
	catalogNamespaceFlag = "catalog-namespace"
	catalogTableNameFlag = "catalog-table-name"

	// Auth flags
	catalogAuthTypeFlag     = "catalog-auth-type"
	dremioTokenEndpointFlag = "dremio-token-endpoint" //nolint:gosec // false positive
	dremioPatFlag           = "dremio-pat"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	// Main
	IntakeId    string
	DisplayName *string
	RunnerId    *string
	Description *string
	Labels      *map[string]string

	// Catalog
	CatalogURI       *string
	CatalogWarehouse *string
	CatalogNamespace *string
	CatalogTableName *string

	// Auth
	CatalogAuthType     *string
	DremioTokenEndpoint *string
	DremioToken         *string
}

func NewCmd(p *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", intakeIdArg),
		Short: "Updates an Intake",
		Long:  "Updates an Intake. Only the specified fields are updated.",
		Args:  args.SingleArg(intakeIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(
				`Update the display name of an Intake with ID "xxx"`,
				`$ stackit beta intake update xxx --runner-id yyy --display-name new-intake-name`),
			examples.NewExample(
				`Update the catalog details for an Intake with ID "xxx"`,
				`$ stackit beta intake update xxx --runner-id yyy --catalog-uri "http://new.uri" --catalog-warehouse "new-warehouse"`),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd, args)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, p.Printer, p.CliVersion, cmd)
			if err != nil {
				p.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update Intake: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p.Printer)
				s.Start("Updating STACKIT Intake Runner instance")
				_, err = wait.CreateOrUpdateIntakeWaitHandler(ctx, apiClient, model.ProjectId, model.Region, model.IntakeId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for STACKIT Instance creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(p.Printer, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	// Top-level flags
	cmd.Flags().String(displayNameFlag, "", "Display name")
	cmd.Flags().Var(flags.UUIDFlag(), runnerIdFlag, "The UUID of the Intake Runner to use")
	cmd.Flags().String(descriptionFlag, "", "Description")
	cmd.Flags().StringToString(labelsFlag, nil, `Labels in key=value format, separated by commas. Example: --labels "key1=value1,key2=value2".`)

	// Catalog flags
	cmd.Flags().String(catalogURIFlag, "", "The URI to the Iceberg catalog endpoint")
	cmd.Flags().String(catalogWarehouseFlag, "", "The Iceberg warehouse to connect to")
	cmd.Flags().String(catalogNamespaceFlag, "", "The namespace to which data shall be written")
	cmd.Flags().String(catalogTableNameFlag, "", "The table name to identify the table in Iceberg")

	// Auth flags
	cmd.Flags().String(catalogAuthTypeFlag, "", "Authentication type for the catalog (e.g., 'none', 'dremio')")
	cmd.Flags().String(dremioTokenEndpointFlag, "", "Dremio OAuth 2.0 token endpoint URL")
	cmd.Flags().String(dremioPatFlag, "", "Dremio personal access token")

	err := flags.MarkFlagsRequired(cmd, runnerIdFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	intakeId := inputArgs[0]

	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := &inputModel{
		GlobalFlagModel:     globalFlags,
		IntakeId:            intakeId,
		DisplayName:         flags.FlagToStringPointer(p, cmd, displayNameFlag),
		RunnerId:            flags.FlagToStringPointer(p, cmd, runnerIdFlag),
		Description:         flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:              flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
		CatalogURI:          flags.FlagToStringPointer(p, cmd, catalogURIFlag),
		CatalogWarehouse:    flags.FlagToStringPointer(p, cmd, catalogWarehouseFlag),
		CatalogNamespace:    flags.FlagToStringPointer(p, cmd, catalogNamespaceFlag),
		CatalogTableName:    flags.FlagToStringPointer(p, cmd, catalogTableNameFlag),
		CatalogAuthType:     flags.FlagToStringPointer(p, cmd, catalogAuthTypeFlag),
		DremioTokenEndpoint: flags.FlagToStringPointer(p, cmd, dremioTokenEndpointFlag),
		DremioToken:         flags.FlagToStringPointer(p, cmd, dremioPatFlag),
	}

	// Check if any optional flag was provided
	if model.DisplayName == nil && model.Description == nil && model.Labels == nil &&
		model.CatalogURI == nil && model.CatalogWarehouse == nil && model.CatalogNamespace == nil &&
		model.CatalogTableName == nil && model.CatalogAuthType == nil &&
		model.DremioTokenEndpoint == nil && model.DremioToken == nil {
		return nil, &cliErr.EmptyUpdateError{}
	}

	p.DebugInputModel(model)
	return model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiUpdateIntakeRequest {
	req := apiClient.UpdateIntake(ctx, model.ProjectId, model.Region, model.IntakeId)

	payload := intake.UpdateIntakePayload{
		IntakeRunnerId: model.RunnerId, // This is required by the API
		DisplayName:    model.DisplayName,
		Description:    model.Description,
		Labels:         model.Labels,
	}

	// Build catalog patch payload only if catalog-related flags are set
	catalogPatch := &intake.IntakeCatalogPatch{}
	catalogNeedsPatching := false

	if model.CatalogURI != nil {
		catalogPatch.Uri = model.CatalogURI
		catalogNeedsPatching = true
	}
	if model.CatalogWarehouse != nil {
		catalogPatch.Warehouse = model.CatalogWarehouse
		catalogNeedsPatching = true
	}
	if model.CatalogNamespace != nil {
		catalogPatch.Namespace = model.CatalogNamespace
		catalogNeedsPatching = true
	}
	if model.CatalogTableName != nil {
		catalogPatch.TableName = model.CatalogTableName
		catalogNeedsPatching = true
	}

	// Build auth patch payload only if auth-related flags are set
	authPatch := &intake.CatalogAuthPatch{}
	authNeedsPatching := false

	if model.CatalogAuthType != nil {
		authType := intake.CatalogAuthType(*model.CatalogAuthType)
		authPatch.Type = &authType
		authNeedsPatching = true
	}

	dremioPatch := &intake.DremioAuthPatch{}
	dremioNeedsPatching := false
	if model.DremioTokenEndpoint != nil {
		dremioPatch.TokenEndpoint = model.DremioTokenEndpoint
		dremioNeedsPatching = true
	}
	if model.DremioToken != nil {
		dremioPatch.PersonalAccessToken = model.DremioToken
		dremioNeedsPatching = true
	}

	if dremioNeedsPatching {
		authPatch.Dremio = dremioPatch
		authNeedsPatching = true
	}

	if authNeedsPatching {
		catalogPatch.Auth = authPatch
		catalogNeedsPatching = true
	}

	if catalogNeedsPatching {
		payload.Catalog = catalogPatch
	}

	req = req.UpdateIntakePayload(payload)
	return req
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, resp *intake.IntakeResponse) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		if resp == nil {
			p.Outputf("Updated Intake for project %q, but no intake ID was returned.\n", projectLabel)
			return nil
		}

		operationState := "Updated"
		if model.Async {
			operationState = "Triggered update of"
		}
		p.Outputf("%s Intake for project %q. Intake ID: %s\n", operationState, projectLabel, utils.PtrString(resp.Id))
		return nil
	})
}
