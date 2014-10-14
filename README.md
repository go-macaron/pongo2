pongo2
======

Middleware pongo2 is a [pongo2](https://github.com/flosch/pongo2) template engine support for [Macaron](https://github.com/Unknwon/macaron).

[API Reference](https://gowalker.org/github.com/macaron-contrib/pongo2)

### Installation

	go get github.com/macaron-contrib/pongo2
	
## Usage

```go
// main.go
package main

import (
  "github.com/Unknwon/macaron"
  "github.com/macaron-contrib/pongo2"
)

func main() {
  m := macaron.Classic()
  m.Use(pongo2.Pongoer())

  m.Get("/", func(ctx *macaron.Context) {
  	ctx.Data["Name"] = "joe"
    ctx.HTML(200, "hello")
  })

  m.Run()
}
```

```html
<!-- templates/hello.tmpl -->
<h2>Hello {{Name}}!</h2>
```

## Options

`pongo2.Options` comes with a variety of configuration options:

```go
// ...
m.Use(pongo2.Pongoer(pongo2.Options{
  Directory: "templates", // Specify what path to load the templates from.
  Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
  Charset: "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
  IndentJSON: true, // Output human readable JSON
  IndentXML: true, // Output human readable XML
  HTMLContentType: "application/xhtml+xml", // Output XHTML content type instead of default "text/html"
}))
// ...
```

## License

This project is under Apache v2 License. See the [LICENSE](LICENSE) file for the full license text.
