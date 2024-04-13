package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/go-playground/form/v4"
	"github.com/sashabaranov/go-openai"
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

func getUnsplashImage(query string) (PhotoResponse, error) {
	resp, err := http.Get(fmt.Sprintf("https://api.unsplash.com/search/photos?per_page=1&query=%s&client_id=%s", query, os.Getenv("UNSPLASH_ACCESS_KEY")))
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

type ImageSearchTerm struct {
	SearchTerm string `json:"image_search_term"`
}

func getImageSearchTermFromTitle(client *openai.Client, title string) (ImageSearchTerm, error) {
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an AI capable of suggesting an image search term for a given course title. Provide the search term in the JSON format as shown: { image_search_term: \"search term here\" }.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("I have a course titled, %s. Provide a good image search term. This search term will be fed into the unsplash API, so make sure it is a good search term that will return good results.", title),
			},
		},
	})
	if err != nil {
		return ImageSearchTerm{}, nil
	}

	var imageSearchTerm ImageSearchTerm
	if err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &imageSearchTerm); err != nil {
		return ImageSearchTerm{}, nil
	}
	return imageSearchTerm, nil
}
