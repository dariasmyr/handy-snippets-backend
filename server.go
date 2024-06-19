package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"log"
	"net/http"
	"os"
	"pastebin/database"
	"pastebin/graph"
	"pastebin/services"
	"time"
)

const defaultPort = "8080"

func main() {

	dbPath := "./data/db.sqlite"

	db, err := database.InitDB(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer db.Close()

	documentService := services.NewDocumentService(db)

	documentService.StartExpiredDocumentsCleaner(1 * time.Hour)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	resolver := graph.NewResolver(documentService)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
