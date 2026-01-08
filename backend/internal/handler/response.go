package handler

import (
	"encoding/json"
	"net/http"

	"github.com/homelab/filemanager/internal/model"
)

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Code    string `json:"code"`
	Details string `json:"details,omitempty"`
}

// writeError writes an error response with the specified message, code, and status
func writeError(w http.ResponseWriter, message, code string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error: message,
		Code:  code,
	})
}

// writeErrorWithDetails writes an error response with additional details
func writeErrorWithDetails(w http.ResponseWriter, message, code, details string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{
		Error:   message,
		Code:    code,
		Details: details,
	})
}

// writeJSON writes a JSON response with the specified data and status code
func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeSuccess writes a simple success response
func writeSuccess(w http.ResponseWriter, message string) {
	writeJSON(w, map[string]string{"message": message}, http.StatusOK)
}

// writeCreated writes a 201 Created response with the specified data
func writeCreated(w http.ResponseWriter, data interface{}) {
	writeJSON(w, data, http.StatusCreated)
}

// writeNoContent writes a 204 No Content response
func writeNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// writeBadRequest writes a 400 Bad Request error response
func writeBadRequest(w http.ResponseWriter, message string) {
	writeError(w, message, model.ErrCodeValidationError, http.StatusBadRequest)
}

// writeUnauthorized writes a 401 Unauthorized error response
func writeUnauthorized(w http.ResponseWriter, message string) {
	writeError(w, message, model.ErrCodeUnauthorized, http.StatusUnauthorized)
}

// writeForbidden writes a 403 Forbidden error response
func writeForbidden(w http.ResponseWriter, message string) {
	writeError(w, message, model.ErrCodePermissionDenied, http.StatusForbidden)
}

// writeNotFound writes a 404 Not Found error response
func writeNotFound(w http.ResponseWriter, message string) {
	writeError(w, message, model.ErrCodeNotFound, http.StatusNotFound)
}

// writeConflict writes a 409 Conflict error response
func writeConflict(w http.ResponseWriter, message string) {
	writeError(w, message, model.ErrCodeConflict, http.StatusConflict)
}

// writeInternalError writes a 500 Internal Server Error response
func writeInternalError(w http.ResponseWriter, message string) {
	writeError(w, message, model.ErrCodeInternalError, http.StatusInternalServerError)
}
