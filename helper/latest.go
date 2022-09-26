package helper

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/ngyewch/asdf-helper/asdf"
	"strings"
)

func Latest(hideLatest bool) error {
	latestVersionMap := make(map[string]string, 0)
	return walk(func(asdfHelper *asdf.Helper, name string, version string, constraint string) error {
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

		latestVersion := ""

		if constraint != "" {
			c, err := semver.NewConstraint(constraint)
			if err != nil {
				fmt.Printf("* %s %s (invalid constraint %s)\n", name, version, constraint)
				return nil
			}
			vers, err := asdfHelper.ListAll(name, versionPrefix)
			if err != nil {
				return err
			}
			for _, ver := range vers {
				testVer := ver
				if versionPrefix != "" {
					parts := strings.SplitN(testVer, "-", 2)
					testVer = parts[1]
				}
				v, err := semver.NewVersion(testVer)
				if err != nil {
					//fmt.Printf("! %s -> %s\n", ver, err.Error())
				} else {
					if c.Check(v) {
						//fmt.Printf("+ %s (%s)\n", ver, testVer)
						latestVersion = ver
					} else {
						//fmt.Printf("- %s\n", ver)
					}
				}
			}
			// TODO
		} else {
			cachedLatestVersion, ok := latestVersionMap[key]
			if !ok {
				ver, err := asdfHelper.Latest(name, versionPrefix)
				if err != nil {
					return err
				}
				latestVersion = ver
				latestVersionMap[key] = latestVersion
			} else {
				latestVersion = cachedLatestVersion
			}
		}

		if latestVersion == "" {
			fmt.Printf("* %s %s (?)\n", name, version)
		} else if version == latestVersion {
			if !hideLatest {
				if constraint != "" {
					fmt.Printf("* %s %s (latest %s)\n", name, version, constraint)
				} else {
					fmt.Printf("* %s %s (latest)\n", name, version)
				}
			}
		} else {
			fmt.Printf("* %s %s => %s\n", name, version, latestVersion)
		}
		return nil
	})
}
