package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"
	"github.com/shariqali-dev/discovery-trail/internal/models"
)

func (app *application) commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonce, err := getNonce(r)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		_ = nonce

		w.Header().Set("Content-Security-Policy",
			fmt.Sprintf("default-src 'self'; script-src 'self' 'nonce-%[1]s'; style-src 'self' fonts.googleapis.com 'nonce-%[1]s'; font-src fonts.gstatic.com;", nonce),
		)
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Go")

		next.ServeHTTP(w, r)
	})
}

func (app *application) generateNonce(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonce := make([]byte, 16)
		_, err := rand.Read(nonce)
		if err != nil {
			app.serverError(w, r, err)
			return
		}

		newContext := context.WithValue(r.Context(), nonceContextKey, base64.StdEncoding.EncodeToString(nonce))
		r = r.WithContext(newContext)
		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieStore, err := app.store.Get(r, "discovery-trail")
		if cookieStore.IsNew || err != nil {
			newContext := context.WithValue(r.Context(), isAuthenticatedContextKey, false)
			r = r.WithContext(newContext)
			next.ServeHTTP(w, r)
			return
		} else {
			sessionToken := cookieStore.Values["token"]
			if sessionToken == nil {
				newContext := context.WithValue(r.Context(), isAuthenticatedContextKey, false)
				r = r.WithContext(newContext)
				next.ServeHTTP(w, r)
				return
			}

			userID, err := app.sessions.GetUserID(sessionToken.(string))
			if err != nil {
				if errors.Is(err, models.ErrorNoRecord) {
					newContext := context.WithValue(r.Context(), isAuthenticatedContextKey, false)
					r = r.WithContext(newContext)
					next.ServeHTTP(w, r)
					return
				} else {
					app.serverError(w, r, err)
					return
				}
			}

			exists, err := app.accounts.Exists(userID)
			if err != nil {
				if errors.Is(err, models.ErrorNoRecord) {
					newContext := context.WithValue(r.Context(), isAuthenticatedContextKey, false)
					r = r.WithContext(newContext)
					next.ServeHTTP(w, r)
					return
				} else {
					app.serverError(w, r, err)
					return
				}
			}

			if exists {
				newContext := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
				r = r.WithContext(newContext)
				next.ServeHTTP(w, r)
			} else {
				newContext := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
				r = r.WithContext(newContext)
				next.ServeHTTP(w, r)
			}
		}

	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return csrfHandler
}
