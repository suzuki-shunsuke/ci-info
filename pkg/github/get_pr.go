package github

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/cenkalti/backoff/v5"
	"github.com/google/go-github/v83/github"
	"github.com/suzuki-shunsuke/ci-info/v2/pkg/domain"
)

var errNoPRFound = errors.New("no associated pull request found")

func (c *Client) getPRNum(ctx context.Context, logger *slog.Logger, params *domain.Params) (int, error) {
	if params.PRNum > 0 {
		return params.PRNum, nil
	}
	logger.Debug("get pull request from SHA", "owner", params.Owner, "repo", params.Repo, "sha", params.SHA)

	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 2 * time.Second
	cnt := 0
	prs, err := backoff.Retry(ctx, func() ([]*github.PullRequest, error) {
		if cnt > 0 {
			logger.Info("retry to get pull request from SHA", "attempt", cnt)
		}
		cnt++
		prs, _, err := c.ListPRsWithCommit(ctx, paramsListPRsWithCommit{
			Owner: params.Owner,
			Repo:  params.Repo,
			SHA:   params.SHA,
		})
		if err != nil {
			return nil, fmt.Errorf("list pull requests with a commit: %w", backoff.Permanent(err))
		}
		if len(prs) == 0 {
			return nil, errNoPRFound
		}
		return prs, nil
	},
		backoff.WithBackOff(b),
		backoff.WithMaxElapsedTime(40*time.Second),
	)
	if err != nil {
		return 0, err //nolint:wrapcheck
	}
	logger.Debug("the number of pull requests assosicated with the commit", "size", len(prs))
	return prs[0].GetNumber(), nil
}

func (c *Client) GetPR(ctx context.Context, logger *slog.Logger, params *domain.Params) (*PullRequest, error) {
	prNum, err := c.getPRNum(ctx, logger, params)
	if err != nil {
		return nil, err
	}
	if prNum <= 0 {
		return nil, nil //nolint:nilnil
	}
	pr, _, err := c.Client.PullRequests.Get(ctx, params.Owner, params.Repo, prNum)
	if err != nil {
		return nil, fmt.Errorf("get a pull request: %w", err)
	}
	return pr, nil
}
