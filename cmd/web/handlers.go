package main

import (
	"net/http"

	"github.com/shariqali-dev/discovery-trail/ui/html/pages"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data, err := app.newTemplateData(r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	app.sessionManager.Put(r.Context(), "flash", "THE SESSION MANAGER IS FUCKING WORKING")
	flash := app.sessionManager.PopString(r.Context(), "flash")
	data.Flash = flash

	homePage := pages.Home(data)

	app.render(w, r, http.StatusOK, homePage)
}
