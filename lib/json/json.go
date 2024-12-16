package json

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

const maxBytes int64 = 1_048_576

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func WriteJSON(w http.ResponseWriter, statusCode int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		err := json.NewEncoder(w).Encode(data)

		if err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	}

	return nil
}

func ReadJSON(w http.ResponseWriter, r *http.Request, data any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBytes)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func WriteJSONError(w http.ResponseWriter, statusCode int, message string) error {
	type errorResponse struct {
		Error string `json:"error"`
	}

	return WriteJSON(w, statusCode, &errorResponse{Error: message})
}

func JSONResponse(w http.ResponseWriter, statusCode int, data any) error {
	type response struct {
		Data any `json:"data"`
	}

	return WriteJSON(w, statusCode, &response{Data: data})
}
