// Package app provides a full featured framework for any web app.
package app

import (
	"flag"
	"github.com/gowww/compress"
	"github.com/gowww/fatal"
	"github.com/gowww/i18n"
	gowwwlog "github.com/gowww/log"
	"github.com/gowww/router"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
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

// A Middleware is a handler that wraps another one.
type Middleware func(http.Handler) http.Handler

// wrapHandler returns handler h wrapped with middlewares mm.
func wrapHandler(h http.Handler, mm ...Middleware) http.Handler {
	for i := len(mm) - 1; i >= 0; i-- {
		h = mm[i](h)
	}
	return h
}

// Route makes a route for method and path.
func Route(method, path string, handler Handler, middlewares ...Middleware) {
	rt.Handle(method, path, wrapHandler(handler, middlewares...))
}

// Get makes a route for GET method.
func Get(path string, handler Handler, middlewares ...Middleware) {
	Route(http.MethodGet, path, handler, middlewares...)
}

// Post makes a route for POST method.
func Post(path string, handler Handler, middlewares ...Middleware) {
	Route(http.MethodPost, path, handler, middlewares...)
}

// Put makes a route for PUT method.
func Put(path string, handler Handler, middlewares ...Middleware) {
	Route(http.MethodPut, path, handler, middlewares...)
}

// Patch makes a route for PATCH method.
func Patch(path string, handler Handler, middlewares ...Middleware) {
	Route(http.MethodPatch, path, handler, middlewares...)
}

// Delete makes a route for DELETE method.
func Delete(path string, handler Handler, middlewares ...Middleware) {
	Route(http.MethodDelete, path, handler, middlewares...)
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
func Run(mm ...Middleware) {
	handler := wrapHandler(rt, mm...)
	if confI18n != nil {
		ll := make(i18n.Locales)
		for lang, trans := range confI18n.Locales {
			ll[lang] = i18n.Translations(trans)
		}
		var pp []i18n.Parser
		for _, parser := range confI18n.Parsers {
			pp = append(pp, i18n.Parser(parser))
		}
		handler = i18n.Handle(handler, ll, confI18n.Fallback, pp...)
	}
	if errorHandler != nil {
		handler = fatal.Handle(handler, &fatal.Options{RecoverHandler: errorHandler})
	} else {
		handler = fatal.Handle(handler, &fatal.Options{RecoverHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		})})
	}
	handler = compress.Handle(handler)
	if !*production {
		handler = gowwwlog.Handle(handler, &gowwwlog.Options{Color: true})
	}
	log.Fatalln(http.ListenAndServe(*address, handler))
}