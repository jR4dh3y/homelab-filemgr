package model

// MountPoint represents a configured filesystem location accessible through the file manager
type MountPoint struct {
	Name         string `json:"name" mapstructure:"name"`
	Path         string `json:"path" mapstructure:"path"`
	ReadOnly     bool   `json:"readOnly" mapstructure:"read_only"`
	AutoDiscover bool   `json:"autoDiscover" mapstructure:"auto_discover"`
}

// ServerConfig contains all server configuration options
type ServerConfig struct {
	Port        int          `mapstructure:"port"`
	Host        string       `mapstructure:"host"`
	MountPoints []MountPoint `mapstructure:"mount_points"`
	JWTSecret   string       `mapstructure:"jwt_secret"`
	MaxUploadMB int          `mapstructure:"max_upload_mb"`
	ChunkSizeMB int          `mapstructure:"chunk_size_mb"`
}

// DefaultServerConfig returns sensible defaults for server configuration
func DefaultServerConfig() ServerConfig {
	return ServerConfig{
		Port:        8080,
		Host:        "0.0.0.0",
		MountPoints: []MountPoint{},
		JWTSecret:   "",
		MaxUploadMB: 10240, // 10GB
		ChunkSizeMB: 5,     // 5MB chunks
	}
}

// IsMountPointReadOnly checks if a mount point is read-only by name
func (c *ServerConfig) IsMountPointReadOnly(name string) bool {
	for _, mp := range c.MountPoints {
		if mp.Name == name {
			return mp.ReadOnly
		}
	}
	return true // Default to read-only if not found
}

// GetMountPoint returns a mount point by name, or nil if not found
func (c *ServerConfig) GetMountPoint(name string) *MountPoint {
	for i := range c.MountPoints {
		if c.MountPoints[i].Name == name {
			return &c.MountPoints[i]
		}
	}
	return nil
}
