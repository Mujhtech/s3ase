package version

import "github.com/spf13/cobra"

func RegisterVersionCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version and exit",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return cmd

}
