package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/homelab/filemanager/internal/model"
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
func Load(configPath string) (*model.ServerConfig, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("port", 80)
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

	var cfg model.ServerConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Parse FM_USERS_* environment variables into the Users map
	// Viper's AutomaticEnv doesn't handle map types from env vars like FM_USERS_username=password
	if cfg.Users == nil {
		cfg.Users = make(map[string]string)
	}
	for _, env := range os.Environ() {
		if strings.HasPrefix(env, "FM_USERS_") {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				username := strings.TrimPrefix(parts[0], "FM_USERS_")
				password := parts[1]
				if username != "" && password != "" {
					cfg.Users[username] = password
				}
			}
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

	return &cfg, nil
}
