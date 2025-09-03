package main

import (
	"Go-final/pkg/db"
	"Go-final/pkg/server"
	"log"
	"net/http"
	"os"
)

func main() {
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = "scheduler.db"
	}
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("failed start database: %v", err)
	}
	webDir := "./web"
	http.Handle("/", http.FileServer(http.Dir(webDir)))
	err := server.StartServer()
	if err != nil {
		log.Fatalf("failed start server: %v", err)
	}
}
