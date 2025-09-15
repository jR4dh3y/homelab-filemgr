package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/homelab/filemanager/internal/service"
)

type FileHandler struct {
	fileService service.FileService
}

func NewFileHandler(fileService service.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

func (h *FileHandler) RegisterRoutes(r chi.Router) {
	r.Get("/*", h.ListFiles)
}

func (h *FileHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	path := chi.URLParam(r, "*")
	
	files, err := h.fileService.List(r.Context(), path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}