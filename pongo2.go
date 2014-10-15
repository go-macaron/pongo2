// Copyright 2014 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

// Package pongo2 is a middleware that provides pongo2 template engine of Macaron.
package pongo2

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Unknwon/macaron"
	"github.com/flosch/pongo2"
)

const (
	ContentType    = "Content-Type"
	ContentLength  = "Content-Length"
	ContentBinary  = "application/octet-stream"
	ContentJSON    = "application/json"
	ContentHTML    = "text/html"
	ContentXHTML   = "application/xhtml+xml"
	ContentXML     = "text/xml"
	defaultCharset = "UTF-8"
)

var (
	renderOpt Options
	tplMap    map[string]*pongo2.Template
	lock      sync.RWMutex // Go map is not safe.
)

func prepareCharset(charset string) string {
	if len(charset) != 0 {
		return "; charset=" + charset
	}

	return "; charset=" + defaultCharset
}

func getExt(s string) string {
	if strings.Index(s, ".") == -1 {
		return ""
	}
	return "." + strings.Join(strings.Split(s, ".")[1:], ".")
}

func compile(options Options) {
	lock.Lock()
	defer lock.Unlock()

	dir := options.Directory
	tplMap = make(map[string]*pongo2.Template)

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		r, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}

		ext := getExt(r)

		for _, extension := range options.Extensions {
			if ext == extension {
				name := (r[0 : len(r)-len(ext)])
				// Bomb out if parse fails. We don't want any silent server starts.
				t, err := pongo2.FromFile(path)
				if err != nil {
					panic(fmt.Errorf("\"%s\": %v", path, err))
				}
				tplMap[strings.Replace(name, "\\", "/", -1)] = t
				break
			}
		}

		return nil
	})
}

// Options represents a struct for specifying configuration options for the Render middleware.
type Options struct {
	// Directory to load templates. Default is "templates"
	Directory string
	// Extensions to parse template files from. Defaults to [".tmpl", ".html"]
	Extensions []string
	// Funcs is a slice of FuncMaps to apply to the template upon compilation. This is useful for helper functions. Defaults to [].
	// Funcs []template.FuncMap
	// Appends the given charset to the Content-Type header. Default is "UTF-8".
	Charset string
	// Outputs human readable JSON
	IndentJSON bool
	// Outputs human readable XML
	IndentXML bool
	// Prefixes the JSON output with the given bytes.
	PrefixJSON []byte
	// Prefixes the XML output with the given bytes.
	PrefixXML []byte
	// Allows changing of output to XHTML instead of HTML. Default is "text/html"
	HTMLContentType string
}

func prepareOptions(options []Options) Options {
	var opt Options
	if len(options) > 0 {
		opt = options[0]
	}

	// Defaults
	if len(opt.Directory) == 0 {
		opt.Directory = "templates"
	}
	if len(opt.Extensions) == 0 {
		opt.Extensions = []string{".tmpl", ".html"}
	}
	if len(opt.HTMLContentType) == 0 {
		opt.HTMLContentType = ContentHTML
	}

	return opt
}

// Pongoer is a Middleware that maps a macaron.Render service into the Macaron handler chain.
// An single variadic pongo2.Options struct can be optionally provided to configure
// HTML rendering. The default directory for templates is "templates" and the default
// file extension is ".tmpl" and ".html".
//
// If MACARON_ENV is set to "" or "development" then templates will be recompiled on every request. For more performance, set the
// MACARON_ENV environment variable to "production".
func Pongoer(options ...Options) macaron.Handler {
	renderOpt = prepareOptions(options)
	cs := prepareCharset(renderOpt.Charset)
	compile(renderOpt)

	return func(ctx *macaron.Context, rw http.ResponseWriter, req *http.Request) {
		if macaron.Env == macaron.DEV {
			compile(renderOpt)
		}
		r := &render{
			TplRender: &macaron.TplRender{
				ResponseWriter: rw,
				Req:            req,
				Opt: macaron.RenderOptions{
					IndentJSON: renderOpt.IndentJSON,
					IndentXML:  renderOpt.IndentXML,
					PrefixJSON: renderOpt.PrefixJSON,
					PrefixXML:  renderOpt.PrefixXML,
				},
				CompiledCharset: cs,
			},
			ResponseWriter:  rw,
			opt:             renderOpt,
			compiledCharset: cs,
		}
		ctx.Render = r
		ctx.MapTo(r, (*macaron.Render)(nil))
	}
}

type render struct {
	*macaron.TplRender
	http.ResponseWriter
	opt             Options
	compiledCharset string

	startTime time.Time
}

func data2Context(data interface{}) pongo2.Context {
	return pongo2.Context(data.(map[string]interface{}))
}

func (r *render) HTML(status int, name string, data interface{}, _ ...macaron.HTMLOptions) {
	r.startTime = time.Now()

	lock.RLock()
	defer lock.RUnlock()

	t := tplMap[name]
	if t == nil {
		http.Error(r, "pongo2: \""+name+"\" is undefined", http.StatusInternalServerError)
		return
	}

	r.Header().Set(ContentType, r.opt.HTMLContentType+r.compiledCharset)
	r.WriteHeader(status)
	if err := t.ExecuteWriter(data2Context(data), r); err != nil {
		http.Error(r, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (r *render) HTMLString(name string, data interface{}, _ ...macaron.HTMLOptions) (string, error) {
	lock.RLock()
	defer lock.RUnlock()

	t := tplMap[name]
	if t == nil {
		http.Error(r, "pongo2: \""+name+"\" is undefined", http.StatusInternalServerError)
		return "", nil
	}

	out, err := t.Execute(data2Context(data))
	if err != nil {
		return "", err
	}

	return out, nil
}

// SetTemplatePath changes templates path.
func (r *render) SetTemplatePath(newPath string) {
	renderOpt.Directory = newPath
	compile(renderOpt)
}
