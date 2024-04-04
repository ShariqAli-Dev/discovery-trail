package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/markbates/goth/gothic"
	"github.com/shariqali-dev/discovery-trail/internal/models"
	"github.com/shariqali-dev/discovery-trail/ui/html/pages"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if app.isAuthenticated(r) {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	data, err := app.newTemplateData(r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	homePage := pages.Home(data)
	app.render(w, r, http.StatusOK, homePage)
}

func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {
	data, err := app.newTemplateData(r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	dashboardPage := pages.Dashboard(data)
	app.render(w, r, http.StatusOK, dashboardPage)
}

func (app *application) callback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	userIDBlob, err := json.Marshal(user.UserID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	cookieStore, err := app.store.Get(r, "discovery-trail")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	currentSessionToken := cookieStore.Values["token"]
	if currentSessionToken != nil {
		err := app.sessions.Destroy(currentSessionToken.(string))
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		cookieStore.Values["token"] = nil
	}

	exists, err := app.users.Exists(user.UserID)

	if err != nil {
		if !errors.Is(err, models.ErrorNoRecord) {
			app.serverError(w, r, err)
			return
		}
	}
	if !exists {
		err = app.users.Insert(user.UserID, fmt.Sprintf("%s %s", user.FirstName, user.LastName), user.Email)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	sessionID, err := app.sessions.Create(userIDBlob)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	cookieStore.Values["token"] = sessionID
	err = cookieStore.Save(r, w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	gothic.Logout(w, r)
	cookieStore, err := app.store.Get(r, "discovery-trail")
	if cookieStore.IsNew || err != nil {
		app.serverError(w, r, err)
		return
	}
	currentSessionToken := cookieStore.Values["token"]
	if currentSessionToken != nil {
		err := app.sessions.Destroy(currentSessionToken.(string))
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		cookieStore.Values["token"] = nil
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		gothic.BeginAuthHandler(w, r)
		return
	}
	userIDBlob, err := json.Marshal(user.UserID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	cookieStore, err := app.store.Get(r, "discovery-trail")
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	currentSessionToken := cookieStore.Values["token"]
	if currentSessionToken != nil {
		err := app.sessions.Destroy(currentSessionToken.(string))
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		cookieStore.Values["token"] = nil
	}

	exists, err := app.users.Exists(user.UserID)

	if err != nil {
		if !errors.Is(err, models.ErrorNoRecord) {
			app.serverError(w, r, err)
			return
		}
	}
	if !exists {
		err = app.users.Insert(user.UserID, fmt.Sprintf("%s %s", user.FirstName, user.LastName), user.Email)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	sessionID, err := app.sessions.Create(userIDBlob)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	cookieStore.Values["token"] = sessionID
	err = cookieStore.Save(r, w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
