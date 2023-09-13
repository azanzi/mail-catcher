package main

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
)

// newTemplateCache loads all the html templates into memory
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	ts, err := template.ParseFiles("templates/emails.tmpl")
	if err != nil {
		return nil, err
	}

	cache["emails"] = ts
	return cache, nil
}

func (app *application) render(w http.ResponseWriter, name string, td interface{}) {
	ts, ok := app.templates[name]
	if !ok {
		// TODO: no template, handle error
		return
	}

	buf := new(bytes.Buffer)

	err := ts.Execute(buf, td)
	if err != nil {
		fmt.Println(err)
		return
	}

	buf.WriteTo(w)
}
