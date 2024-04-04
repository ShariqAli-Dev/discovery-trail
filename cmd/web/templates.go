package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/justinas/nosurf"
	"github.com/shariqali-dev/discovery-trail/internal/types"
)

func (app *application) newTemplateData(r *http.Request) (types.TemplateData, error) {
	nonce, err := getNonce(r)
	if err != nil {
		return types.TemplateData{}, err
	}

	isAuthenticated := app.isAuthenticated(r)

	return types.TemplateData{
		Nonce:           nonce,
		IsAuthenticated: isAuthenticated,
		CSRFToken:       nosurf.Token(r),
	}, nil
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page templ.Component) {
	w.WriteHeader(status)
	err := page.Render(r.Context(), w)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		app.serverError(w, r, err)
		return
	}
}
