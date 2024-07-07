package main

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
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

	documentService.StartExpiredDocumentsCleaner(24 * time.Hour)

	port := os.Getenv("PORT")
	println("port", port)
	if port == "" {
		port = defaultPort
	}

	resolver := graph.NewResolver(documentService)
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{os.Getenv("FRONTEND_URL")},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", c.Handler(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
