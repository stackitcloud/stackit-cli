package dns

import (
	"fmt"

	recordset "github.com/stackitcloud/stackit-cli/internal/cmd/dns/record-set"
	"github.com/stackitcloud/stackit-cli/internal/cmd/dns/zone"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "dns",
	Short:   "Provides functionality for DNS",
	Long:    "Provides functionality for DNS",
	Example: fmt.Sprintf("%s\n%s", zone.Cmd.Example, recordset.Cmd.Example),
}

func init() {
	Cmd.AddCommand(zone.Cmd)
	Cmd.AddCommand(recordset.Cmd)
}
