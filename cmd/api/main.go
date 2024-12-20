package main

import (
	"expvar"
	"runtime"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"tower-defense-api/lib/db"
	"tower-defense-api/lib/env"
	"tower-defense-api/lib/notifications"
	"tower-defense-api/lib/ratelimiter"
	"tower-defense-api/lib/repository"
	"tower-defense-api/lib/repository/cache"
)

const version = "1.3.0"

func main() {
	config := config{
		addr:      env.GetString("ADDR", ":8080"),
		apiURL:    env.GetString("EXTERNAL_URL", "localhost:8080"),
		env:       env.GetString("ENV", "development"),
		authToken: env.GetString("AUTH_TOKEN", ""),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://user:password@localhost/tower_defense?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 10),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisConfig: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false),
		},
		rateLimiter: ratelimiter.Config{
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true),
			RequestsPerTimeFrame: env.GetInt("RATE_LIMITER_REQUESTS_COUNT", 20),
			TimeFrame:            time.Second * 30,
		},
		mailer: mailConfig{
			exp:       time.Hour * 24 * 3,
			fromEmail: env.GetString("FROM_EMAIL", ""),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
		},
	}

	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	db, err := db.New(
		config.db.addr,
		config.db.maxOpenConns,
		config.db.maxIdleConns,
		config.db.maxIdleTime,
	)

	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("Connected to database")

	ratelimiter := ratelimiter.NewFixedWindowLimiter(
		config.rateLimiter.RequestsPerTimeFrame,
		config.rateLimiter.TimeFrame,
	)

	// Mailer
	sendgridMail := notifications.NewSendgrid(config.mailer.sendGrid.apiKey, config.mailer.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	var redisClient *redis.Client

	if config.redisConfig.enabled {
		redisClient = cache.NewRedisClient(
			config.redisConfig.addr,
			config.redisConfig.pw,
			config.redisConfig.db,
		)

		logger.Info("Connected to redis")
	}

	cacheStore := cache.NewRedisStore(redisClient)
	repository := repository.New(db)

	app := &application{
		config:       config,
		repository:   repository,
		cacheStore:   cacheStore,
		logger:       logger,
		rateLimiter:  ratelimiter,
		notification: sendgridMail,
	}

	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
