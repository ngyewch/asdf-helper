package helper

import (
	"fmt"
	"github.com/ngyewch/asdf-helper/asdf"
)

func Install() error {
	return walk(func(asdfHelper *asdf.Helper, name string, version string) error {
		hasInstall, err := asdfHelper.CheckInstall(name, version)
		if err != nil {
			return err
		}
		if hasInstall {
			fmt.Printf("* %s %s already installed\n", name, version)
		} else {
			err = asdfHelper.Install(name, version)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
