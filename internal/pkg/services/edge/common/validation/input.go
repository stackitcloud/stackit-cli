// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package validation

import (
	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/flags"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	commonErr "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/error"
	commonInstance "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/instance"
)

// Struct to model the instance identifier provided by the user (either instance-id or display-name)
type Identifier struct {
	Flag  string
	Value string
}

// GetValidatedInstanceIdentifier gets and validates the instance identifier provided by the user through the command-line flags.
// It checks for either an instance ID or a display name and validates the provided value.
//
// p is the printer used for logging.
// cmd is the cobra command that holds the flags.
//
// Returns an Identifier struct containing the flag and its value if a valid identifier is provided, otherwise returns an error.
// Indirect unit tests of GetValidatedInstanceIdentifier are done within the respective CLI packages.
func GetValidatedInstanceIdentifier(p *print.Printer, cmd *cobra.Command) (*Identifier, error) {
	switch {
	case cmd.Flags().Changed(commonInstance.InstanceIdFlag):
		instanceIdValue := flags.FlagToStringPointer(p, cmd, commonInstance.InstanceIdFlag)
		if err := commonInstance.ValidateInstanceId(instanceIdValue); err != nil {
			return nil, err
		}
		return &Identifier{
			Flag:  commonInstance.InstanceIdFlag,
			Value: *instanceIdValue,
		}, nil
	case cmd.Flags().Changed(commonInstance.DisplayNameFlag):
		displayNameValue := flags.FlagToStringPointer(p, cmd, commonInstance.DisplayNameFlag)
		if err := commonInstance.ValidateDisplayName(displayNameValue); err != nil {
			return nil, err
		}
		return &Identifier{
			Flag:  commonInstance.DisplayNameFlag,
			Value: *displayNameValue,
		}, nil
	default:
		return nil, commonErr.NewNoIdentifierError("")
	}
}
