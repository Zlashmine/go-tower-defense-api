package main

import (
	"net/http"

	"tower-defense-api/lib/json"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("Internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	json.WriteJSONError(w, http.StatusInternalServerError, "Internal server error")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnw("Bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	json.WriteJSONError(w, http.StatusBadRequest, "Bad request")
}

// func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
// 	app.logger.Warnw("forbidden", "method", r.Method, "path", r.URL.Path, "error")

// 	json.WriteJSONError(w, http.StatusForbidden, "forbidden")
// }

// func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
// 	app.logger.Errorf("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error())

// 	json.WriteJSONError(w, http.StatusConflict, err.Error())
// }

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	json.WriteJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	json.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
}

// func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
// 	app.logger.Warnf("unauthorized basic error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

// 	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)

// 	json.WriteJSONError(w, http.StatusUnauthorized, "unauthorized")
// }

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path)

	w.Header().Set("Retry-After", retryAfter)

	json.WriteJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
