package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/phsym/console-slog"

	"github.com/thombashi/eoe"
	"github.com/thombashi/gh-git-describe/pkg/executor"
	"github.com/thombashi/go-gitexec"
)

func newLogger(level slog.Level) *slog.Logger {
	logger := slog.New(
		console.NewHandler(os.Stderr, &console.HandlerOptions{
			Level: level,
		}),
	)

	return logger
}

func getCacheTTL(noCache bool) time.Duration {
	if noCache {
		return 0
	}

	return 24 * time.Hour
}

func main() {
	flags, args, err := setFlags()
	eoe.ExitOnError(err, eoe.NewParams().WithMessage("failed to set flags"))

	var logLevel slog.Level
	err = logLevel.UnmarshalText([]byte(flags.LogLevelStr))
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
		CacheDirPath: flags.CacheDirPath,
		CacheTTL:     getCacheTTL(flags.NoCache),
		LogWithPackage: true,
	}
	gdExecutor, err := executor.New(params)
	eoe.ExitOnError(err, eoeParams.WithMessage("failed to create a git instance"))

	out, err := gdExecutor.RunGitDescribe(&executor.RepoCloneParams{
		RepoID: flags.RepoID,
	}, args...)
	eoe.ExitOnError(err, eoeParams.WithMessage("failed to run git describe"))

	fmt.Println(out)
}
