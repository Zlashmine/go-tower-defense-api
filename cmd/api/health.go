package main

import (
	"net/http"

	"tower-defense-api/lib/json"
)

// HealthCheckHandler godoc
//
//	@Summary		Health Check
//	@Description	Health Check
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Router			/v1/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"version": version,
		"env":     app.config.env,
	}

	if err := json.JSONResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
