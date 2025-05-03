package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ci-info/pkg/controller"
	"github.com/suzuki-shunsuke/ci-info/pkg/domain"
	"github.com/suzuki-shunsuke/ci-info/pkg/github"
	"github.com/suzuki-shunsuke/go-ci-env/v3/cienv"
	"github.com/urfave/cli/v3"
)

func (r *Runner) runCommand() *cli.Command {
	return &cli.Command{
		Name:   "run",
		Usage:  "get CI information",
		Action: r.action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "owner",
				Usage: "repository owner",
			},
			&cli.StringFlag{
				Name:  "repo",
				Usage: "repository name",
			},
			&cli.StringFlag{
				Name:  "sha",
				Usage: "commit SHA",
			},
			&cli.StringFlag{
				Name:  "dir",
				Usage: "directory path where files are created. The directory is created by os.MkdirAll if it doesn't exist. By default the directory is created at Go's ioutil.TempDir",
			},
			&cli.IntFlag{
				Name:  "pr",
				Usage: "pull request number",
			},
			&cli.StringFlag{
				Name:  "github-token",
				Usage: "GitHub Access Token [$GITHUB_TOKEN, $GITHUB_ACCESS_TOKEN]",
			},
			&cli.StringFlag{
				Name:    "github-api-url",
				Usage:   "GitHub API Base URL",
				Sources: cli.EnvVars("GITHUB_API_URL"),
			},
			&cli.StringFlag{
				Name:    "github-graphql-url",
				Usage:   "GitHub GraphQL API URL",
				Sources: cli.EnvVars("GITHUB_GRAPHQL_URL"),
			},
			&cli.StringFlag{
				Name:  "prefix",
				Usage: "The prefix of environment variable name",
				Value: "CI_INFO_",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "log level",
			},
		},
	}
}

func (r *Runner) setCLIArg(cmd *cli.Command, params domain.Params) domain.Params {
	if owner := cmd.String("owner"); owner != "" {
		params.Owner = owner
	}
	if repo := cmd.String("repo"); repo != "" {
		params.Repo = repo
	}
	if token := cmd.String("github-token"); token != "" {
		params.GitHubToken = token
	}
	if logLevel := cmd.String("log-level"); logLevel != "" {
		params.LogLevel = logLevel
	}
	if prefix := cmd.String("prefix"); prefix != "" {
		params.Prefix = prefix
	}
	if sha := cmd.String("sha"); sha != "" {
		params.SHA = sha
	}
	if dir := cmd.String("dir"); dir != "" {
		params.Dir = dir
	}
	if prNum := cmd.Int("pr"); prNum > 0 {
		params.PRNum = prNum
	}
	params.GitHubAPIURL = cmd.String("github-api-url")
	params.GitHubGraphQLURL = cmd.String("github-graphql-url")
	return params
}

func (r *Runner) action(ctx context.Context, c *cli.Command) error {
	params := domain.Params{}
	params = r.setCLIArg(c, params)
	if err := setEnv(&params); err != nil {
		return err
	}
	setLogLevel(params.LogLevel)
	ghClient, err := github.New(ctx, github.ParamsNew{
		Token:      getGitHubToken(params.GitHubToken),
		BaseURL:    params.GitHubAPIURL,
		GraphQLURL: params.GitHubGraphQLURL,
	})
	if err != nil {
		return fmt.Errorf("create a GitHub client: %w", err)
	}

	fs := afero.NewOsFs()

	ctrl := controller.New(ghClient, fs)

	return ctrl.Run(ctx, params) //nolint:wrapcheck
}

func getGitHubToken(token string) string {
	if token != "" {
		return token
	}
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITHUB_ACCESS_TOKEN")
}

func setLogLevel(logLevel string) {
	if logLevel == "" {
		return
	}
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"log_level": logLevel,
		}).WithError(err).Error("the log level is invalid")
	}
	logrus.SetLevel(lvl)
}

func setEnv(params *domain.Params) error {
	platform := cienv.Get(nil)
	if platform == nil {
		return nil
	}
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
			return fmt.Errorf("get the pull request number: %w", err)
		}
		params.PRNum = prNum
	}
	return nil
}
