package controller

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ci-info/pkg/domain"
	"github.com/suzuki-shunsuke/ci-info/pkg/github"
	"github.com/suzuki-shunsuke/ci-info/pkg/output"
	"github.com/suzuki-shunsuke/ci-info/pkg/write"
)

func (c *Controller) Run(ctx context.Context, params domain.Params) error {
	if err := validateParams(params); err != nil {
		return fmt.Errorf("argument is invalid: %w", err)
	}

	isPR := params.PRNum > 0

	pr, err := c.getPR(ctx, params)
	if err != nil {
		return err
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

	dir, err := c.mkDir(params.Dir)
	if err != nil {
		return err
	}

	fmt.Fprintln(c.stdout, output.PREnv(params.Prefix, dir, isPR, params.Owner, params.Repo, pr))

	if err := write.Write(c.fs, dir, pr, files); err != nil {
		return fmt.Errorf("write files: %w", err)
	}
	return nil
}

func (c *Controller) mkDir(dir string) (string, error) {
	if dir == "" {
		d, err := afero.TempDir(c.fs, "", "ci-info")
		if err != nil {
			return "", fmt.Errorf("create a temporal directory: %w", err)
		}
		return d, nil
	}
	if !filepath.IsAbs(dir) {
		d, err := filepath.Abs(dir)
		if err != nil {
			return "", fmt.Errorf("convert -dir %s to absolute path: %w", dir, err)
		}
		dir = d
	}
	if err := c.fs.MkdirAll(dir, 0o755); err != nil { //nolint:mnd
		return "", fmt.Errorf("create a directory %s: %w", dir, err)
	}
	return dir, nil
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
