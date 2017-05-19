package app_test

import (
	"log"
	"net/http"

	"github.com/gowww/app"
	"golang.org/x/text/language"
)

func Example() {
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
		c.View("")
	})

	app.Route("/user", func(c *app.Context) {
		c.Status(http.StatusCreated)
		c.JSON(map[string]interface{}{
			"id":   1,
			"name": "White",
		})
	})

	if !app.EnvProduction() {
		log.Printf("developing app on %s\n", app.Address())
	}

	app.Run()
}
