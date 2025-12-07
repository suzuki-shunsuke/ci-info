package github

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/go-github/v80/github"
	"golang.org/x/oauth2"
)

type (
	PullRequest       = github.PullRequest
	PullRequestBranch = github.PullRequestBranch
	CommitFile        = github.CommitFile
	Response          = github.Response
	Label             = github.Label
	User              = github.User
)

type Client struct {
	Client *github.Client
}

type ParamsNew struct {
	Token      string
	BaseURL    string
	GraphQLURL string
}

func New(ctx context.Context, params ParamsNew) (Client, error) {
	gh := newGitHub(ctx, params.Token)
	if params.BaseURL != "" {
		gh, err := gh.WithEnterpriseURLs(params.BaseURL, params.BaseURL)
		if err != nil {
			return Client{}, fmt.Errorf("configure GitHub API Base URL: %w", err)
		}
		return Client{
			Client: gh,
		}, nil
	}
	return Client{
		Client: gh,
	}, nil
}

func newGitHub(ctx context.Context, token string) *github.Client {
	return github.NewClient(newHTTP(ctx, token))
}

func newHTTP(ctx context.Context, token string) *http.Client {
	if token == "" {
		return http.DefaultClient
	}
	return oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	))
}

type ParamsGetPR struct {
	Owner string
	Repo  string
	PRNum int
}

type ParamsGetPRFiles struct {
	Owner    string
	Repo     string
	PRNum    int
	FileSize int
}

type paramsListPRsWithCommit struct {
	Owner string
	Repo  string
	SHA   string
}

const maxPerPage = 100

func (c *Client) GetPRFiles(ctx context.Context, logger *slog.Logger, params ParamsGetPRFiles) ([]*github.CommitFile, *github.Response, error) {
	ret := []*github.CommitFile{}
	if params.FileSize == 0 {
		logger.Debug("file size is 0")
		return nil, nil, nil
	}
	n := (params.FileSize / maxPerPage) + 1
	var gResp *github.Response
	for i := 1; i <= n; i++ {
		opts := &github.ListOptions{
			Page:    i,
			PerPage: maxPerPage,
		}
		files, resp, err := c.Client.PullRequests.ListFiles(ctx, params.Owner, params.Repo, params.PRNum, opts)
		if err != nil {
			return files, resp, err //nolint:wrapcheck
		}
		gResp = resp
		ret = append(ret, files...)
		if len(files) != maxPerPage {
			return ret, gResp, nil
		}
	}

	return ret, gResp, nil
}

func (c *Client) ListPRsWithCommit(ctx context.Context, params paramsListPRsWithCommit) ([]*github.PullRequest, *github.Response, error) {
	return c.Client.PullRequests.ListPullRequestsWithCommit(ctx, params.Owner, params.Repo, params.SHA, nil) //nolint:wrapcheck
}
