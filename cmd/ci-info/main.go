package main

import (
	"github.com/suzuki-shunsuke/ci-info/v2/pkg/cli"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
)

var version = ""

func main() {
	urfave.Main("ci-info", version, cli.Run)
}
