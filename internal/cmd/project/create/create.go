package create

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/goccy/go-yaml"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/auth"
	"github.com/stackitcloud/stackit-cli/internal/pkg/errors"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/resourcemanager"
)

const (
	parentIdFlag      = "parent-id"
	nameFlag          = "name"
	labelFlag         = "label"
	networkAreaIdFlag = "network-area-id"

	ownerRole        = "project.owner"
	labelKeyRegex    = `[A-ZÄÜÖa-zäüöß0-9_-]{1,64}`
	labelValueRegex  = `^$|[A-ZÄÜÖa-zäüöß0-9_-]{1,64}`
	networkAreaLabel = "networkArea"
)

type inputModel struct {
	*globalflags.GlobalFlagModel
	ParentId      *string
	Name          *string
	Labels        *map[string]string
	NetworkAreaId *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Creates a STACKIT project",
		Long:  "Creates a STACKIT project.",
		Args:  args.NoArgs,
		Example: examples.Build(
			examples.NewExample(
				`Create a STACKIT project`,
				"$ stackit project create --parent-id xxxx --name my-project"),
			examples.NewExample(
				`Create a STACKIT project with a set of labels`,
				"$ stackit project create --parent-id xxxx --name my-project --label key=value --label foo=bar"),
			examples.NewExample(
				`Create a STACKIT project with a network area`,
				"$ stackit project create --parent-id xxxx --name my-project --network-area-id yyyy"),
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

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to create a project under the parent with ID %q?", *model.ParentId)
				err = p.PromptForConfirmation(prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req, err := buildRequest(ctx, model, apiClient)
			if err != nil {
				return fmt.Errorf("build project creation request: %w", err)
			}
			resp, err := req.Execute()
			if err != nil {
				return fmt.Errorf("create project: %w", err)
			}

			return outputResult(p, model, resp)
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(parentIdFlag, "", "Parent resource identifier. Both container ID (user-friendly) and UUID are supported")
	cmd.Flags().String(nameFlag, "", "Project name")
	cmd.Flags().StringToString(labelFlag, nil, "Labels are key-value string pairs which can be attached to a project. A label can be provided with the format key=value and the flag can be used multiple times to provide a list of labels")
	cmd.Flags().Var(flags.UUIDFlag(), networkAreaIdFlag, "Network area ID")

	err := flags.MarkFlagsRequired(cmd, parentIdFlag, nameFlag)
	cobra.CheckErr(err)
}

func parseInput(p *print.Printer, cmd *cobra.Command) (*inputModel, error) {
	globalFlags := globalflags.Parse(p, cmd)

	labels := flags.FlagToStringToStringPointer(p, cmd, labelFlag)
	if labels != nil {
		labelKeyRegex := regexp.MustCompile(labelKeyRegex)
		labelValueRegex := regexp.MustCompile(labelValueRegex)
		for key, value := range *labels {
			if !labelKeyRegex.MatchString(key) {
				return nil, &errors.FlagValidationError{
					Flag:    labelFlag,
					Details: fmt.Sprintf("label key %s didn't match the required regex expression %s", key, labelKeyRegex),
				}
			}

			if !labelValueRegex.MatchString(value) {
				return nil, &errors.FlagValidationError{
					Flag:    labelFlag,
					Details: fmt.Sprintf("label value %s for key %s didn't match the required regex expression %s", value, key, labelValueRegex),
				}
			}
		}
	}

	model := inputModel{
		GlobalFlagModel: globalFlags,
		ParentId:        flags.FlagToStringPointer(p, cmd, parentIdFlag),
		Name:            flags.FlagToStringPointer(p, cmd, nameFlag),
		Labels:          labels,
		NetworkAreaId:   flags.FlagToStringPointer(p, cmd, networkAreaIdFlag),
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

func buildRequest(ctx context.Context, model *inputModel, apiClient *resourcemanager.APIClient) (resourcemanager.ApiCreateProjectRequest, error) {
	req := apiClient.CreateProject(ctx)

	authFlow, err := auth.GetAuthFlow()
	if err != nil {
		return req, fmt.Errorf("get authentication flow: %w", err)
	}
	var email string
	switch authFlow {
	case auth.AUTH_FLOW_SERVICE_ACCOUNT_TOKEN:
		email, err = auth.GetAuthField(auth.SERVICE_ACCOUNT_EMAIL)
		if err != nil {
			return req, fmt.Errorf("get email of the service account that was used to authenticate: %w", err)
		}
	case auth.AUTH_FLOW_SERVICE_ACCOUNT_KEY:
		email, err = auth.GetAuthField(auth.SERVICE_ACCOUNT_EMAIL)
		if err != nil {
			return req, fmt.Errorf("get email of the service account that was used to authenticate: %w", err)
		}
	case auth.AUTH_FLOW_USER_TOKEN:
		email, err = auth.GetAuthField(auth.USER_EMAIL)
		if err != nil {
			return req, fmt.Errorf("get your user email from configuration: %w", err)
		}
	default:
		return req, fmt.Errorf("the configured authentication flow (%s) is not supported, please report this issue", authFlow)
	}

	if email == "" {
		return req, fmt.Errorf("the authenticated subject email cannot be empty, please report this issue")
	}

	labels := model.Labels

	if model.NetworkAreaId != nil {
		if labels == nil {
			labels = &map[string]string{}
		}
		(*labels)[networkAreaLabel] = *model.NetworkAreaId
	}

	req = req.CreateProjectPayload(resourcemanager.CreateProjectPayload{
		ContainerParentId: model.ParentId,
		Name:              model.Name,
		Labels:            labels,
		Members: &[]resourcemanager.Member{
			{
				Role:    utils.Ptr(ownerRole),
				Subject: utils.Ptr(email),
			},
		},
	})

	return req, nil
}

func outputResult(p *print.Printer, model *inputModel, resp *resourcemanager.Project) error {
	switch model.OutputFormat {
	case print.JSONOutputFormat:
		details, err := json.MarshalIndent(resp, "", "  ")
		if err != nil {
			return fmt.Errorf("marshal project: %w", err)
		}
		p.Outputln(string(details))

		return nil
	case print.YAMLOutputFormat:
		details, err := yaml.MarshalWithOptions(resp, yaml.IndentSequence(true))
		if err != nil {
			return fmt.Errorf("marshal project: %w", err)
		}
		p.Outputln(string(details))

		return nil
	default:
		p.Outputf("Created project under the parent with ID %q. Project ID: %s\n", *model.ParentId, *resp.ProjectId)
		return nil
	}
}
