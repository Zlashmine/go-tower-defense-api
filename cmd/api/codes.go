package main

import (
	"net/http"

	"tower-defense-api/lib/json"
	"tower-defense-api/lib/models"
)

func (app *application) createCodeHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.CreateCodePayload

	if err := json.ReadJSON(w, r, &payload); err != nil {
		app.badRequest(w, r, err)
		return
	}

	code := &models.Code{
		Code: payload.Code,
		Item: payload.Item,
	}

	if err := app.repository.Codes.Create(r.Context(), code); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := json.JSONResponse(w, http.StatusCreated, code); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getAllCodesHandler(w http.ResponseWriter, r *http.Request) {
	codes, err := app.repository.Codes.GetAll(r.Context())

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := json.JSONResponse(w, http.StatusOK, codes); err != nil {
		app.internalServerError(w, r, err)
	}
}
