package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/shariqali-dev/discovery-trail/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamic := alice.New()
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("POST /logout/google", dynamic.ThenFunc(app.logout))
	mux.Handle("GET /auth/google", dynamic.ThenFunc(app.login))
	mux.Handle("GET /auth/google/callback", dynamic.ThenFunc(app.callback))

	standard := alice.New(app.recoverPanic, app.logRequest, app.generateNonce)
	return standard.Then(mux)
}
