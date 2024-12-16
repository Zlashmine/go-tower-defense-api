package main

import (
	"net/http"

	"tower-defense-api/lib/json"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"version": version,
	}

	if err := json.JSONResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}
}
