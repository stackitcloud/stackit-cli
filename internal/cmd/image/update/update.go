package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	cliErr "github.com/stackitcloud/stackit-cli/internal/pkg/errors"
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
	Architecture           *string
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

	architectureFlag           = "architecture"
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
	protectedFlag   = "protected"
)

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("update %s", imageIdArg),
		Short: "Updates an image",
		Long:  "Updates an image",
		Args:  args.SingleArg(imageIdArg, utils.ValidateUUID),
		Example: examples.Build(
			examples.NewExample(`Update the name of an image with ID "xxx"`, `$ stackit image update xxx --name my-new-name`),
			examples.NewExample(`Update the labels of an image with ID "xxx"`, `$ stackit image update xxx --labels label1=value1,label2=value2`),
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

			imageLabel, err := iaasUtils.GetImageName(ctx, apiClient, model.ProjectId, model.Id)
			if err != nil {
				params.Printer.Debug(print.WarningLevel, "cannot retrieve image name: %v", err)
				imageLabel = model.Id
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to update the image %q?", imageLabel)
				err = params.Printer.PromptForConfirmation(prompt)
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
			params.Printer.Info("Updated image \"%v\" for %q\n", utils.PtrString(resp.Name), projectLabel)

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "The name of the image.")
	cmd.Flags().String(diskFormatFlag, "", "The disk format of the image. ")

	cmd.Flags().String(architectureFlag, "", "Sets the CPU architecture.")
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
		return nil, &cliErr.ProjectIdError{}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		Id:              cliArgs[0],
		Name:            flags.FlagToStringPointer(p, cmd, nameFlag),

		DiskFormat: flags.FlagToStringPointer(p, cmd, diskFormatFlag),
		Labels:     flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
		Config: &imageConfig{
			Architecture:           flags.FlagToStringPointer(p, cmd, architectureFlag),
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

	if model.Config.isEmpty() {
		model.Config = nil
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiUpdateImageRequest {
	request := apiClient.UpdateImage(ctx, model.ProjectId, model.Id)
	payload := iaas.NewUpdateImagePayload()

	// Config *ImageConfig `json:"config,omitempty"`
	payload.DiskFormat = model.DiskFormat
	payload.Labels = utils.ConvertStringMapToInterfaceMap(model.Labels)
	payload.MinDiskSize = model.MinDiskSize
	payload.MinRam = model.MinRam
	payload.Name = model.Name
	payload.Protected = model.Protected
	payload.Config = nil

	if config := model.Config; config != nil {
		payload.Config = &iaas.ImageConfig{}
		if model.Config.BootMenu != nil {
			payload.Config.BootMenu = model.Config.BootMenu
		}
		if model.Config.CdromBus != nil {
			payload.Config.CdromBus = iaas.NewNullableString(model.Config.CdromBus)
		}
		if model.Config.DiskBus != nil {
			payload.Config.DiskBus = iaas.NewNullableString(model.Config.DiskBus)
		}
		if model.Config.NicModel != nil {
			payload.Config.NicModel = iaas.NewNullableString(model.Config.NicModel)
		}
		if model.Config.OperatingSystem != nil {
			payload.Config.OperatingSystem = model.Config.OperatingSystem
		}
		if model.Config.OperatingSystemDistro != nil {
			payload.Config.OperatingSystemDistro = iaas.NewNullableString(model.Config.OperatingSystemDistro)
		}
		if model.Config.OperatingSystemVersion != nil {
			payload.Config.OperatingSystemVersion = iaas.NewNullableString(model.Config.OperatingSystemVersion)
		}
		if model.Config.RescueBus != nil {
			payload.Config.RescueBus = iaas.NewNullableString(model.Config.RescueBus)
		}
		if model.Config.RescueDevice != nil {
			payload.Config.RescueDevice = iaas.NewNullableString(model.Config.RescueDevice)
		}
		if model.Config.SecureBoot != nil {
			payload.Config.SecureBoot = model.Config.SecureBoot
		}
		if model.Config.Uefi != nil {
			payload.Config.Uefi = model.Config.Uefi
		}
		if model.Config.VideoModel != nil {
			payload.Config.VideoModel = iaas.NewNullableString(model.Config.VideoModel)
		}
		if model.Config.VirtioScsi != nil {
			payload.Config.VirtioScsi = model.Config.VirtioScsi
		}
	}

	request = request.UpdateImagePayload(*payload)

	return request
}
