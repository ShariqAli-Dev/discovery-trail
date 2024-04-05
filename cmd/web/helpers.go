package main

import (
	"errors"
	"net/http"

	"github.com/shariqali-dev/discovery-trail/internal/models"
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

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func (app *application) getAccountFromRequestID(r *http.Request) (models.Account, error) {
	accountID := r.Context().Value(accountIDContextKey)
	if accountID == nil {
		return models.Account{}, nil
	}
	accountIDString, ok := accountID.(string)
	if !ok {
		return models.Account{}, nil
	}

	account, err := app.accounts.Get(accountIDString)
	if err != nil {
		if errors.Is(err, models.ErrorNoRecord) {
			return models.Account{}, nil
		}
		return models.Account{}, err
	}
	return account, nil
}
