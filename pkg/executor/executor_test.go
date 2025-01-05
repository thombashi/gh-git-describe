package executor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thombashi/go-gitexec"
)

func TestRunGitDescribe(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	ctx := context.Background()

	gitExecutor, err := gitexec.New(&gitexec.Params{})
	r.NoError(err)

	params := NewParams()
	params.GitExecutor = gitExecutor
	params.CacheTTL = 60 * time.Second
	executor, err := New(params)
	r.NoError(err)

	const (
		want = "v4.1.7"
		sha  = "692973e3d937129bcbf40652eb9f2f61becf3332"
	)

	rcParams := &RepoCloneParams{
		RepoID:   "actions/checkout",
		CacheTTL: 300,
	}
	got, err := executor.RunGitDescribe(rcParams, "--tags", sha)
	r.NoError(err)
	a.Equal(want, got)

	got, err = executor.RunGitDescribe(rcParams, "--tags", sha)
	r.NoError(err)
	a.Equal(want, got)

	got, err = executor.RunGitDescribeContext(ctx, rcParams, "--tags", sha)
	r.NoError(err)
	a.Equal(want, got)
}

func TestRunGitDescribeInvalidSHA(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	gitExecutor, err := gitexec.New(&gitexec.Params{})
	r.NoError(err)

	params := NewParams()
	params.GitExecutor = gitExecutor
	params.CacheTTL = 60 * time.Second
	executor, err := New(params)
	r.NoError(err)

	const invalidSHA = "0123456789abcdef0123456789abcdef01234567"

	rcParams := &RepoCloneParams{
		RepoID:   "actions/checkout",
		CacheTTL: 300,
	}
	_, err = executor.RunGitDescribe(rcParams, "--tags", invalidSHA)
	r.Error(err)

	_, err = executor.RunGitDescribeContext(ctx, rcParams, "--tags", invalidSHA)
	r.Error(err)
}

func TestRunGitRevParse(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	ctx := context.Background()

	gitExecutor, err := gitexec.New(&gitexec.Params{})
	r.NoError(err)

	params := NewParams()
	params.GitExecutor = gitExecutor
	params.CacheTTL = 30 * time.Second
	executor, err := New(params)
	r.NoError(err)

	want := "692973e3d937129bcbf40652eb9f2f61becf3332"
	tag := "v4.1.7"

	rcParams := &RepoCloneParams{
		RepoID:   "actions/checkout",
		CacheTTL: 300,
	}
	got, err := executor.RunGitRevParse(rcParams, tag)
	r.NoError(err)
	a.Equal(want, got)

	got, err = executor.RunGitRevParseContext(ctx, rcParams, tag)
	r.NoError(err)
	a.Equal(want, got)
}

func TestRunGitRevList(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)
	ctx := context.Background()

	gitExecutor, err := gitexec.New(&gitexec.Params{})
	r.NoError(err)

	params := NewParams()
	params.GitExecutor = gitExecutor
	params.CacheTTL = 30 * time.Second
	executor, err := New(params)
	r.NoError(err)

	const (
		want = "0b496e91ec7ae4428c3ed2eeb4c3a40df431f2cc"
		tag  = "v1.1.0"
	)

	rcParams := &RepoCloneParams{
		RepoID:   "actions/checkout",
		CacheTTL: 300,
	}
	got, err := executor.RunGitRevList(rcParams, "-n", "1", tag)
	r.NoError(err)
	a.Equal(want, got)

	got, err = executor.RunGitRevListContext(ctx, rcParams, "-n", "1", tag)
	r.NoError(err)
	a.Equal(want, got)
}

func TestRunGitRevListInvalidSHA(t *testing.T) {
	r := require.New(t)
	ctx := context.Background()

	gitExecutor, err := gitexec.New(&gitexec.Params{})
	r.NoError(err)

	params := NewParams()
	params.GitExecutor = gitExecutor
	params.CacheTTL = 60 * time.Second
	executor, err := New(params)
	r.NoError(err)

	const invalidSHA = "0123456789abcdef0123456789abcdef01234567"

	rcParams := &RepoCloneParams{
		RepoID:   "actions/checkout",
		CacheTTL: 300,
	}
	_, err = executor.RunGitRevList(rcParams, invalidSHA)
	r.Error(err)

	_, err = executor.RunGitRevListContext(ctx, rcParams, invalidSHA)
	r.Error(err)
}
