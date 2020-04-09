package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alessio-palumbo/linktree-challenge/handlers"
	"github.com/alessio-palumbo/linktree-challenge/middleware"
	"github.com/alessio-palumbo/linktree-challenge/server"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/stdlib"
)

var (
	port     = flag.Int("port", 8080, "port")
	dbSource = flag.String("db_source", "dbname=linktree-dev sslmode=disable", "Db")
	maxDBC   = 5
	nWorkers = 1
	apiURL   = "http://linktr.ee/api"
	authURL  = ""
)

func main() {
	// Parse flags for custom inputs
	flag.Parse()

	// Parse db string and initialise pool
	pgxcfg, err := pgx.ParseConnectionString(*dbSource)
	if err != nil {
		log.Fatal("Failed to parse db-source")
	}

	// TODO Add retry func
	pool := stdlib.OpenDB(pgxcfg)

	pool.SetConnMaxLifetime(time.Duration(10 * time.Minute))
	pool.SetMaxIdleConns(10)
	pool.SetMaxOpenConns(10)

	g := handlers.Group{
		DB:   pool,
		Auth: middleware.NewAuth(pool),
	}

	// Start server
	s := http.Server{
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      server.New(g),
	}

	log.Fatal(s.ListenAndServe())
}
