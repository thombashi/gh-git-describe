package executor

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/thombashi/go-gitexec"
)

func TestRunGitDescribe(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	gitExecutor, err := gitexec.New(&gitexec.Params{})
	r.NoError(err)

	params := NewParams()
	params.GitExecutor = gitExecutor
	params.CacheTTL = 0 * time.Hour
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
}

func TestRunGitDescribeInvalidSHA(t *testing.T) {
	a := assert.New(t)
	r := require.New(t)

	gitExecutor, err := gitexec.New(&gitexec.Params{})
	r.NoError(err)

	params := NewParams()
	params.GitExecutor = gitExecutor
	params.CacheTTL = 0 * time.Hour
	executor, err := New(params)
	r.NoError(err)

	sha := "0123456789abcdef0123456789abcdef01234567"

	rcParams := &RepoCloneParams{
		RepoID:   "actions/checkout",
		CacheTTL: 300,
	}
	_, err = executor.RunGitDescribe(rcParams, "--tags", sha)
	a.Error(err)
}
