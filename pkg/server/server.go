package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func StartServer() error {
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}
	log.Printf("Server start, port:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return fmt.Errorf("could not start server: %w", err)
	}
	return nil
}
