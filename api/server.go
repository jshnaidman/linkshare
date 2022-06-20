package main

import (
	"linkshare_api/graph"
	"linkshare_api/graph/generated"
	"linkshare_api/utils"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func init() {
	// init the random seed
	rand.Seed(time.Now().UnixNano())
	// init the conf object on startup and fail quickly if there's an environment issue
	_ = utils.GetConf()
}

func main() {
	// graphql playground port
	port := "8080"

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
