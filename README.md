pongo2 [![Build Status](https://drone.io/github.com/macaron-contrib/pongo2/status.png)](https://drone.io/github.com/macaron-contrib/pongo2/latest) [![](http://gocover.io/_badge/github.com/macaron-contrib/pongo2)](http://gocover.io/github.com/macaron-contrib/pongo2)
======

Middleware pongo2 is a [pongo2](https://github.com/flosch/pongo2).**v3** template engine support for [Macaron](https://github.com/Unknwon/macaron).

[API Reference](https://gowalker.org/github.com/macaron-contrib/pongo2)

[Full Documentation](http://macaron.gogs.io/docs/middlewares/templating.html)

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
        ctx.Data["Name"] = "jeremy"
        ctx.HTML(200, "hello") // 200 is the response code.
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
    // Directory to load templates. Default is "templates".
    Directory: "templates",
    // Extensions to parse template files from. Defaults are [".tmpl", ".html"].
    Extensions: []string{".tmpl", ".html"},
    // Appends the given charset to the Content-Type header. Default is "UTF-8".
    Charset: "UTF-8",
    // Outputs human readable JSON. Default is false.
    IndentJSON: true,
    // Outputs human readable XML. Default is false.
    IndentXML: true,
    // Allows changing of output to XHTML instead of HTML. Default is "text/html".
    HTMLContentType: "text/html",
}))  
// ...
```

## License

This project is under Apache v2 License. See the [LICENSE](LICENSE) file for the full license text.
