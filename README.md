## Tower Defense API
### *Golang* Api

### Run App
``` bash
# Spin up database
make db

# Run App
make run-app      
```

### ___Db Migrations___
``` bash
make migration create_codes
make migrate-up
make migrate-down
```

### ___Env___
```
# Reloading env
direnv allow .
```

### __Env variables__
`.envrc`
```
export ADDR=:8080
export DB_ADDR=postgres://user:password@localhost/tower_defense?sslmode=disable
export ENV=dev
```

### __Documentation__
Open API
```
# DEV
http://localhost:8080/v1/swagger/index.html

# Production
 TBA
```