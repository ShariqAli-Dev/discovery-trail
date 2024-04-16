package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/go-playground/form/v4"
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

func (app *application) requestGetAccountID(r *http.Request) (string, error) {
	accountID := r.Context().Value(accountIDContextKey)
	if accountID == nil {
		return "", nil
	}
	accountIDString, ok := accountID.(string)
	if !ok {
		return "", nil
	}
	return strings.Trim(accountIDString, `"`), nil
}

func (app *application) requestGetAccount(r *http.Request) (models.Account, error) {
	accountID, err := app.requestGetAccountID(r)
	if err != nil {
		return models.Account{}, nil
	}

	account, err := app.accounts.Get(accountID)
	if err != nil {
		if errors.Is(err, models.ErrorNoRecord) {
			return models.Account{}, nil
		}
		return models.Account{}, err
	}
	return account, nil
}

func (app *application) decodePostForm(r *http.Request, destination any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(destination, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}

type PhotoResponse struct {
	Total      int     `json:"total"`
	TotalPages int     `json:"total_pages"`
	Results    []Photo `json:"results"`
}
type Photo struct {
	ID     string   `json:"id"`
	Images ImageSet `json:"urls"`
}
type ImageSet struct {
	Raw     string `json:"raw"`
	Full    string `json:"full"`
	Regular string `json:"regular"`
	Small   string `json:"small"`
	Thumb   string `json:"thumb"`
}

func unsplashGetImage(query string) (PhotoResponse, error) {
	baseURL := "https://api.unsplash.com/search/photos"
	clientID := os.Getenv("UNSPLASH_ACCESS_KEY")
	params := url.Values{}
	params.Add("per_page", "1")
	params.Add("query", query)
	params.Add("client_id", clientID)
	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	resp, err := http.Get(fullURL)
	if resp.StatusCode != http.StatusOK {
		return PhotoResponse{}, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
	}
	if err != nil {
		return PhotoResponse{}, err
	}
	var photoResponse PhotoResponse
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&photoResponse); err != nil {
		return PhotoResponse{}, err
	}
	return photoResponse, err
}
