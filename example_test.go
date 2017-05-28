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

	app.Get("/", func(c *app.Context) {
		c.View("home")
	})

	app.Post("/user/:id/files/", func(c *app.Context) {
		c.Status(http.StatusCreated)
		c.JSON(map[string]interface{}{
			"id":       c.PathValue("id"),
			"filepath": c.PathValue("*"),
		})
	})

	if !app.EnvProduction() {
		log.Printf("developing app on %s\n", app.Address())
	}

	app.Run()
}
