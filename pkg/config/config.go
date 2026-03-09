package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/EchoJamie/ccs/pkg/util"
)

type Config struct {
	Name                       string `json:"name,omitempty"`
	ANTHROPIC_API_KEY         string `json:"ANTHROPIC_API_KEY,omitempty"`
	ANTHROPIC_BASE_URL        string `json:"ANTHROPIC_BASE_URL,omitempty"`
	ANTHROPIC_MODEL           string `json:"ANTHROPIC_MODEL,omitempty"`
	ANTHROPIC_DEFAULT_OPUS_MODEL    string `json:"ANTHROPIC_DEFAULT_OPUS_MODEL,omitempty"`
	ANTHROPIC_DEFAULT_SONNET_MODEL  string `json:"ANTHROPIC_DEFAULT_SONNET_MODEL,omitempty"`
	ANTHROPIC_DEFAULT_HAIKU_MODEL   string `json:"ANTHROPIC_DEFAULT_HAIKU_MODEL,omitempty"`
	ANTHROPIC_SMALL_FAST_MODEL      string `json:"ANTHROPIC_SMALL_FAST_MODEL,omitempty"`
}

type ConfigStore struct {
	Profiles map[string]Config `json:"profiles"`
}

func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}
	return filepath.Join(home, ".ccs", "config.json")
}

func GetConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = os.Getenv("HOME")
	}
	return filepath.Join(home, ".ccs")
}

func Load() (*ConfigStore, error) {
	path := GetConfigPath()
	if !util.FileExists(path) {
		return &ConfigStore{Profiles: make(map[string]Config)}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	var store ConfigStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	if store.Profiles == nil {
		store.Profiles = make(map[string]Config)
	}
	return &store, nil
}

func (s *ConfigStore) Save() error {
	dir := GetConfigDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	path := GetConfigPath()
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	return nil
}

func (s *ConfigStore) Get(name string) (*Config, bool) {
	cfg, ok := s.Profiles[name]
	return &cfg, ok
}

func (s *ConfigStore) Set(name string, cfg Config) {
	if s.Profiles == nil {
		s.Profiles = make(map[string]Config)
	}
	s.Profiles[name] = cfg
}

func (s *ConfigStore) Delete(name string) {
	delete(s.Profiles, name)
}

func (s *ConfigStore) List() []string {
	names := make([]string, 0, len(s.Profiles))
	for name := range s.Profiles {
		names = append(names, name)
	}
	return names
}
