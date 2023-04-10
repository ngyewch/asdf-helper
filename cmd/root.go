package cmd

import (
	"github.com/ngyewch/go-clibase"
	"github.com/spf13/cobra"
	goVersion "go.hein.dev/go-version"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:  AppName,
		RunE: help,
	}
)

func help(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	clibase.AddVersionCmd(rootCmd, func() *goVersion.Info {
		return VersionInfo
	})
}

func initConfig() {
	// do nothing
}
