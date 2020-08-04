package main

import (
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/korgottt/go-real-world-api/server"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	store := server.ArticleDBStore{}
	if err := store.Init(); err != nil {
		log.Fatalf("unable to access the database: %q", err)
	}
	defer store.Close()
	server := server.NewGlobalServer(&store)

	if err := http.ListenAndServe(":3000", server); err != nil {
		log.Fatalf("could not listen on port 3000 %v", err)
	}
}
