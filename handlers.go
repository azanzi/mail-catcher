package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// landingPageHandler shows the list of all emails in the inbox.
func (app *application) landingPageHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Emails []*Email
	}{
		Emails: app.inbox,
	}
	app.templates["emails"].Execute(w, data)
}

// emailPageHandler shows the email at position 'id'.
func (app *application) emailPageHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 0 || id >= len(app.inbox) {
		// TODO: out of range, handle error
	}

	data := struct {
		Emails []*Email
		Email  *Email
	}{
		Emails: app.inbox,
		Email:  app.inbox[id],
	}
	app.templates["emails"].Execute(w, data)
}
