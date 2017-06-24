// Package app provides a full featured framework for any web app.
package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/alecthomas/kingpin"
	"github.com/gowww/compress"
	"github.com/gowww/fatal"
	"github.com/gowww/i18n"
	gowwwlog "github.com/gowww/log"
	"github.com/gowww/router"
)

var (
	errorHandler Handler

	address    = kingpin.Flag("address", "The address to listen and serve on.").Default(":8080").Short('a').TCP()
	production = kingpin.Flag("production", "Run the server in production environment.").Short('p').Bool()
	rt         = router.New()
)

func init() {
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	// Serve static content.
	rt.Get("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
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
//
// Using Context.Error, you can retrieve the error value stored in request's context during recovering.
func Error(handler Handler) {
	if errorHandler != nil {
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
	return (*address).String()
}

// Run starts the server.
func Run(mm ...Middleware) {
	initViews()

	handler := wrapHandler(rt, mm...)
	handler = contextHandle(handler)

	// gowww/fatal
	if errorHandler != nil {
		handler = fatal.Handle(handler, &fatal.Options{RecoverHandler: errorHandler})
	} else {
		handler = fatal.Handle(handler, &fatal.Options{RecoverHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		})})
	}

	// gowww/i18n
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

	// gowww/compress
	handler = compress.Handle(handler)

	// gowww/log
	if !*production {
		handler = gowwwlog.Handle(handler, &gowwwlog.Options{Color: true})
	}

	// Wait for shut down.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	srv := &http.Server{Addr: (*address).String(), Handler: handler}
	go func() {
		<-quit
		log.Println("Shutting down...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not shut down: %v", err)
		}
	}()

	log.Printf("Running on %v", *address)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	log.Println("Gracefully shut down")
}
