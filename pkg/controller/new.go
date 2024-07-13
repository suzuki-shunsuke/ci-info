package controller

import (
	"context"
	"io"
	"os"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ci-info/pkg/github"
)

type Params struct {
	Owner       string
	Repo        string
	SHA         string
	Dir         string
	PRNum       int
	GitHubToken string
	LogLevel    string
	Prefix      string
}

type GitHub interface {
	GetPR(ctx context.Context, params github.ParamsGetPR) (*github.PullRequest, *github.Response, error)
	GetPRFiles(ctx context.Context, params github.ParamsGetPRFiles) ([]*github.CommitFile, *github.Response, error)
	ListPRsWithCommit(ctx context.Context, params github.ParamsListPRsWithCommit) ([]*github.PullRequest, *github.Response, error)
}

type Controller struct {
	gh     GitHub
	stdout io.Writer
	stderr io.Writer
	fs     afero.Fs
}

func New(ghClient github.Client, fs afero.Fs) Controller {
	return Controller{
		gh:     &ghClient,
		fs:     fs,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}
