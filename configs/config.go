package configs

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

//go:embed config.json
var defaultConfig []byte

const (
	// Configuration constants
	storageDirectory = "build" // TODO: fix
	configFileName   = "config.json"
)

// AppConfig holds the application's configuration settings as defined in the
// "config.json" file. This includes settings such as the server address and
// port. It does not contain any device-specific configurations like
// IP addresses or serial numbers.
type AppConfig struct {
	Server struct {
		HTTP struct {
			Address string `json:"address"`
			Port    string `json:"port"`
		} `json:"http"`
	} `json:"server"`
	Database struct {
		Path string `json:"path"`
	} `json:"database"`
}

// Manager is a struct that manages the configuration of the application.
// It includes a mutex for read/write locking, a path to the configuration file,
// and the application configuration itself.
type Manager struct {
	mutex  sync.RWMutex
	path   string
	config AppConfig
}

// ConfigManager initializes and returns a new configuration manager.
// It retrieves the current working directory, constructs the configuration
// file path, ensures the configuration file exists, and loads the configuration.
// If any step fails, an error is returned.
//
// Returns:
//   - *Manager: A pointer to the initialized configuration manager.
//   - error: An error if any step in the initialization process fails.
func ConfigManager() (*Manager, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get working directory: %w", err)
	}

	configPath := filepath.Join(cwd, storageDirectory, configFileName)

	manager := &Manager{path: configPath}

	if err := manager.ensureConfigFile(); err != nil {
		return nil, err
	}

	if err := manager.load(); err != nil {
		return nil, err
	}

	slog.Debug("created config manager",
		slog.String("cwd", cwd),
		slog.String("path", configPath))
	return manager, nil
}

// Get retrieves the current application configuration in a thread-safe manner.
// It acquires a read lock on the mutex to ensure that the configuration is not
// modified while being accessed, and releases the lock once the configuration
// is returned.
//
// Returns:
//   - AppConfig: The current application configuration.
func (cm *Manager) Get() AppConfig {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.config
}

// Update applies the provided update function to the AppConfig instance
// managed by the Manager. It ensures that the update is performed in a
// thread-safe manner by acquiring a lock before applying the update and
// releasing it afterward. After the update function is executed, the
// updated configuration is saved.
//
// Parameters:
//   - updateFunc: A function that takes a pointer to an AppConfig instance
//     and performs the desired updates.
//
// Returns:
//   - error: An error if the save operation fails, otherwise nil.
func (cm *Manager) Update(updateFunc func(*AppConfig)) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	updateFunc(&cm.config)
	return cm.save()
}

// ensureConfigFile ensures that the configuration file exists at the
// specified path. If the file already exists, it logs a debug message and
// returns nil. If the file does not exist, it creates the necessary directories
// and writes the default configuration to the file. It returns an error if any
// of these operations fail.
//
// Returns:
//   - error: An error if the file creation or writing fails.
func (cm *Manager) ensureConfigFile() error {
	if _, err := os.Stat(cm.path); err == nil {
		slog.Debug("config file already exists",
			slog.String("path", cm.path))
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(cm.path), 0755); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	if err := os.WriteFile(cm.path, defaultConfig, 0644); err != nil {
		return fmt.Errorf("create config file: %w", err)
	}

	slog.Debug("created new config file",
		slog.String("path", cm.path))
	return nil
}

// load reads the configuration file from the specified path and unmarshals its content
// into the Manager's config field. It returns an error if the file cannot be read or
// if the content cannot be parsed as JSON.
//
// Returns:
//   - error: An error if there is an issue reading the file or parsing its content.
func (cm *Manager) load() error {
	data, err := os.ReadFile(cm.path)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	if err := json.Unmarshal(data, &cm.config); err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}

	slog.Debug("loaded config",
		slog.String("path", cm.path))
	return nil
}

// save serializes the Manager's configuration to JSON format and writes it to a file.
// It returns an error if the serialization or file writing fails.
//
// Returns:
//   - error: An error if the serialization or file writing fails.
func (cm *Manager) save() error {
	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	if err := os.WriteFile(cm.path, data, 0644); err != nil {
		return fmt.Errorf("write config file: %w", err)
	}

	slog.Debug("saved config",
		slog.String("path", cm.path))
	return nil
}
