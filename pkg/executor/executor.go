package executor

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cli/go-gh/v2"
	"github.com/cli/go-gh/v2/pkg/repository"

	"github.com/thombashi/go-gitexec"
)

const extensionName = "gh-git-describe"

// RepoCloneParams represents the parameters for the RunRepoClone function.
type RepoCloneParams struct {
	// RepoID is the repository ID (OWNER/NAME) of the repository to clone.
	RepoID string

	// CacheTTL is the cache TTL duration.
	CacheTTL time.Duration
}

// Executor is an interface for the command executors.
type Executor interface {
	// RunRepoClone clones the specified GitHub repository.
	RunRepoClone(params *RepoCloneParams) (string, error)

	// RunRepoCloneContext clones the specified GitHub repository with the specified context.
	RunRepoCloneContext(ctx context.Context, params *RepoCloneParams) (string, error)

	// RunGitDescribe runs the 'git describe' command for the specified GitHub repository.
	RunGitDescribe(params *RepoCloneParams, args ...string) (string, error)

	// RunGitDescribe runs the 'git describe' command for the specified GitHub repository with the specified context.
	RunGitDescribeContext(ctx context.Context, params *RepoCloneParams, args ...string) (string, error)
}

type executor struct {
	gitExecutor  gitexec.GitExecutor
	logger       *slog.Logger
	cacheDirPath string
	cacheTTL     time.Duration
}

type Params struct {
	// GitExecutor is an instance of gitexec.GitExecutor
	GitExecutor gitexec.GitExecutor

	// Logger is the slog.Logger instance.
	Logger *slog.Logger

	// CacheDirPath is the directory path to store the cache.
	CacheDirPath string

	// CacheTTL is the cache retention duration for the cloned repository.
	CacheTTL time.Duration
}

func NewParams() *Params {
	return &Params{
		CacheTTL: 1 * time.Hour,
	}
}

// WithLogger sets the logger instance.
func (p *Params) WithLogger(logger *slog.Logger) *Params {
	p.Logger = logger
	return p
}

// New creates a new Executor instance.
func New(params *Params) (Executor, error) {
	var err error

	gitExecutor := params.GitExecutor
	if gitExecutor == nil {
		gitExecutor, err = gitexec.New(&gitexec.Params{
			Logger: params.Logger,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create a GitExecutor instance: %w", err)
		}
	}

	logger := params.Logger
	if logger == nil {
		logger = slog.Default()
	}

	var cacheDirPath string
	if params.CacheDirPath == "" {
		userCacheDir, err := os.UserCacheDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get the user cache directory: %w", err)
		}

		cacheDirPath = filepath.Join(userCacheDir, extensionName)
	} else {
		cacheDirPath = params.CacheDirPath
	}

	if err := os.MkdirAll(cacheDirPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create a cache directory: %w", err)
	}

	return &executor{
		gitExecutor:  gitExecutor,
		logger:       logger,
		cacheDirPath: cacheDirPath,
		cacheTTL:     params.CacheTTL,
	}, nil
}

// RunRepoClone clones the specified GitHub repository.
func (e executor) RunRepoClone(params *RepoCloneParams) (string, error) {
	return e.RunRepoCloneContext(context.Background(), params)
}

// RunRepoCloneContext clones the specified GitHub repository with the specified context.
func (e executor) RunRepoCloneContext(ctx context.Context, params *RepoCloneParams) (string, error) {
	logger := e.logger.With(slog.String("repo", params.RepoID))

	if params.RepoID == "" {
		return "", fmt.Errorf("repository ID must be specified")
	}

	repo, err := repository.Parse(params.RepoID)
	if err != nil {
		return "", fmt.Errorf("failed to parse the repository ID: %w", err)
	}

	outputDir := filepath.Join(e.cacheDirPath, repo.Owner, repo.Name)
	logger.Debug("cloning a repository", slog.String("path", outputDir))

	// find the cache directory
	plock.RLock(outputDir)
	info, err := os.Stat(outputDir)
	plock.RUnlock(outputDir)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to get the information of the directory: %w", err)
	}
	if err == nil {
		cacheTTL := e.cacheTTL
		if params.CacheTTL > 0 {
			cacheTTL = params.CacheTTL
		}

		if time.Since(info.ModTime()) < cacheTTL {
			logger.Debug("repo cache found", slog.String("path", outputDir))
			return outputDir, nil
		}
	}

	tempDir, err := os.MkdirTemp("", extensionName)
	if err != nil {
		return "", fmt.Errorf("failed to create a temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	if err := os.MkdirAll(filepath.Dir(outputDir), 0755); err != nil {
		return "", fmt.Errorf("failed to create the parent directory: %w", err)
	}

	repoID := fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
	_, _, err = gh.ExecContext(ctx, "repo", "clone", repoID, tempDir, "--", "--bare")
	if err != nil {
		return "", fmt.Errorf("failed to clone the repository: %w", err)
	}

	plock.Lock(outputDir)
	defer plock.Unlock(outputDir)

	if err := os.RemoveAll(outputDir); err != nil {
		return "", fmt.Errorf("failed to remove the directory: %w", err)
	}
	if err := os.Rename(tempDir, outputDir); err != nil {
		return "", fmt.Errorf("failed to rename the directory: %w", err)
	}

	return outputDir, nil
}

// RunGitDescribe runs the 'git describe' command for the specified GitHub repository.
func (e executor) RunGitDescribe(params *RepoCloneParams, args ...string) (string, error) {
	return e.RunGitDescribeContext(context.Background(), params, args...)
}

// RunGitDescribeContext runs the 'git describe' command for the specified GitHub repository with the specified context.
func (e executor) RunGitDescribeContext(ctx context.Context, params *RepoCloneParams, args ...string) (string, error) {
	clonedDir, err := e.RunRepoCloneContext(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to clone the repository: %w", err)
	}

	plock.RLock(clonedDir)
	defer plock.RUnlock(clonedDir)

	gitArgs := []string{"-C", clonedDir, "describe"}
	gitArgs = append(gitArgs, args...)
	result, err := e.gitExecutor.RunGitContext(ctx, gitArgs...)
	if err != nil {
		return "", fmt.Errorf("failed to run git: error=%w, stderr=%s", err, result.Stderr.String())
	}

	stdout := strings.TrimSpace(result.Stdout.String())

	return stdout, nil
}
