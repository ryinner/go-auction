package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

const internalErrorMsg = "internal server error"

type errorResponse struct {
	Error string `json:"error"`
}

// Можно прокидывать логи и сразу ошибки писать, вместо реализации в каждом handler
type responder struct {
	log *slog.Logger
}

func (rs responder) JSON(w http.ResponseWriter, r *http.Request, code int, v any) {
	body, err := json.Marshal(v)
	if err != nil {
		rs.log.ErrorContext(r.Context(), "encode response", "err", err)
		code = http.StatusInternalServerError
		body = []byte(`{"error":"` + internalErrorMsg + `"}`)
	}

	// Если задать до WriteHeader, ничего не будет
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if _, err := w.Write(body); err != nil {
		rs.log.WarnContext(r.Context(), "write response", "err", err)
	}
}

func (rs responder) Error(w http.ResponseWriter, r *http.Request, code int, msg string) {
	rs.JSON(w, r, code, errorResponse{Error: msg})
}

func (rs responder) Internal(w http.ResponseWriter, r *http.Request, msg string, err error) {
	rs.log.ErrorContext(r.Context(), msg, "err", err)
	rs.Error(w, r, http.StatusInternalServerError, internalErrorMsg)
}
