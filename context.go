package app

import (
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/gowww/i18n"
	"github.com/gowww/router"
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

// Text writes a string to the response.
func (c *Context) Text(s string) {
	c.Write([]byte(s))
}

// Bytes writes a bytes slice to the response.
func (c *Context) Bytes(b []byte) {
	c.Write(b)
}

// Status sets the HTTP status of the response.
func (c *Context) Status(code int) {
	c.Res.WriteHeader(code)
}

// View writes a rendered view to the response.
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

// JSON writes a marshalled JSON to the response.
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

// NotFound only responds with a status.
func (c *Context) NotFound() {
	http.NotFound(c.Res, c.Req) // TODO: Send custom 404 template.
}

func (c *Context) Error(err error) {
	log.Println(err)
	c.Res.WriteHeader(http.StatusInternalServerError)
	c.Res.Write([]byte(http.StatusText(http.StatusInternalServerError))) // TODO: Send custom 500 template.
}
