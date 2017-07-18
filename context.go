package app

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"net/http"

	"golang.org/x/text/language"

	"github.com/gowww/check"
	"github.com/gowww/fatal"
	"github.com/gowww/i18n"
	"github.com/gowww/router"
)

type responseType int

// Response content types.
const (
	HTML responseType = iota
	JSON
)

// A Context contains the data for a handler.
type Context struct {
	Res http.ResponseWriter
	Req *http.Request
}

// contextHandle wraps the router for setting headers and deferring their write.
func contextHandle(h http.Handler) http.Handler {
	return Handler(func(c *Context) {
		cw := &contextWriter{ResponseWriter: c.Res}
		defer func() {
			if cw.status != 0 {
				c.Res.WriteHeader(cw.status)
			}
		}()
		cw.Header().Set("Connection", "keep-alive")
		h.ServeHTTP(cw, c.Req)
	})
}

// logWriter keeps the status code from WriteHeader to allow setting headers after a Context.Status call.
// Required when using Context.Status with Context.JSON, for example.
type contextWriter struct {
	http.ResponseWriter
	status int
}

func (cw *contextWriter) WriteHeader(status int) {
	cw.status = status
}

func (cw *contextWriter) Write(b []byte) (int, error) {
	if cw.status != 0 {
		cw.ResponseWriter.WriteHeader(cw.status)
		cw.status = 0
	}
	return cw.ResponseWriter.Write(b)
}

// CloseNotify implements the http.CloseNotifier interface.
// No channel is returned if CloseNotify is not implemented by an upstream response writer.
func (cw *contextWriter) CloseNotify() <-chan bool {
	n, ok := cw.ResponseWriter.(http.CloseNotifier)
	if !ok {
		return nil
	}
	return n.CloseNotify()
}

// Flush implements the http.Flusher interface.
// Nothing is done if Flush is not implemented by an upstream response writer.
func (cw *contextWriter) Flush() {
	f, ok := cw.ResponseWriter.(http.Flusher)
	if ok {
		f.Flush()
	}
}

// Hijack implements the http.Hijacker interface.
// Error http.ErrNotSupported is returned if Hijack is not implemented by an upstream response writer.
func (cw *contextWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := cw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return h.Hijack()
}

// Push implements the http.Pusher interface.
// http.ErrNotSupported is returned if Push is not implemented by an upstream response writer or not supported by the client.
func (cw *contextWriter) Push(target string, opts *http.PushOptions) error {
	p, ok := cw.ResponseWriter.(http.Pusher)
	if !ok {
		return http.ErrNotSupported
	}
	return p.Push(target, opts)
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
//	.errors	the translated errors map
func (c *Context) View(name string, data ...ViewData) {
	mdata := mergeViewData(data)
	mdata["c"] = c
	switch errs := mdata["errors"].(type) {
	case check.TranslatedErrors:
		break
	case check.Errors:
		mdata["errors"] = c.TErrors(errs)
	default:
		mdata["errors"] = make(check.TranslatedErrors)
	}
	err := views.ExecuteTemplate(c, name, mdata)
	if err != nil {
		c.Panic(err)
	}
}

// JSON writes the response with a marshalled JSON.
// If v has a JSON() interface{} method, it will be used.
func (c *Context) JSON(v interface{}) {
	c.Res.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(c.Res)
	if vjson, ok := v.(interface {
		JSON() interface{}
	}); ok {
		if err := enc.Encode(vjson.JSON()); err != nil {
			c.Panic(err)
		}
		return
	}
	if err := enc.Encode(v); err != nil {
		c.Panic(err)
	}
}

// Check uses a check.Checker to validate request's data and always returns the non-nil errors map.
func (c *Context) Check(checker check.Checker) check.Errors {
	return checker.CheckRequest(c.Req)
}

// TErrors returns translated checking errors.
func (c *Context) TErrors(errs check.Errors) check.TranslatedErrors {
	return errs.T(i18n.RequestTranslator(c.Req))
}

// BadRequest uses a check.Checker to validate request form data, and a view name to execute on fail.
// If you don't provide a view name (empty string), the response will be a JSON errors map.
//
// If the check fails, it sets the status to "400 Bad Request" and returns true, allowing you to exit from the handler.
func (c *Context) BadRequest(checker check.Checker, view string, data ...ViewData) bool {
	errs := c.Check(checker)
	if errs.Empty() {
		return false
	}
	c.Status(http.StatusBadRequest)
	if view == "" {
		c.JSON(errs)
	} else {
		data = append(data, ViewData{"errors": c.TErrors(errs)})
		c.View(view, data...)
	}
	return true
}

// Redirect redirects the client to the url with status code.
func (c *Context) Redirect(url string, status int) {
	http.Redirect(c.Res, c.Req, url, status)
}

// Cookie returns the value of the named cookie.
// If multiple cookies match the given name, only one cookie value will be returned.
// If the secret key is set for app, value will be decrypted before returning.
// If cookie is not found or the decryption fails, an empty string is returned.
func (c *Context) Cookie(name string) string {
	ck, _ := c.Req.Cookie(name)
	if ck == nil {
		return ""
	}
	if encrypter == nil {
		return ck.Value
	}
	v, err := encrypter.DecryptBase64([]byte(ck.Value))
	if err != nil {
		c.DeleteCookie(name)
		return ""
	}
	return string(v)
}

// SetCookie sets a cookie to the response.
// If the secret key is set for app, value will be encrypted.
// If the app is not in a production environment, the "secure" flag will be set to false.
func (c *Context) SetCookie(cookie *http.Cookie) {
	if !production {
		cookie.Secure = false
	}
	if encrypter != nil {
		v, err := encrypter.EncryptBase64([]byte(cookie.Value))
		if err != nil {
			c.Panic(err)
		}
		cookie.Value = string(v)
	}
	http.SetCookie(c.Res, cookie)
}

// DeleteCookie removes a cookie from the client.
func (c *Context) DeleteCookie(name string) {
	http.SetCookie(c.Res, &http.Cookie{Name: name, MaxAge: -1})
}

// translator returns the request translator.
func (c *Context) translator() *i18n.Translator {
	rt := i18n.RequestTranslator(c.Req)
	if rt == nil {
		panic("app: no locales, no translator set")
	}
	return rt
}

// Locale returns the locale used for the client.
func (c *Context) Locale() language.Tag {
	return c.translator().Locale()
}

// T returns the translation associated to key, for the client locale.
func (c *Context) T(key string, a ...interface{}) string {
	return c.translator().T(key, a...)
}

// Tn returns the translation associated to key, for the client locale.
// If the translation defines plural forms (zero, one, other), it uses the most appropriate.
// All i18n.TnPlaceholder in the translation are replaced with number n.
// If translation is not found, an empty string is returned.
func (c *Context) Tn(key string, n int, args ...interface{}) string {
	return c.translator().Tn(key, n, args...)
}

// THTML works like T but returns an HTML unescaped translation. An "nl2br" function is applied to the result.
func (c *Context) THTML(key string, a ...interface{}) template.HTML {
	return c.translator().THTML(key, a...)
}

// TnHTML works like Tn but returns an HTML unescaped translation. An "nl2br" function is applied to the result.
func (c *Context) TnHTML(key string, n int, args ...interface{}) template.HTML {
	return c.translator().TnHTML(key, n, args...)
}

// FmtNumber returns a formatted number with decimal and thousands marks.
func (c *Context) FmtNumber(n interface{}) string {
	return i18n.FmtNumber(c.translator().Locale(), n)
}

// Push initiates an HTTP/2 server push if supported.
// See net/http.Pusher for documentation.
func (c *Context) Push(target string, opts *http.PushOptions) {
	if pusher, ok := c.Res.(http.Pusher); ok {
		pusher.Push(target, opts)
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

// Panic logs error with stack trace and responds with the error handler if set.
func (c *Context) Panic(err error) {
	panic(fmt.Errorf("Failed serving %s: %v", c.Req.RemoteAddr, err))
}

// Error returns the error value stored in request's context after a recovering or a Context.Error call.
func (c *Context) Error() error {
	return fmt.Errorf("%v", fatal.Error(c.Req))
}
