package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

type APIError struct {
	StatusCode int `json:"statusCode"`
	Msg        any `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("api error: %d", e.StatusCode)
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
}

func InvalidRequestData(errors map[string]string) APIError {
	return APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Msg:        errors,
	}
}

func InvalidJSON() APIError {
	return NewAPIError(
		http.StatusBadRequest,
		fmt.Errorf("invalid JSON request data"),
	)
}

func ServerError() APIError {
	return NewAPIError(
		http.StatusInternalServerError,
		fmt.Errorf("internal server error"),
	)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func Make(h apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			apiErr, ok := err.(APIError)
			if !ok {
				slog.Error(
					"HTTP API error",
					"err",
					err.Error(),
					"path",
					r.URL.Path,
				)
				apiErr = ServerError()
			}
			WriteJSON(w, apiErr.StatusCode, apiErr)
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
