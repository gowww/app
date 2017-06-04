package app_test

import (
	"github.com/gowww/app"
	"golang.org/x/text/language"
	"log"
	"net/http"
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

	v1 := app.Group("/v1")
	{
		v1.Get("/user", func(c *app.Context) { c.Text("User for V1") })
		v1.Get("/item", func(c *app.Context) { c.Text("Item for V1") })
	}

	v2 := app.Group("/v2")
	{
		v2.Get("/user", func(c *app.Context) { c.Text("User for V2") })
		v2.Get("/item", func(c *app.Context) { c.Text("Item for V2") })
	}

	if !app.EnvProduction() {
		log.Printf("developing app on %s\n", app.Address())
	}

	app.Run()
}
