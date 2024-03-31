package main

import (
	"crypto/tls"
	"database/sql"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger         *slog.Logger
	sessionManager *scs.SessionManager
}

func main() {
	// authToken := os.Getenv("AUTH_TOKEN")
	// dbName := os.Getenv("DB_NAME")
	addr := flag.String("addr", ":4000", "HTTP network address")
	prod := flag.Bool("prod", false, "Enable production")
	dataSourceName := flag.String("dsn", "web:password@/discovery_trail?parseTime=true", "MySQL data source name")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
	}))

	db, err := openDB(*dataSourceName)
	// db, err := openTursoDB(dbName, authToken)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	app := &application{
		logger:         logger,
		sessionManager: sessionManager,
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

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

// func openTursoDB(dbName, authToken string) (*sql.DB, error) {
// 	url := fmt.Sprintf("libsql://%s.turso.io?authToken=%s", dbName, authToken)

// 	db, err := sql.Open("libsql", url)
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

// func init() {
// 	if err := godotenv.Load(); err != nil {
// 		log.Fatal(err)
// 	}
// }
