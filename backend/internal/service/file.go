// Package service provides business logic for the file manager.
package service

import (
	"context"
	"errors"

	"github.com/homelab/filemanager/internal/model"
)

// File service errors
var (
	ErrPathNotFound     = errors.New("path not found")
	ErrPathExists       = errors.New("path already exists")
	ErrNotDirectory     = errors.New("path is not a directory")
	ErrPermissionDenied = errors.New("permission denied")
)

// FileService defines the file operations service interface
type FileService interface {
	// List returns files in a directory
	List(ctx context.Context, path string) ([]model.FileInfo, error)
	// GetInfo returns metadata for a file or directory
	GetInfo(ctx context.Context, path string) (*model.FileInfo, error)
}

// fileService implements FileService
type fileService struct {
	basePath string
}

// NewFileService creates a new file service
func NewFileService(basePath string) FileService {
	return &fileService{basePath: basePath}
}

// List returns files in a directory
func (s *fileService) List(ctx context.Context, path string) ([]model.FileInfo, error) {
	// TODO: Implement directory listing
	return nil, nil
}

// GetInfo returns metadata for a file or directory
func (s *fileService) GetInfo(ctx context.Context, path string) (*model.FileInfo, error) {
	// TODO: Implement file info retrieval
	return nil, nil
}