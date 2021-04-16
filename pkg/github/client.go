package github

import (
	"context"

	"github.com/google/go-github/v35/github"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type Client struct {
	Client *github.Client
}

type ParamsNew struct {
	Token string
}

func New(ctx context.Context, params ParamsNew) Client {
	tc := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: params.Token},
	))
	return Client{
		Client: github.NewClient(tc),
	}
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

type ParamsListPRsWithCommit struct {
	Owner string
	Repo  string
	SHA   string
}

func (client *Client) GetPR(ctx context.Context, params ParamsGetPR) (*github.PullRequest, *github.Response, error) {
	return client.Client.PullRequests.Get(ctx, params.Owner, params.Repo, params.PRNum)
}

func (client *Client) getPRFiles(ctx context.Context, params ParamsGetPRFiles, opts *github.ListOptions) ([]*github.CommitFile, *github.Response, error) {
	return client.Client.PullRequests.ListFiles(ctx, params.Owner, params.Repo, params.PRNum, opts)
}

const maxPerPage = 100

func (client *Client) GetPRFiles(ctx context.Context, params ParamsGetPRFiles) ([]*github.CommitFile, *github.Response, error) {
	ret := []*github.CommitFile{}
	if params.FileSize == 0 {
		logrus.Debug("file size is 0")
		return nil, nil, nil
	}
	n := (params.FileSize / maxPerPage) + 1
	var gResp *github.Response
	for i := 1; i <= n; i++ {
		opts := &github.ListOptions{
			Page:    i,
			PerPage: maxPerPage,
		}
		files, resp, err := client.getPRFiles(ctx, params, opts)
		if err != nil {
			return files, resp, err
		}
		gResp = resp
		ret = append(ret, files...)
		if len(files) != maxPerPage {
			return ret, gResp, nil
		}
	}

	return ret, gResp, nil
}

func (client *Client) ListPRsWithCommit(ctx context.Context, params ParamsListPRsWithCommit) ([]*github.PullRequest, *github.Response, error) {
	return client.Client.PullRequests.ListPullRequestsWithCommit(ctx, params.Owner, params.Repo, params.SHA, nil)
}
