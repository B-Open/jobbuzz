package main

import (
	"log"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/b-open/jobbuzz/graph"
	"github.com/b-open/jobbuzz/graph/generated"
	"github.com/b-open/jobbuzz/internal/config"
	"github.com/b-open/jobbuzz/pkg/controller"
	"github.com/b-open/jobbuzz/pkg/service"
	"github.com/gin-gonic/gin"
)

// Defining the Graphql handler
func graphqlHandler(service service.Servicer) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{Service: service}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	configuration, err := config.LoadConfig("../../")

	if err != nil {
		log.Fatal("Fail to load db config", err)
	}

	db, err := configuration.GetDb()

	if err != nil {
		log.Fatal("Fail to get db connection", err)
	}

	service := service.Service{DB: db}
	controller := controller.Controller{Service: &service}

	r := gin.Default()

	if gin.IsDebugging() {
	    r.GET("/", playgroundHandler())
	}
	r.POST("/graphql", graphqlHandler(&service))

	apiV1 := r.Group("/api/v1")
	apiV1.GET("/job", controller.GetJobs)
	r.Run()
}
