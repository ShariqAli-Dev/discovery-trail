package main

import (
	"crypto/tls"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/shariqali-dev/discovery-trail/internal/auth"
	"github.com/shariqali-dev/discovery-trail/internal/database"
	"github.com/shariqali-dev/discovery-trail/internal/models"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type application struct {
	logger   *slog.Logger
	accounts models.AccountModelInterface
	sessions models.SessionModelInterface
	store    *sessions.CookieStore
}

func main() {
	dbAuthToken := os.Getenv("DB_AUTH_TOKEN")
	dbURL := os.Getenv("DB_URL")
	addr := flag.String("addr", ":4000", "HTTP network address")
	authKey := os.Getenv("AUTH_KEY")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	db, err := database.OpenDB(dbURL, dbAuthToken)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	store := sessions.NewCookieStore([]byte(authKey))
	store.MaxAge(int(time.Hour.Seconds() * 24 * 7))
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = true

	auth.NewAuth(addr, store)

	app := &application{
		logger:   logger,
		accounts: &models.AccountModel{DB: db},
		sessions: &models.SessionModel{DB: db},
		store:    store,
	}

	tlsConfig := &tls.Config{
		CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256},
	}
	server := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		TLSConfig:    tlsConfig,
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logger.Info("starting server", "addr", server.Addr)

	err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	logger.Error(err.Error())
	os.Exit(1)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
