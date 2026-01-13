// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package edge

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/edge/instance"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/edge/kubeconfig"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/edge/plans"
	"github.com/stackitcloud/stackit-cli/internal/cmd/beta/edge/token"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edge-cloud",
		Short: "Provides functionality for edge services.",
		Long:  "Provides functionality for STACKIT Edge Cloud (STEC) services.",
		Args:  args.NoArgs,
		Run:   utils.CmdHelp,
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(instance.NewCmd(params))
	cmd.AddCommand(plans.NewCmd(params))
	cmd.AddCommand(kubeconfig.NewCmd(params))
	cmd.AddCommand(token.NewCmd(params))
}
