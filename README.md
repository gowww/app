<p align="center">
	<img src="https://cloud.githubusercontent.com/assets/9503891/26803819/ebb4bf86-4a45-11e7-9f1e-a8c2e2b39717.png" alt="gowww/app">
</p>

Gowww App is a full featured HTTP framework for any web app.  
It greatly increases productivity by providing helpers at all levels while maintaining best performance.

- [Start](#start)
- [Routing](#routing)
  - [Path parameters](#path-parameters)
  - [Groups](#groups)
  - [Errors](#errors)
- [Context](#context)
  - [Request](#request)
  - [Response](#response)
  - [Values](#values)
- [Views](#views)
- [Static files](#static-files)
- [Running](#running)
- [Middlewares](#middlewares)
- [Internationalization](#internationalization)

## Start

1. [Install Go](https://golang.org/doc/install)

2. Install gowww/app:

   ```Shell
   go get https://github.com/gowww/app
   ```

3. Import it in your new app:

   ```Go
   import github.com/gowww/app
   ```

## Routing

You can add a route using [`Route`](https://godoc.org/github.com/gowww/app#Route) and providing a method, a path and a handler:

```Go
app.Route("GET", "/", func(c *app.Context) {
	// Write response for GET /
})
```

There are also shortcuts for common HTTP methods:

```Go
app.Get("/", func(c *app.Context) {
	// Write response for GET /
})

app.Post("/", func(c *app.Context) {
	// Write response for POST /
})

app.Put("/", func(c *app.Context) {
	// Write response for PUT /
})

app.Patch("/", func(c *app.Context) {
	// Write response for PATCH /
})

app.Delete("/", func(c *app.Context) {
	// Write response for DELETE /
})
```

### Path parameters

Use `:` to set a named parameter in the path and matching any value in this path part (between two `/`):

```Go
app.Get("/user/:id", func(c *app.Context) {
	// Write response for GET /user/1 or GET /user/2, etc.
})
```

To access a parameter value, call [`Context.PathValue`](https://godoc.org/github.com/gowww/app#Context.PathValue):

```Go
app.Get("/users/:id", func(c *app.Context) {
	userID := c.PathValue("id")
})
```

Use a trailing `/` to match the beginning of a path and retreive the end of the path under `*`:

```Go
app.Get("/users/:id/files/", func(c *app.Context) {
	userID := c.PathValue("id")
	userFile := c.PathValue("*")
})
```

### Groups

A routing group works like the top-level router but prefixes all subroutes paths:

```Go
api := app.Group("/api")
{
	v1 := app.Group("/v1")
	{
		v1.Get("/user", func(c *app.Context) { // Write response for GET /api/v1/user })
		v1.Get("/item", func(c *app.Context) { // Write response for GET /api/v1/item })
	}

	v2 := app.Group("/v2")
	{
		v2.Get("/user", func(c *app.Context) { // Write response for GET /api/v2/user })
		v2.Get("/item", func(c *app.Context) { // Write response for GET /api/v2/item })
	}
}
```

### Errors

You can set a custom handler for "404 Not Found" error with [`NotFound`](https://godoc.org/github.com/gowww/app#NotFound):

```Go
app.NotFound(func(c *app.Context) {
	// Write response for "404 Not Found"
})
```

The app is also recovered from panics so you can set a custom handler (which is used only when the response is not already written) for "500 Internal Server Error" with [`Error`](https://godoc.org/github.com/gowww/app#NotFound):

```Go
app.Error(func(c *app.Context) {
	// Write response for "500 Internal Server Error"
})
```

## Context

A [`Context`](https://godoc.org/github.com/gowww/app#Context) is always used inside a [`Handler`](https://godoc.org/github.com/gowww/app#Handler).  
It contains the original request and response writer but provides all the necessary helpers to access them:

### Request

Use [`Context.Req`](https://godoc.org/github.com/gowww/app#Context.Req) to access the original request:

```Go
app.Get("/", func(c *app.Context) {
	r := c.Req
})
```

Use [`Context.FormValue`](https://godoc.org/github.com/gowww/app#Context.FormValue) to access a value from URL or body.  
You can also use [`Context.HasFormValue`](https://godoc.org/github.com/gowww/app#Context.HasFormValue) to check its existence:

```Go
app.Get("/", func(c *app.Context) {
	if c.HasFormValue("id") {
		id := c.FormValue("id")
	}
})
```

### Response

Use [`Context.Res`](https://godoc.org/github.com/gowww/app#Context.Res) to access the original response writer:

```Go
app.Get("/", func(c *app.Context) {
	w := c.Res
})
```

Use [`Context.Text`](https://godoc.org/github.com/gowww/app#Context.Text) or [`Context.Bytes`](https://godoc.org/github.com/gowww/app#Context.Bytes) to send a string:

```Go
app.Get("/", func(c *app.Context) {
	c.Text("Hello")
	c.Bytes([]byte("World"))
})
```

Use [`Context.JSON`](https://godoc.org/github.com/gowww/app#Context.JSON) to send a JSON formatted response:

```Go
app.Get("/", func(c *app.Context) {
	c.JSON(map[string]interface{}{
		"id":       c.PathValue("id"),
		"filepath": c.PathValue("*"),
	})
})
```

Use [`Context.Status`](https://godoc.org/github.com/gowww/app#Context.Status) to set the response status code:

```Go
app.Get("/", func(c *app.Context) {
	c.Status(http.StatusCreated)
})
```

Use [`Context.NotFound`](https://godoc.org/github.com/gowww/app#Context.NotFound) to send a "404 Not Found" response:

```Go
app.Get("/", func(c *app.Context) {
	c.NotFound()
})
```

Use [`Context.Error`](https://godoc.org/github.com/gowww/app#Context.Error) to log an error and send a "500 Internal Server Error" response:

```Go
app.Get("/", func(c *app.Context) {
	c.Error("database connection failed")
})
```

Use [`Context.Redirect`](https://godoc.org/github.com/gowww/app#Context.Redirect) to redirect the client:

```Go
app.Get("/old", func(c *app.Context) {
	c.Redirect("/new", http.StatusMovedPermanently)
})
```

Use [`Context.Push`](https://godoc.org/github.com/gowww/app#Context.Push) to initiate an HTTP/2 server push:

```Go
app.Get("/", func(c *app.Context) {
	c.Push("/static/main.css")
})
```

### Values

You can use context values kept inside the context for future usage downstream (like views or subhandlers).

Use [`Context.Set`](https://godoc.org/github.com/gowww/app#Context.Set) to set a value:

```Go
app.Get("/", func(c *app.Context) {
	c.Set("clientCountry", "UK")
})
```

Use [`Context.Get`](https://godoc.org/github.com/gowww/app#Context.Get) to retreive a value:

```Go
app.Get("/", func(c *app.Context) {
	clientCountry := c.Get("clientCountry")
})
```

## Views

Views are standard [Go HTML templates](https://golang.org/pkg/html/template/) and must be stored inside `.gohtml` files within `views` directory.  
They are automatically parsed on running.

Use [`Context.View`](https://godoc.org/github.com/gowww/app#Context.View) to format and send a view template:

```Go
app.Get("/", func(c *app.Context) {
	c.View("home")
})
```

## Static files

Static files must be stored inside the `static` directory.  
They are automatically accessible from the `/static/` path prefix.

## Running

Call [`Run`](https://godoc.org/github.com/gowww/app#Run) at the end of your main function:

```Go
app.Run()
```

By default, your app will listend and serve on `:8080`.  
But you can change this address by using flag `-a` when running your app:

```Shell
./myapp -a :1234
```

## Middlewares

Custom middlewares can be used if they are compatible with standard interface [`net/http.Handler`](https://golang.org/pkg/net/http/#Handler).  
They can be set for:

- The entire app:

  ```Go
  app.Run(hand1, hand2, hand3)
  ```

- A group:

  ```Go
  api := app.Group("/api", hand1, hand2, hand3)
  ```

- A single route:

  ```Go
  api := app.Get("/", func(c *app.Context) {
	  // Write response for GET /
  }, hand1, hand2, hand3)
  ```

First handler wraps the second and so on.

## Internationalization

To have translations accessible all over your app, use [`Localize`](https://godoc.org/github.com/gowww/app#Localize) with your locales, their translations (a map of string to string) and the default locale:

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
```

In your views, use function `t` (or its variants: `tn`, `thtml`, `tnhtml`) to get a translation:

```HTML
<h>{{t .c "hello"}}</h1>
```

<p align="center">
	<br><br>
	<a href="https://godoc.org/github.com/gowww/app"><img src="https://godoc.org/github.com/gowww/app?status.svg" alt="GoDoc"></a>
	<a href="https://travis-ci.org/gowww/app"><img src="https://travis-ci.org/gowww/app.svg?branch=master" alt="Build"></a>
	<a href="https://coveralls.io/github/gowww/app?branch=master"><img src="https://coveralls.io/repos/github/gowww/app/badge.svg?branch=master" alt="Coverage"></a>
	<a href="https://goreportcard.com/report/github.com/gowww/app"><img src="https://goreportcard.com/badge/github.com/gowww/app" alt="Go Report"></a>
	<br><br>
</p>