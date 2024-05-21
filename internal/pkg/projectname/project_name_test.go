package projectname

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stackitcloud/stackit-cli/internal/pkg/config"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

var testProjectId = uuid.NewString()

func TestGetProjectName(t *testing.T) {
	tests := []struct {
		description string
		projectName string
		projectId   string
		isValid     bool
	}{
		{
			description: "Project name from config",
			projectName: "project-name",
			projectId:   testProjectId,
			isValid:     true,
		},
		{
			description: "empty project name and id",
			projectName: "",
			projectId:   "",
			isValid:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			viper.Set(config.ProjectNameKey, tt.projectName)
			viper.Set(config.ProjectIdKey, tt.projectId)
			defer viper.Reset()
			p := print.NewPrinter()
			cmd := &cobra.Command{}

			projectName, err := GetProjectName(context.Background(), p, cmd)
			if err != nil {
				if tt.isValid {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if !tt.isValid {
				t.Fatalf("expected error, got project name %q", projectName)
			}

			if projectName != tt.projectName {
				t.Fatalf("expected project name %q, got %q", tt.projectName, projectName)
			}
		})
	}
}
