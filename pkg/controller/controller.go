package controller

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/b-open/jobbuzz/pkg/graph"
	"github.com/b-open/jobbuzz/pkg/graph/generated"
	"github.com/b-open/jobbuzz/pkg/service"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	Service service.Servicer
}

func (controller *Controller) GetJobs(c *gin.Context) {
	jobs, err := controller.Service.GetJobs()

	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{
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

