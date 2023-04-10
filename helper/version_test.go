package helper

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAsdfVersion(t *testing.T) {
	{
		asdfVersion := NewAsdfVersion("adoptopenjdk-19.0.0-beta+36.0.202208190932")
		expectAsdfVersion(t, asdfVersion, true, "adoptopenjdk", "19.0.0-beta+36.0.202208190932", "beta")
	}
	{
		asdfVersion := NewAsdfVersion("adoptopenjdk-19.0.1+10")
		expectAsdfVersion(t, asdfVersion, true, "adoptopenjdk", "19.0.1+10", "")
	}
	{
		asdfVersion := NewAsdfVersion("adoptopenjdk-openj9-8.0.192+12.OpenJDK8U-jdk_x64_linux_openj9_8u192b12.tar.gz")
		expectAsdfVersion(t, asdfVersion, false, "adoptopenjdk-openj9", "8.0.192+12.OpenJDK8U-jdk_x64_linux_openj9_8u192b12.tar.gz", "")
	}
	{
		asdfVersion := NewAsdfVersion("adoptopenjdk-openj9-8.0.192+12.openj9-0.11.0")
		expectAsdfVersion(t, asdfVersion, true, "adoptopenjdk-openj9", "8.0.192+12.openj9-0.11.0", "")
	}
	{
		asdfVersion := NewAsdfVersion("7.27.0-0")
		expectAsdfVersion(t, asdfVersion, true, "", "7.27.0-0", "0")
	}
}

func expectAsdfVersion(t *testing.T, asdfVersion *AsdfVersion, expectedValid bool, expectedVersionPrefix string, expectedVersionNumber string, expectedPrerelease string) {
	assert.Equalf(t, expectedValid, asdfVersion.Valid(), `asdfVersion.VersionPrefix = %v, expected = %v`,
		asdfVersion.Valid(), expectedValid)
	assert.Equalf(t, expectedVersionPrefix, asdfVersion.VersionPrefix,
		`asdfVersion.VersionPrefix = "%s", expected = "%s"`, asdfVersion.VersionPrefix, expectedVersionPrefix)
	assert.Equalf(t, expectedVersionNumber, asdfVersion.VersionNumber,
		`asdfVersion.VersionNumber = "%s", expected = "%s"`, asdfVersion.VersionNumber, expectedVersionNumber)
	if asdfVersion.Valid() {
		assert.Equalf(t, expectedPrerelease, asdfVersion.Prerelease(),
			`asdfVersion.Prerelease() = "%s", expected = "%s"`, asdfVersion.Prerelease(), expectedPrerelease)
	}
}
