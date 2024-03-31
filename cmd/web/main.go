package main

import (
	"crypto/tls"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	// _ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger *slog.Logger
}

func main() {
	// googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	// googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	addr := flag.String("addr", ":4000", "HTTP network address")
	// dataSourceName := flag.String("dsn", "web:password@/discovery_trail?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	// db, err := openDB(*dataSourceName)
	// if err != nil {
	// 	logger.Error(err.Error())
	// 	os.Exit(1)
	// }
	// defer db.Close()

	// sessionManager := scs.New()
	// sessionManager.Store = mysqlstore.New(db)
	// sessionManager.Lifetime = 12 * time.Hour
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

	// err = server.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	err := server.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

// func openDB(dsn string) (*sql.DB, error) {
// 	db, err := sql.Open("mysql", dsn)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		db.Close()
// 		return nil, err
// 	}

// 	return db, nil
// }

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
