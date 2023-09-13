package main

import "net/http"

// landingPageHandler shows the list of all emails in the inbox.
func (app *application) landingPageHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Emails []*Email
	}{
		Emails: app.inbox,
	}
	app.templates["emails"].Execute(w, data)
}

func (app *application) emailPageHandler(w http.ResponseWriter, r *http.Request) {

}
