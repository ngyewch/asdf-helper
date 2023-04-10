package helper

import (
	"github.com/Masterminds/semver/v3"
	"regexp"
)

var (
	versionPrefixLocatorRegex1 = regexp.MustCompile(`^[0-9]`)
	versionPrefixLocatorRegex2 = regexp.MustCompile(`-[0-9]`)
)

type AsdfVersion struct {
	Version       string
	VersionPrefix string
	VersionNumber string
	SemVer        *semver.Version
}

func NewAsdfVersion(version string) *AsdfVersion {
	versionPrefixLocation := versionPrefixLocatorRegex1.FindStringIndex(version)
	if versionPrefixLocation != nil {
		v, _ := semver.NewVersion(version)
		return &AsdfVersion{
			Version:       version,
			VersionPrefix: "",
			VersionNumber: version,
			SemVer:        v,
		}
	}
	versionPrefixLocation = versionPrefixLocatorRegex2.FindStringIndex(version)
	if versionPrefixLocation != nil {
		versionPrefix := version[0:versionPrefixLocation[0]]
		versionNumber := version[versionPrefixLocation[0]+1:]
		v, _ := semver.NewVersion(versionNumber)
		return &AsdfVersion{
			Version:       version,
			VersionPrefix: versionPrefix,
			VersionNumber: versionNumber,
			SemVer:        v,
		}
	}
	return &AsdfVersion{
		Version:       version,
		VersionPrefix: version,
		VersionNumber: "",
		SemVer:        nil,
	}
}

func (v *AsdfVersion) Valid() bool {
	return v.SemVer != nil
}

func (v *AsdfVersion) Prerelease() string {
	return v.SemVer.Prerelease()
}

func (v *AsdfVersion) CheckConstraints(c *semver.Constraints) bool {
	if v.SemVer == nil {
		return false
	}
	if c != nil {
		return c.Check(v.SemVer)
	} else {
		return true
	}
}
