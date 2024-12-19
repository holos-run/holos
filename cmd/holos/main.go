package main

import (
	"os"

	"github.com/holos-run/holos/cmd"
)

func main() {
	os.Exit(cmd.MakeMain()())
}
