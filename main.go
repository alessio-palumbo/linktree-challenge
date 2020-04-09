package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alessio-palumbo/linktree-challenge/server"
	"github.com/jackc/pgx"
)

var (
	port     = flag.Int("port", 80, "port")
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
	dbPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     pgxcfg,
		MaxConnections: maxDBC,
	})
	if err != nil {
		log.Fatal("Failed to create new connection pool")
	}
	defer dbPool.Close()

	// TODO Add middleware (auth)

	// Start server
	s := http.Server{
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      server.New(dbPool),
	}

	s.ListenAndServe()
}
