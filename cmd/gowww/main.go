// gowww is the CLI of the gowww/app framework.
package main

import (
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/gowww/cli"
)

var (
	flagBuildDocker bool
	flagBuildName   string

	watcher     *fsnotify.Watcher
	runningProc *os.Process
)

func main() {
	defer clean()
	atexit(clean)

	cli.Description = "The CLI of the gowww/app framework."

	cli.Command("build", build, "Create binary for app.").
		Bool(&flagBuildDocker, "docker", false, `Use Docker's "golang:latest" image to build for Linux.`).
		String(&flagBuildName, "name", getwd(false), "The file name used for build.")

	cli.Command("watch", watch, "Detect changes and rerun app.")

	cli.Parse()

	watch()
}

func run() {
	if runningProc != nil {
		if err := runningProc.Kill(); err != nil {
			panic(err)
		}
	}
	cmd := exec.Command("./"+buildName(), cli.SubArgs()...)
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
