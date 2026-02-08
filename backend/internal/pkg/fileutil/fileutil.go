// Package fileutil provides shared file utility functions.
package fileutil

import (
	"io/fs"
	"mime"
	"path/filepath"

	"github.com/homelab/filemanager/internal/model"
)

// ToFileInfo converts fs.FileInfo to model.FileInfo
// This is a centralized utility function used by file service and search service
func ToFileInfo(name, path string, info fs.FileInfo) model.FileInfo {
	fileInfo := model.FileInfo{
		Name:        name,
		Path:        path,
		Size:        info.Size(),
		IsDir:       info.IsDir(),
		ModTime:     info.ModTime(),
		Permissions: info.Mode().String(),
	}

	// Set MIME type for files
	if !info.IsDir() {
		ext := filepath.Ext(name)
		if ext != "" {
			mimeType := mime.TypeByExtension(ext)
			if mimeType != "" {
				fileInfo.MimeType = mimeType
			}
		}
	}

	return fileInfo
}

// DetectMimeType returns the MIME type for a file based on its extension
// Returns "application/octet-stream" if the type cannot be determined
func DetectMimeType(filename string) string {
	ext := filepath.Ext(filename)
	if ext != "" {
		mimeType := mime.TypeByExtension(ext)
		if mimeType != "" {
			return mimeType
		}
	}
	return "application/octet-stream"
}
