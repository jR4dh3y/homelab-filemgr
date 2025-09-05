package model

import "time"

// FileInfo represents metadata for a file or directory
type FileInfo struct {
	Name        string    `json:"name"`
	Path        string    `json:"path"`
	Size        int64     `json:"size"`
	IsDir       bool      `json:"isDir"`
	ModTime     time.Time `json:"modTime"`
	Permissions string    `json:"permissions"`
	MimeType    string    `json:"mimeType,omitempty"`
}

// FileList represents a paginated list of files in a directory
type FileList struct {
	Path       string     `json:"path"`
	Items      []FileInfo `json:"items"`
	TotalCount int        `json:"totalCount"`
	Page       int        `json:"page"`
	PageSize   int        `json:"pageSize"`
}

// ListOptions contains options for listing directory contents
type ListOptions struct {
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	SortBy   string `json:"sortBy"`
	SortDir  string `json:"sortDir"`
	Filter   string `json:"filter"`
}

// DefaultListOptions returns sensible defaults for listing
func DefaultListOptions() ListOptions {
	return ListOptions{
		Page:     1,
		PageSize: 50,
		SortBy:   "name",
		SortDir:  "asc",
		Filter:   "",
	}
}
