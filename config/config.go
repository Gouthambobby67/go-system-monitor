package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// AppConfig holds the application configuration
type AppConfig struct {
	CPUThreshold       float64 `json:"cpu_threshold"`
	MemoryThreshold    float64 `json:"memory_threshold"`
	DiskThreshold      float64 `json:"disk_threshold"`
	SwapThreshold      float64 `json:"swap_threshold"`
	RefreshInterval    int     `json:"refresh_interval_ms"`
	MaxProcesses       int     `json:"max_processes"`
	MaxAlertsToKeep    int     `json:"max_alerts_to_keep"`
	DefaultSortingMode string  `json:"default_sorting_mode"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() AppConfig {
	return AppConfig{
		CPUThreshold:       85.0,
		MemoryThreshold:    85.0,
		DiskThreshold:      90.0,
		SwapThreshold:      80.0,
		RefreshInterval:    1000, // 1 second in milliseconds
		MaxProcesses:       15,
		MaxAlertsToKeep:    100,
		DefaultSortingMode: "cpu",
	}
}

// LoadConfig loads the configuration from file or returns default if not found
func LoadConfig() (AppConfig, error) {
	config := DefaultConfig()
	
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return config, fmt.Errorf("couldn't get home directory: %v", err)
	}
	
	// Create config directory if it doesn't exist
	configDir := filepath.Join(homeDir, ".config", "sysmon")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return config, fmt.Errorf("couldn't create config directory: %v", err)
	}
	
	configFile := filepath.Join(configDir, "config.json")
	
	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Create default config file
		if err := SaveConfig(configFile, config); err != nil {
			return config, fmt.Errorf("couldn't create default config file: %v", err)
		}
		return config, nil
	}
	
	// Read existing config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return config, fmt.Errorf("couldn't read config file: %v", err)
	}
	
	// Unmarshal JSON into config struct
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("invalid config file format: %v", err)
	}
	
	return config, nil
}

// SaveConfig saves the configuration to the specified file
func SaveConfig(filePath string, config AppConfig) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("couldn't marshal config: %v", err)
	}
	
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("couldn't write config file: %v", err)
	}
	
	return nil
}
