package cmd

import (
	"github.com/ngyewch/asdf-helper/helper"
	"github.com/spf13/cobra"
)

var (
	latestCmd = &cobra.Command{
		Use:  "latest",
		RunE: latest,
	}
)

func latest(cmd *cobra.Command, args []string) error {
	hideLatest, err := cmd.Flags().GetBool("hide-latest")
	if err != nil {
		return err
	}
	return helper.Latest(hideLatest)
}

func init() {
	latestCmd.Flags().Bool("hide-latest", false, "do not print tools already at the latest version")

	rootCmd.AddCommand(latestCmd)
}
