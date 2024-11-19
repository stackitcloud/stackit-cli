package create

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/goccy/go-yaml"
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

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a server",
		Long:  "Creates a server.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a server with machine type "t1.1", name "server1" and image with id xxx`,
				`$ stackit beta server create --machine-type t1.1 --name server1 --image-id xxx`,
			),
			examples.NewExample(
				`Create a server with machine type "t1.1", name "server1", image with id xxx and labels`,
				`$ stackit beta server create --machine-type t1.1 --name server1 --image-id xxx --labels key=value,foo=bar`,
			),
			examples.NewExample(
				`Create a server with machine type "t1.1", name "server1", boot volume source id "xxx", type "image" and size 64GB`,
				`$ stackit beta server create --machine-type t1.1 --name server1 --boot-volume-source-id xxx --boot-volume-source-type image --boot-volume-size 64`,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(p, cmd)
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a server for project %q?", projectLabel)
				err = p.PromptForConfirmation(prompt)
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
				s := spinner.New(p)
				s.Start("Creating server")
				_, err = wait.CreateServerWaitHandler(ctx, apiClient, model.ProjectId, serverId).WaitWithContext(ctx)
				if err != nil {
					return fmt.Errorf("wait for server creation: %w", err)
				}
				s.Stop()
			}

			return outputResult(p, model, projectLabel, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(nameFlag, "n", "", "Server name")
	cmd.Flags().String(machineTypeFlag, "", "Machine type the server shall belong to")
	cmd.Flags().String(affinityGroupFlag, "", "The affinity group the server is assigned to")
	cmd.Flags().String(availabilityZoneFlag, "", "Availability zone")
	cmd.Flags().String(bootVolumeSourceIdFlag, "", "ID of the source object of boot volume. It can be either 'image-id' or 'volume-id'")
	cmd.Flags().String(bootVolumeSourceTypeFlag, "", "Type of the source object of boot volume. It can be either  'image' or 'volume'")
	cmd.Flags().Int64(bootVolumeSizeFlag, 0, "Boot volume size (GB). Size is required for the image type boot volumes")
	cmd.Flags().String(bootVolumePerformanceClassFlag, "", "Boot volume performance class")
	cmd.Flags().Bool(bootVolumeDeleteOnTerminationFlag, false, "Delete the volume during the termination of the server. Defaults to false")
	cmd.Flags().String(imageIdFlag, "", "ID of the image. Either image-id or boot volume is required")
	cmd.Flags().String(keypairNameFlag, "", "The SSH keypair used during the server creation")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a server. E.g. '--labels key1=value1,key2=value2,...'")
	cmd.Flags().String(networkIdFlag, "", "ID of the network for the initial networking setup for the server creation")
	cmd.Flags().StringSlice(networkInterfaceIdsFlag, []string{}, "List of network interface IDs for the initial networking setup for the server creation")
	cmd.Flags().StringSlice(securityGroupsFlag, []string{}, "The initial security groups for the server creation")
	cmd.Flags().StringSlice(serviceAccountEmailsFlag, []string{}, "List of the service account mails")
	cmd.Flags().String(userDataFlag, "", "User data that is provided to the server")
	cmd.Flags().StringSlice(volumesFlag, []string{}, "The list of volumes attached to the server")

	err := flags.MarkFlagsRequired(cmd, nameFlag, machineTypeFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	if globalFlags.ProjectId == "" {
		return nil, &cliErr.ProjectIdError{}
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *iaas.APIClient) iaas.ApiCreateServerRequest {
	req := apiClient.CreateServer(ctx, model.ProjectId)
	var labelsMap *map[string]interface{}
	if model.Labels != nil && len(*model.Labels) > 0 {
		// convert map[string]string to map[string]interface{}
		labelsMap = utils.Ptr(map[string]interface{}{})
		for k, v := range *model.Labels {
			(*labelsMap)[k] = v
		}
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
		UserData:            model.UserData,
		Volumes:             model.Volumes,
		Labels:              labelsMap,
	}

	if model.BootVolumePerformanceClass != nil || model.BootVolumeSize != nil || model.BootVolumeDeleteOnTermination != nil || model.BootVolumeSourceId != nil || model.BootVolumeSourceType != nil {
		payload.BootVolume = &iaas.CreateServerPayloadBootVolume{
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
		payload.Networking = &iaas.CreateServerPayloadNetworking{}

		if model.NetworkInterfaceIds != nil {
			payload.Networking.CreateServerNetworkingWithNics = &iaas.CreateServerNetworkingWithNics{
				NicIds: model.NetworkInterfaceIds,
			}
		} else if model.NetworkId != nil {
			payload.Networking.CreateServerNetworking = &iaas.CreateServerNetworking{
				NetworkId: model.NetworkId,
			}
		}
	}

	return req.CreateServerPayload(payload)
}

func outputResult(p *print.Printer, model *inputModel, projectLabel string, server *iaas.Server) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(server, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal server: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(server, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal server: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created server for project %q.\nServer ID: %s\n", projectLabel, *server.Id)
		return nil
	}
}
