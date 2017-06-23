// gowww is the CLI of the gowww/app framework.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"regexp"
)

var (
	reFilenameGo              = regexp.MustCompile(`^[0-9A-Za-z_-]+[^_test].go$`)
	reFilenameScriptsGopherJS = regexp.MustCompile(`^scripts/[0-9A-Za-z_-]+[^_test].go$`)
	reFilenameViews           = regexp.MustCompile(`^views/[0-9A-Za-z_-]+.gohtml$`)
)

func main() {
	flag.Usage = cmdHelp
	flag.Parse()
	switch flag.Arg(0) {
	case "", "watch":
		cmdWatch()
	default:
		cmdHelp()
	}
}

func cleanLines(n int) {
	fmt.Printf("\033[%dA\033[0K", n)
}

func atexit(f func()) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, os.Interrupt)
	go func() {
		<-s
		f()
		os.Exit(0)
	}()
}
