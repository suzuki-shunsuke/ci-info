package cli

import (
	"github.com/suzuki-shunsuke/ci-info/pkg/controller"
	"github.com/urfave/cli/v2"
)

func (runner Runner) setCLIArg(c *cli.Context, params controller.Params) controller.Params {
	if owner := c.String("owner"); owner != "" {
		params.Owner = owner
	}
	if repo := c.String("repo"); repo != "" {
		params.Repo = repo
	}
	if token := c.String("github-token"); token != "" {
		params.GitHubToken = token
	}
	if logLevel := c.String("log-level"); logLevel != "" {
		params.LogLevel = logLevel
	}
	if prefix := c.String("prefix"); prefix != "" {
		params.Prefix = prefix
	}
	if sha := c.String("sha"); sha != "" {
		params.SHA = sha
	}
	if dir := c.String("dir"); dir != "" {
		params.Dir = dir
	}
	if prNum := c.Int("pr"); prNum > 0 {
		params.PRNum = prNum
	}
	return params
}

func (runner Runner) action(c *cli.Context) error {
	params := controller.Params{}
	params = runner.setCLIArg(c, params)

	ctrl, params, err := controller.New(c.Context, params)
	if err != nil {
		return err
	}

	return ctrl.Run(c.Context, params)
}
