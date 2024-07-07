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

func toRepoID(repo repository.Repository) string {
	return fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
}

func makeCacheDir(cacheDirPath string) (string, error) {
	cacheDirPath = strings.TrimSpace(cacheDirPath)

	if cacheDirPath == "" {
		userCacheDir, err := os.UserCacheDir()
		if err != nil {
			return "", fmt.Errorf("failed to get the user cache directory: %w", err)
		}

		cacheDirPath = filepath.Join(userCacheDir, extensionName)
	} else {
		cacheDirPath = filepath.Clean(cacheDirPath)
		cacheDirPath = filepath.Join(cacheDirPath, extensionName)
	}

	if err := os.MkdirAll(cacheDirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create a cache directory: %w", err)
	}

	return cacheDirPath, nil
}

// RepoCloneParams represents the parameters for the RunRepoClone function.
type RepoCloneParams struct {
	// RepoID is the repository ID (OWNER/NAME) of the repository to clone.
	// If not specified, the current repository is used.
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

	// RunGit runs the specified git command.
	RunGit(params *RepoCloneParams, command string, args ...string) (string, error)

	// RunGitContext runs the specified git command with the specified context.
	RunGitContext(ctx context.Context, params *RepoCloneParams, command string, args ...string) (string, error)

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

// Params represents the parameters for the New function.
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

// NewParams creates a new Params instance.
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

	cacheDirPath, err := makeCacheDir(params.CacheDirPath)
	if err != nil {
		return nil, err
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
	var repo repository.Repository
	var err error

	if params.RepoID == "" {
		repo, err = repository.Current()
		if err != nil {
			return "", fmt.Errorf("failed to get the current repository: %w", err)
		}
	} else {
		repo, err = repository.Parse(params.RepoID)
		if err != nil {
			return "", fmt.Errorf("failed to parse the repository ID: %w", err)
		}
	}

	repoID := toRepoID(repo)
	logger := e.logger.With(slog.String("repo", repoID))

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

	// reaching here means that the cache is not found or expired.
	// clone the GitHub repository and replace the local cache directory.

	tempDir, err := os.MkdirTemp("", extensionName)
	if err != nil {
		return "", fmt.Errorf("failed to create a temporary directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	if err := os.MkdirAll(filepath.Dir(outputDir), 0755); err != nil {
		return "", fmt.Errorf("failed to create the parent directory: %w", err)
	}

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

// RunGit runs the specified git command.
func (e executor) RunGit(params *RepoCloneParams, command string, args ...string) (string, error) {
	return e.RunGitContext(context.Background(), params, command, args...)
}

// RunGitContext runs the specified git command with the specified context.
func (e executor) RunGitContext(ctx context.Context, params *RepoCloneParams, command string, args ...string) (string, error) {
	command = strings.TrimSpace(command)
	if command == "" {
		return "", fmt.Errorf("require a git subcommand")
	}

	clonedDir, err := e.RunRepoCloneContext(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to clone the repository: %w", err)
	}

	plock.RLock(clonedDir)
	defer plock.RUnlock(clonedDir)

	gitArgs := []string{"-C", clonedDir, command}
	gitArgs = append(gitArgs, args...)
	result, err := e.gitExecutor.RunGitContext(ctx, gitArgs...)
	if err != nil {
		return "", fmt.Errorf("failed to run git: error=%w, stderr=%s", err, result.Stderr.String())
	}

	stdout := strings.TrimSpace(result.Stdout.String())

	return stdout, nil
}

// RunGitDescribe runs the 'git describe' command for the specified GitHub repository.
func (e executor) RunGitDescribe(params *RepoCloneParams, args ...string) (string, error) {
	return e.RunGitDescribeContext(context.Background(), params, args...)
}

// RunGitDescribeContext runs the 'git describe' command for the specified GitHub repository with the specified context.
func (e executor) RunGitDescribeContext(ctx context.Context, params *RepoCloneParams, args ...string) (string, error) {
	subcommand := "describe"
	stdout, err := e.RunGitContext(ctx, params, subcommand, args...)
	if err != nil {
		return "", fmt.Errorf("failed to run git-%s: %w", subcommand, err)
	}

	return stdout, nil
}

