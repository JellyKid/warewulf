package export

import "github.com/spf13/cobra"

var (
	baseCmd = &cobra.Command{
		Use:                "export",
		Short:              "Kernel Image Export",
		Long:               "Export kernel image",
		RunE:				CobraRunE,
		Args: 				cobra.ExactArgs(1),

	}
)

func init() {
}

// GetRootCommand returns the root cobra.Command for the application.
func GetCommand() *cobra.Command {
	return baseCmd
}