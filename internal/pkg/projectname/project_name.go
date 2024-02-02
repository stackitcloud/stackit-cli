package projectname

import (
	"context"
	"fmt"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Returns the project name associated to the project ID set in config
//
// Uses the one stored in config if it's valid, otherwise gets it from the API
func GetProjectName(ctx context.Context, cmd *cobra.Command) (string, error) {
	// If we can use the project name from config, return it
	if useProjectNameFromConfig(cmd) {
		return viper.GetString(config.ProjectNameKey), nil
	}

	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return "", fmt.Errorf("found empty project ID and name")
	}

	apiClient, err := client.ConfigureClient(cmd)
	if err != nil {
		return "", fmt.Errorf("configure resource manager client: %w", err)
	}
	req := apiClient.GetProject(ctx, projectId)
	resp, err := req.Execute()
	if err != nil {
		return "", fmt.Errorf("read project details: %w", err)
	}
	projectName := *resp.Name

	// If project ID is set in config, we store the project name in config
	// (So next time we can just pull it from there)
	if !isProjectIdSetInFlags(cmd) {
		viper.Set(config.ProjectNameKey, projectName)
		err = viper.WriteConfig()
		if err != nil {
			return "", fmt.Errorf("write new config to file: %w", err)
		}
	}

	return projectName, nil
}

// Returns True if project name from config should be used, False otherwise
func useProjectNameFromConfig(cmd *cobra.Command) bool {
	// We use the project name from the config file, if:
	// - Project id is not set to a different value than the one in the config file
	// - Project name in the config file is not empty
	projectIdSet := isProjectIdSetInFlags(cmd)
	projectName := viper.GetString(config.ProjectNameKey)
	projectNameSet := false
	if projectName != "" {
		projectNameSet = true
	}
	return !projectIdSet && projectNameSet
}

func isProjectIdSetInFlags(cmd *cobra.Command) bool {
	// FlagToStringPointer pulls the projectId from passed flags
	// viper.GetString uses the flags, and fallsback to config file
	// To check if projectId was passed, we use the first rather than the second
	projectIdFromFlag := flags.FlagToStringPointer(cmd, globalflags.ProjectIdFlag)
	projectIdSetInFlag := false
	if projectIdFromFlag != nil {
		projectIdSetInFlag = true
	}
	return projectIdSetInFlag
}
