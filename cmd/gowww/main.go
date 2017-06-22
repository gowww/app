// Gowww helps developing with gowww/app.
package main

import (
	"flag"
	"os"
	"os/signal"
)

func main() {
	flag.Usage = cmdHelp
	flag.Parse()
	switch flag.Arg(0) {
	case "":
		cmdRun()
	case "help":
		cmdHelp()
	case "run":
		cmdRun()
	default:
		cmdHelp()
	}
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
