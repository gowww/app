/*
Package app provides a full featured framework for any web apps.
*/
package app

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	gowwwlog "github.com/gowww/log"
)

var (
	address    = flag.String("a", ":8080", "the address to listen and serving on")
	production = flag.Bool("p", false, "run the server in production environment")
	mux        = http.NewServeMux()
)

func init() {
	flag.Parse()

	// Serve static content
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Parse views
	files, _ := ioutil.ReadDir("views")
	for _, f := range files {
		if !f.IsDir() && filepath.Ext(f.Name()) == ".gohtml" {
			parseViews()
			return
		}
	}
}

// A Handler handles a request.
type Handler func(*Context)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(&Context{w, r})
}

// Route adds a new route to the muxer.
func Route(pattern string, handler Handler) {
	mux.Handle(pattern, handler)
}

// EnvProduction tells if the app is run with the production flag.
func EnvProduction() bool {
	return *production
}

// Address gives the address on which the app is running.
func Address() string {
	return *address
}

// Run starts the server.
func Run() {
	handler := http.Handler(mux)
	if confI18n != nil {
		confI18n.handleI18n(&handler)
	}
	if !*production {
		handler = gowwwlog.Handle(handler, &gowwwlog.Options{Color: true})
	}
	log.Fatalln(http.ListenAndServe(*address, handler))
}
