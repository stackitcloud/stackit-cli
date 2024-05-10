package options

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/ske/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/tables"

	"github.com/spf13/cobra"
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

func NewCmd(p *print.Printer) *cobra.Command {
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

			// Call API
			req := buildRequest(ctx, apiClient)
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("get SKE provider options: %w", err)
			}

			return outputResult(p, model, resp)
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
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(options, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal SKE options: %w", err)
		}
		p.Outputln(string(details))
		return nil
	default:
		return outputResultAsTable(p, model, options)
	}
}

func outputResultAsTable(p *print.Printer, model *inputModel, options *ske.ProviderOptions) error {
	content := ""
	if model.AvailabilityZones {
		content += renderAvailabilityZones(options)
	}
	if model.KubernetesVersions {
		kubernetesVersionsRendered, err := renderKubernetesVersions(options)
		if err != nil {
			return fmt.Errorf("render Kubernetes versions: %w", err)
		}
		content += kubernetesVersionsRendered
	}
	if model.MachineImages {
		content += renderMachineImages(options)
	}
	if model.MachineTypes {
		content += renderMachineTypes(options)
	}
	if model.VolumeTypes {
		content += renderVolumeTypes(options)
	}

	err := p.PagerDisplay(content)
	if err != nil {
		return fmt.Errorf("display output: %w", err)
	}

	return nil
}

func renderAvailabilityZones(resp *ske.ProviderOptions) string {
	zones := *resp.AvailabilityZones

	table := tables.NewTable()
	table.SetTitle("Availability Zones")
	table.SetHeader("ZONE")
	for i := range zones {
		z := zones[i]
		table.AddRow(*z.Name)
	}
	return table.Render()
}

func renderKubernetesVersions(resp *ske.ProviderOptions) (string, error) {
	versions := *resp.KubernetesVersions

	table := tables.NewTable()
	table.SetTitle("Kubernetes Versions")
	table.SetHeader("VERSION", "STATE", "EXPIRATION DATE", "FEATURE GATES")
	for i := range versions {
		v := versions[i]
		featureGate, err := json.Marshal(*v.FeatureGates)
		if err != nil {
			return "", fmt.Errorf("marshal featureGates of Kubernetes version %q: %w", *v.Version, err)
		}
		expirationDate := ""
		if v.ExpirationDate != nil {
			expirationDate = *v.ExpirationDate
		}
		table.AddRow(*v.Version, *v.State, expirationDate, string(featureGate))
	}
	return table.Render(), nil
}

func renderMachineImages(resp *ske.ProviderOptions) string {
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
				expirationDate = *version.ExpirationDate
			}
			table.AddRow(*image.Name, *version.Version, *version.State, expirationDate, criNamesString)
		}
	}
	table.EnableAutoMergeOnColumns(1)
	return table.Render()
}

func renderMachineTypes(resp *ske.ProviderOptions) string {
	types := *resp.MachineTypes

	table := tables.NewTable()
	table.SetTitle("Machine Types")
	table.SetHeader("TYPE", "CPU", "MEMORY")
	for i := range types {
		t := types[i]
		table.AddRow(*t.Name, *t.Cpu, *t.Memory)
	}
	return table.Render()
}

func renderVolumeTypes(resp *ske.ProviderOptions) string {
	types := *resp.VolumeTypes

	table := tables.NewTable()
	table.SetTitle("Volume Types")
	table.SetHeader("TYPE")
	for i := range types {
		z := types[i]
		table.AddRow(*z.Name)
	}
	return table.Render()
}
