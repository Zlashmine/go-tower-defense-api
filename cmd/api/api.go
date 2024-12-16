package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"tower-defense-api/lib/repository"
)

type application struct {
	config     config
	repository repository.Repository
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(middleware.Timeout(60 * time.Second))

	router.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

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
	})

	return router
}

func (app *application) run(mux http.Handler) error {
	server := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 60,
	}

	fmt.Printf("Starting server on %s\n", app.config.addr)

	return server.ListenAndServe()
}
