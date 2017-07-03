// Package app provides a full featured framework for any web app.
package app

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gowww/compress"
	"github.com/gowww/fatal"
	"github.com/gowww/i18n"
	gowwwlog "github.com/gowww/log"
	"github.com/gowww/router"
	"github.com/gowww/secure"
)

var (
	errorHandler Handler
	encrypter    secure.Encrypter

	address    = flag.String("a", ":8080", "The address to listen and serve on.")
	production = flag.Bool("p", false, "Run the server in production environment.")
	rt         = router.New()
)

func init() {
	flag.Parse()

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

// Secret sets the secret key used for encryption.
// The key must be 32 bytes long.
func Secret(key string) {
	if encrypter != nil {
		panic("app: secret key set multiple times")
	}
	var err error
	encrypter, err = secure.NewEncrypter(key)
	if err != nil {
		panic(fmt.Errorf("app: %v", err))
	}
}

// Encrypter returns the global encrypter.
func Encrypter() secure.Encrypter {
	if encrypter == nil {
		panic("app: no secret key set, no encrypter")
	}
	return encrypter
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
	srv := &http.Server{Addr: *address, Handler: handler}
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
