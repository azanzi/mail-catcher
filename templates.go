package main

import "html/template"

// loadTemplates loads all the html templates into memory
func (app *application) loadTemplates() {
	app.templates["emails"] = template.Must(template.ParseFiles("templates/emails.tmpl"))
}
