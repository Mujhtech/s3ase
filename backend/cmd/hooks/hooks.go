package hooks

import "github.com/spf13/cobra"

var skipHook = map[string]struct{}{
	"up":      {},
	"down":    {},
	"version": {},
}

func PreHook() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if _, ok := skipHook[cmd.Use]; ok {
			return nil
		}

		return nil
	}
}

func PostHook() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {

		return nil
	}
}
