package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// MountPoint represents a configured filesystem location accessible through the file manager
type MountPoint struct {
	Name     string `json:"name" mapstructure:"name"`
	Path     string `json:"path" mapstructure:"path"`
	ReadOnly bool   `json:"readOnly" mapstructure:"read_only"`
}

// ServerConfig holds all server configuration
type ServerConfig struct {
	Port          int          `mapstructure:"port"`
	Host          string       `mapstructure:"host"`
	MountPoints   []MountPoint `mapstructure:"mount_points"`
	JWTSecret     string       `mapstructure:"jwt_secret"`
	MaxUploadMB   int          `mapstructure:"max_upload_mb"`
	ChunkSizeMB   int          `mapstructure:"chunk_size_mb"`
	AdminUsername string       `mapstructure:"admin_username"`
	AdminPassword string       `mapstructure:"admin_password"`
}

// Load reads configuration from file and environment variables
func Load(configPath string) (*ServerConfig, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("port", 8080)
	v.SetDefault("host", "0.0.0.0")
	v.SetDefault("max_upload_mb", 10240) // 10GB default
	v.SetDefault("chunk_size_mb", 5)     // 5MB chunks
	v.SetDefault("admin_username", "admin")
	v.SetDefault("admin_password", "") // Must be set via env var

	// Config file settings
	if configPath != "" {
		v.SetConfigFile(configPath)
	} else {
		v.SetConfigName("config")
		v.SetConfigType("yaml")
		v.AddConfigPath(".")
		v.AddConfigPath("./config")
		v.AddConfigPath("/etc/filemanager")
	}

	// Environment variable settings
	v.SetEnvPrefix("FM")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found is okay, we'll use defaults and env vars
	}

	var cfg ServerConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Validate checks that the configuration is valid
func (c *ServerConfig) Validate() error {
	if c.JWTSecret == "" {
		return fmt.Errorf("jwt_secret is required")
	}

	if len(c.MountPoints) == 0 {
		return fmt.Errorf("at least one mount_point is required")
	}

	for i, mp := range c.MountPoints {
		if mp.Name == "" {
			return fmt.Errorf("mount_point[%d].name is required", i)
		}
		if mp.Path == "" {
			return fmt.Errorf("mount_point[%d].path is required", i)
		}
	}

	if c.Port < 1 || c.Port > 65535 {
		return fmt.Errorf("port must be between 1 and 65535")
	}

	if c.MaxUploadMB < 1 {
		return fmt.Errorf("max_upload_mb must be at least 1")
	}

	if c.ChunkSizeMB < 1 {
		return fmt.Errorf("chunk_size_mb must be at least 1")
	}

	if c.AdminUsername == "" {
		return fmt.Errorf("admin_username is required")
	}

	if c.AdminPassword == "" {
		return fmt.Errorf("admin_password is required (set FM_ADMIN_PASSWORD environment variable)")
	}

	return nil
}

// GetMountPoint returns the mount point for a given name, or nil if not found
func (c *ServerConfig) GetMountPoint(name string) *MountPoint {
	for i := range c.MountPoints {
		if c.MountPoints[i].Name == name {
			return &c.MountPoints[i]
		}
	}
	return nil
}

// IsMountPointReadOnly checks if a mount point is read-only
func (c *ServerConfig) IsMountPointReadOnly(name string) bool {
	mp := c.GetMountPoint(name)
	if mp == nil {
		return true // Default to read-only for unknown mounts
	}
	return mp.ReadOnly
}
