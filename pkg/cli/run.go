package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/ci-info/v2/pkg/controller"
	"github.com/suzuki-shunsuke/ci-info/v2/pkg/domain"
	"github.com/suzuki-shunsuke/ci-info/v2/pkg/github"
	"github.com/suzuki-shunsuke/go-ci-env/v3/cienv"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/urfave/cli/v3"
)

func (r *Runner) runCommand(logger *slogutil.Logger) *cli.Command { //nolint:funlen
	params := &domain.Params{}
	return &cli.Command{
		Name:  "run",
		Usage: "get CI information",
		Action: func(ctx context.Context, _ *cli.Command) error {
			return r.action(ctx, logger, params)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "owner",
				Usage:       "repository owner",
				Destination: &params.Owner,
			},
			&cli.StringFlag{
				Name:        "repo",
				Usage:       "repository name",
				Destination: &params.Repo,
			},
			&cli.StringFlag{
				Name:        "sha",
				Usage:       "commit SHA",
				Destination: &params.SHA,
			},
			&cli.StringFlag{
				Name:        "dir",
				Usage:       "directory path where files are created. The directory is created by os.MkdirAll if it doesn't exist. By default the directory is created at Go's ioutil.TempDir",
				Destination: &params.Dir,
			},
			&cli.IntFlag{
				Name:        "pr",
				Usage:       "pull request number",
				Destination: &params.PRNum,
			},
			&cli.StringFlag{
				Name:        "github-token",
				Usage:       "GitHub Access Token [$GITHUB_TOKEN, $GITHUB_ACCESS_TOKEN]",
				Destination: &params.GitHubToken,
			},
			&cli.StringFlag{
				Name:        "github-api-url",
				Usage:       "GitHub API Base URL",
				Sources:     cli.EnvVars("GITHUB_API_URL"),
				Destination: &params.GitHubAPIURL,
			},
			&cli.StringFlag{
				Name:        "github-graphql-url",
				Usage:       "GitHub GraphQL API URL",
				Sources:     cli.EnvVars("GITHUB_GRAPHQL_URL"),
				Destination: &params.GitHubGraphQLURL,
			},
			&cli.StringFlag{
				Name:        "prefix",
				Usage:       "The prefix of environment variable name",
				Value:       "CI_INFO_",
				Destination: &params.Prefix,
			},
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "log level",
				Destination: &params.LogLevel,
			},
		},
	}
}

func (r *Runner) action(ctx context.Context, logger *slogutil.Logger, params *domain.Params) error {
	if err := setEnv(params); err != nil {
		return err
	}
	if err := logger.SetLevel(params.LogLevel); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}
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

	l := logger.Logger
	if params.Owner != "" {
		l = l.With("owner", params.Owner)
	}
	if params.Repo != "" {
		l = l.With("repo", params.Repo)
	}
	if params.Prefix != "" {
		l = l.With("prefix", params.Prefix)
	}
	if params.SHA != "" {
		l = l.With("sha", params.SHA)
	}
	if params.Dir != "" {
		l = l.With("dir", params.Dir)
	}
	if params.PRNum > 0 {
		l = l.With("pr", params.PRNum)
	}

	return ctrl.Run(ctx, l, params) //nolint:wrapcheck
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
