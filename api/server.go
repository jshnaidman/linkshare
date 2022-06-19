package main

import (
	b64 "encoding/base64"
	"linkshare_api/conf"
	"linkshare_api/graph"
	"linkshare_api/graph/generated"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "8080"

var URL_encoding = b64.URLEncoding.WithPadding(b64.NoPadding)

func init() {
	// init the random seed
	rand.Seed(time.Now().UnixNano())
	// init the conf object on startup and fail quickly if there's an environment issue
	_ = conf.GetConf()
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
