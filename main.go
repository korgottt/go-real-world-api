package main

import (
	"log"
	"net/http"
	"realworldapi/server"
)

func main() {
	server := server.NewGlobalServer()

	if err := http.ListenAndServe(":3000", server); err != nil {
		log.Fatalf("could not listen on port 3000 %v", err)
	}
}
