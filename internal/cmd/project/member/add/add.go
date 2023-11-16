package add

import (
	"context"
	"fmt"

	"stackit/internal/pkg/args"
	"stackit/internal/pkg/confirm"
	"stackit/internal/pkg/errors"
	"stackit/internal/pkg/examples"
	"stackit/internal/pkg/flags"
	"stackit/internal/pkg/globalflags"
	"stackit/internal/pkg/projectname"
	"stackit/internal/pkg/services/membership/client"
	"stackit/internal/pkg/utils"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-sdk-go/services/membership"
)

const (
	roleFlag = "role"

	subjectArg = "SUBJECT"

	projectResourceType = "project"
)

type inputModel struct {
	*globalflags.GlobalFlagModel

	Subject string
	Role    *string
}

func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("add %s", subjectArg),
		Short: "Add a member to a project",
		Long: fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
			"Add a member to a project.",
			"A member is a combination of a subject (user, service account or client) and a role.",
			"The subject is usually email address for users or name in case of clients",
			"For more details on the available roles, run:",
			"  $ stackit project role list --project-id <PROJECT ID>",
		),
		Args: args.SingleArg(subjectArg, nil),
		Example: examples.Build(
			examples.NewExample(
				`Add a member to a project with the "reader" role`,
				"$ stackit project member add someone@domain.com --project-id xxx --role reader"),
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

			projectLabel, err := projectname.GetProjectName(ctx, cmd)
			if err != nil {
				projectLabel = model.ProjectId
			}

			if !model.AssumeYes {
				prompt := fmt.Sprintf("Are you sure you want to add the role %s to %s on project %s?", *model.Role, model.Subject, projectLabel)
				err = confirm.PromptForConfirmation(cmd, prompt)
				if err != nil {
					return err
				}
			}

			// Call API
			req := buildRequest(ctx, model, apiClient)
			_, err = req.Execute()
			if err != nil {
				return fmt.Errorf("add member: %w", err)
			}

			cmd.Printf("Added the role %s to %s on project %s\n", *model.Role, model.Subject, projectLabel)
			return nil
		},
	}
	configureFlags(cmd)
	return cmd
}

func configureFlags(cmd *cobra.Command) {
	cmd.Flags().String(roleFlag, "", "The role to add to the subject")

	err := flags.MarkFlagsRequired(cmd, roleFlag)
	cobra.CheckErr(err)
}

func parseInput(cmd *cobra.Command, inputArgs []string) (*inputModel, error) {
	subject := inputArgs[0]

	globalFlags := globalflags.Parse(cmd)
	if globalFlags.ProjectId == "" {
		return nil, &errors.ProjectIdError{}
	}

	return &inputModel{
		GlobalFlagModel: globalFlags,
		Subject:         subject,
		Role:            flags.FlagToStringPointer(cmd, roleFlag),
	}, nil
}

func buildRequest(ctx context.Context, model *inputModel, apiClient *membership.APIClient) membership.ApiAddMembersRequest {
	req := apiClient.AddMembers(ctx, model.GlobalFlagModel.ProjectId)
	req = req.AddMembersPayload(membership.AddMembersPayload{
		Members: utils.Ptr([]membership.Member{
			{
				Subject: utils.Ptr(model.Subject),
				Role:    model.Role,
			},
		}),
		ResourceType: utils.Ptr(projectResourceType),
	})
	return req
}
