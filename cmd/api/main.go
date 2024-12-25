package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/danilshap/domains-generator/internal/api"
	db "github.com/danilshap/domains-generator/internal/db/sqlc"
	"github.com/danilshap/domains-generator/pkg/config"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Load(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(cfg.DBDriver, cfg.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	log.Printf("Server is running on %s", cfg.ServerAddres)
	err = http.ListenAndServe(cfg.ServerAddres, server)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}
