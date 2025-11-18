package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/homelab/filemanager/internal/model"
	"github.com/homelab/filemanager/internal/service"
)

// SearchHandler handles search-related HTTP requests
type SearchHandler struct {
	searchService service.SearchService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(searchService service.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

// RegisterRoutes registers search routes on the given router
func (h *SearchHandler) RegisterRoutes(r chi.Router) {
	r.Get("/", h.Search)
}

// SearchResponse represents the search results response
type SearchResponse struct {
	Path    string           `json:"path"`
	Query   string           `json:"query"`
	Results []model.FileInfo `json:"results"`
	Count   int              `json:"count"`
}

// Search handles search requests
// GET /api/v1/search?path=&q=
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	path := r.URL.Query().Get("path")
	query := r.URL.Query().Get("q")

	// Validate query is not empty
	if query == "" {
		writeError(w, "Search query is required", model.ErrCodeValidationError, http.StatusBadRequest)
		return
	}

	// Default path to empty string (will search from root mount points)
	if path == "" {
		writeError(w, "Search path is required", model.ErrCodeValidationError, http.StatusBadRequest)
		return
	}

	// Perform search
	results, err := h.searchService.Search(r.Context(), path, query)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	// Return results
	resp := SearchResponse{
		Path:    path,
		Query:   query,
		Results: results,
		Count:   len(results),
	}

	writeJSON(w, resp, http.StatusOK)
}

// handleServiceError converts service errors to HTTP responses
func (h *SearchHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrEmptyQuery):
		writeError(w, "Search query cannot be empty", model.ErrCodeValidationError, http.StatusBadRequest)
	case errors.Is(err, service.ErrPathNotFound):
		writeError(w, "Path not found", model.ErrCodeNotFound, http.StatusNotFound)
	case errors.Is(err, service.ErrNotDirectory):
		writeError(w, "Path is not a directory", model.ErrCodeValidationError, http.StatusBadRequest)
	case errors.Is(err, service.ErrMountPointNotFound):
		writeError(w, "Mount point not found or access denied", model.ErrCodeAccessDenied, http.StatusForbidden)
	default:
		writeError(w, "Internal server error", model.ErrCodeInternalError, http.StatusInternalServerError)
	}
}
