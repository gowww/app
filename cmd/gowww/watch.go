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
				if strings.HasPrefix(event.Name, "scripts/") {
					if strings.HasSuffix(event.Name, ".go") && !strings.HasSuffix(event.Name, "_test.go") {
						buildScriptsGopherJS()
					}
					// TODO: Babel, CoffeeScript, TypeScript...
				} else if strings.HasPrefix(event.Name, "styles/") {
					if !strings.Contains(event.Name, "partial") && !strings.Contains(event.Name, "mixin") && filepath.Base(event.Name)[0] != '_' {
						if strings.HasSuffix(event.Name, ".styl") {
							buildStylesStylus(event.Name)
						}
						// TODO: LESS, SASS, SCSS...
					}
				} else if strings.HasPrefix(event.Name, "views/") {
					run()
				} else if strings.HasSuffix(event.Name, ".go") && !strings.HasSuffix(event.Name, "_test.go") {
					if buildGo() == nil {
						run()
					}
				}
			}
		case err := <-watcher.Errors:
			if err != nil {
				panic(err)
			}
		}
	}
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
