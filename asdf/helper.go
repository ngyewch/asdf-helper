package asdf

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type Helper struct {
	asdfDir         string
	asdfPluginsDir  string
	asdfInstallsDir string
}

func NewHelper() (*Helper, error) {
	asdfDir, err := GetAsdfDir()
	if err != nil {
		return nil, err
	}
	return &Helper{
		asdfDir:         asdfDir,
		asdfPluginsDir:  filepath.Join(asdfDir, "plugins"),
		asdfInstallsDir: filepath.Join(asdfDir, "installs"),
	}, nil
}

func GetAsdfDir() (string, error) {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(userHomeDir, ".asdf"), nil
}

func (helper *Helper) CheckPlugin(name string) (bool, error) {
	asdfPluginDir := filepath.Join(helper.asdfPluginsDir, name)
	_, err := os.Stat(asdfPluginDir)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (helper *Helper) CheckInstall(name string, version string) (bool, error) {
	asdfInstallDir := filepath.Join(helper.asdfInstallsDir, name, version)
	_, err := os.Stat(asdfInstallDir)
	if os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else {
		return true, nil
	}
}

func (helper *Helper) AddPlugin(name string) error {
	cmd := exec.Command("asdf", "plugin", "add", name)
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

func (helper *Helper) AddCustomPlugin(name string, gitUrl string) error {
	cmd := exec.Command("asdf", "plugin", "add", name, gitUrl)
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

func (helper *Helper) Install(name string, version string) error {
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

func (helper *Helper) Latest(name string, prefix string) (string, error) {
	cmd := exec.Command("asdf", "latest", name, prefix)
	combinedOutput, err := cmd.CombinedOutput()
	if (cmd.ProcessState != nil) && (cmd.ProcessState.ExitCode() != 0) {
		fmt.Printf("exit code = %d\n", cmd.ProcessState.ExitCode())
	} else {
		if err != nil {
			return "", err
		}
	}
	return strings.TrimSpace(string(combinedOutput)), nil
}

func (helper *Helper) ListAll(name string, prefix string) ([]string, error) {
	cmd := exec.Command("asdf", "list", "all", name, prefix)
	combinedOutput, err := cmd.CombinedOutput()
	if (cmd.ProcessState != nil) && (cmd.ProcessState.ExitCode() != 0) {
		fmt.Printf("exit code = %d\n", cmd.ProcessState.ExitCode())
	} else {
		if err != nil {
			return nil, err
		}
	}
	re := regexp.MustCompile(`\r?\n`)
	return re.Split(strings.TrimSpace(string(combinedOutput)), -1), nil
}
