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
	watcherAddRecur("scripts")
	watcherAddRecur("styles")
	watcherAddRecur("views")
}

// watcherAddRecur adds a directory and its subdirectories to the watcher.
func watcherAddRecur(dir string) {
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

// watcherAddCreated adds a directory and its subdirectories to the watcher if event is "create".
func watcherAddCreated(e fsnotify.Event) {
	if eventIs(e, fsnotify.Create) { // Watch newly created directories.
		watcherAddRecur(e.Name)
	}
}

// eventIs check that event e is one of ops.
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
		watcherAddCreated(e)

		// GopherJS
		if strings.HasSuffix(e.Name, ".go") {
			if strings.HasSuffix(e.Name, "_test.go") {
				return
			}
			buildScriptsGopherJS()
			return
		}

		// TODO: Babel, CoffeeScript, TypeScript...
		return
	}

	if strings.HasPrefix(e.Name, "styles/") {
		if strings.Contains(e.Name, "mixin") || strings.Contains(e.Name, "partial") || filepath.Base(e.Name)[0] == '_' {
			return
		}
		watcherAddCreated(e)

		// Stylus
		if strings.HasSuffix(e.Name, ".styl") {
			if eventIs(e, fsnotify.Write) {
				buildStylesStylus(e.Name)
				return
			}
			name := filepath.Base(e.Name)
			name = strings.TrimSuffix(name, filepath.Ext(name))
			os.Remove("static/styles/" + name + ".css")
			os.Remove("static/styles/" + name + ".css.map")
			return
		}

		// TODO: LESS, SASS, SCSS...
		return
	}

	if strings.HasPrefix(e.Name, "views/") {
		watcherAddCreated(e)
		run()
		return
	}

	if strings.HasSuffix(e.Name, ".go") {
		if strings.HasSuffix(e.Name, "_test.go") {
			return
		}
		if buildGo() == nil {
			run()
			return
		}
		return
	}
}
