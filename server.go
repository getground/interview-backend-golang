package main

import (
	"context"
	"net/http"

	"github.com/getground/interview-backend-golang/handlers"
	"github.com/getground/interview-backend-golang/internal/app/example"
	"github.com/getground/interview-backend-golang/internal/pkg/config"
	"github.com/getground/interview-backend-golang/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.Load,
			models.NewExampleRepository,
			example.NewService,
			handlers.NewExampleHandler,
			newRouter,
			newHTTPServer,
		),
		fx.Invoke(startServer),
	)
	app.Run()
}

func newRouter(
	exampleHandler *handlers.ExampleHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api/v1")
	{
		examples := api.Group("/examples")
		{
			examples.POST("/", exampleHandler.CreateExample)
			examples.GET("/", exampleHandler.GetAllExamples)
			examples.GET("/:id", exampleHandler.GetExampleByID)
			examples.PUT("/:id", exampleHandler.UpdateExample)
			examples.DELETE("/:id", exampleHandler.DeleteExample)
		}
	}
	return router
}

func newHTTPServer(
	cfg *config.Config,
	router *gin.Engine,
) *http.Server {
	return &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
}

func startServer(
	lifecycle fx.Lifecycle,
	server *http.Server,
) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					panic(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
	})
}
