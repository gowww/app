package app

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

var (
	views           *template.Template
	globalViewData  ViewData
	globalViewFuncs ViewFuncs
)

// ViewData represents data for a view rendering.
type ViewData map[string]interface{}

// ViewFuncs is a map of functions passed to all view renderings.
type ViewFuncs map[string]interface{}

// GlobalViewFuncs sets a map of functions passed to all view renderings.
func GlobalViewFuncs(funcs ViewFuncs) {
	if globalViewFuncs != nil {
		panic(`app: view functions set multiple times`)
	}
	globalViewFuncs = funcs
}

// GlobalViewData sets a map of data passed to all view renderings.
func GlobalViewData(data ViewData) {
	if globalViewData != nil {
		panic(`app: view data set multiple times`)
	}
	globalViewData = data
}

func parseViews() {
	ff := template.FuncMap{
		"t": func(c *Context, key string, a ...interface{}) string {
			return c.T(key, a...)
		},
		"tn": func(c *Context, key string, n interface{}, a ...interface{}) string {
			return c.Tn(key, n, a...)
		},
		"thtml": func(c *Context, key string, a ...interface{}) template.HTML {
			return c.THTML(key, a...)
		},
		"tnhtml": func(c *Context, key string, n interface{}, a ...interface{}) template.HTML {
			return c.TnHTML(key, n, a...)
		},
		"fmtn": func(c *Context, n interface{}) string {
			return c.Fmtn(n)
		},
		"safehtml": func(s string) template.HTML {
			return template.HTML(s)
		},
		"nl2br": func(s string) template.HTML {
			return template.HTML(strings.Replace(s, "\n", "<br>", -1))
		},
		"styles": viewFuncStyles,
		"scripts": func(scripts ...string) (h template.HTML) {
			for _, script := range scripts {
				h += template.HTML(`<script src="` + script + `"></script>`)
			}
			return
		},
		"googlefonts": func(fonts ...string) template.HTML {
			return viewFuncStyles("https://fonts.googleapis.com/css?family=" + strings.Join(fonts, "|"))
		},
		"copyright": func(name string) string {
			return fmt.Sprintf("© %s %d", name, time.Now().Year())
		},
		"envproduction": func(name string) bool {
			return EnvProduction()
		},
	}

	for n, f := range globalViewFuncs {
		if _, ok := ff[n]; ok {
			panic(`app: view function "` + n + `" already exists`)
		}
		ff[n] = f
	}

	views = template.Must(template.New("main").Funcs(ff).ParseGlob("views/*.gohtml"))
}

func viewFuncStyles(styles ...string) (h template.HTML) {
	for _, style := range styles {
		h += template.HTML(`<link rel="stylesheet" href="` + style + `">`)
	}
	return
}
