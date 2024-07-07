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

	want := "v4"
	sha := "692973e3d937129bcbf40652eb9f2f61becf3332"

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

	sha := "0123456789abcdef0123456789abcdef01234567"

	rcParams := &RepoCloneParams{
		RepoID:   "actions/checkout",
		CacheTTL: 300,
	}
	_, err = executor.RunGitDescribe(rcParams, "--tags", sha)
	r.Error(err)

	_, err = executor.RunGitDescribeContext(ctx, rcParams, "--tags", sha)
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
