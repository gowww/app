package app

import "net/http"

// RouterGroup contains the first path part for a routes group.
type RouterGroup struct {
	path        string
	middlewares []Middleware
}

// Group initiates a routing group.
func Group(path string, middlewares ...Middleware) *RouterGroup {
	return &RouterGroup{path, middlewares}
}

// Group contains the first path part for a routes subgroup.
func (rg *RouterGroup) Group(path string, middlewares ...Middleware) *RouterGroup {
	return &RouterGroup{rg.path + path, middlewares}
}

// Route makes a route for method and path.
func (rg *RouterGroup) Route(method, path string, handler Handler, middlewares ...Middleware) {
	if path == "/" {
		path = ""
	}
	rt.Handle(method, rg.path+path, wrapHandler(wrapHandler(handler, middlewares...), rg.middlewares...))
}

// Get makes a route for GET method.
func (rg *RouterGroup) Get(path string, handler Handler, middlewares ...Middleware) {
	rg.Route(http.MethodGet, path, handler, middlewares...)
}

// Post makes a route for POST method.
func (rg *RouterGroup) Post(path string, handler Handler, middlewares ...Middleware) {
	rg.Route(http.MethodPost, path, handler, middlewares...)
}

// Put makes a route for PUT method.
func (rg *RouterGroup) Put(path string, handler Handler, middlewares ...Middleware) {
	rg.Route(http.MethodPut, path, handler, middlewares...)
}

// Patch makes a route for PATCH method.
func (rg *RouterGroup) Patch(path string, handler Handler, middlewares ...Middleware) {
	rg.Route(http.MethodPatch, path, handler, middlewares...)
}

// Delete makes a route for DELETE method.
func (rg *RouterGroup) Delete(path string, handler Handler, middlewares ...Middleware) {
	rg.Route(http.MethodDelete, path, handler, middlewares...)
}
