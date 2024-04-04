package main

import (
	"net/http"

	"github.com/justinas/alice"
	"github.com/shariqali-dev/discovery-trail/ui"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=604800") // 7 days
		http.FileServer(http.FS(ui.Files)).ServeHTTP(w, r)
	}))

	base := alice.New()
	mux.Handle("GET /auth/{provider}", base.ThenFunc(app.login))
	mux.Handle("GET /auth/{provider}/callback", base.ThenFunc(app.callback))

	dynamic := base.Append(noSurf, app.authenticate)
	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))

	protected := dynamic.Append(app.requireAuthentication)
	mux.Handle("GET /dashboard", protected.ThenFunc(app.dashboard))
	mux.Handle("POST /logout/{provider}", protected.ThenFunc(app.logout))

	standard := alice.New(app.recoverPanic, app.logRequest, app.generateNonce, app.commonHeaders)
	return standard.Then(mux)
}
