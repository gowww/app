package app_test

import (
	"log"
	"net/http"

	"github.com/gowww/app"
	"github.com/gowww/check"
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
		log.Printf("Developing app on %s\n", app.Address())
	}

	app.Run()
}

func ExampleBadRequest() {
	userChecker := check.Checker{
		"email": {check.Required, check.Email},
		"phone": {check.Phone},
	}

	app.Post("/users", func(c *app.Context) {
		if c.BadRequest(userChecker) {
			return
		}
		c.Status(http.StatusCreated)
	})
}
