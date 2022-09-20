package cmd

import (
	"fmt"
	"github.com/denormal/go-gitignore"
	slog "github.com/go-eden/slf4go"
	"github.com/ngyewch/asdf-helper/asdf"
	"github.com/ngyewch/asdf-helper/util"
	"github.com/ngyewch/go-clibase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	goVersion "go.hein.dev/go-version"
	"os"
	"path/filepath"
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Use:  AppName,
		RunE: run,
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	ignore, err := gitignore.NewRepository(".")
	if err != nil {
		return err
	}

	helper, err := asdf.NewHelper()
	if err != nil {
		return err
	}
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		match := ignore.Relative(path, info.IsDir())
		if match != nil {
			if match.Ignore() {
				if info.IsDir() {
					return filepath.SkipDir
				} else {
					return nil
				}
			}
		}
		if info.Name() == ".tool-versions" {
			fmt.Println()
			fmt.Println(path)

			pluginMap := make(map[string]string, 0)
			dir := filepath.Dir(path)
			err := util.ScanFile(filepath.Join(dir, ".plugin-versions"), func(line string) error {
				if line == "" {
					return nil
				}
				parts := strings.Split(line, " ")
				if len(parts) == 2 {
					name := parts[0]
					gitUrl := parts[1]
					pluginMap[name] = gitUrl
				}
				return nil
			})
			if err != nil && !os.IsNotExist(err) {
				return err
			}

			err = util.ScanFile(path, func(line string) error {
				if line == "" {
					return nil
				}
				parts := strings.Split(line, " ")
				if len(parts) == 2 {
					name := parts[0]
					version := parts[1]
					hasPlugin, err := helper.CheckPlugin(name)
					if err != nil {
						return err
					}
					if !hasPlugin {
						gitUrl, ok := pluginMap[name]
						if ok {
							err = helper.AddCustomPlugin(name, gitUrl)
							if err != nil {
								return err
							}
						} else {
							err = helper.AddPlugin(name)
							if err != nil {
								return err
							}
						}
					}
					hasInstall, err := helper.CheckInstall(name, version)
					if err != nil {
						return err
					}
					if hasInstall {
						fmt.Printf("* %s %s already installed\n", name, version)
					} else {
						err = helper.Install(name, version)
						if err != nil {
							return err
						}
					}
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
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
