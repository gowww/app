package main

import (
	"fmt"
	"os"
)

func help() {
	fmt.Print(`The CLI of the gowww/app framework.

Usage:

	gowww [command] [flags]

Commands:

	build  Create binary for app.
	watch  Detect changes and rerun app.

Flags:

	-name  The file name used for build. Default: ` + getwd(false) + `.

`)
	os.Exit(0)
}
