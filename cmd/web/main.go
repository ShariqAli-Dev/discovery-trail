package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type application struct {
	logger *slog.Logger
}

func main() {
	authToken := os.Getenv("AUTH_TOKEN")
	dbName := os.Getenv("DB_NAME")
	addr := flag.String("addr", ":4000", "HTTP network address")
	prod := flag.Bool("prod", false, "Enable production")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	db, err := openDB(dbName, authToken)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	app := &application{
		logger: logger,
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

	if *prod {
		err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	} else {
		err = server.ListenAndServe()
	}
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dbName, authToken string) (*sql.DB, error) {
	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", dbName, authToken)

	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, fmt.Errorf("failed to open db %s: %s", url, err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
