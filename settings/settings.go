package settings

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

//go:embed settings.json
var defaultSettings []byte

const (
	storageDirectory = "build" // TODO: fix
	settingsFileName = "settings.json"
)

type ApplicationSettings struct {
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

type SettingsManager struct {
	mutex    sync.RWMutex
	path     string
	settings ApplicationSettings
}

func NewSettingsManager() (*SettingsManager, error) {
	cwd, err := os.Getwd() // get cwd first
	if err != nil {
		return nil, fmt.Errorf("get working directory: %w", err)
	}

	settingsPath := filepath.Join(cwd, storageDirectory, settingsFileName)

	manager := &SettingsManager{path: settingsPath}

	if err := manager.ensureSettingsFile(); err != nil {
		return nil, err
	}

	if err := manager.load(); err != nil {
		return nil, err
	}

	slog.Info("created settings manager",
		slog.String("cwd", cwd),
		slog.String("path", settingsPath))
	return manager, nil
}

func (cm *SettingsManager) Get() ApplicationSettings {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return cm.settings
}

func (cm *SettingsManager) Update(updateFunc func(*ApplicationSettings)) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	updateFunc(&cm.settings)
	return cm.save()
}

func (cm *SettingsManager) ensureSettingsFile() error {
	if _, err := os.Stat(cm.path); err == nil {
		slog.Info("settings file already exists",
			slog.String("path", cm.path))
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(cm.path), 0755); err != nil {
		return fmt.Errorf("create settings directory: %w", err)
	}

	if err := os.WriteFile(cm.path, defaultSettings, 0644); err != nil {
		return fmt.Errorf("create settings file: %w", err)
	}

	slog.Info("created new settings file",
		slog.String("path", cm.path))
	return nil
}

func (cm *SettingsManager) load() error {
	data, err := os.ReadFile(cm.path)
	if err != nil {
		return fmt.Errorf("read settings file: %w", err)
	}

	if err := json.Unmarshal(data, &cm.settings); err != nil {
		return fmt.Errorf("parse settings file: %w", err)
	}

	slog.Info("loaded settings",
		slog.String("path", cm.path))
	return nil
}

func (cm *SettingsManager) save() error {
	data, err := json.MarshalIndent(cm.settings, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal settings: %w", err)
	}

	if err := os.WriteFile(cm.path, data, 0644); err != nil {
		return fmt.Errorf("write settings file: %w", err)
	}

	slog.Info("saved settings",
		slog.String("path", cm.path))
	return nil
}
