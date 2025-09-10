package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"test_kode/internal/config"
	"test_kode/internal/db"
	"test_kode/internal/server"
	"test_kode/internal/service"
)

func main() {
	// конфиг
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// база
	// будем использовать драйвер для psql
    pdb, err := sql.Open("postgres", cfg.DatabaseURL)
    if err != nil {
        log.Fatal(err)
    }
    defer pdb.Close()

	dataBaseInit := db.New(pdb)
	service  := service.New(dataBaseInit, cfg)
	server  := server.New(service, cfg)

	// http-сервер
	httpSrv := &http.Server{
		Addr:    cfg.Port,
		Handler: server,
	}
	log.Fatal(httpSrv.ListenAndServe())
}