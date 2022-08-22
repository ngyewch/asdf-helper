package cmd

import (
	slog "github.com/go-eden/slf4go"
	goVersion "go.hein.dev/go-version"
)

const (
	AppName = "asdf-helper"

	ConfigLogLevel = "log.level"

	DefaultLogLevel = "INFO"
)

var (
	VersionInfo *goVersion.Info
	log         = slog.GetLogger()
)
