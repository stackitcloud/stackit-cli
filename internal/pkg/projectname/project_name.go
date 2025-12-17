package projectname

import (
	"context"
	"fmt"
	"os"

	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/globalflags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/client"
	"github.com/stackitcloud/stackit-cli/internal/pkg/services/resourcemanager/utils"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Returns the project name associated to the project ID set in config
//
// Uses the one stored in config if it's valid, otherwise gets it from the API
func GetProjectName(ctx context.Context, p *print.Printer, cliVersion string, cmd *cobra.Command) (string, error) {
	// If we can use the project name from config, return it
	if useProjectNameFromConfig(p, cmd) {
		return viper.GetString(config.ProjectNameKey), nil
	}

	projectId := viper.GetString(config.ProjectIdKey)
	if projectId == "" {
		return "", fmt.Errorf("found empty project ID and name")
	}

	apiClient, err := client.ConfigureClient(p, cliVersion)
	if err != nil {
		return "", fmt.Errorf("configure resource manager client: %w", err)
	}

	projectName, err := utils.GetProjectName(ctx, apiClient, projectId)
	if err != nil {
		return "", fmt.Errorf("get project name: %w", err)
	}

	// If project ID is set in config, we store the project name in config
	// (So next time we can just pull it from there)
	if !isProjectIdSetInFlags(p, cmd) && !isProjectIdSetInEnvVar() {
		viper.Set(config.ProjectNameKey, projectName)
		err = config.Write()
		if err != nil {
			return "", fmt.Errorf("write new config to file: %w", err)
		}
	}

	return projectName, nil
}

// Returns True if project name from config should be used, False otherwise
func useProjectNameFromConfig(p *print.Printer, cmd *cobra.Command) bool {
	// We use the project name from the config file, if:
	// - Project id is not set to a different value than the one in the config file
	// - Project name in the config file is not empty
	projectIdSetInFlags := isProjectIdSetInFlags(p, cmd)
	projectIdSetInEnv := isProjectIdSetInEnvVar()
	projectName := viper.GetString(config.ProjectNameKey)
	projectNameSet := projectName != ""
	return !projectIdSetInFlags && !projectIdSetInEnv && projectNameSet
}

func isProjectIdSetInFlags(p *print.Printer, cmd *cobra.Command) bool {
	// FlagToStringPointer pulls the projectId from passed flags
	// viper.GetString uses the flags, and fallsback to config file
	// To check if projectId was passed, we use the first rather than the second
	projectIdFromFlag := flags.FlagToStringPointer(p, cmd, globalflags.ProjectIdFlag)
	projectIdSetInFlag := projectIdFromFlag != nil
	return projectIdSetInFlag
}

func isProjectIdSetInEnvVar() bool {
	// Reads the project Id from the environment variable PROJECT_ID
	_, projectIdSetInEnv := os.LookupEnv("STACKIT_PROJECT_ID")
	return projectIdSetInEnv
}
