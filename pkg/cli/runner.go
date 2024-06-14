package cli

import (
	"context"
	"io"

	"github.com/urfave/cli/v2"
)

type LDFlags struct {
	Version string
	Commit  string
	Date    string
}

func (flags *LDFlags) AppVersion() string {
	return flags.Version + " (" + flags.Commit + ")"
}

type Runner struct {
	Stdin   io.Reader
	Stdout  io.Writer
	Stderr  io.Writer
	LDFlags *LDFlags
}

func (runner *Runner) Run(ctx context.Context, args ...string) error {
	app := cli.App{
		Name:    "ci-info",
		Usage:   "get CI information. https://github.com/suzuki-shunsuke/ci-info",
		Version: runner.LDFlags.AppVersion(),
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
					&cli.BoolFlag{
						Name:  "wait-mergeable",
						Usage: "wait until the pull request's 'mergeable' becomes not null",
					},
					&cli.IntFlag{
						Name:  "wait-mergeable-timeout",
						Usage: "timeout of wait-mergeable (second)",
						Value: 60, //nolint:gomnd
					},
				},
			},
		},
	}

	return app.RunContext(ctx, args) //nolint:wrapcheck
}
