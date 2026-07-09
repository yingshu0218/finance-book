package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type BookInfo struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	CreatedAt string `json:"created_at"`
}

type Config struct {
	DefaultBook string     `json:"default_book"`
	Books       []BookInfo `json:"books"`
}

func GetConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ledger")
}

func GetConfigPath() string {
	return filepath.Join(GetConfigDir(), "config.json")
}

func GetBooksDir() string {
	return filepath.Join(GetConfigDir(), "books")
}

func Load() (*Config, error) {
	path := GetConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return createDefaultConfig()
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	if err := os.MkdirAll(GetConfigDir(), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(GetConfigPath(), data, 0644)
}

func createDefaultConfig() (*Config, error) {
	cfg := &Config{
		DefaultBook: "default",
		Books: []BookInfo{
			{
				Name:      "default",
				Path:      filepath.Join(GetConfigDir(), "books", "default.db"),
				CreatedAt: "2024-01-01",
			},
		},
	}
	if err := cfg.Save(); err != nil {
		return nil, err
	}
	return cfg, nil
}
