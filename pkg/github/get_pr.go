package github

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ci-info/pkg/domain"
)

func (c *Client) getPRNum(ctx context.Context, params domain.Params) (int, error) {
	if params.PRNum > 0 {
		return params.PRNum, nil
	}
	logrus.WithFields(logrus.Fields{
		"owner": params.Owner,
		"repo":  params.Repo,
		"sha":   params.SHA,
	}).Debug("get pull request from SHA")
	prs, _, err := c.ListPRsWithCommit(ctx, paramsListPRsWithCommit{
		Owner: params.Owner,
		Repo:  params.Repo,
		SHA:   params.SHA,
	})
	if err != nil {
		return 0, fmt.Errorf("list pull requests with a commit: %w", err)
	}
	logrus.WithFields(logrus.Fields{
		"size": len(prs),
	}).Debug("the number of pull requests assosicated with the commit")
	if len(prs) == 0 {
		return 0, nil
	}
	return prs[0].GetNumber(), nil
}

func (c *Client) GetPR(ctx context.Context, params domain.Params) (*PullRequest, error) {
	prNum, err := c.getPRNum(ctx, params)
	if err != nil {
		return nil, err
	}
	if prNum <= 0 {
		return nil, nil //nolint:nilnil
	}
	pr, _, err := c.Client.PullRequests.Get(ctx, params.Owner, params.Repo, params.PRNum)
	if err != nil {
		return nil, fmt.Errorf("get a pull request: %w", err)
	}
	return pr, nil
}
