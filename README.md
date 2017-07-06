<p align="center">
	<img src="https://cloud.githubusercontent.com/assets/9503891/26803819/ebb4bf86-4a45-11e7-9f1e-a8c2e2b39717.png" alt="gowww/app">
</p>

**gowww/app** is a full featured HTTP framework for any web app.  
It greatly increases productivity by providing helpers at all levels while maintaining best performance.

- [Start](#start)
- [Routing](#routing)
	- [Path parameters](#path-parameters)
		- [Named](#named)
		- [Regular expressions](#regular-expressions)
		- [Wildcard](#wildcard)
	- [Groups](#groups)
	- [Errors](#errors)
- [Context](#context)
	- [Request](#request)
	- [Response](#response)
	- [Values](#values)
- [Internationalization](#internationalization)
- [Views](#views)
	- [Data](#data)
	- [Functions](#functions)
	- [Built-in](#built-in)
- [Static files](#static-files)
- [Running](#running)
- [Middlewares](#middlewares)

## Start

1. [Install Go](https://golang.org/doc/install)

2. Install gowww/app:

	```Shell
	go get github.com/gowww/app
	```

3. Import it in your new app:

	```Go
	import github.com/gowww/app
	```

## Routing

There are methods for common HTTP methods:

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

#### Named

A named parameter begins with `:` and matches any value until the next `/` in path.

To retrieve the value, ask [`Context.PathValue`](https://godoc.org/github.com/gowww/app#Context.PathValue).  
It will return the value as a string (empty if the parameter doesn't exist).

Example, with a parameter `id`:

```Go
app.Get("/users/:id", func(c *app.Context) {
	id := c.PathValue("id")
	fmt.Fprintf(w, "Page of user #%s", id)
}))
```

#### Regular expressions

If a parameter must match an exact pattern (digits only, for example), you can also set a [regular expression](https://golang.org/pkg/regexp/syntax) constraint just after the parameter name and another `:`:

```Go
app.Get(`/users/:id:^\d+$`, func(c *app.Context) {
	id := c.PathValue("id")
	fmt.Fprintf(w, "Page of user #%s", id)
}))
```

If you don't need to retrieve the parameter value but only use a regular expression, you can omit the parameter name.

#### Wildcard

A trailing slash behaves like a wildcard by matching the beginning of the request path and keeping the rest as a parameter value, under `*`:

```Go
rt.Get("/files/", func(c *app.Context) {
	filepath := c.PathValue("*")
	fmt.Fprintf(w, "Get file %s", filepath)
}))
```

For more details, see [gowww/router](https://github.com/gowww/router).

### Groups

A routing group works like the top-level router but prefixes all subroute paths:

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

You can set a custom "not found" handler with [`NotFound`](https://godoc.org/github.com/gowww/app#NotFound):

```Go
app.NotFound(func(c *app.Context) {
	c.Status(http.StatusNotFound)
	c.View("notFound")
})
```

The app is also recovered from panics so you can set a custom "serving error" handler (which is used only when the response is not already written) with [`Error`](https://godoc.org/github.com/gowww/app#NotFound) and retrive the recovered error value with [`Context.Error`](https://godoc.org/github.com/gowww/app#Context.Error):

```Go
app.Error(func(c *app.Context) {
	c.Status(http.StatusInternalServerError)
	if c.Error() == ErrCannotOpenFile" {
		c.View("errorStorage")
		return
	}
	c.View("error")
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
app.Get(`/users/:id:^\d+$/files/`, func(c *app.Context) {
	c.JSON(map[string]interface{}{
		"userID":   c.PathValue("id"),
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

Use [`Context.NotFound`](https://godoc.org/github.com/gowww/app#Context.NotFound) to send a "not found" response:

```Go
app.Get("/", func(c *app.Context) {
	c.NotFound()
})
```

Use [`Context.Panic`](https://godoc.org/github.com/gowww/app#Context.Panic) to log an error and send a "serving error" response:

```Go
app.Get("/", func(c *app.Context) {
	c.Panic("database connection failed")
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
	c.Push("/static/main.css", nil)
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

Use [`Context.Get`](https://godoc.org/github.com/gowww/app#Context.Get) to retrieve a value:

```Go
app.Get("/", func(c *app.Context) {
	clientCountry := c.Get("clientCountry")
})
```

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
<h1>{{t .c "hello"}}</h1>
```

For more details, see [gowww/i18n](https://github.com/gowww/i18n).

## Views

Views are standard [Go HTML templates](https://golang.org/pkg/html/template/) and must be stored inside the `views` directory, within `.gohtml` files.  
They are automatically parsed during launch.

Use [`Context.View`](https://godoc.org/github.com/gowww/app#Context.View) to send a view:

```Go
app.Get("/", func(c *app.Context) {
	c.View("home")
})
```

### Data

Use a [`ViewData`](https://godoc.org/github.com/gowww/app#ViewData) map to pass data to a view.  
Note that the context is automatically stored in the view data under key `c`.

You can also use [`GlobalViewData`](https://godoc.org/github.com/gowww/app#GlobalViewData) to set data for all views:


```Go
app.GlobalViewData(app.ViewData{
	"appName": "My app",
})

app.Get("/", func(c *app.Context) {
	user := &User{
		ID:   1,
		Name: "John Doe",
	}
	c.View("home", app.ViewData{
		"user": user,
	})
})
```

In *views/home.gohtml*:

```HTML
{{define "home"}}
	<h1>Hello {{.user.Name}} ({{.c.Req.RemoteAddr}}) and welcome on {{.appName}}!</h1>
{{end}}
```

### Functions

Use [`GlobalViewFuncs`](https://godoc.org/github.com/gowww/app#GlobalViewFuncs) to set functions for all views:

```Go
app.GlobalViewFuncs(app.ViewFuncs{
	"pathescape": url.PathEscape,
})

app.Get("/posts/new", func(c *app.Context) {
	c.View("postsNew")
})
```

In *views/posts.gohtml*:

```HTML
{{define "postsNew"}}
	<a href="/sign-in?return-to={{pathescape "/posts/new"}}">Sign in</a>
{{end}}
```

#### Built-in

In addition to the functions provided by the standard [template](https://golang.org/pkg/text/template/#hdr-Functions) package, these function are also available out of the box:

| Function        | Description                                                                                                                      | Usage                                              |
------------------|----------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------|
| `envproduction` | Tells if the app is run with the production flag.                                                                                | `{{if envproduction}}Live{{else}}Testing{{end}}`   |
| `fmtn`          | Returns a formatted number with decimal and thousands marks.                                                                     | `{{fmtn 123456.123456}}`                           |
| `googlefonts`   | Sets an HTML link to [Google Fonts](https://fonts.google.com)'s stylesheet of the given font(s).                                 | `{{googlefonts "Open+Sans:400,700\|Spectral"}}`    |
| `nl2br`         | Converts `\n` to HTML `<br>`.                                                                                                    | `{{nl2br "line one\nline two"}}`                   |
| `safehtml`      | Prevents string to be escaped. Be careful.                                                                                       | `{{safehtml "<strong>word</strong>"}}`             |
| `scripts`       | Sets HTML script tags for the given script sources.                                                                              | `{{scripts "/static/main.js" "/static/user.js"}}`  |
| `styles`        | Sets HTML link tags for the given stylesheets.                                                                                   | `{{styles "/static/main.css" "/static/user.css"}}` |
| `t`             | Returns the translation associated to key, for the client locale.                                                                | `{{t .c "hello"}}`                                 |
| `thtml`         | Works like `t` but returns an HTML unescaped translation. `nl2br` is applied to the result.                                      | `{{t .c "hello"}}`                                 |
| `tn`            | Works like `t` with plural variations (zero, one, other). See [`Context.Tn`](https://godoc.org/github.com/gowww/app#Context.Tn). | `{{tn .c "item" 12}}`                              |
| `tnhtml`        | Works like `tn` + `thml`. See [`Context.TnHTML`](https://godoc.org/github.com/gowww/app#Context.TnHTML).                         | `{{tnhtml .c "item" 12}}`                          |

## Static files

Static files must be stored inside the `static` directory.  
They are automatically accessible from the `/static/` path prefix.

## Running

Call [`Run`](https://godoc.org/github.com/gowww/app#Run) at the end of your main function:

```Go
app.Run()
```

By default, your app will listen and serve on `:8080`.  
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

<p align="center">
<br><br>
<a href="https://godoc.org/github.com/gowww/app"><img src="https://godoc.org/github.com/gowww/app?status.svg" alt="GoDoc"></a>
<a href="https://travis-ci.org/gowww/app"><img src="https://travis-ci.org/gowww/app.svg?branch=master" alt="Build"></a>
<a href="https://coveralls.io/github/gowww/app?branch=master"><img src="https://coveralls.io/repos/github/gowww/app/badge.svg?branch=master" alt="Coverage"></a>
<a href="https://goreportcard.com/report/github.com/gowww/app"><img src="https://goreportcard.com/badge/github.com/gowww/app" alt="Go Report"></a>
<br><br>
</p>
