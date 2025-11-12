package create

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
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
	"github.com/stackitcloud/stackit-sdk-go/services/iaas"
	"github.com/stackitcloud/stackit-sdk-go/services/iaas/wait"

	"github.com/spf13/cobra"
)

const (
	nameFlag                          = "name"
	machineTypeFlag                   = "machine-type"
	affinityGroupFlag                 = "affinity-group"
	availabilityZoneFlag              = "availability-zone"
	bootVolumeSourceIdFlag            = "boot-volume-source-id"
	bootVolumeSourceTypeFlag          = "boot-volume-source-type"
	bootVolumeSizeFlag                = "boot-volume-size"
	bootVolumePerformanceClassFlag    = "boot-volume-performance-class"
	bootVolumeDeleteOnTerminationFlag = "boot-volume-delete-on-termination"
	imageIdFlag                       = "image-id"
	keypairNameFlag                   = "keypair-name"
	labelFlag                         = "labels"
	networkIdFlag                     = "network-id"
	networkInterfaceIdsFlag           = "network-interface-ids"
	securityGroupsFlag                = "security-groups"
	serviceAccountEmailsFlag          = "service-account-emails"
	userDataFlag                      = "user-data"
	volumesFlag                       = "volumes"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	Name                          *string
	MachineType                   *string
	AffinityGroup                 *string
	AvailabilityZone              *string
	BootVolumeSourceId            *string
	BootVolumeSourceType          *string
	BootVolumeSize                *int64
	BootVolumePerformanceClass    *string
	BootVolumeDeleteOnTermination *bool
	ImageId                       *string
	KeypairName                   *string
	Labels                        *map[string]string
	NetworkId                     *string
	NetworkInterfaceIds           *[]string
	SecurityGroups                *[]string
	ServiceAccountMails           *[]string
	UserData                      *string
	Volumes                       *[]string
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a server",
		Long:  "Creates a server.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a server from an image with id xxx`,
				`$ stackit server create --machine-type t1.1 --name server1 --image-id xxx`,
			),
			examples.NewExample(
				`Create a server with labels from an image with id xxx`,
				`$ stackit server create --machine-type t1.1 --name server1 --image-id xxx --labels key=value,foo=bar`,
			),
			examples.NewExample(
				`Create a server with a boot volume`,
				`$ stackit server create --machine-type t1.1 --name server1 --boot-volume-source-id xxx --boot-volume-source-type image --boot-volume-size 64`,
			),
			examples.NewExample(
				`Create a server with a boot volume from an existing volume`,
				`$ stackit server create --machine-type t1.1 --name server1 --boot-volume-source-id xxx --boot-volume-source-type volume`,
			),
			examples.NewExample(
				`Create a server with a keypair`,
				`$ stackit server create --machine-type t1.1 --name server1 --image-id xxx --keypair-name example`,
			),
			examples.NewExample(
				`Create a server with a network`,
				`$ stackit server create --machine-type t1.1 --name server1 --image-id xxx --network-id yyy`,
			),
			examples.NewExample(
				`Create a server with a network interface`,
				`$ stackit server create --machine-type t1.1 --name server1 --boot-volume-source-id xxx --boot-volume-source-type image --boot-volume-size 64 --network-interface-ids yyy`,
			),
			examples.NewExample(
				`Create a server with an attached volume`,
				`$ stackit server create --machine-type t1.1 --name server1 --boot-volume-source-id xxx --boot-volume-source-type image --boot-volume-size 64 --volumes yyy`,
			),
			examples.NewExample(
				`Create a server with user data (cloud-init)`,
				`$ stackit server create --machine-type t1.1 --name server1 --boot-volume-source-id xxx --boot-volume-source-type image --boot-volume-size 64 --user-data @path/to/file.yaml")`,
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a server for project %q?", projectLabel)
				err = params.Printer.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create server : %w", err)
			}
			serverId := *resp.Id

			// Wait for async operation, if async mode not enabled
			if !model.Async {
				s := spinner.New(params.Printer)
				s.Start("Creating server")
				_, err = wait.CreateServerWaitHandler(ctx, apiClient, model.ProjectId, model.Region, serverId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for server creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(params.Printer, model.OutputFormat, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "Server name")
	cmd.Flags().String(machineTypeFlag, "", "Name of the type of the machine for the server. Possible values are documented in https://docs.stackit.cloud/stackit/en/virtual-machine-flavors-75137231.html")
	cmd.Flags().String(affinityGroupFlag, "", "The affinity group the server is assigned to")
	cmd.Flags().String(availabilityZoneFlag, "", "The availability zone of the server")
	cmd.Flags().String(bootVolumeSourceIdFlag, "", "ID of the source object of boot volume. It can be either an image or volume ID")
	cmd.Flags().String(bootVolumeSourceTypeFlag, "", "Type of the source object of boot volume. It can be either  'image' or 'volume'")
	cmd.Flags().Int64(bootVolumeSizeFlag, 0, "The size of the boot volume in GB. Must be provided when 'boot-volume-source-type' is 'image'")
	cmd.Flags().String(bootVolumePerformanceClassFlag, "", "Boot volume performance class")
	cmd.Flags().Bool(bootVolumeDeleteOnTerminationFlag, false, "Delete the volume during the termination of the server. Defaults to false")
	cmd.Flags().String(imageIdFlag, "", "The image ID to be used for an ephemeral disk on the server. Either 'image-id' or 'boot-volume-...' flags are required")
	cmd.Flags().String(keypairNameFlag, "", "The name of the SSH keypair used during the server creation")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a server. E.g. '--labels key1=value1,key2=value2,...'")
	cmd.Flags().String(networkIdFlag, "", "ID of the network for the initial networking setup for the server creation")
	cmd.Flags().StringSlice(networkInterfaceIdsFlag, []string{}, "List of network interface IDs for the initial networking setup for the server creation")
	cmd.Flags().StringSlice(securityGroupsFlag, []string{}, "The initial security groups for the server creation")
	cmd.Flags().StringSlice(serviceAccountEmailsFlag, []string{}, "List of the service account mails")
	cmd.Flags().Var(flags.ReadFromFileFlag(), userDataFlag, "User data that is passed via cloud-init to the server")
	cmd.Flags().StringSlice(volumesFlag, []string{}, "The list of volumes attached to the server")

	err := flags.MarkFlagsRequired(cmd, nameFlag, machineTypeFlag)
	cmd.MarkFlagsMutuallyExclusive(imageIdFlag, bootVolumeSourceIdFlag)
	cmd.MarkFlagsMutuallyExclusive(imageIdFlag, bootVolumeSourceTypeFlag)
	cmd.MarkFlagsMutuallyExclusive(networkIdFlag, networkInterfaceIdsFlag)
	cmd.MarkFlagsOneRequired(networkIdFlag, networkInterfaceIdsFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command, _ []string) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
	}

	bootVolumeSourceId := flags.FlagToStringPointer(p, cmd, bootVolumeSourceIdFlag)
	bootVolumeSourceType := flags.FlagToStringPointer(p, cmd, bootVolumeSourceTypeFlag)
	bootVolumeSize := flags.FlagToInt64Pointer(p, cmd, bootVolumeSizeFlag)
	imageId := flags.FlagToStringPointer(p, cmd, imageIdFlag)

	if imageId == nil && bootVolumeSourceId == nil && bootVolumeSourceType == nil {
		return nil, &cliErr.ServerCreateMissingFlagsError{
			Cmd: cmd,
		}
	}

	if imageId == nil {
		err := flags.MarkFlagsRequired(cmd, bootVolumeSourceIdFlag, bootVolumeSourceTypeFlag)
		cobra.CheckErr(err)
	}

	if bootVolumeSourceId != nil && bootVolumeSourceType == nil {
		err := cmd.MarkFlagRequired(bootVolumeSourceTypeFlag)
		cobra.CheckErr(err)

		return nil, &cliErr.ServerCreateMissingVolumeTypeError{
			Cmd: cmd,
		}
	}

	if bootVolumeSourceType != nil {
		if bootVolumeSourceId == nil {
			err := cmd.MarkFlagRequired(bootVolumeSourceIdFlag)
			cobra.CheckErr(err)

			return nil, &cliErr.ServerCreateMissingVolumeIdError{
				Cmd: cmd,
			}
		}

		if *bootVolumeSourceType == "image" && bootVolumeSize == nil {
			err := cmd.MarkFlagRequired(bootVolumeSizeFlag)
			cobra.CheckErr(err)
			return nil, &cliErr.ServerCreateError{
				Cmd: cmd,
			}
		}
	}

	if bootVolumeSourceId == nil && bootVolumeSourceType == nil {
		err := cmd.MarkFlagRequired(imageIdFlag)
		cobra.CheckErr(err)
	}

	model := inputModel{
		GlobalFlagModel:               globalFlags,
		Name:                          flags.FlagToStringPointer(p, cmd, nameFlag),
		MachineType:                   flags.FlagToStringPointer(p, cmd, machineTypeFlag),
		AffinityGroup:                 flags.FlagToStringPointer(p, cmd, affinityGroupFlag),
		AvailabilityZone:              flags.FlagToStringPointer(p, cmd, availabilityZoneFlag),
		BootVolumeSourceId:            flags.FlagToStringPointer(p, cmd, bootVolumeSourceIdFlag),
		BootVolumeSourceType:          flags.FlagToStringPointer(p, cmd, bootVolumeSourceTypeFlag),
		BootVolumeSize:                flags.FlagToInt64Pointer(p, cmd, bootVolumeSizeFlag),
		BootVolumePerformanceClass:    flags.FlagToStringPointer(p, cmd, bootVolumePerformanceClassFlag),
		BootVolumeDeleteOnTermination: flags.FlagToBoolPointer(p, cmd, bootVolumeDeleteOnTerminationFlag),
		ImageId:                       flags.FlagToStringPointer(p, cmd, imageIdFlag),
		KeypairName:                   flags.FlagToStringPointer(p, cmd, keypairNameFlag),
		Labels:                        flags.FlagToStringToStringPointer(p, cmd, labelFlag),
		NetworkId:                     flags.FlagToStringPointer(p, cmd, networkIdFlag),
		NetworkInterfaceIds:           flags.FlagToStringSlicePointer(p, cmd, networkInterfaceIdsFlag),
		SecurityGroups:                flags.FlagToStringSlicePointer(p, cmd, securityGroupsFlag),
		ServiceAccountMails:           flags.FlagToStringSlicePointer(p, cmd, serviceAccountEmailsFlag),
		UserData:                      flags.FlagToStringPointer(p, cmd, userDataFlag),
		Volumes:                       flags.FlagToStringSlicePointer(p, cmd, volumesFlag),
	}

	p.DebugInputModel(model)
	return &model, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateServerRequest {
	req := apiClient.CreateServer(ctx, model.ProjectId, model.Region)

	var userData *[]byte
	if model.UserData != nil {
		userData = utils.Ptr([]byte(*model.UserData))
	}

	payload := iaas.CreateServerPayload{
		Name:             model.Name,
		MachineType:      model.MachineType,
		AffinityGroup:    model.AffinityGroup,
		AvailabilityZone: model.AvailabilityZone,

		ImageId:             model.ImageId,
		KeypairName:         model.KeypairName,
		SecurityGroups:      model.SecurityGroups,
		ServiceAccountMails: model.ServiceAccountMails,
		UserData:            userData,
		Volumes:             model.Volumes,
		Labels:              utils.ConvertStringMapToInterfaceMap(model.Labels),
	}

	if model.BootVolumePerformanceClass != nil || model.BootVolumeSize != nil || model.BootVolumeDeleteOnTermination != nil || model.BootVolumeSourceId != nil || model.BootVolumeSourceType != nil {
		payload.BootVolume = &iaas.ServerBootVolume{
			PerformanceClass:    model.BootVolumePerformanceClass,
			Size:                model.BootVolumeSize,
			DeleteOnTermination: model.BootVolumeDeleteOnTermination,
			Source: &iaas.BootVolumeSource{
				Id:   model.BootVolumeSourceId,
				Type: model.BootVolumeSourceType,
			},
		}
	}

	if model.NetworkInterfaceIds != nil || model.NetworkId != nil {
		payload.Networking = &iaas.CreateServerPayloadAllOfNetworking{}

		if model.NetworkInterfaceIds != nil {
			payload.Networking.CreateServerNetworkingWithNics = &iaas.CreateServerNetworkingWithNics{
				NicIds: model.NetworkInterfaceIds,
			}
		}
		if model.NetworkId != nil {
			payload.Networking.CreateServerNetworking = &iaas.CreateServerNetworking{
				NetworkId: model.NetworkId,
			}
		}
	}

	return req.CreateServerPayload(payload)
}

func outputResult(p *print.Printer, outputFormat, projectLabel string, server *iaas.Server) error {
	if server == nil {
		return fmt.Errorf("server response is empty")
	}
	return p.OutputResult(outputFormat, server, func() error {
		p.Outputf("Created server for project %q.\nServer ID: %s\n", projectLabel, utils.PtrString(server.Id))
		return nil
	})
}
