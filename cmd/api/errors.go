package main

import (
	"fmt"
	"net/http"

	"tower-defense-api/lib/json"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Printf("Internal server error: %s - %s - %s\n", r.Method, r.URL.Path, err)

	json.WriteJSONError(w, http.StatusInternalServerError, "Internal server error")
}

// func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
// 	fmt.Printf("Not found: %s - %s\n", r.Method, r.URL.Path)

// 	json.WriteJSONError(w, http.StatusNotFound, "Resource not found")
// }

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Printf("Bad request: %s - %s - %s\n", r.Method, r.URL.Path, err)

	json.WriteJSONError(w, http.StatusBadRequest, "Bad request")
}
