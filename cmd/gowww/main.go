// gowww is the CLI of the gowww/app framework.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"

	"github.com/fsnotify/fsnotify"
)

var (
	flagName        = flag.String("name", getwd(false), "The file name used for build.")
	flagBuildDocker = flag.Bool("docker", false, `Use Docker's "golang:latest" image to build for Linux.`)

	subprocArgs []string
	watcher     *fsnotify.Watcher
	runningProc *os.Process

	reFilenameGo              = regexp.MustCompile(`^[0-9A-Za-z_-]+[^_test].go$`)
	reFilenameScriptsGopherJS = regexp.MustCompile(`^scripts/[0-9A-Za-z_-]+[^_test].go$`)
	reFilenameViews           = regexp.MustCompile(`^views/[0-9A-Za-z_-]+.gohtml$`)
)

func main() {
	defer clean()
	atexit(clean)

	flag.Usage = help
	flag.Parse()

	// Pass args after "--" to subprocess.
	for i := 0; i < len(os.Args); i++ {
		if os.Args[i] == "--" {
			subprocArgs = os.Args[i+1:]
		}
	}

	switch flag.Arg(0) {
	case "", "watch":
		watch()
	case "build":
		build()
	default:
		help()
	}
}

func run() {
	if runningProc != nil {
		if err := runningProc.Kill(); err != nil {
			panic(err)
		}
	}
	cmd := exec.Command("./"+buildName(), subprocArgs...)
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
