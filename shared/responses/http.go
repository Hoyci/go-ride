package responses

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func ParseValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range validationErrors {
			errors[fieldErr.Field()] = fieldErr.Tag()
		}
	}
	return errors
}
