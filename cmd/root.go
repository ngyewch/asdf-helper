package cmd

import (
	"bufio"
	"fmt"
	slog "github.com/go-eden/slf4go"
	"github.com/ngyewch/go-clibase"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	goVersion "go.hein.dev/go-version"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == ".tool-versions" {
			f, err := os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()

			fmt.Println()
			fmt.Println(path)
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := scanner.Text()
				idx := strings.Index(line, "#")
				if idx >= 0 {
					line = line[0:idx]
				}
				space := regexp.MustCompile(`\s+`)
				line = space.ReplaceAllString(line, " ")
				line = strings.TrimSpace(line)
				if line != "" {
					parts := strings.Split(line, " ")
					if len(parts) == 2 {
						name := parts[0]
						version := parts[1]
						hasPlugin, err := checkPlugin(name)
						if err != nil {
							return err
						}
						if hasPlugin {
							hasInstall, err := checkInstall(name, version)
							if err != nil {
								return err
							}
							if hasInstall {
								fmt.Printf("* %s %s already installed\n", name, version)
							} else {

							}
						} else {
							fmt.Printf("* plugin not added: %s\n", name)
						}
					}
				}
			}
			err = scanner.Err()
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

func getAsdfDir() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userHomeDir, ".asdf"), nil
}

func checkPlugin(name string) (bool, error) {
	asdfDir, err := getAsdfDir()
	if err != nil {
		return false, err
	}
	asdfPluginsDir := filepath.Join(asdfDir, "plugins")
	asdfPluginDir := filepath.Join(asdfPluginsDir, name)
	_, err = os.Stat(asdfPluginDir)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func checkInstall(name string, version string) (bool, error) {
	asdfDir, err := getAsdfDir()
	if err != nil {
		return false, err
	}
	asdfInstallsDir := filepath.Join(asdfDir, "installs")
	asdfInstallDir := filepath.Join(asdfInstallsDir, name, version)
	_, err = os.Stat(asdfInstallDir)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func install(name string, version string) error {
	cmd := exec.Command("asdf", "install", name, version)
	combinedOutput, err := cmd.CombinedOutput()
	if (cmd.ProcessState != nil) && (cmd.ProcessState.ExitCode() != 0) {
		fmt.Printf("exit code = %d\n", cmd.ProcessState.ExitCode())
	} else {
		if err != nil {
			return err
		}
	}
	fmt.Println(string(combinedOutput))
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
