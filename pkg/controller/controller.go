package controller

import (
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/b-open/jobbuzz/pkg/graph"
	"github.com/b-open/jobbuzz/pkg/graph/generated"
	"github.com/b-open/jobbuzz/pkg/graph/graphmodel"

	"github.com/b-open/jobbuzz/pkg/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Controller struct {
	Service service.Servicer
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func (controller *Controller) GetJobs(c *gin.Context) {
	jobs, err := controller.Service.GetJobs(graphmodel.PaginationInput{})
	if err != nil {
		panic(err)
	}

	log.Debug().Msgf("Found %d jobs", len(jobs))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"jobs":    jobs,
	})
}

// Defining the Graphql handler
func GraphqlHandler(service service.Servicer) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	h := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{Service: service}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func PlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
