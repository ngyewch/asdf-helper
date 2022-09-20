package cmd

import (
	slog "github.com/go-eden/slf4go"
	"github.com/ngyewch/go-clibase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	rootCmd.PersistentFlags().String("config", "", "config file")
	rootCmd.PersistentFlags().String("log-level", DefaultLogLevel, "log level")

	_ = viper.BindPFlag(ConfigLogLevel, rootCmd.Flag("log-level"))

	clibase.AddVersionCmd(rootCmd, func() *goVersion.Info {
		return VersionInfo
	})
}

func initConfig() {
	cfgFile, err := rootCmd.PersistentFlags().GetString("config")
	if cfgFile != "" && err == nil {
		viper.SetConfigFile(cfgFile)
		err := viper.ReadInConfig()
		if err != nil {
			log.Warnf("error reading config file: %s", err)
		}
	}

	logLevel := viper.GetString(ConfigLogLevel)
	slog.SetLevel(clibase.ToSlogLevel(logLevel, slog.InfoLevel))
}
