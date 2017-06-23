package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

var (
	watchBuildName string
	runningProc    *os.Process
	watcher        *fsnotify.Watcher
)

func setWatchBuildName() {
	var err error
	watchBuildName, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	watchBuildName = filepath.Base(watchBuildName)
}

func initWatcher() {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	watcher.Add(".")
	watcher.Add("scripts")
	watcher.Add("views")
}

func cmdWatch() {
	setWatchBuildName()
	atexit(watchClean)
	initWatcher()

	var args []string
	if len(flag.Args()) > 1 {
		args = flag.Args()[1:]
	}

	if build() == nil {
		run(args)
	}
	// TODO: Build scripts (Babel, GopherJS, TypeScript...) and styles (LESS, SCSS, Stylus...) before and during watching.
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Remove == fsnotify.Remove ||
				event.Op&fsnotify.Rename == fsnotify.Rename {
				if reFilenameGo.MatchString(event.Name) {
					if build() == nil {
						run(args)
					}
				} else if reFilenameScriptsGopherJS.MatchString(event.Name) {
					buildScriptsGopherJS()
				} else if reFilenameViews.MatchString(event.Name) {
					run(args)
				}
			}
		case err := <-watcher.Errors:
			if err != nil {
				panic(err)
			}
		}
	}
}

func build() error {
	log.Println("Building...")
	cmd := exec.Command("go", "build", "-o", watchBuildName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		cleanLines(1)
	}
	return err
}

func run(args []string) {
	if runningProc != nil {
		if err := runningProc.Kill(); err != nil {
			panic(err)
		}
	}
	cmd := exec.Command("./"+watchBuildName, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	runningProc = cmd.Process
}

func buildScriptsGopherJS() error {
	log.Println("Building scripts with GopherJS...")
	cmd := exec.Command("gopherjs", "build", "./scripts", "-o", "static/main.js", "-m")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		cleanLines(1)
	}
	return err
}

func watchClean() {
	if watcher != nil {
		watcher.Close()
	}
	if runningProc != nil {
		runningProc.Kill()
	}
	if watchBuildName != "" {
		os.Remove(watchBuildName)
	}
}
