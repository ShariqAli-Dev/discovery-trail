package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func customGetProviderName(r *http.Request) (string, error) {
	provider := r.PathValue("provider")
	if provider == "" {
		return "", fmt.Errorf("expected provider, got %s", provider)
	}

	return provider, nil
}

func NewAuth(addr *string, store *sessions.CookieStore) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	isProd := os.Getenv("ENVIROMENT") == "production"

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	gothic.Store = store
	gothic.GetProviderName = customGetProviderName
	fmt.Println(*addr)
	if isProd {
		goth.UseProviders(
			google.New(googleClientId, googleClientSecret, "https://discovery-trail.shariqapps.dev/auth/google/callback"),
		)
	} else {
		goth.UseProviders(
			google.New(googleClientId, googleClientSecret, fmt.Sprintf("http://localhost%s/auth/google/callback", *addr)),
		)
	}
}
