# [![gowww](https://avatars.githubusercontent.com/u/18078923?s=20)](https://github.com/gowww) app [![GoDoc](https://godoc.org/github.com/gowww/app?status.svg)](https://godoc.org/github.com/gowww/app) [![Build](https://travis-ci.org/gowww/app.svg?branch=master)](https://travis-ci.org/gowww/app) [![Coverage](https://coveralls.io/repos/github/gowww/app/badge.svg?branch=master)](https://coveralls.io/github/gowww/app?branch=master) [![Go Report](https://goreportcard.com/badge/github.com/gowww/app)](https://goreportcard.com/report/github.com/gowww/app)

Package [app](https://godoc.org/github.com/gowww/app) provides a full featured framework for any web app.

## Example

```Go
var locales = app.Locales{
	language.English: {
		"hello": "Hello!",
	},
	language.French: {
		"hello": "Bonjour !",
	},
}

app.Localize(locales, language.English)

app.Route("/", func(c *app.Context) {
        c.View("home")
})

app.Route("/user", func(c *app.Context) {
        c.Status(http.StatusCreated)
        c.JSON(map[string]interface{}{
                "id":  1,
                "name": "White",
        })
})

if !app.EnvProduction() {
        log.Printf("developing app on %s\n", app.Address())
}

app.Run()
```
