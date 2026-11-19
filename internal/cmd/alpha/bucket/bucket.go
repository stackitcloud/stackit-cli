package bucket

import (
	"github.com/spf13/cobra"

	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/bucket/cp"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/bucket/ls"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/bucket/mb"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/bucket/mv"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/bucket/presign"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/bucket/rb"
	"github.com/stackitcloud/stackit-cli/internal/cmd/alpha/bucket/rm"
	"github.com/stackitcloud/stackit-cli/internal/pkg/args"
	"github.com/stackitcloud/stackit-cli/internal/pkg/examples"
	"github.com/stackitcloud/stackit-cli/internal/pkg/types"
	"github.com/stackitcloud/stackit-cli/internal/pkg/utils"
)

// The bucket command group is in alpha because the Object Storage data-plane CLI integration is still experimental.
func NewCmd(params *types.CmdParams) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bucket",
		Short: "Manages objects in STACKIT Object Storage buckets (data plane)",
		Long: `Manages objects in STACKIT Object Storage buckets.

Buckets and objects can be addressed with s3 URIs (s3://BUCKET[/KEY]) or
STACKIT Object Storage URLs in virtual hosted style
(https://BUCKET.object.storage.REGION.onstackit.cloud[/KEY]) or path style
(https://object.storage.REGION.onstackit.cloud/BUCKET[/KEY]).

Credentials (access key ID / secret access key) are required and can be provided
with the --access-key-id and --secret-access-key flags, the STACKIT_OBJECT_STORAGE_ACCESS_KEY_ID
and STACKIT_OBJECT_STORAGE_SECRET_ACCESS_KEY environment variables, or the config file.
Create credentials with 'stackit object-storage credentials create'.`,
		Args: args.NoArgs,
		Run:  utils.CmdHelp,
		Example: examples.Build(
			examples.NewExample(
				"List all buckets",
				"$ stackit alpha bucket ls"),
			examples.NewExample(
				"Upload a file to a bucket",
				`$ stackit alpha bucket cp local/file.txt s3://my-bucket/file.txt`),
		),
	}
	addSubcommands(cmd, params)
	return cmd
}

func addSubcommands(cmd *cobra.Command, params *types.CmdParams) {
	cmd.AddCommand(cp.NewCmd(params))
	cmd.AddCommand(ls.NewCmd(params))
	cmd.AddCommand(mb.NewCmd(params))
	cmd.AddCommand(mv.NewCmd(params))
	cmd.AddCommand(presign.NewCmd(params))
	cmd.AddCommand(rb.NewCmd(params))
	cmd.AddCommand(rm.NewCmd(params))
}
