package main

import (
	"errors"
	"net/http"
)

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func getNonce(r *http.Request) (string, error) {
	nonce, ok := r.Context().Value(nonceContextKey).(string)
	if !ok {
		return "", errors.New("could not convert nonce to string")
	}
	return nonce, nil
}
