// Package app provides a full featured framework for any web app.
package app

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	gowwwlog "github.com/gowww/log"
	"github.com/gowww/router"
)

var (
	address      = flag.String("a", ":8080", "the address to listen and serving on")
	production   = flag.Bool("p", false, "run the server in production environment")
	rt           = router.New()
	errorHandler Handler
)

func init() {
	flag.Parse()

	// Serve static content
	rt.Get("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

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

// Route makes a route for method and path.
func Route(method, path string, handler Handler) {
	rt.Handle(method, path, handler)
}

// Get makes a route for GET method.
func Get(path string, handler Handler) {
	rt.Get(path, handler)
}

// Post makes a route for POST method.
func Post(path string, handler Handler) {
	rt.Post(path, handler)
}

// Put makes a route for PUT method.
func Put(path string, handler Handler) {
	rt.Put(path, handler)
}

// Patch makes a route for PATCH method.
func Patch(path string, handler Handler) {
	rt.Patch(path, handler)
}

// Delete makes a route for DELETE method.
func Delete(path string, handler Handler) {
	rt.Delete(path, handler)
}

// NotFound registers the "not found" handler.
func NotFound(handler Handler) {
	if rt.NotFoundHandler != nil {
		panic(`app: "not found" handler set multiple times`)
	}
	rt.NotFoundHandler = handler
}

// Error registers the "internal error" handler.
func Error(handler Handler) {
	if rt.NotFoundHandler != nil {
		panic(`app: "internal error" handler set multiple times`)
	}
	errorHandler = handler
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
	handler := http.Handler(rt)
	if confI18n != nil {
		confI18n.handleI18n(&handler)
	}
	if !*production {
		handler = gowwwlog.Handle(handler, &gowwwlog.Options{Color: true})
	}
	log.Fatalln(http.ListenAndServe(*address, handler))
}
