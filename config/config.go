// Package config provides configuration management for Gerry
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

var (
	// ErrNoConfigFile is returned when the config file does not exist
	ErrNoConfigFile = errors.New("configuration file not found")

	// ErrNoGerritCredentials is returned when Gerrit credentials are missing
	ErrNoGerritCredentials = errors.New("no Gerrit credentials found. Please set gerritUsername and gerritPassword in ~/.config/gerry.json")
)

// Config represents the configuration for Gerry
type Config struct {
	GerritUsername string `json:"gerritUsername,omitempty"`
	GerritPassword string `json:"gerritPassword,omitempty"`
}

// getPath returns the path to the configuration file
func getPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not determine home directory: %w", err)
	}
	return filepath.Join(home, ".config", "gerry.json"), nil
}

// Load loads the configuration from the config file
func Load() (*Config, error) {
	configPath, err := getPath()
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrNoConfigFile
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(content, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate that credentials are present
	if cfg.GerritUsername == "" || cfg.GerritPassword == "" {
		return nil, ErrNoGerritCredentials
	}

	return &cfg, nil
}
