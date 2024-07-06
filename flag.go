package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

// flag variables
var (
	logLevelStr  string
	repoID       string
	cacheDirPath string
	noCache      bool
)

func setFlags() ([]string, error) {
	pflag.StringVarP(
		&repoID,
		"repo",
		"R",
		"",
		"[required] GitHub repository ID",
	)
	pflag.StringVar(
		&logLevelStr,
		"log-level",
		"info",
		"log level (debug, info, warn, error)",
	)
	pflag.StringVar(
		&cacheDirPath,
		"cache-dir",
		"",
		"cache directory path. If not specified, use the system's temporary directory.",
	)
	pflag.BoolVar(
		&noCache,
		"no-cache",
		false,
		"disable cache",
	)
	pflag.Parse()

	if repoID == "" {
		return nil, fmt.Errorf("--repo flag must be specified")
	}

	return pflag.Args(), nil
}
