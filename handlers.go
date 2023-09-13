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
		Email  *Email
		Id     int
	}{
		Emails: app.inbox,
		Email:  nil,
		Id:     -1,
	}
	app.render(w, "emails", data)
}

// emailPageHandler shows the email at position 'id'.
func (app *application) emailPageHandler(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 0 || id >= len(app.inbox) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	app.inbox[id].Seen = true

	data := struct {
		Emails []*Email
		Email  *Email
		Id     int
	}{
		Emails: app.inbox,
		Email:  app.inbox[id],
		Id:     id,
	}
	app.render(w, "emails", data)
}
