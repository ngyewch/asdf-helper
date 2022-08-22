package main

import (
	"github.com/ngyewch/asdf-helper/cmd"
	goVersion "go.hein.dev/go-version"
)

var (
	version string
	commit  string
	date    string
)

func main() {
	versionInfo := goVersion.New(version, commit, date)
	cmd.VersionInfo = versionInfo

	cmd.Execute()
}
