package app

import (
	"html/template"
	"os"

	"github.com/gowww/view"
)

const viewsDir = "views"

var views = view.New()

// ViewData represents data for a view rendering.
type ViewData map[string]interface{}

// ViewFuncs is a map of functions passed to all view renderings.
type ViewFuncs map[string]interface{}

// GlobalViewData adds global data for view templates.
func GlobalViewData(data ViewData) {
	views.Data(view.Data(data))
}

// GlobalViewFuncs adds functions for view templates.
func GlobalViewFuncs(funcs ViewFuncs) {
	views.Funcs(view.Funcs(funcs))
}

func initViews() {
	if _, err := os.Stat(viewsDir); err != nil { // viewsDir not found: nothing to parse.
		return
	}

	view.AllHelpers["envproduction"] = func() bool {
		return EnvProduction()
	}
	view.AllHelpers["fmtnumber"] = func(c *Context, n interface{}) string {
		return c.FmtNumber(n)
	}
	view.AllHelpers["t"] = func(c *Context, key string, a ...interface{}) string {
		return c.T(key, a...)
	}
	view.AllHelpers["thtml"] = func(c *Context, key string, a ...interface{}) template.HTML {
		return c.THTML(key, a...)
	}
	view.AllHelpers["tn"] = func(c *Context, key string, n int, args ...interface{}) string {
		return c.Tn(key, n, args...)
	}
	view.AllHelpers["tnhtml"] = func(c *Context, key string, n int, args ...interface{}) template.HTML {
		return c.TnHTML(key, n, args...)
	}

	views.Funcs(view.AllHelpers).ParseDir(viewsDir)
}

func mergeViewData(dd []ViewData) view.Data {
	data := make(view.Data, len(dd))
	for _, d := range dd {
		for k, v := range d {
			data[k] = v
		}
	}
	return data
}
