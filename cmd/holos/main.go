package main

import (
	"os"

	"github.com/holos-run/holos/internal/cli"
)

func main() {
	os.Exit(cli.MakeMain()())
}
