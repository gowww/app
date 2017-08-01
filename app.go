// Package app provides a full featured framework for any web app.
package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gowww/cli"
	"github.com/gowww/compress"
	"github.com/gowww/crypto"
	"github.com/gowww/fatal"
	"github.com/gowww/i18n"
	gowwwlog "github.com/gowww/log"
	"github.com/gowww/router"
	"github.com/gowww/secure"
)

var (
	errorHandler    Handler
	encrypter       crypto.Encrypter
	securityOptions *secure.Options

	address    string
	production bool
	rt         = router.New()
)

func init() {
	cli.String(&address, "a", ":8080", "The address to listen and serve on.")
	cli.Bool(&production, "p", false, "Run the server in production environment.")

	// Set route for static content.
	rt.Get("/"+staticDir+"/", staticHandler)
}

// A Handler handles a request.
type Handler func(*Context)

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h(&Context{Res: w, Req: r})
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

// Secure sets security options.
func Secure(o *secure.Options) {
	if encrypter != nil {
		panic("app: security options set multiple times")
	}
	securityOptions = o
}

// Secret sets the secret key used for encryption.
// The key must be 32 bytes long.
func Secret(key string) {
	if encrypter != nil {
		panic("app: secret key set multiple times")
	}
	var err error
	encrypter, err = crypto.NewEncrypter([]byte(key))
	if err != nil {
		panic(fmt.Errorf("app: %v", err))
	}
}

// Encrypter returns the global encrypter.
func Encrypter() crypto.Encrypter {
	if encrypter == nil {
		panic("app: no secret key set, no encrypter")
	}
	return encrypter
}

// EnvProduction tells if the app is run with the production flag.
// It ensures that flags are parsed so don't use this function before setting your own flags with gowww/cli or they will be ignored.
func EnvProduction() bool {
	if !cli.Parsed() {
		cli.Parse()
	}
	return production
}

// Address gives the address on which the app is running.
// It ensures that flags are parsed so don't use this function before setting your own flags with gowww/cli or they will be ignored.
func Address() string {
	if !cli.Parsed() {
		cli.Parse()
	}
	return address
}

// Run ensures that flags are parsed, sets the middlewares and starts the server.
func Run(mm ...Middleware) {
	if !cli.Parsed() {
		cli.Parse()
	}

	initViews()

	handler := wrapHandler(rt, mm...)
	handler = contextHandle(handler)

	// gowww/secure
	if securityOptions != nil {
		securityOptions.EnvDevelopment = !production
		handler = secure.Handle(handler, securityOptions)
	} else {
		handler = secure.Handle(handler, &secure.Options{EnvDevelopment: !production})
	}

	// gowww/fatal
	if errorHandler != nil {
		handler = fatal.Handle(handler, &fatal.Options{RecoverHandler: errorHandler})
	} else {
		handler = fatal.Handle(handler, &fatal.Options{RecoverHandler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		})})
	}

	// gowww/i18n
	if confI18n.Locales != nil {
		handler = i18n.Handle(handler, confI18n.Locales, confI18n.Fallback, confI18n.Parsers...)
	}

	// gowww/compress
	handler = compress.Handle(handler)

	// gowww/log
	if !production {
		handler = gowwwlog.Handle(handler, &gowwwlog.Options{Color: true})
	}

	// Wait for shut down.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	srv := &http.Server{Addr: address, Handler: handler}
	go func() {
		<-quit
		log.Println("Shutting down...")
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Fatalf("Could not shut down: %v", err)
		}
	}()

	log.Printf("Running on %v", address)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	log.Println("Gracefully shut down")
}
