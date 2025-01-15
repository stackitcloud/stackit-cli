package rotate

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

const (
	clusterNameArg = "CLUSTER_NAME"
)

func NewCmd(_ *print.Printer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("rotate %s", clusterNameArg),
		Short: "Rotates credentials associated to a SKE cluster",
		Long:  "Rotates credentials associated to a STACKIT Kubernetes Engine (SKE) cluster. The old credentials will be invalid after the operation.",
		Args:  args.NoArgs,
		Deprecated: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n",
			"and was removed.",
			"Please use the 2-step credential rotation flow instead, by running the commands:",
			" $ stackit ske credentials start-rotation CLUSTER_NAME",
			" $ stackit ske credentials complete-rotation CLUSTER_NAME",
			"For more information, visit: https://docs.stackit.cloud/stackit/en/how-to-rotate-ske-credentials-200016334.html",
		),
		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}
	return cmd
}
