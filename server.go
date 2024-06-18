package main

import (
	"log"
	"net/http"
	"os"
	"pastebin/database"
	"pastebin/graph"
	"pastebin/services"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
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

	go func() {
		ticker := time.NewTicker(24 * time.Hour) // Пример: каждые сутки
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				err := documentService.DeleteExpiredDocuments()
				if err != nil {
					log.Printf("failed to delete expired documents: %v", err)
				}
			}
		}
	}()

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
