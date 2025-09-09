// Package service provides business logic for the file manager.
package service

import (
	"context"
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/homelab/filemanager/internal/model"
	"github.com/homelab/filemanager/internal/pkg/filesystem"
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
	List(ctx context.Context, path string) ([]model.FileInfo, error)
	GetInfo(ctx context.Context, path string) (*model.FileInfo, error)
	CreateDir(ctx context.Context, path string) error
	Delete(ctx context.Context, path string) error
}

type fileService struct {
	fs       filesystem.FS
	basePath string
}

func NewFileService(fsys filesystem.FS, basePath string) FileService {
	return &fileService{fs: fsys, basePath: basePath}
}

func (s *fileService) List(ctx context.Context, path string) ([]model.FileInfo, error) {
	fullPath := filepath.Join(s.basePath, path)
	
	info, err := s.fs.Stat(fullPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, ErrPathNotFound
		}
		return nil, err
	}
	if !info.IsDir() {
		return nil, ErrNotDirectory
	}

	entries, err := s.fs.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	items := make([]model.FileInfo, 0, len(entries))
	for _, entry := range entries {
		entryInfo, _ := entry.Info()
		items = append(items, model.FileInfo{
			Name:    entry.Name(),
			Path:    filepath.Join(path, entry.Name()),
			Size:    entryInfo.Size(),
			IsDir:   entry.IsDir(),
			ModTime: entryInfo.ModTime(),
		})
	}

	return items, nil
}

func (s *fileService) GetInfo(ctx context.Context, path string) (*model.FileInfo, error) {
	fullPath := filepath.Join(s.basePath, path)
	info, err := s.fs.Stat(fullPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, ErrPathNotFound
		}
		return nil, err
	}
	return &model.FileInfo{
		Name:    info.Name(),
		Path:    path,
		Size:    info.Size(),
		IsDir:   info.IsDir(),
		ModTime: info.ModTime(),
	}, nil
}

func (s *fileService) CreateDir(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)
	return s.fs.MkdirAll(fullPath, 0755)
}

func (s *fileService) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)
	return s.fs.RemoveAll(fullPath)
}