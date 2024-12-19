package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"tower-defense-api/docs"
	"tower-defense-api/lib/ratelimiter"
	"tower-defense-api/lib/repository"
	"tower-defense-api/lib/repository/cache"
)

type application struct {
	config      config
	logger      *zap.SugaredLogger
	repository  repository.Repository
	cacheStore  cache.Store
	rateLimiter ratelimiter.Limiter
}

type config struct {
	addr        string
	db          dbConfig
	env         string
	apiURL      string
	redisConfig redisConfig
	rateLimiter ratelimiter.Config
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type redisConfig struct {
	addr string
	pw   string
	db   int
	enabled bool
}

//	@Title			Tower Defense API
//	@Description	Documentation for the Tower Defense API

// @securityDefinitions.apiKey	ApiKeyAuth
// @in							header
// @name						Authorization
func (app *application) mount() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	if app.config.rateLimiter.Enabled {
		router.Use(app.RateLimiterMiddleware)
	}

	docsUrl := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)

	router.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsUrl)))

		r.Route("/codes", func(r chi.Router) {
			r.Get("/", app.getAllCodesHandler)
			r.Post("/", app.createCodeHandler)
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/", app.createUserHandler)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", app.getUserHandler)
			})
		})

		r.Route("/messages", func(r chi.Router) {
			r.Post("/", app.createMessageHandler)

			r.Route("/{id}", func(r chi.Router) {
				r.Patch("/", app.setReadHandler)
			})
		})
	})

	return router
}

func (app *application) run(mux http.Handler) error {
	docs.SwaggerInfo.Title = "Tower Defense API"
	docs.SwaggerInfo.Description = "Documentation for the Tower Defense API"
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL

	return app.withGracefulShutdown(&http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 60,
	})
}

func (app *application) withGracefulShutdown(server *http.Server) error {
	shutdown := make(chan error)

	go func() {
		killChannel := make(chan os.Signal, 1)

		signal.Notify(killChannel, syscall.SIGINT, syscall.SIGTERM)
		killSignal := <-killChannel

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.Infow("signal caught", "signal", killSignal.String())

		shutdown <- server.Shutdown(ctx)
	}()

	app.logger.Infow("Server started", "addr", app.config.addr, "env", app.config.env)

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	app.logger.Infow("Server has stopped", "addr", app.config.addr, "env", app.config.env)

	return nil
}
