package helper

import (
	"fmt"
	"github.com/ngyewch/asdf-helper/asdf"
	"strings"
)

func Latest() error {
	latestVersionMap := make(map[string]string, 0)
	return walk(func(asdfHelper *asdf.Helper, name string, version string) error {
		versionPrefix := ""
		//versionNumber := version
		if version[0] < '0' || version[0] > '9' {
			parts := strings.SplitN(version, "-", 2)
			versionPrefix = parts[0]
			//versionNumber = parts[1]
		}
		key := name
		if versionPrefix != "" {
			key = fmt.Sprintf("%s %s", name, versionPrefix)
		}
		latestVersion, ok := latestVersionMap[key]
		if !ok {
			ver, err := asdfHelper.Latest(name, versionPrefix)
			if err != nil {
				return err
			}
			latestVersion = ver
			latestVersionMap[key] = latestVersion
		}
		if version == latestVersion {
			fmt.Printf("* %s %s (latest)\n", name, version)
		} else {
			fmt.Printf("* %s %s => %s\n", name, version, latestVersion)
		}
		return nil
	})
}
