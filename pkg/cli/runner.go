package cli

import (
	"context"
	"io"

	"github.com/suzuki-shunsuke/ci-info/pkg/constant"
	"github.com/urfave/cli/v2"
)

type Runner struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (runner *Runner) Run(ctx context.Context, args ...string) error {
	app := cli.App{
		Name:    "ci-info",
		Usage:   "get CI information. https://github.com/suzuki-shunsuke/ci-info",
		Version: constant.Version,
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "get CI information",
				Action: runner.action,
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
						Name:  "prefix",
						Usage: "The prefix of environment variable name",
						Value: "CI_INFO_",
					},
					&cli.StringFlag{
						Name:  "log-level",
						Usage: "log level",
					},
				},
			},
		},
	}

	return app.RunContext(ctx, args)
}
