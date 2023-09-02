package helper

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
	"github.com/ngyewch/asdf-helper/asdf"
)

func getAllVersions(asdfHelper *asdf.Helper, name string, versionPrefix string, includePrereleases bool) ([]*AsdfVersion, error) {
	candidateVersions, err := asdfHelper.ListAll(name, versionPrefix)
	if err != nil {
		return nil, err
	}
	allVersions := make([]*AsdfVersion, 0)
	for _, candidateVersion := range candidateVersions {
		candidateAsdfVersion := NewAsdfVersion(candidateVersion)
		if (candidateAsdfVersion.VersionPrefix == versionPrefix) && candidateAsdfVersion.Valid() {
			if includePrereleases || (candidateAsdfVersion.SemVer.Prerelease() == "") {
				allVersions = append(allVersions, candidateAsdfVersion)
			}
		}
	}
	return allVersions, nil
}

func Latest(hideLatest bool, includePrereleases bool, recursive bool) error {
	allVersionsMap := make(map[string][]*AsdfVersion, 0)
	return walk(recursive, func(asdfHelper *asdf.Helper, name string, version string, constraint string) error {
		var c *semver.Constraints = nil
		if constraint != "" {
			c1, err := semver.NewConstraint(constraint)
			if err != nil {
				fmt.Printf("* %s %s (invalid constraint %s)\n", name, version, constraint)
				return nil
			}
			c = c1
		}

		asdfVersion := NewAsdfVersion(version)
		allVersionsMapKey := name
		if asdfVersion.VersionPrefix != "" {
			allVersionsMapKey = fmt.Sprintf("%s %s", name, asdfVersion.VersionPrefix)
		}

		allVersions, ok := allVersionsMap[allVersionsMapKey]
		if !ok {
			allVersions1, err := getAllVersions(asdfHelper, name, asdfVersion.VersionPrefix, includePrereleases)
			if err != nil {
				return err
			}
			allVersionsMap[allVersionsMapKey] = allVersions1
			allVersions = allVersions1
		}

		var latestVersion *AsdfVersion = nil
		for _, candidateVersion := range allVersions {
			if candidateVersion.Valid() {
				if candidateVersion.CheckConstraints(c) && ((latestVersion == nil) || candidateVersion.SemVer.GreaterThan(latestVersion.SemVer)) {
					//fmt.Printf("[%s] + %s (%s)\n", name, candidateVersion.SemVer, c)
					latestVersion = candidateVersion
				} else {
					//fmt.Printf("[%s] - %s (%s)\n", name, candidateVersion.SemVer, c)
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
