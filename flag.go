package main

import (
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
		"GitHub repository ID. If not specified, use the current repository.",
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

	return pflag.Args(), nil
}
