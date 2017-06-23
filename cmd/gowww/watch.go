package main

import "github.com/fsnotify/fsnotify"

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
	watcher.Add("scripts")
	watcher.Add("views")
}

func watch() {
	initWatcher()
	if build() == nil {
		run()
	}
	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Create == fsnotify.Create ||
				event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Remove == fsnotify.Remove ||
				event.Op&fsnotify.Rename == fsnotify.Rename {
				if reFilenameGo.MatchString(event.Name) {
					if build() == nil {
						run()
					}
				} else if reFilenameScriptsGopherJS.MatchString(event.Name) {
					buildScriptsGopherJS()
				} else if reFilenameViews.MatchString(event.Name) {
					run()
				}
			}
		case err := <-watcher.Errors:
			if err != nil {
				panic(err)
			}
		}
	}
}
