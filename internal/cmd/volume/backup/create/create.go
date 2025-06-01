package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"
)

const (
	sourceIdFlag   = "source-id"
	sourceTypeFlag = "source-type"
	nameFlag       = "name"
	labelsFlag     = "labels"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	SourceID   string
	SourceType string
	Name       *string
	Labels     map[string]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a backup from a specific source",
		Long:  "Creates a backup from a specific source (volume or snapshot).",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a backup from a volume`,
				"$ stackit volume backup create --source-id xxx --source-type volume --project-id xxx"),
			examples.NewExample(
				`Create a backup from a snapshot with a name`,
				"$ stackit volume backup create --source-id xxx --source-type snapshot --name my-backup --project-id xxx"),
			examples.NewExample(
				`Create a backup with labels`,
				"$ stackit volume backup create --source-id xxx --source-type volume --labels key1=value1,key2=value2 --project-id xxx"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer, params.CliVersion)
			if err != nil {
				return err
			}

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			// Get source name for label (use ID if name not available)
			sourceLabel := model.SourceID
			if model.SourceType == "volume" {
				volume, err := apiClient.GetVolume(ctx, model.ProjectId, model.SourceID).Execute()
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get volume name: %v", err)
				} else if volume != nil && volume.Name != nil {
					sourceLabel = *volume.Name
				}
			} else if model.SourceType == "snapshot" {
				snapshot, err := apiClient.GetSnapshot(ctx, model.ProjectId, model.SourceID).Execute()
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get snapshot name: %v", err)
				} else if snapshot != nil && snapshot.Name != nil {
					sourceLabel = *snapshot.Name
				}
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create backup from %s? (This cannot be undone)", sourceLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(model, apiClient, ctx)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create volume backup: %w", err)
			}

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating backup")
				resp, err = wait.CreateBackupWaitHandler(ctx, apiClient, model.ProjectId, *resp.Id).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for backup creation: %w", err)
				}
				s.Stop()
			}

			if model.Async {
				params.Printer.Info("Triggered backup of %q in %q. Backup ID: %s\n", sourceLabel, projectLabel, *resp.Id)
			} else {
				params.Printer.Info("Created backup of %q in %q. Backup ID: %s\n", sourceLabel, projectLabel, *resp.Id)
			}
			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(sourceIdFlag, "", "ID of the source from which a backup should be created")
	cmd.Flags().String(sourceTypeFlag, "", "Source type of the backup (volume or snapshot)")
	cmd.Flags().String(nameFlag, "", "Name of the backup")
	cmd.Flags().StringToString(labelsFlag, nil, "Key-value string pairs as labels")

	err := flags.MarkFlagsRequired(cmd, sourceIdFlag, sourceTypeFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	sourceID := flags.FlagToStringValue(p, cmd, sourceIdFlag)
	if sourceID == "" {
		return nil, fmt.Errorf("source-id is required")
	}

	sourceType := flags.FlagToStringValue(p, cmd, sourceTypeFlag)
	if sourceType != "volume" && sourceType != "snapshot" {
		return nil, fmt.Errorf("source-type must be either 'volume' or 'snapshot'")
	}

	name := flags.FlagToStringPointer(p, cmd, nameFlag)
	labels := flags.FlagToStringToStringPointer(p, cmd, labelsFlag)
	if labels == nil {
		labels = &map[string]string{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		SourceID:        sourceID,
		SourceType:      sourceType,
		Name:            name,
		Labels:          *labels,
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

// TODO: Enough?
func buildRequest(model *inputModel, apiClient *iaas.APIClient, ctx context.Context) iaas.ApiCreateBackupRequest {
	// TODO: doc says if createeBackup func provides snapshot-id but isnt in the func-signature?
	req := apiClient.CreateBackup(ctx, model.ProjectId)
	return req
}

// TODO: create(volume)BackupResponse or createBackupResponse needs to be created
func outputResult(p *print.Printer, outputFormat string, async bool, sourceLabel, projectLabel string, resp *iaas.Backup) error {
	if resp == nil {
		return fmt.Errorf("create backup response is empty")
	}

	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal backup: %w", err)
		}
		p.Outputln(string(details))
		return nil

	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal backup: %w", err)
		}
		p.Outputln(string(details))
		return nil

	default:
		if async {
			p.Outputf("Triggered backup of %s in %s. Backup ID: %s\n", sourceLabel, projectLabel, *resp.Id)
		} else {
			p.Outputf("Created backup of %s in %s. Backup ID: %s\n", sourceLabel, projectLabel, *resp.Id)
		}
		return nil
	}
}
