package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"tower-defense-api/lib/json"
	"tower-defense-api/lib/models"
	"tower-defense-api/lib/repository"
)

// CreateUserHandler godoc
//
//	@Summary		Create a new user
//	@Description	Create a new user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		models.CreateUserPayload	true	"User payload"
//	@Success		201		{object}	models.User
//	@Security		ApiKeyAuth
//	@Router			/v1/users [post]
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

// GetUserHandler godoc
//
//	@Summary		Get a user by ID
//	@Description	Get a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	models.User
//	@Security		ApiKeyAuth
//	@Router			/v1/users/{id} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)

	if err != nil {
		json.WriteJSONError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	ctx := r.Context()
	user, err := app.getUserFromCacheOr(ctx, int64(id))

	if err != nil {
		switch err {
		case repository.ErrNotFound:
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	messages, err := app.repository.Messages.GetByPlayerId(ctx, int64(user.ID))

	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	user.Messages = messages

	if err := app.setUserToCache(ctx, user); err != nil {
		app.internalServerError(w, r, err)
	}

	if err := json.JSONResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getUserFromCacheOr(context context.Context, id int64) (*models.User, error) {
	if !app.config.redisConfig.enabled {
		return app.repository.Users.GetById(context, id)
	}

	user, _ := app.cacheStore.Users.Get(context, id)

	if user != nil {
		return user, nil
	}

	return app.repository.Users.GetById(context, id)
}

func (app *application) setUserToCache(context context.Context, user *models.User) error {
	if !app.config.redisConfig.enabled {
		return nil
	}

	return app.cacheStore.Users.Set(context, user)
}
