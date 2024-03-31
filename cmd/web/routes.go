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

	standard := alice.New(app.recoverPanic, app.logRequest, app.generateNonce, app.commonHeaders)
	return standard.Then(mux)
}
