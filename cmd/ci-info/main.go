package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"github.com/suzuki-shunsuke/ci-info/v2/pkg/cli"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
)

var (
	version = ""
	commit  = "" //nolint:gochecknoglobals
	date    = "" //nolint:gochecknoglobals
)

func main() {
	if code := core(); code != 0 {
		os.Exit(code)
	}
}

func core() int {
	logLevelVar := &slog.LevelVar{}
	logger := slogutil.New(&slogutil.InputNew{
		Name:    "ci-info",
		Version: version,
		Out:     os.Stderr,
		Level:   logLevelVar,
	})
	runner := cli.Runner{
		LDFlags: &cli.LDFlags{
			Version: version,
			Commit:  commit,
			Date:    date,
		},
		LogLevelVar: logLevelVar,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	if err := runner.Run(ctx, logger, os.Args...); err != nil {
		slogerr.WithError(logger, err).Error("ci-info failed")
		return 1
	}
	return 0
}
