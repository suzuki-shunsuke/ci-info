package controller

import (
	"context"
	"io"
	"os"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ci-info/pkg/domain"
	"github.com/suzuki-shunsuke/ci-info/pkg/github"
)

type GitHub interface {
	GetPR(ctx context.Context, params domain.Params) (*github.PullRequest, error)
	GetPRFiles(ctx context.Context, params github.ParamsGetPRFiles) ([]*github.CommitFile, *github.Response, error)
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
