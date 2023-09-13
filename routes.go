package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/", app.landingPageHandler)
	router.HandlerFunc(http.MethodGet, "/:id", app.emailPageHandler)

	return router
}
