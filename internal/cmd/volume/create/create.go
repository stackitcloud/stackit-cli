package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/types"

	iaas "github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/v2api/wait"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/spinner"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
)

const (
	availabilityZoneFlag = "availability-zone"
	nameFlag             = "name"
	descriptionFlag      = "description"
	labelFlag            = "labels"
	performanceClassFlag = "performance-class"
	sizeFlag             = "size"
	sourceIdFlag         = "source-id"
	sourceTypeFlag       = "source-type"

	encryptionKekKeyIdFlag       = "encryption-kek-key-id"
	encryptionKekKeyringIdFlag   = "encryption-kek-keyring-id"
	encryptionKekKeyVersionFlag  = "encryption-kek-key-version"
	encryptionServiceAccountFlag = "encryption-service-account"
	encryptionKeyPayloadFlag     = "encryption-key-payload-base64"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	AvailabilityZone string
	Name             *string
	Description      *string
	Labels           map[string]any
	PerformanceClass *string
	Size             *int64
	SourceId         *string
	SourceType       *string

	EncryptionKekKeyId       *string
	EncryptionKekKeyringId   *string
	EncryptionKekKeyVersion  *int64
	EncryptionServiceAccount *string
	EncryptionKeyPayload     *string `json:"-"`
}

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a volume",
		Long:  "Creates a volume.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a volume with availability zone "eu01-1" and size 64 GB`,
				`$ stackit volume create --availability-zone eu01-1 --size 64`,
			),
			examples.NewExample(
				`Create a volume with availability zone "eu01-1", size 64 GB and labels`,
				`$ stackit volume create --availability-zone eu01-1 --size 64 --labels key=value,foo=bar`,
			),
			examples.NewExample(
				`Create a volume with name "volume-1", from a source image with ID "xxx"`,
				`$ stackit volume create --availability-zone eu01-1 --name volume-1 --source-id xxx --source-type image`,
			),
			examples.NewExample(
				`Create a volume with availability zone "eu01-1", performance class "storage_premium_perf1" and size 64 GB`,
				`$ stackit volume create --availability-zone eu01-1 --performance-class storage_premium_perf1 --size 64`,
			),
			examples.NewExample(
				`Create an encrypted volume using a STACKIT KMS key`,
				`$ stackit volume create --availability-zone eu01-1 --size 64 --encryption-kek-key-id xxx --encryption-kek-keyring-id yyy --encryption-kek-key-version 1 --encryption-service-account service-account@example.com`,
			),
			examples.NewExample(
				`Create an encrypted volume with a predefined base64-encoded secret read from a file`,
				`$ stackit volume create --availability-zone eu01-1 --size 64 --encryption-kek-key-id xxx --encryption-kek-keyring-id yyy --encryption-kek-key-version 1 --encryption-service-account service-account@example.com --encryption-key-payload-base64 @./secret.txt`,
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

			projectLabel, err := projectname.GetProjectName(ctx, params.Printer, params.CliVersion, cmd)
			if err != nil {
				params.Printer.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			prompt := fmt.Sprintf("Are you sure you want to create a volume for project %q?", projectLabel)
			err = params.Printer.PromptForConfirmation(prompt)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient.DefaultAPI)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create volume : %w", err)
			}
			volumeId := *resp.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				err := spinner.Run(params.Printer, "Creating volume", func() error {
					_, err = wait.CreateVolumeWaitHandler(ctx, apiClient.DefaultAPI, model.ProjectId, model.Region, volumeId).WaitWithContext(ctx)
					return err
				})
				if err != nil {
					return fmt.Errorf("wait for volume creation: %w", err)
				}
			}

			return outputResult(params.Printer, model, projectLabel, resp)
		},
	}
	configureFlags(cmd, params)
	return cmd
}

func configureFlags(cmd *cobra.Command, params *types.CmdParams) {
	cmd.Flags().String(availabilityZoneFlag, "", "Availability zone")
	cmd.Flags().StringP(nameFlag, "n", "", "Volume name")
	cmd.Flags().String(descriptionFlag, "", "Volume description")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a volume. E.g. '--labels key1=value1,key2=value2,...'")
	cmd.Flags().String(performanceClassFlag, "", "Performance class")
	cmd.Flags().Int64(sizeFlag, 0, "Volume size (GB). Either 'size' or the 'source-id' and 'source-type' flags must be given")
	cmd.Flags().String(sourceIdFlag, "", "ID of the source object of volume. Either 'size' or the 'source-id' and 'source-type' flags must be given")
	cmd.Flags().String(sourceTypeFlag, "", "Type of the source object of volume. Either 'size' or the 'source-id' and 'source-type' flags must be given")

	cmd.Flags().Var(flags.UUIDFlag(), encryptionKekKeyIdFlag, "UUID of the STACKIT KMS key used for encryption")
	cmd.Flags().Var(flags.UUIDFlag(), encryptionKekKeyringIdFlag, "UUID of the STACKIT KMS keyring containing the key")
	cmd.Flags().Int64(encryptionKekKeyVersionFlag, 0, "Version of the STACKIT KMS key used for encryption")
	cmd.Flags().String(encryptionServiceAccountFlag, "", "Service account linked to the STACKIT KMS key")

	keyPayload := flags.SecretFlag(encryptionKeyPayloadFlag, params)
	cmd.Flags().Var(keyPayload, encryptionKeyPayloadFlag, keyPayload.Usage()+" Must be base64 encoded. If omitted, the API generates a random secret.")

	cmd.MarkFlagsRequiredTogether(encryptionKekKeyIdFlag, encryptionKekKeyringIdFlag, encryptionKekKeyVersionFlag, encryptionServiceAccountFlag)
	err := flags.MarkFlagsRequired(cmd, availabilityZoneFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel:  globalFlags,
		AvailabilityZone: flags.FlagToStringValue(p, cmd, availabilityZoneFlag),
		Name:             flags.FlagToStringPointer(p, cmd, nameFlag),
		Description:      flags.FlagToStringPointer(p, cmd, descriptionFlag),
		Labels:           flags.FlagToStringToAny(p, cmd, labelFlag),
		PerformanceClass: flags.FlagToStringPointer(p, cmd, performanceClassFlag),
		Size:             flags.FlagToInt64Pointer(p, cmd, sizeFlag),
		SourceId:         flags.FlagToStringPointer(p, cmd, sourceIdFlag),
		SourceType:       flags.FlagToStringPointer(p, cmd, sourceTypeFlag),

		EncryptionKekKeyId:       flags.FlagToStringPointer(p, cmd, encryptionKekKeyIdFlag),
		EncryptionKekKeyringId:   flags.FlagToStringPointer(p, cmd, encryptionKekKeyringIdFlag),
		EncryptionKekKeyVersion:  flags.FlagToInt64Pointer(p, cmd, encryptionKekKeyVersionFlag),
		EncryptionServiceAccount: flags.FlagToStringPointer(p, cmd, encryptionServiceAccountFlag),
	}

	// Cobra requires the encryption parameters together, so checking one is sufficient.
	if cmd.Flags().Changed(encryptionKeyPayloadFlag) && model.EncryptionKekKeyId == nil {
		return nil, &cliErr.FlagValidationError{
			Flag:    encryptionKekKeyIdFlag,
			Details: fmt.Sprintf("when using flag %q, flag %q is required", encryptionKeyPayloadFlag, encryptionKekKeyIdFlag),
		}
	}

	keyPayload, err := flags.ResolveSecretFlag(p, cmd, encryptionKeyPayloadFlag)
	if err != nil {
		return nil, err
	}
	model.EncryptionKeyPayload = keyPayload

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient iaas.DefaultAPI) iaas.ApiCreateVolumeRequest {
	req := apiClient.CreateVolume(ctx, model.ProjectId, model.Region)

	payload := iaas.CreateVolumePayload{
		AvailabilityZone: model.AvailabilityZone,
		Name:             model.Name,
		Description:      model.Description,
		Labels:           model.Labels,
		PerformanceClass: model.PerformanceClass,
		Size:             model.Size,
	}

	if model.SourceId != nil && model.SourceType != nil {
		payload.Source = &iaas.VolumeSource{
			Id:   *model.SourceId,
			Type: *model.SourceType,
		}
	}

	if model.EncryptionKekKeyId != nil && model.EncryptionKekKeyringId != nil && model.EncryptionKekKeyVersion != nil && model.EncryptionServiceAccount != nil {
		payload.EncryptionParameters = &iaas.VolumeEncryptionParameter{
			KekKeyId:       *model.EncryptionKekKeyId,
			KekKeyringId:   *model.EncryptionKekKeyringId,
			KekKeyVersion:  *model.EncryptionKekKeyVersion,
			ServiceAccount: *model.EncryptionServiceAccount,
			KeyPayload:     model.EncryptionKeyPayload,
		}
	}

	return req.CreateVolumePayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, volume *iaas.Volume) error {
	if volume == nil {
		return fmt.Errorf("volume response is empty")
	}
	return p.OutputResult(model.OutputFormat, volume, func() error {
		operationState := "Created"
		if model.Async {
			operationState = "Triggered creation of"
		}
		p.Outputf("%s volume for project %q.\nVolume ID: %s\n", operationState, projectLabel, utils.PtrString(volume.Id))
		return nil
	})
}
