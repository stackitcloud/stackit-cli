package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/projectname"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	iaasUtils "github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/utils"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

type imageConfig struct {
	BootMenu               *bool
	CdromBus               *string
	DiskBus                *string
	NicModel               *string
	OperatingSystem        *string
	OperatingSystemDistro  *string
	OperatingSystemVersion *string
	RescueBus              *string
	RescueDevice           *string
	SecureBoot             *bool
	Uefi                   *bool
	VideoModel             *string
	VirtioScsi             *bool
}

func (ic *imageConfig) isEmpty() bool {
	return ic.BootMenu == nil &&
		ic.CdromBus == nil &&
		ic.DiskBus == nil &&
		ic.NicModel == nil &&
		ic.OperatingSystem == nil &&
		ic.OperatingSystemDistro == nil &&
		ic.OperatingSystemVersion == nil &&
		ic.RescueBus == nil &&
		ic.RescueDevice == nil &&
		ic.SecureBoot == nil &&
		ic.Uefi == nil &&
		ic.VideoModel == nil &&
		ic.VirtioScsi == nil
}

type inputModel struct {
	*globalflags.GlobalFlagModel

	Id          string
	Name        *string
	DiskFormat  *string
	Labels      *map[string]string
	Config      *imageConfig
	MinDiskSize *int64
	MinRam      *int64
	Protected   *bool
}

func (im *inputModel) isEmpty() bool {
	return im.Name == nil &&
		im.DiskFormat == nil &&
		im.Labels == nil &&
		(im.Config == nil || im.Config.isEmpty()) &&
		im.MinDiskSize == nil &&
		im.MinRam == nil &&
		im.Protected == nil
}

const imageIdArg = "IMAGE_ID"

const (
	nameFlag       = "name"
	diskFormatFlag = "disk-format"

	bootMenuFlag               = "boot-menu"
	cdromBusFlag               = "cdrom-bus"
	diskBusFlag                = "disk-bus"
	nicModelFlag               = "nic-model"
	operatingSystemFlag        = "os"
	operatingSystemDistroFlag  = "os-distro"
	operatingSystemVersionFlag = "os-version"
	rescueBusFlag              = "rescue-bus"
	rescueDeviceFlag           = "rescue-device"
	secureBootFlag             = "secure-boot"
	uefiFlag                   = "uefi"
	videoModelFlag             = "video-model"
	virtioScsiFlag             = "virtio-scsi"

	labelsFlag = "labels"

	minDiskSizeFlag = "min-disk-size"
	minRamFlag      = "min-ram"
	ownerFlag       = "owner"
	protectedFlag   = "protected"
)

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", imageIdArg),
		Short: "Updates an image",
		Long:  "Updates an image",
		Args:  args.SingleArg(imageIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(`Update the name of an image with ID "xxx"`, `$ stackit beta image update xxx --name my-new-name`),
			examples.NewExample(`Update the labels of an image with ID "xxx"`, `$ stackit beta image update xxx --labels label1=value1,label2=value2`),
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

			projectLabel, err := projectname.GetProjectName(ctx, p, cmd)
			if err != nil {
				p.Debug(print.ErrorLevel, "get project name: %v", err)
				projectLabel = model.ProjectId
			}

			imageLabel, err := iaasUtils.GetImageName(ctx, apiClient, model.ProjectId, model.Id)
			if err != nil {
				p.Debug(print.WarningLevel, "cannot retrieve image name: %v", err)
				imageLabel = model.Id
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update the image %q?", imageLabel)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)

			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("update image: %w", err)
			}
			p.Info("Updated image \"%v\" for %q\n", utils.PtrString(resp.Name), projectLabel)

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "The name of the image.")
	cmd.Flags().String(diskFormatFlag, "", "The disk format of the image. ")

	cmd.Flags().Bool(bootMenuFlag, false, "Enables the BIOS bootmenu.")
	cmd.Flags().String(cdromBusFlag, "", "Sets CDROM bus controller type.")
	cmd.Flags().String(diskBusFlag, "", "Sets Disk bus controller type.")
	cmd.Flags().String(nicModelFlag, "", "Sets virtual nic model.")
	cmd.Flags().String(operatingSystemFlag, "", "Enables OS specific optimizations.")
	cmd.Flags().String(operatingSystemDistroFlag, "", "Operating System Distribution.")
	cmd.Flags().String(operatingSystemVersionFlag, "", "Version of the OS.")
	cmd.Flags().String(rescueBusFlag, "", "Sets the device bus when the image is used as a rescue image.")
	cmd.Flags().String(rescueDeviceFlag, "", "Sets the device when the image is used as a rescue image.")
	cmd.Flags().Bool(secureBootFlag, false, "Enables Secure Boot.")
	cmd.Flags().Bool(uefiFlag, false, "Enables UEFI boot.")
	cmd.Flags().String(videoModelFlag, "", "Sets Graphic device model.")
	cmd.Flags().Bool(virtioScsiFlag, false, "Enables the use of VirtIO SCSI to provide block device access. By default instances use VirtIO Block.")

	cmd.Flags().StringToString(labelsFlag, nil, "Labels are key-value string pairs which can be attached to a network-interface. E.g. '--labels key1=value1,key2=value2,...'")

	cmd.Flags().Int64(minDiskSizeFlag, 0, "Size in Gigabyte.")
	cmd.Flags().Int64(minRamFlag, 0, "Size in Megabyte.")
	cmd.Flags().Bool(protectedFlag, false, "Protected VM.")

	cmd.MarkFlagsRequiredTogether(rescueBusFlag, rescueDeviceFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command, cliArgs []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Id:              cliArgs[0],
		Name:            flags.FlagToStringPointer(p, cmd, nameFlag),

		DiskFormat: flags.FlagToStringPointer(p, cmd, diskFormatFlag),
		Labels:     flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
		Config: &imageConfig{
			BootMenu:               flags.FlagToBoolPointer(p, cmd, bootMenuFlag),
			CdromBus:               flags.FlagToStringPointer(p, cmd, cdromBusFlag),
			DiskBus:                flags.FlagToStringPointer(p, cmd, diskBusFlag),
			NicModel:               flags.FlagToStringPointer(p, cmd, nicModelFlag),
			OperatingSystem:        flags.FlagToStringPointer(p, cmd, operatingSystemFlag),
			OperatingSystemDistro:  flags.FlagToStringPointer(p, cmd, operatingSystemDistroFlag),
			OperatingSystemVersion: flags.FlagToStringPointer(p, cmd, operatingSystemVersionFlag),
			RescueBus:              flags.FlagToStringPointer(p, cmd, rescueBusFlag),
			RescueDevice:           flags.FlagToStringPointer(p, cmd, rescueDeviceFlag),
			SecureBoot:             flags.FlagToBoolPointer(p, cmd, secureBootFlag),
			Uefi:                   flags.FlagToBoolPointer(p, cmd, uefiFlag),
			VideoModel:             flags.FlagToStringPointer(p, cmd, videoModelFlag),
			VirtioScsi:             flags.FlagToBoolPointer(p, cmd, virtioScsiFlag),
		},
		MinDiskSize: flags.FlagToInt64Pointer(p, cmd, minDiskSizeFlag),
		MinRam:      flags.FlagToInt64Pointer(p, cmd, minRamFlag),
		Protected:   flags.FlagToBoolPointer(p, cmd, protectedFlag),
	}

	if model.isEmpty() {
		return nil, fmt.Errorf("no flags have been passed")
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiUpdateImageRequest {
	request := apiClient.UpdateImage(ctx, model.ProjectId, model.Id)
	payload := iaas.NewUpdateImagePayload()
	var labelsMap *map[string]any
	if model.Labels != nil && len(*model.Labels) > 0 {
		// convert map[string]string to map[string]interface{}
		labelsMap = utils.Ptr(map[string]interface{}{})
		for k, v := range *model.Labels {
			(*labelsMap)[k] = v
		}
	}
	// Config *ImageConfig `json:"config,omitempty"`
	payload.DiskFormat = model.DiskFormat
	payload.Labels = labelsMap
	payload.MinDiskSize = model.MinDiskSize
	payload.MinRam = model.MinRam
	payload.Name = model.Name
	payload.Protected = model.Protected

	if model.Config != nil {
		payload.Config = &iaas.ImageConfig{
			BootMenu:               model.Config.BootMenu,
			CdromBus:               iaas.NewNullableString(model.Config.CdromBus),
			DiskBus:                iaas.NewNullableString(model.Config.DiskBus),
			NicModel:               iaas.NewNullableString(model.Config.NicModel),
			OperatingSystem:        model.Config.OperatingSystem,
			OperatingSystemDistro:  iaas.NewNullableString(model.Config.OperatingSystemDistro),
			OperatingSystemVersion: iaas.NewNullableString(model.Config.OperatingSystemVersion),
			RescueBus:              iaas.NewNullableString(model.Config.RescueBus),
			RescueDevice:           iaas.NewNullableString(model.Config.RescueDevice),
			SecureBoot:             model.Config.SecureBoot,
			Uefi:                   model.Config.Uefi,
			VideoModel:             iaas.NewNullableString(model.Config.VideoModel),
			VirtioScsi:             model.Config.VirtioScsi,
		}
	}

	request = request.UpdateImagePayload(*payload)

	return request
}
