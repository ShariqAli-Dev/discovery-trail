package auth

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

const (
	key    = "9013g3thnohd@#OJKJq0"
	MaxAge = 86400 * 40
	IsProd = true
)

func CustomGetProviderName(r *http.Request) (string, error) {
	return "google", nil
}

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store
	gothic.GetProviderName = CustomGetProviderName

	goth.UseProviders(
		google.New(googleClientId, googleClientSecret, "http://localhost:4000/auth/google/callback"),
	)
}
