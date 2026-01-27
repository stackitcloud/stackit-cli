package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/intake/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/intake"
	"github.com/stackitcloud/stackit-sdk-go/services/intake/wait"
)

const (
	displayNameFlag         = "display-name"
	runnerIdFlag            = "runner-id"
	descriptionFlag         = "description"
	labelsFlag              = "labels"
	catalogURIFlag          = "catalog-uri"
	catalogWarehouseFlag    = "catalog-warehouse"
	catalogNamespaceFlag    = "catalog-namespace"
	catalogTableNameFlag    = "catalog-table-name"
	catalogPartitioningFlag = "catalog-partitioning"
	catalogPartitionByFlag  = "catalog-partition-by"
	catalogAuthTypeFlag     = "catalog-auth-type"
	dremioTokenEndpointFlag = "dremio-token-endpoint" //nolint:gosec // false positive
	dremioPatFlag           = "dremio-pat"
)

// inputModel struct holds all the input parameters for the command
type inputModel struct {
	*globalflags.GlobalFlagModel

	// Top-level fields
	DisplayName *string
	RunnerId    *string
	Description *string
	Labels      *map[string]string

	// Catalog fields
	CatalogURI          *string
	CatalogWarehouse    *string
	CatalogNamespace    *string
	CatalogTableName    *string
	CatalogPartitioning *string
	CatalogPartitionBy  *[]string

	// Auth fields
	CatalogAuthType     *string
	DremioTokenEndpoint *string
	DremioToken         *string
}

func NewCmd(p *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a new Intake",
		Long:  "Creates a new Intake.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a new Intake with required parameters`,
				`$ stackit beta intake create --display-name my-intake --runner-id xxx --catalog-uri "http://dremio.example.com" --catalog-warehouse "my-warehouse" --catalog-auth-type none`),
			examples.NewExample(
				`Create a new Intake with a description, labels, and Dremio authentication`,
				`$ stackit beta intake create --display-name my-intake --runner-id xxx --description "Production intake" --labels "env=prod,team=billing" --catalog-uri "http://dremio.example.com" --catalog-warehouse "my-warehouse" --catalog-auth-type "dremio" --dremio-token-endpoint "https://auth.dremio.cloud/oauth/token" --dremio-pat "MY_TOKEN"`),
			examples.NewExample(
				`Create a new Intake with manual partitioning by a date field`,
				`$ stackit beta intake create --display-name my-partitioned-intake --runner-id xxx --catalog-uri "http://dremio.example.com" --catalog-warehouse "my-warehouse" --catalog-partitioning "manual" --catalog-partition-by "day(__intake_ts)" --catalog-auth-type "none"`),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(p.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(p.Printer, p.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, p.Printer, p.CliVersion, cmd)
			if err != nil {
				p.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			prompt := fmt.Sprintf("Are you sure you want to create an Intake for project %q?", projectLabel)
			err = p.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create Intake: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(p.Printer)
				s.Start("Creating STACKIT Intake instance")
				_, err = wait.CreateOrUpdateIntakeWaitHandler(ctx, apiClient, model.ProjectId, model.Region, resp.GetId()).WaitWithContext(ctx)
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
	cmd.Flags().StringToString(labelsFlag, nil, "Labels in key=value format, separated by commas. Example: --labels \"key1=value1,key2=value2\"")

	// Catalog flags
	cmd.Flags().String(catalogURIFlag, "", "The URI to the Iceberg catalog endpoint")
	cmd.Flags().String(catalogWarehouseFlag, "", "The Iceberg warehouse to connect to")
	cmd.Flags().String(catalogNamespaceFlag, "", "The namespace to which data shall be written (default: 'intake')")
	cmd.Flags().String(catalogTableNameFlag, "", "The table name to identify the table in Iceberg")
	cmd.Flags().String(catalogPartitioningFlag, "", "The target table's partitioning. One of 'none', 'intake-time', 'manual'")
	cmd.Flags().StringSlice(catalogPartitionByFlag, nil, "List of Iceberg partitioning expressions. Only used when --catalog-partitioning is 'manual'")

	// Auth flags
	cmd.Flags().String(catalogAuthTypeFlag, "", "Authentication type for the catalog (e.g., 'none', 'dremio')")
	cmd.Flags().String(dremioTokenEndpointFlag, "", "Dremio OAuth 2.0 token endpoint URL. Required if auth-type is 'dremio'")
	cmd.Flags().String(dremioPatFlag, "", "Dremio personal access token. Required if auth-type is 'dremio'")

	err := flags.MarkFlagsRequired(cmd, displayNameFlag, runnerIdFlag, catalogURIFlag, catalogWarehouseFlag, catalogAuthTypeFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		// Top-level fields
		DisplayName: flags.FlagToStringPointer(p, cmd, displayNameFlag),
		RunnerId:    flags.FlagToStringPointer(p, cmd, runnerIdFlag),
		Description: flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:      flags.FlagToStringToStringPointer(p, cmd, labelsFlag),

		// Catalog fields
		CatalogURI:          flags.FlagToStringPointer(p, cmd, catalogURIFlag),
		CatalogWarehouse:    flags.FlagToStringPointer(p, cmd, catalogWarehouseFlag),
		CatalogNamespace:    flags.FlagToStringPointer(p, cmd, catalogNamespaceFlag),
		CatalogTableName:    flags.FlagToStringPointer(p, cmd, catalogTableNameFlag),
		CatalogPartitioning: flags.FlagToStringPointer(p, cmd, catalogPartitioningFlag),
		CatalogPartitionBy:  flags.FlagToStringSlicePointer(p, cmd, catalogPartitionByFlag),

		// Auth fields
		CatalogAuthType:     flags.FlagToStringPointer(p, cmd, catalogAuthTypeFlag),
		DremioTokenEndpoint: flags.FlagToStringPointer(p, cmd, dremioTokenEndpointFlag),
		DremioToken:         flags.FlagToStringPointer(p, cmd, dremioPatFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *intake.APIClient) intake.ApiCreateIntakeRequest {
	req := apiClient.CreateIntake(ctx, model.ProjectId, model.Region)

	// Build catalog authentication
	var catalogAuth *intake.CatalogAuth
	if model.CatalogAuthType != nil {
		authType := intake.CatalogAuthType(*model.CatalogAuthType)
		catalogAuth = &intake.CatalogAuth{
			Type: &authType,
		}
		if *model.CatalogAuthType == "dremio" {
			catalogAuth.Dremio = &intake.DremioAuth{
				TokenEndpoint:       model.DremioTokenEndpoint,
				PersonalAccessToken: model.DremioToken,
			}
		}
	}

	var partitioning *intake.PartitioningType
	if model.CatalogPartitioning != nil {
		partitioning = utils.Ptr(intake.PartitioningType(*model.CatalogPartitioning))
	}

	// Build catalog
	catalogPayload := intake.IntakeCatalog{
		Uri:          model.CatalogURI,
		Warehouse:    model.CatalogWarehouse,
		Namespace:    model.CatalogNamespace,
		TableName:    model.CatalogTableName,
		Partitioning: partitioning,
		PartitionBy:  model.CatalogPartitionBy,
		Auth:         catalogAuth,
	}

	// Build main payload
	payload := intake.CreateIntakePayload{
		DisplayName:    model.DisplayName,
		IntakeRunnerId: model.RunnerId,
		Description:    model.Description,
		Labels:         model.Labels,
		Catalog:        &catalogPayload,
	}
	req = req.CreateIntakePayload(payload)

	return req
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, resp *intake.IntakeResponse) error {
	return p.OutputResult(model.OutputFormat, resp, func() error {
		if resp == nil {
			p.Outputf("Triggered creation of Intake for project %q, but no intake ID was returned.\n", projectLabel)
			return nil
		}

		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s Intake for project %q. Intake ID: %s\n", operationState, projectLabel, utils.PtrString(resp.Id))
		return nil
	})
}
