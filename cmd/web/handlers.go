package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/markbates/goth/gothic"
	"github.com/shariqali-dev/discovery-trail/internal/models"
	"github.com/shariqali-dev/discovery-trail/internal/types"
	"github.com/shariqali-dev/discovery-trail/internal/validator"
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

func (app *application) create(w http.ResponseWriter, r *http.Request) {
	data, err := app.newTemplateData(r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	createPage := pages.Create(data)
	app.render(w, r, http.StatusOK, createPage)
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	var form types.CourseCreateForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	unitCountAsInt, err := strconv.Atoi(r.PostForm.Get("unit-count"))
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
	}
	form.UnitCount = unitCountAsInt
	form.UnitValues = make(map[string]string)
	form.CheckField(validator.NotBlank(form.Title), "title", "This field cannot be blank")
	form.CheckField(validator.MinCount(form.UnitCount, 1), "unit-count", "Cant have less than 1 unit")
	form.CheckField(validator.MaxCount(form.UnitCount, 5), "unit-count", "Cant have more than 5 unit")
	for unit := range form.UnitCount {
		unit++
		unitString := fmt.Sprintf("unit-%d", unit)
		unitFormValue := r.PostForm.Get(unitString)
		form.UnitValues[unitString] = unitFormValue

		form.CheckField(validator.NotBlank(unitFormValue), unitString, "This field cannot be blank")
		form.CheckField(validator.MaxChars(unitFormValue, 30), unitString, "This field cannot be more than 30 characters long")
	}
	if !form.Valid() {
		courseCreateFormInputs := pages.CourseCreateFormInputs(form)
		app.render(w, r, http.StatusOK, courseCreateFormInputs)
		return
	}

	unsplashSearchTerm, err := getImageSearchTermFromTitle(app.openAiClient, form.Title)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	unsplashResult, err := getUnsplashImage(strings.TrimSpace(unsplashSearchTerm.SearchTerm))
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	courseID, err := app.courses.Insert(form.Title, unsplashResult.Results[0].Images.Regular)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	for unit := 1; unit <= form.UnitCount; unit++ {
		go func(unit int) {
			_ = courseID
			// start the openai course creation here
			// update the databsae with the course creation test
		}(unit)
	}

	courseCreateInputs := pages.CourseCreateFormInputs(form)
	app.render(w, r, http.StatusOK, courseCreateInputs)
}

func (app *application) callback(w http.ResponseWriter, r *http.Request) {
	account, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	accountIDBlob, err := json.Marshal(account.UserID)
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

	exists, err := app.accounts.Exists(account.UserID)

	if err != nil {
		if !errors.Is(err, models.ErrorNoRecord) {
			app.serverError(w, r, err)
			return
		}
	}
	if !exists {
		err = app.accounts.Insert(account.UserID, fmt.Sprintf("%s %s", account.FirstName, account.LastName), account.Email)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	sessionID, err := app.sessions.Create(accountIDBlob)
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
	if err != nil {
		app.serverError(w, r, err)
		return
	}
	currentSessionToken := cookieStore.Values["token"]
	if currentSessionToken != nil {
		err := app.sessions.Destroy(currentSessionToken.(string))
		if err != nil {
			if !errors.Is(err, models.ErrorNoRecord) {
				app.serverError(w, r, err)
				return
			}
		}
	}

	cookieStore.Options.MaxAge = -1
	err = cookieStore.Save(r, w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {
	account, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		gothic.BeginAuthHandler(w, r)
		return
	}
	accountIDBlob, err := json.Marshal(account.UserID)
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

	exists, err := app.accounts.Exists(account.UserID)

	if err != nil {
		if !errors.Is(err, models.ErrorNoRecord) {
			app.serverError(w, r, err)
			return
		}
	}
	if !exists {
		err = app.accounts.Insert(account.UserID, fmt.Sprintf("%s %s", account.FirstName, account.LastName), account.Email)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
	}

	sessionID, err := app.sessions.Create(accountIDBlob)
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
