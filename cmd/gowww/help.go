package main

import (
	"fmt"
	"os"
)

func cmdHelp() {
	fmt.Print(`gowww is the CLI of the gowww/app framework.

Usage:

	gowww

Inside a gowww/app project, run command "gowww" to watch for changes and rerun your app.
`)
	os.Exit(0)
}
