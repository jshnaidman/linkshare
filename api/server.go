package main

import (
	"linkshare_api/auth"
	"linkshare_api/graph"
	"linkshare_api/graph/generated"
	"linkshare_api/utils"
	"math/rand"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
)

var CSRFHandler utils.Middleware

var conf *utils.Conf

func init() {
	// init the random seed
	rand.Seed(time.Now().UnixNano())
	// init the conf object on startup and fail quickly if there's an environment issue
	conf = utils.GetConf()
}

func main() {
	router := chi.NewRouter()

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	// corsMiddleware := cors.New(cors.Options{
	// 	AllowedOrigins:   []string{conf.AllowedOrigins},
	// 	AllowCredentials: true,
	// 	Debug:            conf.DebugMode,
	// }).Handler

	// router.Use(corsMiddleware)
	//router.Use(auth.AuthMiddleware())

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	router.Handle("/", playground.Handler("Linkshare", "/query"))
	router.Handle("/query", srv)
	router.HandleFunc("/loginGoogleJWT", auth.GoogleLoginHandleFunc)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}
