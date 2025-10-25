package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/suzuki-shunsuke/ci-info/v2/pkg/domain"
	"github.com/suzuki-shunsuke/ci-info/v2/pkg/github"
	"github.com/suzuki-shunsuke/ci-info/v2/pkg/output"
	"github.com/suzuki-shunsuke/ci-info/v2/pkg/write"
)

func (c *Controller) Run(ctx context.Context, params domain.Params) error {
	if err := validateParams(params); err != nil {
		return fmt.Errorf("argument is invalid: %w", err)
	}

	isPR := params.PRNum > 0

	pr, err := c.gh.GetPR(ctx, params)
	if err != nil {
		return fmt.Errorf("get an associated pull request: %w", err)
	}

	if pr == nil {
		fmt.Fprintln(c.stdout, output.NonPREnv(params))
		return nil
	}

	files, _, err := c.gh.GetPRFiles(ctx, github.ParamsGetPRFiles{
		Owner:    params.Owner,
		Repo:     params.Repo,
		PRNum:    pr.GetNumber(),
		FileSize: pr.GetChangedFiles(),
	})
	if err != nil {
		return fmt.Errorf("get pull request files: %w", err)
	}

	dir, err := write.MkDir(c.fs, params.Dir)
	if err != nil {
		return fmt.Errorf("create a directory: %w", err)
	}

	fmt.Fprintln(c.stdout, output.PREnv(params.Prefix, dir, isPR, params.Owner, params.Repo, pr))

	if err := write.Write(c.fs, dir, pr, files); err != nil {
		return fmt.Errorf("write files: %w", err)
	}
	return nil
}

var (
	errOwnerRequired      = errors.New("owner is required")
	errRepoRequired       = errors.New("repo is required")
	errSHAOrPRNumRequired = errors.New("sha or pr number is required")
)

func validateParams(params domain.Params) error {
	if params.Owner == "" {
		return errOwnerRequired
	}
	if params.Repo == "" {
		return errRepoRequired
	}
	if params.PRNum <= 0 && params.SHA == "" {
		return errSHAOrPRNumRequired
	}
	return nil
}
