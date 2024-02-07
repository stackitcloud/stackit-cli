package remove

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/confirm"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/authorization/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/authorization"
)

const (
	subjectArg = "SUBJECT"

	organizationIdFlag = "organization-id"
	roleFlag           = "role"
	forceFlag          = "force"

	organizationResourceType = "organization"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	OrganizationId *string
	Subject        string
	Role           *string
	Force          bool
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("remove %s", subjectArg),
		Short: "Removes a member from an organization",
		Long: fmt.Sprintf("%s\n%s\n%s",
			"Removes a member from an organization.",
			"A member is a combination of a subject (user, service account or client) and a role.",
			"The subject is usually email address (for users) or name (for clients).",
		),
		Args: args.SingleArg(subjectArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Remove a member (user "someone@domain.com" with an "editor" role) from an organization`,
				"$ stackit organization member remove someone@domain.com --organization-id xxx --role editor"),
			examples.NewExample(
				`Remove a member (user "someone@domain.com" with a "reader" role) from an organization, along with all other roles of the subject that would stop the removal of the "reader" role`,
				"$ stackit organization member remove someone@domain.com --organization-id xxx --role reader --force"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to remove the %s role from %s on organization with ID %s?", *model.Role, model.Subject, *model.OrganizationId)
				if model.Force {
					prompt = fmt.Sprintf("%s This will also remove other roles of the subject that would stop the removal of the requested role", prompt)
				}
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("remove member: %w", err)
			}

			cmd.Println("Member removed")
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(organizationIdFlag, "", "The organization ID")
	cmd.Flags().String(roleFlag, "", "The role to be removed from the subject")
	cmd.Flags().Bool(forceFlag, false, "When true, removes other roles of the subject that would stop the removal of the requested role")

	err := flags.MarkFlagsRequired(cmd, organizationIdFlag, roleFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	subject := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)

	return &inputModel{
		GlobalFlagModel: globalFlags,
		OrganizationId:  flags.FlagToStringPointer(cmd, organizationIdFlag),
		Subject:         subject,
		Role:            flags.FlagToStringPointer(cmd, roleFlag),
		Force:           flags.FlagToBoolValue(cmd, forceFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *authorization.APIClient) authorization.ApiRemoveMembersRequest {
	req := apiClient.RemoveMembers(ctx, *model.OrganizationId)
	payload := authorization.RemoveMembersPayload{
		Members: utils.Ptr([]authorization.Member{
			{
				Subject: utils.Ptr(model.Subject),
				Role:    model.Role,
			},
		}),
		ResourceType: utils.Ptr(organizationResourceType),
	}
	payload.ForceRemove = &model.Force
	req = req.RemoveMembersPayload(payload)
	return req
}
