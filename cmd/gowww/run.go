package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func cmdRun() {
	var args []string
	if len(flag.Args()) > 1 {
		args = flag.Args()[1:]
	}

	name, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	name = filepath.Base(name)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	var prevProc *os.Process
	atexit(func() {
		watcher.Close()
		os.Remove(name)
	})

	if err = watcher.Add("."); err != nil {
		panic(err)
	}
	if err = watcher.Add("views"); err != nil {
		panic(err)
	}

	if build(name) == nil {
		prevProc = run(prevProc, name, args)
	}
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Remove == fsnotify.Remove ||
				event.Op&fsnotify.Rename == fsnotify.Rename {
				switch filepath.Ext(event.Name) {
				case ".go":
					if build(name) == nil {
						prevProc = run(prevProc, name, args)
					}
				case ".gohtml":
					prevProc = run(prevProc, name, args)
				}
			}
		case err := <-watcher.Errors:
			if err != nil {
				panic(err)
			}
		}
	}
}

func build(name string) error {
	log.Println("Rebuilding...")
	buildCmd := exec.Command("go", "build", "-o", name)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	return buildCmd.Run()
}

func run(prevProc *os.Process, name string, args []string) *os.Process {
	log.Println("Rerunning...")
	if prevProc != nil {
		if err := prevProc.Kill(); err != nil {
			panic(err)
		}
	}

	runCmd := exec.Command("./"+name, args...)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	if err := runCmd.Start(); err != nil {
		panic(err)
	}
	return runCmd.Process
}
