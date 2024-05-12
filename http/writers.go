package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

func WriteJson(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var err error
	if buf, ok := data.([]byte); ok {
		_, err = w.Write(buf)
	} else {
		err = json.NewEncoder(w).Encode(data)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func WriteStatus(w http.ResponseWriter, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
}

type ErrorResponse struct {
	Error string `json:"error" example:"an error occurred"`
}

func WriteError(w http.ResponseWriter, r *http.Request, status int, err error) {
	log := r.Context().Value("logger").(*slog.Logger)

	if status == 0 {
		status = http.StatusBadRequest
	}

	var msg string
	if err != nil {
		msg = err.Error()
		if strings.Contains(msg, "error: couldn't find ") {
			status = http.StatusNotFound
		}
	} else {
		msg = "unspecified error"
	}

	resp := ErrorResponse{
		Error: msg,
	}

	// suppress internal errors
	if status >= http.StatusInternalServerError {
		log.Error("Error", "error", resp.Error)
		resp.Error = "internal server error"
	} else {
		log.Warn("Error", "error", resp.Error)
	}

	WriteJson(w, status, resp)
}
