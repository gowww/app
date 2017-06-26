// gowww is the CLI of the gowww/app framework.
package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"

	"github.com/arthurwhite/kingpin"
	"github.com/fsnotify/fsnotify"
)

var (
	cmd         = kingpin.New("gowww", "The CLI for the gowww/app framework.")
	cmdFlagName = cmd.Flag("name", "The file name used for build.").Default(getwd(false)).Short('n').String()

	cmdBuild           = cmd.Command("build", "Create binary for app.").Alias("b")
	cmdBuildFlagDocker = cmdBuild.Flag("docker", "User Docker's \"golang:latest\" image to build for Linux.").Short('d').Bool()

	cmdWatch = cmd.Command("watch", "Detect changes and rerun app.").Alias("w")

	watcher     *fsnotify.Watcher
	runningProc *os.Process

	reFilenameGo              = regexp.MustCompile(`^[0-9A-Za-z_-]+[^_test].go$`)
	reFilenameScriptsGopherJS = regexp.MustCompile(`^scripts/[0-9A-Za-z_-]+[^_test].go$`)
	reFilenameViews           = regexp.MustCompile(`^views/[0-9A-Za-z_-]+.gohtml$`)
)

func init() {
	cmd.HelpFlag.Short('h')
}

func main() {
	defer clean()
	atexit(clean)

	switch kingpin.MustParse(cmd.Parse(os.Args[1:])) {
	case cmdBuild.FullCommand():
		if *cmdBuildFlagDocker {
			buildDocker()
		} else {
			build()
		}
	case cmdWatch.FullCommand():
		watch()
	}
}

func run() {
	if runningProc != nil {
		if err := runningProc.Kill(); err != nil {
			panic(err)
		}
	}
	cmd := exec.Command("./" + buildName())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	runningProc = cmd.Process
}

func clean() {
	if watcher != nil {
		watcher.Close()
	}
	if runningProc != nil {
		runningProc.Kill()
	}
}

func cleanLines(n int) {
	fmt.Printf("\033[%dA\033[0K", n)
}

func getwd(fullpath bool) string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if fullpath {
		return wd
	}
	return filepath.Base(wd)
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
