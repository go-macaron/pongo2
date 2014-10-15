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

package pongo2

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/Unknwon/macaron"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Test_Render_Bad_HTML(t *testing.T) {
	m := macaron.Classic()
	m.Use(Pongoer(Options{
		Directory: "fixtures/basic",
	}))

	// routing
	m.Get("/foobar", func(r macaron.Render) {
		r.HTML(200, "nope", nil)
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 500)
	expect(t, res.Body.String(), "pongo2: \"nope\" is undefined\n")
}

func Test_Render_HTML(t *testing.T) {
	m := macaron.Classic()
	m.Use(Pongoer(Options{
		Directory: "fixtures/basic",
	}))

	// routing
	m.Get("/foobar", func(r macaron.Render) {
		r.HTML(200, "hello", map[string]interface{}{
			"Name": "jeremy",
		})
		r.SetTemplatePath("fixtures/basic2")
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML+"; charset=UTF-8")
	expect(t, res.Body.String(), "<h1>Hello jeremy</h1>\n")

	// Change templates path.
	res = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML+"; charset=UTF-8")
	expect(t, res.Body.String(), "<h1>What's up, jeremy</h1>\n")
}

func Test_Render_XHTML(t *testing.T) {
	m := macaron.Classic()
	m.Use(Pongoer(Options{
		Directory:       "fixtures/basic",
		HTMLContentType: ContentXHTML,
	}))

	m.Get("/foobar", func(r macaron.Render) {
		r.HTML(200, "hello", map[string]interface{}{
			"Name": "jeremy",
		})
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentXHTML+"; charset=UTF-8")
	expect(t, res.Body.String(), "<h1>Hello jeremy</h1>\n")
}

func Test_Render_Extensions(t *testing.T) {
	m := macaron.Classic()
	m.Use(Pongoer(Options{
		Directory:  "fixtures/basic",
		Extensions: []string{".tmpl", ".html"},
	}))

	// routing
	m.Get("/foobar", func(r macaron.Render) {
		r.HTML(200, "hypertext", map[string]interface{}{})
	})

	res := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/foobar", nil)

	m.ServeHTTP(res, req)

	expect(t, res.Code, 200)
	expect(t, res.Header().Get(ContentType), ContentHTML+"; charset=UTF-8")
	expect(t, res.Body.String(), "Hypertext!\n")
}

func Test_Render_NoRace(t *testing.T) {
	// This test used to fail if run with -race
	m := macaron.Classic()
	m.Use(Pongoer(Options{
		Directory: "fixtures/basic",
	}))

	// routing
	m.Get("/foobar", func(r macaron.Render) {
		r.HTML(200, "hello", map[string]interface{}{
			"Name": "world",
		})
	})

	done := make(chan bool)
	doreq := func() {
		res := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/foobar", nil)

		m.ServeHTTP(res, req)

		expect(t, res.Code, 200)
		expect(t, res.Header().Get(ContentType), ContentHTML+"; charset=UTF-8")
		// ContentLength should be deferred to the ResponseWriter and not Render
		expect(t, res.Header().Get(ContentLength), "")
		expect(t, res.Body.String(), "<h1>Hello world</h1>\n")
		done <- true
	}
	// Run two requests to check there is no race condition
	go doreq()
	go doreq()
	<-done
	<-done
}
