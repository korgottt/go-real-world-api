package main

import (
	"log"
	"net/http"

	"github.com/korgottt/go-real-world-api/server"
)

func main() {
	server := server.NewGlobalServer(&server.InMemoryStore{})

	if err := http.ListenAndServe(":3000", server); err != nil {
		log.Fatalf("could not listen on port 3000 %v", err)
	}
}
