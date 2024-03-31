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
	homePage := pages.Home(data)

	app.render(w, r, http.StatusOK, homePage)

}
