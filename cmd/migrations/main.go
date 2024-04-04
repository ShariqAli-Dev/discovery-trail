package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/shariqali-dev/discovery-trail/internal/database"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	up := flag.Bool("up", false, "migrate up database")
	down := flag.Bool("down", false, "migrate down database")
	reset := flag.Bool("reset", false, "reset database")
	flag.Parse()

	if !*up && !*down && !*reset {
		fmt.Println("Usage: up, down, reset")
		os.Exit(1)
	}

	dbAuthToken := os.Getenv("DB_AUTH_TOKEN")
	dbURL := os.Getenv("DB_URL")
	db, err := database.OpenDB(dbURL, dbAuthToken)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if *up {
		err := database.MigrateUp(db)
		if err != nil {
			log.Fatal(err)
		}
	} else if *down {
		err := database.MigrateDown(db)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := database.MigrateDown(db)
		if err != nil {
			log.Fatal(err)
		}
		err = database.MigrateUp(db)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}

type Session struct {
	Token string
}
