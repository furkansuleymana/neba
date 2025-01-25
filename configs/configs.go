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

//go:embed configs.json
var defaultConfig []byte

const (
	storageDirectory = "build" // TODO: fix
	configFileName   = "config.json"
)

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

type Manager struct {
	mutex  sync.RWMutex
	path   string
	config AppConfig
}

func ConfigManager() (*Manager, error) {
	cwd, err := os.Getwd() // get cwd first
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

func (cm *Manager) Get() AppConfig {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.config
}

func (cm *Manager) Update(updateFunc func(*AppConfig)) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	updateFunc(&cm.config)
	return cm.save()
}

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
