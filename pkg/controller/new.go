package controller

import (
	"context"
	"io"

	"github.com/google/go-github/v32/github"
	gh "github.com/suzuki-shunsuke/ci-info/pkg/github"
)

type Controller struct {
	GitHub GitHub
	Stdout io.Writer
	Stderr io.Writer
}

type GitHub interface {
	GetPR(ctx context.Context, params gh.ParamsGetPR) (*github.PullRequest, *github.Response, error)
	GetPRFiles(ctx context.Context, params gh.ParamsGetPRFiles) ([]*github.CommitFile, *github.Response, error)
	ListPRsWithCommit(ctx context.Context, params gh.ParamsListPRsWithCommit) ([]*github.PullRequest, *github.Response, error)
	GetBranch(ctx context.Context, owner, repo, branch string) (*github.Branch, *github.Response, error)
}
