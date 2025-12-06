package cli

import (
	"context"

	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
	"github.com/urfave/cli/v3"
)

func Run(ctx context.Context, logger *slogutil.Logger, env *urfave.Env) error {
	r := &Runner{}
	return urfave.Command(env, &cli.Command{ //nolint:wrapcheck
		Name:  "ci-info",
		Usage: "get CI information. https://github.com/suzuki-shunsuke/ci-info",
		Flags: []cli.Flag{},
		Commands: []*cli.Command{
			r.runCommand(logger),
		},
	}).Run(ctx, env.Args)
}

type Runner struct{}
