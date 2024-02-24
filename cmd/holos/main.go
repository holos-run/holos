package main

import (
	"github.com/holos-run/holos/pkg/cli"
	"os"
)

func main() {
	os.Exit(cli.MakeMain()())
}
