package options

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/params"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
	"github.com/stackitcloud/stackit-sdk-go/services/ske"
)

const (
	availabilityZonesFlag  = "availability-zones"
	kubernetesVersionsFlag = "kubernetes-versions"
	machineImagesFlag      = "machine-images"
	machineTypesFlag       = "machine-types"
	volumeTypesFlag        = "volume-types"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	AvailabilityZones  bool
	KubernetesVersions bool
	MachineImages      bool
	MachineTypes       bool
	VolumeTypes        bool
}

func NewCmd(params *params.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "options",
		Short: "Lists SKE provider options",
		Long: fmt.Sprintf("%s\n%s",
			"Lists STACKIT Kubernetes Engine (SKE) provider options (availability zones, Kubernetes versions, machine images and types, volume types).",
			"Pass one or more flags to filter what categories are shown.",
		),
		Args: args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`List SKE options for all categories`,
				"$ stackit ske options"),
			examples.NewExample(
				`List SKE options regarding Kubernetes versions only`,
				"$ stackit ske options --kubernetes-versions"),
			examples.NewExample(
				`List SKE options regarding Kubernetes versions and machine images`,
				"$ stackit ske options --kubernetes-versions --machine-images"),
		),
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := context.Background()
			model, err := parseInput(params.Printer, cmd)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(params.Printer)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get SKE provider options: %w", err)
			}

			return outputResult(params.Printer, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().Bool(availabilityZonesFlag, false, "Lists availability zones")
	cmd.Flags().Bool(kubernetesVersionsFlag, false, "Lists supported kubernetes versions")
	cmd.Flags().Bool(machineImagesFlag, false, "Lists supported machine images")
	cmd.Flags().Bool(machineTypesFlag, false, "Lists supported machine types")
	cmd.Flags().Bool(volumeTypesFlag, false, "Lists supported volume types")
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)
	availabilityZones := flags.FlagToBoolValue(p, cmd, availabilityZonesFlag)
	kubernetesVersions := flags.FlagToBoolValue(p, cmd, kubernetesVersionsFlag)
	machineImages := flags.FlagToBoolValue(p, cmd, machineImagesFlag)
	machineTypes := flags.FlagToBoolValue(p, cmd, machineTypesFlag)
	volumeTypes := flags.FlagToBoolValue(p, cmd, volumeTypesFlag)

	// If no flag was passed, take it as if every flag were passed
	if !availabilityZones && !kubernetesVersions && !machineImages && !machineTypes && !volumeTypes {
		availabilityZones = true
		kubernetesVersions = true
		machineImages = true
		machineTypes = true
		volumeTypes = true
	}

	model := inputModel{
		GlobalFlagModel:    globalFlags,
		AvailabilityZones:  availabilityZones,
		KubernetesVersions: kubernetesVersions,
		MachineImages:      machineImages,
		MachineTypes:       machineTypes,
		VolumeTypes:        volumeTypes,
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

func buildRequest(ctx context.Context, apiClient *ske.APIClient) ske.ApiListProviderOptionsRequest {
	req := apiClient.ListProviderOptions(ctx)
	return req
}

func outputResult(p *print.Printer, model *inputModel, options *ske.ProviderOptions) error {
	if model == nil || model.GlobalFlagModel == nil {
		return fmt.Errorf("model is nil")
	} else if options == nil {
		return fmt.Errorf("options is nil")
	}

	// filter output based on the flags
	if !model.AvailabilityZones {
		options.AvailabilityZones = nil
	}

	if !model.KubernetesVersions {
		options.KubernetesVersions = nil
	}

	if !model.MachineImages {
		options.MachineImages = nil
	}

	if !model.MachineTypes {
		options.MachineTypes = nil
	}

	if !model.VolumeTypes {
		options.VolumeTypes = nil
	}

	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(options, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SKE options: %w", err)
		}
		p.Outputln(string(details))
		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(options, yaml.IndentSequence(true), yaml.UseJSONMarshaler())
		if err != nil {
			return fmt.Errorf("marshal SKE options: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		return outputResultAsTable(p, options)
	}
}

func outputResultAsTable(p *print.Printer, options *ske.ProviderOptions) error {
	if options == nil {
		return fmt.Errorf("options is nil")
	}

	content := []tables.Table{}
	if options.AvailabilityZones != nil && len(*options.AvailabilityZones) != 0 {
		content = append(content, buildAvailabilityZonesTable(options))
	}

	if options.KubernetesVersions != nil && len(*options.KubernetesVersions) != 0 {
		kubernetesVersionsTable, err := buildKubernetesVersionsTable(options)
		if err != nil {
			return fmt.Errorf("build Kubernetes versions table: %w", err)
		}
		content = append(content, kubernetesVersionsTable)
	}

	if options.MachineImages != nil && len(*options.MachineImages) != 0 {
		content = append(content, buildMachineImagesTable(options))
	}

	if options.MachineTypes != nil && len(*options.MachineTypes) != 0 {
		content = append(content, buildMachineTypesTable(options))
	}

	if options.VolumeTypes != nil && len(*options.VolumeTypes) != 0 {
		content = append(content, buildVolumeTypesTable(options))
	}

	err := tables.DisplayTables(p, content)
	if err != nil {
		return fmt.Errorf("display output: %w", err)
	}

	return nil
}

func buildAvailabilityZonesTable(resp *ske.ProviderOptions) tables.Table {
	zones := *resp.AvailabilityZones

	table := tables.NewTable()
	table.SetTitle("Availability Zones")
	table.SetHeader("ZONE")
	for i := range zones {
		z := zones[i]
		table.AddRow(*z.Name)
	}
	return table
}

func buildKubernetesVersionsTable(resp *ske.ProviderOptions) (tables.Table, error) {
	versions := *resp.KubernetesVersions

	table := tables.NewTable()
	table.SetTitle("Kubernetes Versions")
	table.SetHeader("VERSION", "STATE", "EXPIRATION DATE", "FEATURE GATES")
	for i := range versions {
		v := versions[i]
		featureGate, err := json.Marshal(*v.FeatureGates)
		if err != nil {
			return table, fmt.Errorf("marshal featureGates of Kubernetes version %q: %w", *v.Version, err)
		}
		expirationDate := ""
		if v.ExpirationDate != nil {
			expirationDate = v.ExpirationDate.Format(time.RFC3339)
		}
		table.AddRow(
			utils.PtrString(v.Version),
			utils.PtrString(v.State),
			expirationDate,
			string(featureGate))
	}
	return table, nil
}

func buildMachineImagesTable(resp *ske.ProviderOptions) tables.Table {
	images := *resp.MachineImages

	table := tables.NewTable()
	table.SetTitle("Machine Images")
	table.SetHeader("NAME", "VERSION", "STATE", "EXPIRATION DATE", "SUPPORTED CRI")
	for i := range images {
		image := images[i]
		versions := *image.Versions
		for j := range versions {
			version := versions[j]
			criNames := make([]string, 0)
			for i := range *version.Cri {
				cri := (*version.Cri)[i]
				criNames = append(criNames, *cri.Name)
			}
			criNamesString := strings.Join(criNames, ", ")

			expirationDate := "-"
			if version.ExpirationDate != nil {
				expirationDate = version.ExpirationDate.Format(time.RFC3339)
			}
			table.AddRow(
				utils.PtrString(image.Name),
				utils.PtrString(version.Version),
				utils.PtrString(version.State),
				expirationDate,
				criNamesString,
			)
		}
	}
	table.EnableAutoMergeOnColumns(1)
	return table
}

func buildMachineTypesTable(resp *ske.ProviderOptions) tables.Table {
	types := *resp.MachineTypes

	table := tables.NewTable()
	table.SetTitle("Machine Types")
	table.SetHeader("TYPE", "CPU", "MEMORY")
	for i := range types {
		t := types[i]
		table.AddRow(
			utils.PtrString(t.Name),
			utils.PtrString(t.Cpu),
			utils.PtrString(t.Memory),
		)
	}
	return table
}

func buildVolumeTypesTable(resp *ske.ProviderOptions) tables.Table {
	types := *resp.VolumeTypes

	table := tables.NewTable()
	table.SetTitle("Volume Types")
	table.SetHeader("TYPE")
	for i := range types {
		z := types[i]
		table.AddRow(utils.PtrString(z.Name))
	}
	return table
}
