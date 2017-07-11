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

func watch() {
	initWatcher()
	if buildGo() == nil {
		run()
	}
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Remove == fsnotify.Remove ||
				event.Op&fsnotify.Rename == fsnotify.Rename {
				watchEvent(event.Name)
			}
		case err := <-watcher.Errors:
			if err != nil {
				panic(err)
			}
		}
	}
}

func watchEvent(name string) {
	if strings.HasPrefix(name, "scripts/") {
		if strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
			buildScriptsGopherJS()
		}
		// TODO: Babel, CoffeeScript, TypeScript...
	} else if strings.HasPrefix(name, "styles/") &&
		!strings.Contains(name, "mixin") &&
		!strings.Contains(name, "partial") &&
		filepath.Base(name)[0] != '_' {
		if strings.HasSuffix(name, ".styl") {
			buildStylesStylus(name)
		}
		// TODO: LESS, SASS, SCSS...
	} else if strings.HasPrefix(name, "views/") {
		run()
	} else if strings.HasSuffix(name, ".go") && !strings.HasSuffix(name, "_test.go") {
		if buildGo() == nil {
			run()
		}
	}
}
