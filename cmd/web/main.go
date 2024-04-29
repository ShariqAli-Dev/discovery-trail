package main

import (
	"crypto/tls"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
	"github.com/shariqali-dev/discovery-trail/internal/auth"
	"github.com/shariqali-dev/discovery-trail/internal/database"
	"github.com/shariqali-dev/discovery-trail/internal/models"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type application struct {
	logger       *slog.Logger
	accounts     models.AccountModelInterface
	sessions     models.SessionModelInterface
	courses      models.CourseModelInterface
	units        models.UnitModelInterface
	chapters     models.ChapterModelInterface
	questions    models.QuestionModelInterface
	store        *sessions.CookieStore
	formDecoder  *form.Decoder
	openAiClient *openai.Client
}

func main() {
	dbAuthToken := os.Getenv("DB_AUTH_TOKEN")
	dbURL := os.Getenv("DB_URL")
	openAiAPIKey := os.Getenv("OPEN_AI_API_KEY")
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
	formDecoder := form.NewDecoder()

	openAiClient := openai.NewClient(openAiAPIKey)

	app := &application{
		logger:       logger,
		accounts:     &models.AccountModel{DB: db},
		sessions:     &models.SessionModel{DB: db},
		courses:      &models.CourseModel{DB: db},
		units:        &models.UnitModel{DB: db},
		chapters:     &models.ChapterModel{DB: db},
		questions:    &models.QuestionModel{DB: db},
		store:        store,
		formDecoder:  formDecoder,
		openAiClient: openAiClient,
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
		WriteTimeout: 15 * time.Second,
	}

	err = server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
