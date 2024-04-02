package main

import (
	"net/http"

	"github.com/markbates/goth/gothic"
	"github.com/shariqali-dev/discovery-trail/ui/html/pages"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data, err := app.newTemplateData(r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	// app.sessionManager.Put(r.Context(), "flash", "THE SESSION MANAGER IS FUCKING WORKING")
	// flash := app.sessionManager.PopString(r.Context(), "flash")
	// data.Flash = flash

	homePage := pages.Home(data)

	app.render(w, r, http.StatusOK, homePage)
}

func (app *application) callback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.logger.Info("callbacked user", "user", user)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		gothic.BeginAuthHandler(w, r)
		return
	}
	app.logger.Info("logged in user", "user", user)
	http.Redirect(w, r, "", http.StatusSeeOther)
}
