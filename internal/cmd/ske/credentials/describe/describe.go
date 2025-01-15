package describe

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
		Use:   fmt.Sprintf("describe %s", clusterNameArg),
		Short: "Shows details of the credentials associated to a SKE cluster",
		Long:  "Shows details of the credentials associated to a STACKIT Kubernetes Engine (SKE) cluster",
		Args:  args.NoArgs,
		Deprecated: fmt.Sprintf("%s\n%s\n%s\n%s\n",
			"and was removed.",
			"Please use the following command to obtain a kubeconfig file instead:",
			" $ stackit ske kubeconfig create CLUSTER_NAME",
			"For more information, visit: https://docs.stackit.cloud/stackit/en/how-to-rotate-ske-credentials-200016334.html",
		),

		RunE: func(_ *cobra.Command, _ []string) error {
			return nil
		},
	}
	return cmd
}
