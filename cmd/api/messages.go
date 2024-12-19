package main

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"tower-defense-api/lib/json"
	"tower-defense-api/lib/models"
)

// CreateMessageHandler godoc
//
//	@Summary		Create a new message
//	@Description	Create a new message
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.CreateMessagePayload	true	"Message payload"
//	@Success		201		{object}	models.Message
//	@Router			/v1/messages [post]
//	@Security		ApiKeyAuth
func (app *application) createMessageHandler(w http.ResponseWriter, r *http.Request) {
	var payload models.CreateMessagePayload

	if err := json.ReadJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := json.Validate.Struct(payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	message := &models.Message{
		UserID:  payload.UserID,
		Content: payload.Content,
		Sender:  payload.Sender,
	}

	if err := app.repository.Messages.Create(r.Context(), message); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := json.JSONResponse(w, http.StatusCreated, message); err != nil {
		app.internalServerError(w, r, err)
	}
}

// GetMessagesHandler godoc
//
//	@Summary		Set Message as Read
//	@Description	Set Message as Read
//	@Tags			messages
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"Message ID"
//	@Success		200	
//	@Router			/v1/messages/{id} [put]
//	@Security		ApiKeyAuth
func (app *application) setReadHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		json.WriteJSONError(w, http.StatusBadRequest, "Invalid message ID")
		return
	}

	if err := app.repository.Messages.SetRead(r.Context(), int64(id)); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
