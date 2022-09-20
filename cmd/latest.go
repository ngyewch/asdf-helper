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
	return helper.Latest()
}

func init() {
	rootCmd.AddCommand(latestCmd)
}
