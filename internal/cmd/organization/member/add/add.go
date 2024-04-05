package add

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
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/authorization"
)

const (
	subjectArg = "SUBJECT"

	organizationIdFlag = "organization-id"
	roleFlag           = "role"

	organizationResourceType = "organization"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	OrganizationId *string
	Subject        string
	Role           *string
}

func NewCmd(p *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("add %s", subjectArg),
		Short: "Adds a member to an organization",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
			"Adds a member to an organization.",
			"A member is a combination of a subject (user, service account or client) and a role.",
			"The subject is usually email address for users or name in case of clients",
			"For more details on the available roles, run:",
			"  $ stackit organization role list --organization-id <RESOURCE ID>",
		),
		Args: args.SingleArg(subjectArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Add a member to an organization with the "reader" role`,
				"$ stackit organization member add someone@domain.com --organization-id xxx --role reader"),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			model, err := parseInput(cmd, args)
			if err != nil {
				return err
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to add the %s role to %s on organization with ID %q?", *model.Role, model.Subject, *model.OrganizationId)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Configure API client
			apiClient, err := client.ConfigureClient(cmd)
			if err != nil {
				return err
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("add member: %w", err)
			}

			cmd.Println("Member added")
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(organizationIdFlag, "", "The organization ID")
	cmd.Flags().String(roleFlag, "", "The role to add to the subject")

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
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *authorization.APIClient) authorization.ApiAddMembersRequest {
	req := apiClient.AddMembers(ctx, *model.OrganizationId)
	req = req.AddMembersPayload(authorization.AddMembersPayload{
		Members: utils.Ptr([]authorization.Member{
			{
				Subject: utils.Ptr(model.Subject),
				Role:    model.Role,
			},
		}),
		ResourceType: utils.Ptr(organizationResourceType),
	})
	return req
}
