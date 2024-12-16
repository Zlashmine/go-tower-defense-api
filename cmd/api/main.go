package main

import (
	"fmt"
	"log"

	"tower-defense-api/lib/db"
	"tower-defense-api/lib/env"
	"tower-defense-api/lib/repository"
)

const version = "v0.0.1"

func main() {
	config := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:password@localhost/tower_defense?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 10),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	db, err := db.New(
		config.db.addr,
		config.db.maxOpenConns,
		config.db.maxIdleConns,
		config.db.maxIdleTime,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	fmt.Println("Connected to database")

	repository := repository.New(db)

	app := &application{
		config:     config,
		repository: repository,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
