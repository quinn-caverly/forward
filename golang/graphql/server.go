package main

import (
	"log"
	"net/http"
	"os"

	"github.com/quinn-caverly/forward-graphql/internal/db"
	"github.com/rs/cors"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/quinn-caverly/forward-graphql/graph"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	err := db.CreateConnToMongo()
	if err != nil {
		log.Fatal("Could not connect to the mongodb, ", err)
	}
	defer db.CloseConnToMongo()

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Replace with your React app's URL
		AllowedMethods: []string{"GET", "POST"},
	})

	// Wrap the http.DefaultServeMux with the CORS middleware
	handler := c.Handler(http.DefaultServeMux)

	log.Println("Running on: ", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
