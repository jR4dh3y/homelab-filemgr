package model

// DriveNameMapping represents a custom name for a mount point
type DriveNameMapping struct {
	MountPoint string `json:"mountPoint"`
	CustomName string `json:"customName"`
}

// DriveNamesRequest represents a request to set a custom drive name
type DriveNamesRequest struct {
	MountPoint string `json:"mountPoint"`
	CustomName string `json:"customName"`
}

// DriveNamesResponse represents the response with all custom drive names
type DriveNamesResponse struct {
	Mappings []DriveNameMapping `json:"mappings"`
}
