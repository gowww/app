package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func initWatcher() {
	if watcher != nil {
		return
	}
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	watcher.Add(".")
	watcherAdd("scripts")
	watcherAdd("styles")
	watcherAdd("views")
}

// watcherAdd adds a directory and its subdirectories to the watcher.
func watcherAdd(dir string) {
	if watcher.Add(dir) == nil {
		filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !f.IsDir() {
				return nil
			}
			return watcher.Add(path)
		})
	}
}

func eventIs(e fsnotify.Event, ops ...fsnotify.Op) bool {
	for _, op := range ops {
		if e.Op&op == op {
			return true
		}
	}
	return false
}

func watch() {
	initWatcher()
	if buildGo() == nil {
		run()
	}
	for {
		select {
		case e := <-watcher.Events:
			watchEvent(e)
		case err := <-watcher.Errors:
			if err != nil {
				panic(err)
			}
		}
	}
}

func watchEvent(e fsnotify.Event) {
	if eventIs(e, fsnotify.Chmod) {
		return
	}
	if strings.HasPrefix(e.Name, "scripts/") {
		if strings.HasSuffix(e.Name, ".go") && !strings.HasSuffix(e.Name, "_test.go") {
			buildScriptsGopherJS()
		}
		// TODO: Babel, CoffeeScript, TypeScript...
	} else if strings.HasPrefix(e.Name, "styles/") &&
		!strings.Contains(e.Name, "mixin") &&
		!strings.Contains(e.Name, "partial") &&
		filepath.Base(e.Name)[0] != '_' {
		if strings.HasSuffix(e.Name, ".styl") {
			buildStylesStylus(e.Name)
		}
		// TODO: LESS, SASS, SCSS...
	} else if strings.HasPrefix(e.Name, "views/") {
		run()
	} else if strings.HasSuffix(e.Name, ".go") && !strings.HasSuffix(e.Name, "_test.go") {
		if buildGo() == nil {
			run()
		}
	}
}
