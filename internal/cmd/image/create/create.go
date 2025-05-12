package create

import (
	"bufio"
	"context"
	"encoding/json"
	goerrors "errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/iaas/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
)

const (
	nameFlag                = "name"
	diskFormatFlag          = "disk-format"
	localFilePathFlag       = "local-file-path"
	noProgressIndicatorFlag = "no-progress"

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
	Uefi                   bool
	VideoModel             *string
	VirtioScsi             *bool
}
type inputModel struct {
	*globalflags.GlobalFlagModel

	Id                  *string
	Name                string
	DiskFormat          string
	LocalFilePath       string
	Labels              *map[string]string
	Config              *imageConfig
	MinDiskSize         *int64
	MinRam              *int64
	Protected           *bool
	NoProgressIndicator *bool
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates images",
		Long:  "Creates images.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create an image with name 'my-new-image' from a raw disk image located in '/my/raw/image'`,
				`$ stackit image create --name my-new-image --disk-format=raw --local-file-path=/my/raw/image`,
			),
			examples.NewExample(
				`Create an image with name 'my-new-image' from a qcow2 image read from '/my/qcow2/image' with labels describing its contents`,
				`$ stackit image create --name my-new-image --disk-format=qcow2 --local-file-path=/my/qcow2/image --labels os=linux,distro=alpine,version=3.12`,
			),
		),
		RunE: func(cmd *cobra.Command, _ []string) (err error) {
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

			// we open input file first to fail fast, if it is not readable
			file, err := os.Open(model.LocalFilePath)
			if err != nil {
				return fmt.Errorf("create image: file %q is not readable: %w", model.LocalFilePath, err)
			}
			defer func() {
				if inner := file.Close(); inner != nil {
					err = fmt.Errorf("error closing input file: %w (%w)", inner, err)
				}
			}()

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create the image %q?", model.Name)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			request := buildRequest(ctx, model, apiClient)

			result, err := request.Execute()
			if err != nil {
				return fmt.Errorf("create image: %w", err)
			}
			model.Id = result.Id
			url, ok := result.GetUploadUrlOk()
			if !ok {
				return fmt.Errorf("create image: no upload URL has been provided")
			}
			if err := uploadAsync(ctx, params.Printer, model, file, url); err != nil {
				return err
			}

			if err := outputResult(params.Printer, model, result); err != nil {
				return err
			}

			return nil
		},
	}

	configureFlags(cmd)
	return cmd
}

func uploadAsync(ctx context.Context, p *print.Printer, model *inputModel, file *os.File, url string) error {
	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("upload file: %w", err)
	}

	var reader io.Reader
	if model.NoProgressIndicator != nil && *model.NoProgressIndicator {
		reader = file
	} else {
		var ch <-chan int
		reader, ch = newProgressReader(file)
		go func() {
			ticker := time.NewTicker(2 * time.Second)
			var uploaded int
		done:
			for {
				select {
				case <-ticker.C:
					p.Info("uploaded %3.1f%%\r", 100.0/float64(stat.Size())*float64(uploaded))
				case n, ok := <-ch:
					if !ok {
						break done
					}
					if n >= 0 {
						uploaded += n
					}
				}
			}
		}()
	}

	if err = uploadFile(ctx, p, reader, stat.Size(), url); err != nil {
		return fmt.Errorf("upload file: %w", err)
	}

	return nil
}

var _ io.Reader = (*progressReader)(nil)

type progressReader struct {
	delegate io.Reader
	ch       chan int
}

func newProgressReader(delegate io.Reader) (reader io.Reader, result <-chan int) {
	ch := make(chan int)
	return &progressReader{
		delegate: delegate,
		ch:       ch,
	}, ch
}

// Read implements io.Reader.
func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.delegate.Read(p)
	if goerrors.Is(err, io.EOF) && n <= 0 {
		close(pr.ch)
	} else {
		pr.ch <- n
	}
	return n, err
}

func uploadFile(ctx context.Context, p *print.Printer, reader io.Reader, filesize int64, url string) error {
	p.Debug(print.DebugLevel, "uploading image to %s", url)

	start := time.Now()
	// pass the file contents as stream, as they can get arbitrarily large. We do
	// _not_ want to load them into an internal buffer. The downside is, that we
	// have to set the content-length header manually
	uploadRequest, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bufio.NewReader(reader))
	if err != nil {
		return fmt.Errorf("create image: cannot create request: %w", err)
	}
	uploadRequest.Header.Add("Content-Type", "application/octet-stream")
	uploadRequest.ContentLength = filesize

	uploadResponse, err := http.DefaultClient.Do(uploadRequest)
	if err != nil {
		return fmt.Errorf("create image: error contacting server for upload: %w", err)
	}
	defer func() {
		if inner := uploadResponse.Body.Close(); inner != nil {
			err = fmt.Errorf("error closing file: %w (%w)", inner, err)
		}
	}()
	if uploadResponse.StatusCode != http.StatusOK {
		return fmt.Errorf("create image: server rejected image upload with %s", uploadResponse.Status)
	}
	delay := time.Since(start)
	p.Debug(print.DebugLevel, "uploaded %d bytes in %v", filesize, delay)

	return nil
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(nameFlag, "", "The name of the image.")
	cmd.Flags().String(diskFormatFlag, "", "The disk format of the image. ")
	cmd.Flags().String(localFilePathFlag, "", "The path to the local disk image file.")
	cmd.Flags().Bool(noProgressIndicatorFlag, false, "Show no progress indicator for upload.")

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
	cmd.Flags().Bool(uefiFlag, true, "Enables UEFI boot.")
	cmd.Flags().String(videoModelFlag, "", "Sets Graphic device model.")
	cmd.Flags().Bool(virtioScsiFlag, false, "Enables the use of VirtIO SCSI to provide block device access. By default instances use VirtIO Block.")

	cmd.Flags().StringToString(labelsFlag, nil, "Labels are key-value string pairs which can be attached to a network-interface. E.g. '--labels key1=value1,key2=value2,...'")

	cmd.Flags().Int64(minDiskSizeFlag, 0, "Size in Gigabyte.")
	cmd.Flags().Int64(minRamFlag, 0, "Size in Megabyte.")
	cmd.Flags().Bool(protectedFlag, false, "Protected VM.")

	if err := flags.MarkFlagsRequired(cmd, nameFlag, diskFormatFlag, localFilePathFlag); err != nil {
		cobra.CheckErr(err)
	}
	cmd.MarkFlagsRequiredTogether(rescueBusFlag, rescueDeviceFlag)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}
	name := flags.FlagToStringValue(p, cmd, nameFlag)

	model := inputModel{
		GlobalFlagModel:     globalFlags,
		Name:                name,
		DiskFormat:          flags.FlagToStringValue(p, cmd, diskFormatFlag),
		LocalFilePath:       flags.FlagToStringValue(p, cmd, localFilePathFlag),
		Labels:              flags.FlagToStringToStringPointer(p, cmd, labelsFlag),
		NoProgressIndicator: flags.FlagToBoolPointer(p, cmd, noProgressIndicatorFlag),
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
			Uefi:                   flags.FlagToBoolValue(p, cmd, uefiFlag),
			VideoModel:             flags.FlagToStringPointer(p, cmd, videoModelFlag),
			VirtioScsi:             flags.FlagToBoolPointer(p, cmd, virtioScsiFlag),
		},
		MinDiskSize: flags.FlagToInt64Pointer(p, cmd, minDiskSizeFlag),
		MinRam:      flags.FlagToInt64Pointer(p, cmd, minRamFlag),
		Protected:   flags.FlagToBoolPointer(p, cmd, protectedFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateImageRequest {
	request := apiClient.CreateImage(ctx, model.ProjectId).
		CreateImagePayload(createPayload(ctx, model))
	return request
}

func createPayload(_ context.Context, model *inputModel) iaas.CreateImagePayload {
	var labelsMap *map[string]any
	if model.Labels != nil && len(*model.Labels) > 0 {
		// convert map[string]string to map[string]interface{}
		labelsMap = utils.Ptr(map[string]interface{}{})
		for k, v := range *model.Labels {
			(*labelsMap)[k] = v
		}
	}
	payload := iaas.CreateImagePayload{
		DiskFormat:  &model.DiskFormat,
		Name:        &model.Name,
		Labels:      labelsMap,
		MinDiskSize: model.MinDiskSize,
		MinRam:      model.MinRam,
		Protected:   model.Protected,
	}
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
			Uefi:                   utils.Ptr(model.Config.Uefi),
			VideoModel:             iaas.NewNullableString(model.Config.VideoModel),
			VirtioScsi:             model.Config.VirtioScsi,
		}
	}

	return payload
}

func outputResult(p *print.Printer, model *inputModel, resp *iaas.ImageCreateResponse) error {
	if model == nil {
		return fmt.Errorf("input model is nil")
	}
	var outputFormat string
	if model.GlobalFlagModel != nil {
		outputFormat = model.OutputFormat
	}
	switch outputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal image: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal image: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created image %q with id %s\n", model.Name, utils.PtrString(model.Id))
		return nil
	}
}
