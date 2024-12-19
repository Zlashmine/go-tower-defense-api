package main

import (
	"net/http"

	"tower-defense-api/lib/json"
	"tower-defense-api/lib/models"
)

// CreateCodeHandler godoc
//
//	@Summary		Create a new code
//	@Description	Create a new code
//	@Tags			codes
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.CreateCodePayload	true	"Code payload"
//	@Success		201		{object}	models.Code
//	@Router			/v1/codes [post]
//	@Security		ApiKeyAuth
func (app *application) createCodeHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.CreateCodePayload

	if err := json.ReadJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		app.internalServerError(w, r, err)
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

// GetAllCodesHandler godoc
//
//	@Summary		Get all codes
//	@Description	Get all codes
//	@Tags			codes
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.Code
//	@Router			/v1/codes [get]
//	@Security		ApiKeyAuth
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
