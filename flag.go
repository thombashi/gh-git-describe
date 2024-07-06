package main

import (
	"fmt"

	"github.com/spf13/pflag"
)

type Flags struct {
	RepoID       string
	LogLevelStr  string
	CacheDirPath string
	NoCache      bool
}

var flags = Flags{}

func setFlags() ([]string, error) {
	pflag.StringVarP(
		&flags.RepoID,
		"repo",
		"R",
		"",
		"[required] GitHub repository ID",
	)
	pflag.StringVar(
		&flags.LogLevelStr,
		"log-level",
		"info",
		"log level (debug, info, warn, error)",
	)
	pflag.StringVar(
		&flags.CacheDirPath,
		"cache-dir",
		"",
		"cache directory path. If not specified, use the system's temporary directory.",
	)
	pflag.BoolVar(
		&flags.NoCache,
		"no-cache",
		false,
		"disable cache",
	)
	pflag.Parse()

	if flags.RepoID == "" {
		return nil, fmt.Errorf("--repo flag must be specified")
	}

	return pflag.Args(), nil
}
