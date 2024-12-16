package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"tower-defense-api/lib/json"
	"tower-defense-api/lib/models"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.CreateUserPayload

	if err := json.ReadJSON(w, r, &payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	user := &models.User{
		Username: payload.Username,
	}

	if err := app.repository.Users.Create(r.Context(), user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := json.JSONResponse(w, http.StatusCreated, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		json.WriteJSONError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	user, err := app.repository.Users.GetById(r.Context(), int64(id))

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := json.JSONResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}
