package app

import (
	"net/http"
)

// RouterGroup contains the first path part for a routes group.
type RouterGroup struct {
	path string
}

// Group initiates a routing group.
func Group(path string) *RouterGroup {
	return &RouterGroup{path}
}

// Group contains the first path part for a routes subgroup.
func (rg *RouterGroup) Group(path string) *RouterGroup {
	return &RouterGroup{rg.path + path}
}

// Route makes a route for method and path.
func (rg *RouterGroup) Route(method, path string, handler Handler) {
	if path == "/" {
		path = ""
	}
	rt.Handle(method, rg.path+path, handler)
}

// Get makes a route for GET method.
func (rg *RouterGroup) Get(path string, handler Handler) {
	rg.Route(http.MethodGet, path, handler)
}

// Post makes a route for POST method.
func (rg *RouterGroup) Post(path string, handler Handler) {
	rg.Route(http.MethodPost, path, handler)
}

// Put makes a route for PUT method.
func (rg *RouterGroup) Put(path string, handler Handler) {
	rg.Route(http.MethodPut, path, handler)
}

// Patch makes a route for PATCH method.
func (rg *RouterGroup) Patch(path string, handler Handler) {
	rg.Route(http.MethodPatch, path, handler)
}

// Delete makes a route for DELETE method.
func (rg *RouterGroup) Delete(path string, handler Handler) {
	rg.Route(http.MethodDelete, path, handler)
}
