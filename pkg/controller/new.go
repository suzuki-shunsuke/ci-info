package controller

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/google/go-github/v39/github"
	"github.com/sirupsen/logrus"
	gh "github.com/suzuki-shunsuke/ci-info/pkg/github"
	"github.com/suzuki-shunsuke/go-ci-env/v2/cienv"
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
	GetPR(ctx context.Context, params gh.ParamsGetPR) (*github.PullRequest, *github.Response, error)
	GetPRFiles(ctx context.Context, params gh.ParamsGetPRFiles) ([]*github.CommitFile, *github.Response, error)
	ListPRsWithCommit(ctx context.Context, params gh.ParamsListPRsWithCommit) ([]*github.PullRequest, *github.Response, error)
}

type Controller struct {
	GitHub GitHub
	Stdout io.Writer
	Stderr io.Writer
}

func New(ctx context.Context, params Params) (Controller, Params, error) { //nolint:cyclop
	if params.LogLevel != "" {
		lvl, err := logrus.ParseLevel(params.LogLevel)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"log_level": params.LogLevel,
			}).WithError(err).Error("the log level is invalid")
		}
		logrus.SetLevel(lvl)
	}

	if params.GitHubToken == "" {
		params.GitHubToken = os.Getenv("GITHUB_TOKEN")
		if params.GitHubToken == "" {
			params.GitHubToken = os.Getenv("GITHUB_ACCESS_TOKEN")
		}
	}

	//nolint:nestif
	if platform := cienv.Get(); platform != nil {
		if params.Owner == "" {
			params.Owner = platform.RepoOwner()
		}
		if params.Repo == "" {
			params.Repo = platform.RepoName()
		}
		if params.SHA == "" {
			params.SHA = platform.SHA()
		}
		if params.PRNum <= 0 {
			prNum, err := platform.PRNumber()
			if err != nil {
				return Controller{}, params, fmt.Errorf("get the pull request number: %w", err)
			}
			params.PRNum = prNum
		}
	}

	ghClient := gh.New(ctx, gh.ParamsNew{
		Token: params.GitHubToken,
	})

	return Controller{
		GitHub: &ghClient,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}, params, nil
}
