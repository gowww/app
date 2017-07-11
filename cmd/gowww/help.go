package main

import (
	"flag"
	"fmt"
)

func help() {
	switch flag.Arg(0) {
	case "build":
		helpBuild()
	case "help":
		switch flag.Arg(1) {
		case "build":
			helpBuild()
		case "watch":
			helpWatch()
		default:
			helpMain()
		}
	case "watch":
		helpWatch()
	default:
		helpMain()
	}
}

func helpMain() {
	fmt.Print(`The CLI of the gowww/app framework.

Usage:

	gowww [command] [flags]

Commands:

	build  Create binary for app.
	watch  Detect changes and rerun app.

Flags:

	-name  The file name used for build. Default: ` + getwd(false) + `.

`)
}

func helpBuild() {
	fmt.Print(`Create binary for app.

Usage:

	gowww build [flags]

Flags:

	-docker  Use Docker's "golang:latest" image to build for Linux.
	-name    The file name used for build.

`)
}

func helpWatch() {
	fmt.Print(`Detect changes and rerun app.

Usage:

	gowww watch

`)
}
