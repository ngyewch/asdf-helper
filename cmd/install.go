package cmd

import (
	"github.com/ngyewch/asdf-helper/helper"
	"github.com/spf13/cobra"
)

var (
	installCmd = &cobra.Command{
		Use:  "install",
		RunE: install,
	}
)

func install(cmd *cobra.Command, args []string) error {
	recursive, err := cmd.Flags().GetBool("recursive")
	if err != nil {
		return err
	}

	return helper.Install(recursive)
}

func init() {
	installCmd.Flags().Bool("recursive", true, "run recursively")
	rootCmd.AddCommand(installCmd)
}
