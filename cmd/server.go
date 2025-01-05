package main

import (
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"log"
	"net/http"
	"os"
	"pastebin/database"
	"pastebin/graph"
	"pastebin/services"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/rs/cors"
)

const defaultPort = "4000"

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
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	frontendURL := os.Getenv("FRONTEND_URL")
	log.Println("Frontend URL:", frontendURL)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{frontendURL},
		AllowCredentials: true,
		AllowedMethods:   []string{"POST", "GET", "OPTIONS", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", c.Handler(srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
