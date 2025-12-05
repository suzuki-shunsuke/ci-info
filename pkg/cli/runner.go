package cli

import (
	"context"
	"io"
	"log/slog"

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
	Stdin       io.Reader
	Stdout      io.Writer
	Stderr      io.Writer
	LDFlags     *LDFlags
	LogLevelVar *slog.LevelVar
}

func (r *Runner) Run(ctx context.Context, logger *slog.Logger, args ...string) error {
	cmd := cli.Command{
		Name:    "ci-info",
		Usage:   "get CI information. https://github.com/suzuki-shunsuke/ci-info/v2",
		Version: r.LDFlags.AppVersion(),
		Commands: []*cli.Command{
			r.runCommand(logger),
		},
	}

	return cmd.Run(ctx, args) //nolint:wrapcheck
}
