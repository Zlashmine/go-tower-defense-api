## Tower Defense API
### *Go* Api

#### Google Cloud:
`https://go-tower-defense-api-485626356297.us-central1.run.app/v1/health`
#### Run App
``` bash
# Spin up database
make db

# Run App
make run-app      
```

#### ___Db Migrations___
``` bash
make migration create_codes
make migrate-up
make migrate-down
```

#### ___Env___
```
# Reloading env
direnv allow .
```

#### __Env variables__

``` bash
# /.envrc
export ADDR=:8080
export EXTERNAL_URL=localhost:8080
export ENV=development
export AUTH_TOKEN=abc123

export DB_ADDR=postgres://user:password@localhost/tower_defense?sslmode=disable

export DB_MAX_OPEN_CONNS=10
export DB_MAX_IDLE_CONNS=10
export DB_MAX_IDLE_TIME=5m

export REDIS_ADDR=localhost:6379
# export REDIS_PW=
export REDIS_DB=0
export REDIS_ENABLED=false

export RATE_LIMITER_REQUESTS_COUNT=20
export RATE_LIMITER_ENABLED=true

export FROM_EMAIL="td.mazing@gmail.com"
export SENDGRID_API_KEY=""
```

## __Documentation__
### Open API

#### `Development`
``` bash
http://localhost:8080/v1/swagger/index.html
```
#### `Production`
``` bash
# UI
https://go-tower-defense-api-485626356297.us-central1.run.app/v1/swagger/index.html

# Json
https://go-tower-defense-api-485626356297.us-central1.run.app/v1/swagger/0.0.0.0:8080/swagger/doc.json
```