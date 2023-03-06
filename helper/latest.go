package helper

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/ngyewch/asdf-helper/asdf"
	"regexp"
	"strings"
)

type AsdfVersion struct {
	Version       string
	VersionPrefix string
	VersionNumber string
}

var (
	versionPrefixLocatorRegex1 = regexp.MustCompile(`^[0-9]`)
	versionPrefixLocatorRegex2 = regexp.MustCompile(`-[0-9]`)
)

func toAsdfVersion(version string) *AsdfVersion {
	versionPrefixLocation := versionPrefixLocatorRegex1.FindStringIndex(version)
	if versionPrefixLocation != nil {
		return &AsdfVersion{
			Version:       version,
			VersionPrefix: "",
			VersionNumber: version,
		}
	}
	versionPrefixLocation = versionPrefixLocatorRegex2.FindStringIndex(version)
	if versionPrefixLocation != nil {
		return &AsdfVersion{
			Version:       version,
			VersionPrefix: version[0:versionPrefixLocation[0]],
			VersionNumber: version[versionPrefixLocation[0]+1:],
		}
	}
	return &AsdfVersion{
		Version:       version,
		VersionPrefix: version,
		VersionNumber: "",
	}
}

func getAllVersions(asdfHelper *asdf.Helper, name string, versionPrefix string, excludes []string) ([]*AsdfVersion, error) {
	allVersions := make([]*AsdfVersion, 0)
	candidateVersions, err := asdfHelper.ListAll(name, versionPrefix)
	if err != nil {
		return nil, err
	}
	for _, candidateVersion := range candidateVersions {
		candidateAsdfVersion := toAsdfVersion(candidateVersion)
		if candidateAsdfVersion.VersionPrefix == versionPrefix {
			v, _ := semver.NewVersion(candidateAsdfVersion.VersionNumber)
			if v != nil {
				if v.Prerelease() != "" {
					continue
				}
			} else {
				excludeVersion := false
				for _, exclude := range excludes {
					if strings.Contains(candidateAsdfVersion.VersionNumber, exclude) {
						excludeVersion = true
						break
					}
				}
				if excludeVersion {
					continue
				}
			}
			allVersions = append(allVersions, candidateAsdfVersion)
		}
	}
	return allVersions, nil
}

func Latest(hideLatest bool, excludes []string) error {
	allVersionsMap := make(map[string][]*AsdfVersion, 0)
	return walk(func(asdfHelper *asdf.Helper, name string, version string, constraint string) error {
		var c *semver.Constraints = nil
		if constraint != "" {
			c1, err := semver.NewConstraint(constraint)
			if err != nil {
				fmt.Printf("* %s %s (invalid constraint %s)\n", name, version, constraint)
				return nil
			}
			c = c1
		}

		asdfVersion := toAsdfVersion(version)
		allVersionsMapKey := name
		if asdfVersion.VersionPrefix != "" {
			allVersionsMapKey = fmt.Sprintf("%s %s", name, asdfVersion.VersionPrefix)
		}

		allVersions, ok := allVersionsMap[allVersionsMapKey]
		if !ok {
			allVersions1, err := getAllVersions(asdfHelper, name, asdfVersion.VersionPrefix, excludes)
			if err != nil {
				return err
			}
			allVersionsMap[allVersionsMapKey] = allVersions1
			allVersions = allVersions1
		}

		var latestVersion *AsdfVersion = nil
		if (allVersions != nil) && (len(allVersions) > 0) {
			if c == nil {
				latestVersion = allVersions[len(allVersions)-1]
			} else {
				for _, candidateVersion := range allVersions {
					v, err := semver.NewVersion(candidateVersion.VersionNumber)
					if err != nil {
						//fmt.Printf("! %s -> %s\n", v, err.Error())
					} else {
						if c.Check(v) {
							//fmt.Printf("+ %s (%s)\n", v, c)
							latestVersion = candidateVersion
						} else {
							//fmt.Printf("- %s\n", v)
						}
					}
				}
			}
		}

		if latestVersion == nil {
			fmt.Printf("* %s %s (?)\n", name, version)
		} else if version == latestVersion.Version {
			if !hideLatest {
				if c != nil {
					fmt.Printf("* %s %s (latest %s)\n", name, version, constraint)
				} else {
					fmt.Printf("* %s %s (latest)\n", name, version)
				}
			}
		} else {
			if c != nil {
				fmt.Printf("* %s %s => %s (constraint %s)\n", name, version, latestVersion.Version, constraint)
			} else {
				fmt.Printf("* %s %s => %s\n", name, version, latestVersion.Version)
			}
		}
		return nil
	})
}
