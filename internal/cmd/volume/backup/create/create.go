package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasutils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

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

var sourceTypeFlagOptions = []string{"volume", "snapshot"}

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
				"$ stackit volume backup create --source-id xxx --source-type volume"),
			examples.NewExample(
				`Create a backup from a snapshot with a name`,
				"$ stackit volume backup create --source-id xxx --source-type snapshot --name my-backup"),
			examples.NewExample(
				`Create a backup with labels`,
				"$ stackit volume backup create --source-id xxx --source-type volume --labels key1=value1,key2=value2"),
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			// Get source name for label (use ID if name not available)
			sourceLabel := model.SourceID
			if model.SourceType == "volume" {
				name, err := iaasutils.GetVolumeName(ctx, apiClient, model.ProjectId, model.Region, model.SourceID)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get volume name: %v", err)
				} else if name != "" {
					sourceLabel = name
				}
			} else if model.SourceType == "snapshot" {
				name, err := iaasutils.GetSnapshotName(ctx, apiClient, model.ProjectId, model.Region, model.SourceID)
				if err != nil {
					params.Printer.Debug(print.ErrorLevel, "get snapshot name: %v", err)
				} else if name != "" {
					sourceLabel = name
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
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create volume backup: %w", err)
			}
			if resp == nil || resp.Id == nil {
				return fmt.Errorf("create volume: empty response")
			}
			volumeId := *resp.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating backup")
				resp, err = wait.CreateBackupWaitHandler(ctx, apiClient, model.ProjectId, model.Region, volumeId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for backup creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model.OutputFormat, model.Async, sourceLabel, projectLabel, resp)
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(sourceIdFlag, "", "ID of the source from which a backup should be created")
	cmd.Flags().Var(flags.EnumFlag(false, "", sourceTypeFlagOptions...), sourceTypeFlag, fmt.Sprintf("Source type of the backup, one of %q", sourceTypeFlagOptions))
	cmd.Flags().String(nameFlag, "", "Name of the backup")
	cmd.Flags().StringToString(labelsFlag, nil, "Key-value string pairs as labels")

	err := flags.MarkFlagsRequired(cmd, sourceIdFlag, sourceTypeFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	sourceID := flags.FlagToStringValue(p, cmd, sourceIdFlag)
	if sourceID == "" {
		return nil, fmt.Errorf("source-id is required")
	}

	sourceType := flags.FlagToStringValue(p, cmd, sourceTypeFlag)

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

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateBackupRequest {
	req := apiClient.CreateBackup(ctx, model.ProjectId, model.Region)

	payload := iaas.CreateBackupPayload{
		Name:   model.Name,
		Labels: utils.ConvertStringMapToInterfaceMap(utils.Ptr(model.Labels)),
		Source: &iaas.BackupSource{
			Id:   &model.SourceID,
			Type: &model.SourceType,
		},
	}

	return req.CreateBackupPayload(payload)
}

func outputResult(p *print.Printer, outputFormat string, async bool, sourceLabel, projectLabel string, resp *iaas.Backup) error {
	if resp == nil {
		return fmt.Errorf("create backup response is empty")
	}

	return p.OutputResult(outputFormat, resp, func() error {
		if async {
			p.Outputf("Triggered backup of %s in %s. Backup ID: %s\n", sourceLabel, projectLabel, utils.PtrString(resp.Id))
		} else {
			p.Outputf("Created backup of %s in %s. Backup ID: %s\n", sourceLabel, projectLabel, utils.PtrString(resp.Id))
		}
		return nil
	})
}
