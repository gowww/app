package main

import (
	"bufio"
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
	watcherAddRecur(dirScripts)
	watcherAddRecur(dirStyles)
	watcherAddRecur(dirViews)
}

// watcherAddRecur adds a directory and its subdirectories to the watcher.
func watcherAddRecur(dir string) {
	if watcher.Add(dir) == nil {
		err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !f.IsDir() {
				return nil
			}
			return watcher.Add(path)
		})
		if err != nil {
			panic(err)
		}
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

	if strings.HasPrefix(e.Name, dirScripts+"/") {
		watcherAddCreated(e)

		// GopherJS
		if strings.HasSuffix(e.Name, ".go") {
			if strings.HasSuffix(e.Name, "_test.go") || !packageIsMain(e.Name) {
				return
			}
			outFile := filepath.Dir(strings.TrimPrefix(e.Name, dirScripts+"/"))
			outFile += "/main.js"
			buildScriptsGopherJS(e.Name, outFile)
			return
		}

		// TODO: Babel, CoffeeScript, TypeScript...
		return
	}

	if strings.HasPrefix(e.Name, dirStyles+"/") {
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
			os.Remove("static/" + dirStyles + "/" + name + ".css")
			os.Remove("static/" + dirStyles + "/" + name + ".css.map")
			return
		}

		// TODO: LESS, SASS, SCSS...
		return
	}

	if strings.HasPrefix(e.Name, dirViews+"/") {
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

// packageIsMain checks that package defined in Go file is "main".
func packageIsMain(file string) bool {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if scanner.Text() == "package main" {
			return true
		}
	}
	return false
}
