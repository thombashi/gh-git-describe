package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/phsym/console-slog"
	"github.com/spf13/pflag"

	"github.com/thombashi/eoe"
	"github.com/thombashi/gh-git-describe/pkg/executor"
	"github.com/thombashi/go-gitexec"
)

// const toolName = "gh-git-describe"

// flag variables
var (
	logLevelStr  string
	repoID       string
	cacheDirPath string
	noCache      bool
)

func newLogger(level slog.Level) *slog.Logger {
	logger := slog.New(
		console.NewHandler(os.Stderr, &console.HandlerOptions{
			Level: level,
		}),
	)

	return logger
}

func getCacheTTL() time.Duration {
	if noCache {
		return 0
	}

	return 24 * time.Hour
}

func setFlags() error {
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
		return fmt.Errorf("--repo flag must be specified")
	}

	return nil
}

func main() {
	err := setFlags()
	eoe.ExitOnError(err, eoe.NewParams().WithMessage("failed to set flags"))

	var logLevel slog.Level
	err = logLevel.UnmarshalText([]byte(logLevelStr))
	eoe.ExitOnError(err, eoe.NewParams().WithMessage("failed to get a slog level"))

	logger := newLogger(logLevel)
	eoeParams := eoe.NewParams().WithLogger(logger)

	gitExecutor, err := gitexec.New(&gitexec.Params{
		Logger: logger,
		Env:    os.Environ(),
	})
	eoe.ExitOnError(err, eoe.NewParams().WithMessage("failed to create a GitExecutor instance"))

	params := &executor.Params{
		GitExecutor:  gitExecutor,
		Logger:       logger,
		CacheDirPath: cacheDirPath,
		CacheTTL:     getCacheTTL(),
	}
	gdExecutor, err := executor.New(params)
	eoe.ExitOnError(err, eoeParams.WithMessage("failed to create a git instance"))

	out, err := gdExecutor.RunGitDescribe(&executor.RepoCloneParams{
		RepoID: repoID,
	}, pflag.Args()...)
	eoe.ExitOnError(err, eoeParams.WithMessage("failed to run git describe"))

	fmt.Println(out)
}
