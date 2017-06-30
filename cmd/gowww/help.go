package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	helpMain = `The CLI of the gowww/app framework.

Usage:

	gowww [command] [flags]

Commands:

	build  Create binary for app.
	watch  Detect changes and rerun app.

Flags:

	-name  The file name used for build. Default: ` + getwd(false) + `.

`

	helpBuild = `Create binary for app.

Usage:

	gowww build [flags]

Flags:

	-docker  Use Docker's "golang:latest" image to build for Linux.

`

	helpWatch = `Detect changes and rerun app.

Usage:

	gowww watch

`
)

func help() {
	switch flag.Arg(0) {
	case "build":
		fmt.Print(helpBuild)
	case "watch":
		fmt.Print(helpWatch)
	default:
		fmt.Print(helpMain)
	}
	os.Exit(0)
}
