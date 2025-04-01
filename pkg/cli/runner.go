package cli

import (
	"context"
	"io"

	"github.com/urfave/cli/v3"
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

func (r *Runner) Run(ctx context.Context, args ...string) error {
	app := cli.App{
		Name:    "ci-info",
		Usage:   "get CI information. https://github.com/suzuki-shunsuke/ci-info",
		Version: r.LDFlags.AppVersion(),
		Commands: []*cli.Command{
			r.runCommand(),
		},
	}

	return app.RunContext(ctx, args) //nolint:wrapcheck
}
