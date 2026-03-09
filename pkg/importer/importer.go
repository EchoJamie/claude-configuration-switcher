package importer

import (
	"database/sql"
	"encoding/json"
	"fmt"

	_ "modernc.org/sqlite"
)

type CCSwitchConfig struct {
	ID        int
	Name      string
	APIKey    string
	BaseURL   string
	Model     string
	Enabled   bool
	CreatedAt string
	UpdatedAt string
	Env       map[string]interface{}
}

// SettingsConfig cc-switch 的 settings_config JSON 结构
type SettingsConfig struct {
	Env map[string]interface{} `json:"env"`
}

// DefaultModels 默认模型
var DefaultModels = map[string]string{
	"opus":   "opus-4-5-20250514",
	"sonnet": "sonnet-4-20250514",
	"haiku":  "haiku-3-20250514",
}

func LoadFromCCSwitch(dbPath string) ([]CCSwitchConfig, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}
	defer db.Close()

	// 查询 providers 表获取配置
	rows, err := db.Query("SELECT name, settings_config, is_current FROM providers WHERE app_type = 'claude'")
	if err != nil {
		return nil, fmt.Errorf("查询配置失败: %w", err)
	}
	defer rows.Close()

	var configs []CCSwitchConfig
	for rows.Next() {
		var c CCSwitchConfig
		var settingsConfig string
		var isCurrent int
		if err := rows.Scan(&c.Name, &settingsConfig, &isCurrent); err != nil {
			return nil, fmt.Errorf("扫描配置失败: %w", err)
		}

		// 解析 settings_config JSON
		var sc SettingsConfig
		if err := json.Unmarshal([]byte(settingsConfig), &sc); err != nil {
			continue
		}
		if sc.Env == nil {
			continue
		}

		// 从 env 中提取配置
		c.APIKey = getEnvString(sc.Env, "ANTHROPIC_AUTH_TOKEN")
		c.BaseURL = getEnvStringDefault(sc.Env, "ANTHROPIC_BASE_URL", "https://api.anthropic.com")
		c.Model = getEnvStringDefault(sc.Env, "ANTHROPIC_MODEL", "opus")
		c.Enabled = isCurrent == 1
		c.Env = sc.Env

		configs = append(configs, c)
	}
	return configs, nil
}

func getEnvString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getEnvStringDefault(m map[string]interface{}, key, def string) string {
	s := getEnvString(m, key)
	if s != "" {
		return s
	}
	return def
}
