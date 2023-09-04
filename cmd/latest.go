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

	includePrereleases, err := cmd.Flags().GetBool("include-prereleases")
	if err != nil {
		return err
	}

	recursive, err := cmd.Flags().GetBool("recursive")
	if err != nil {
		return err
	}

	return helper.Latest(hideLatest, includePrereleases, recursive)
}

func init() {
	latestCmd.Flags().Bool("hide-latest", false, "do not print tools already at the latest version")
	latestCmd.Flags().Bool("include-prereleases", false, "include prereleases")
	latestCmd.Flags().Bool("recursive", true, "run recursively")

	rootCmd.AddCommand(latestCmd)
}
