package main

import (
	"authentication/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var count int64

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	log.Println("Starting authentication service")

	// Connect to db
	conn := connectToDB()
	if conn == nil {
		log.Panic("Could not connect to DB!")
	}

	// Set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	// Define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	// Start he server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dns string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dns)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready yet...")
			count++
		} else {
			fmt.Println("Connected to PostgresDB")
			return connection
		}

		if count > 10 {
			log.Panicln(err)
			return nil
		}

		log.Println("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
	}
}
