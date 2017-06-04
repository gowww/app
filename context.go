package app

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gowww/i18n"
	"github.com/gowww/router"
	"html/template"
	"log"
	"net/http"
)

// A Context contains the data for a handler.
type Context struct {
	Res http.ResponseWriter
	Req *http.Request
}

// Get returns a context value.
func (c *Context) Get(key interface{}) interface{} {
	return c.Req.Context().Value(key)
}

// Set sets a context value.
func (c *Context) Set(key, val interface{}) {
	c.Req = c.Req.WithContext(context.WithValue(c.Req.Context(), key, val))
}

// PathValue returns the value of path parameter.
func (c *Context) PathValue(key string) string {
	return router.Parameter(c.Req, key)
}

// FormValue gets the form value from the request.
func (c *Context) FormValue(key string) string {
	return c.Req.FormValue(key)
}

// HasFormValue checks if the form value exists in the request.
func (c *Context) HasFormValue(key string) bool {
	return c.Req.FormValue(key) != ""
}

// Write writes the response.
func (c *Context) Write(b []byte) (int, error) {
	return c.Res.Write(b)
}

// Text writes the response with a string.
func (c *Context) Text(s string) {
	c.Write([]byte(s))
}

// Textf writes the response with a formatted string.
func (c *Context) Textf(s string, a ...interface{}) {
	fmt.Fprintf(c.Res, s, a...)
}

// Bytes writes the response with a bytes slice.
func (c *Context) Bytes(b []byte) {
	c.Write(b)
}

// Status sets the HTTP status of the response.
func (c *Context) Status(code int) {
	c.Res.WriteHeader(code)
}

// View writes the response with a rendered view.
// This data is always part of the rendering:
//	.	the GlobalViewData
//	.c	the Context
func (c *Context) View(name string, data ...ViewData) {
	d := make(ViewData)
	for k, v := range GlobalViewData {
		d[k] = v
	}
	for _, dt := range data {
		for k, v := range dt {
			d[k] = v
		}
	}
	d["c"] = c
	err := views.ExecuteTemplate(c, name, d)
	if err != nil {
		log.Println(err)
	}
}

// JSON writes the response with a marshalled JSON.
func (c *Context) JSON(v interface{}) {
	enc := json.NewEncoder(c.Res)
	enc.Encode(v)
}

// Redirect redirects the client to the url with status code.
func (c *Context) Redirect(url string, status int) {
	http.Redirect(c.Res, c.Req, url, status)
}

// T returns the translation associated to key, for the client locale.
func (c *Context) T(key string, a ...interface{}) string {
	return i18n.RequestTranslator(c.Req).T(key, a...)
}

// Tn returns the translation associated to key, for the client locale.
// If the translation defines plural forms (zero, one, other), it uses the most appropriate.
// All i18n.TnPlaceholder in the translation are replaced with number n.
// When translation is not found, an empty string is returned.
func (c *Context) Tn(key string, n interface{}, a ...interface{}) string {
	return i18n.RequestTranslator(c.Req).Tn(key, n, a...)
}

// THTML works like T but returns an HTML unescaped translation. An "nl2br" function is applied to the result.
func (c *Context) THTML(key string, a ...interface{}) template.HTML {
	return i18n.RequestTranslator(c.Req).THTML(key, a...)
}

// TnHTML works like Tn but returns an HTML unescaped translation. An "nl2br" function is applied to the result.
func (c *Context) TnHTML(key string, n interface{}, a ...interface{}) template.HTML {
	return i18n.RequestTranslator(c.Req).TnHTML(key, n, a...)
}

// Fmtn returns a formatted number with decimal and thousands marks.
func (c *Context) Fmtn(n interface{}) string {
	return i18n.Fmtn(i18n.RequestTranslator(c.Req).Locale, n)
}

// Push initiates an HTTP/2 server push with an Accept-Encoding header.
// See net/http.Pusher for documentation.
func (c *Context) Push(target string) {
	if pusher, ok := c.Res.(http.Pusher); ok {
		pusher.Push(target, nil)
	}
}

// NotFound responds with the "not found" handler.
func (c *Context) NotFound() {
	if rt.NotFoundHandler != nil {
		rt.NotFoundHandler.ServeHTTP(c.Res, c.Req)
	} else {
		http.NotFound(c.Res, c.Req)
	}
}

// Error logs error and responds with the error handler.
func (c *Context) Error(err error) {
	log.Println(err)
	if errorHandler != nil {
		errorHandler.ServeHTTP(c.Res, c.Req)
	} else {
		http.Error(c.Res, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}