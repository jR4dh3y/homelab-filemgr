package config

import (
	"fmt"
	"strings"

	"github.com/homelab/filemanager/internal/model"
	"github.com/spf13/viper"
)

// Load reads configuration from file and environment variables
func Load(configPath string) (*model.ServerConfig, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("port", 8080)
	v.SetDefault("host", "0.0.0.0")
	v.SetDefault("max_upload_mb", 10240) // 10GB default
	v.SetDefault("chunk_size_mb", 5)     // 5MB chunks

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

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &cfg, nil
}
