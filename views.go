package app

import (
	"html/template"
	"os"
	"strings"

	"github.com/gowww/static"
	"github.com/gowww/view"
)

const (
	viewsDir  = "views"
	staticDir = "static"
)

var (
	staticHandler = static.Handle("/"+staticDir+"/", staticDir)
	views         = view.New()
)

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

	GlobalViewData(ViewData{
		"envProduction": production,
	})

	GlobalViewFuncs(ViewFuncs{
		"asset": func(path string) string {
			return staticHandler.Hash(path)
		},
		"script": func(src string) template.HTML {
			return view.HelperScript(staticHandler.Hash("scripts/" + strings.TrimPrefix(src, "/")))
		},
		"style": func(href string) template.HTML {
			return view.HelperStyle(staticHandler.Hash("styles/" + strings.TrimPrefix(href, "/")))
		},
	})

	views.ParseDir(viewsDir)
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
