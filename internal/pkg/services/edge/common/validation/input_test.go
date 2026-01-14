// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2025 STACKIT GmbH & Co. KG

package validation

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
	commonErr "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/error"
	commonInstance "github.com/stackitcloud/stackit-cli/internal/pkg/services/edge/common/instance"
	testUtils "github.com/stackitcloud/stackit-cli/internal/pkg/testutils"
)

func TestGetValidatedInstanceIdentifier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(*cobra.Command)
		want    *Identifier
		wantErr any
	}{
		{
			name: "instance id success",
			setup: func(cmd *cobra.Command) {
				cmd.Flags().String(commonInstance.InstanceIdFlag, "", "")
				_ = cmd.Flags().Set(commonInstance.InstanceIdFlag, "edgesvc01")
			},
			want: &Identifier{Flag: commonInstance.InstanceIdFlag, Value: "edgesvc01"},
		},
		{
			name: "display name success",
			setup: func(cmd *cobra.Command) {
				cmd.Flags().String(commonInstance.DisplayNameFlag, "", "")
				_ = cmd.Flags().Set(commonInstance.DisplayNameFlag, "edge01")
			},
			want: &Identifier{Flag: commonInstance.DisplayNameFlag, Value: "edge01"},
		},
		{
			name: "instance id validation error",
			setup: func(cmd *cobra.Command) {
				cmd.Flags().String(commonInstance.InstanceIdFlag, "", "")
				_ = cmd.Flags().Set(commonInstance.InstanceIdFlag, "id")
			},
			wantErr: "too short",
		},
		{
			name: "display name validation error",
			setup: func(cmd *cobra.Command) {
				cmd.Flags().String(commonInstance.DisplayNameFlag, "", "")
				_ = cmd.Flags().Set(commonInstance.DisplayNameFlag, "x")
			},
			wantErr: "too short",
		},
		{
			name:    "no identifier",
			setup:   func(_ *cobra.Command) {},
			wantErr: &commonErr.NoIdentifierError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			printer := print.NewPrinter()
			cmd := &cobra.Command{Use: "test"}
			tt.setup(cmd)

			got, err := GetValidatedInstanceIdentifier(printer, cmd)
			if !testUtils.AssertError(t, err, tt.wantErr) {
				return
			}
			if tt.want != nil {
				testUtils.AssertValue(t, got, tt.want)
			}
		})
	}
}
