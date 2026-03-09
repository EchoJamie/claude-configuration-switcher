package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/EchoJamie/ccs/pkg/config"
	"github.com/EchoJamie/ccs/pkg/importer"
	"github.com/EchoJamie/ccs/pkg/util"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [path]",
	Short: "从 cc-switch 导入配置",
	Args:  cobra.RangeArgs(0, 1),
	Run:   runImport,
}

func init() {
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) {
	logger := slog.Default()
	var dbPath string
	if len(args) > 0 {
		dbPath = util.ExpandPath(args[0])
	} else {
		var input string
		qs := &survey.Input{Message: "请输入 cc-switch 数据库路径:", Default: "~/.cc-switch/cc-switch.db"}
		if err := survey.AskOne(qs, &input); err != nil {
			logger.Error("输入失败", "error", err)
			os.Exit(1)
		}
		dbPath = util.ExpandPath(input)
	}
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Printf("数据库文件不存在: %s\n", dbPath)
		os.Exit(1)
	}
	configs, err := importer.LoadFromCCSwitch(dbPath)
	if err != nil {
		logger.Error("加载配置失败", "error", err)
		os.Exit(1)
	}
	if len(configs) == 0 {
		fmt.Println("未找到可导入的配置")
		os.Exit(0)
	}
	var selected []string
	qs := &survey.MultiSelect{Message: "选择要导入的配置:", Options: buildImportOptions(configs)}
	if err := survey.AskOne(qs, &selected); err != nil {
		logger.Error("选择失败", "error", err)
		os.Exit(1)
	}
	cfg, err := config.Load()
	if err != nil {
		logger.Error("加载配置失败", "error", err)
		os.Exit(1)
	}
	importSelectedConfigs(cfg, configs, selected)
	if err := cfg.Save(); err != nil {
		logger.Error("保存配置失败", "error", err)
		os.Exit(1)
	}
	fmt.Printf("已导入 %d 个配置\n", len(selected))
	fmt.Println("当前可用配置:")
	for _, name := range cfg.List() {
		fmt.Printf("  - %s\n", name)
	}
}

func buildImportOptions(configs []importer.CCSwitchConfig) []string {
	names := make([]string, len(configs))
	for i, c := range configs {
		suffix := ""
		if c.APIKey == "" || c.BaseURL == "" {
			suffix = " (缺失: "
			if c.APIKey == "" {
				suffix += "APIKey"
			}
			if c.BaseURL == "" {
				if c.APIKey == "" {
					suffix += ", "
				}
				suffix += "BaseURL"
			}
			suffix += ")"
		}
		names[i] = c.Name + suffix
	}
	return names
}

func importSelectedConfigs(cfg *config.ConfigStore, configs []importer.CCSwitchConfig, selected []string) {
	for _, name := range selected {
		actualName := name
		if idx := findConfigByName(configs, name); idx >= 0 {
			actualName = configs[idx].Name
		}

		var sc *importer.CCSwitchConfig
		for i := range configs {
			if configs[i].Name == actualName {
				sc = &configs[i]
				break
			}
		}
		if sc == nil {
			continue
		}
		if sc.APIKey == "" {
			fmt.Printf("配置 '%s' 缺少 APIKey，跳过导入\n", actualName)
			continue
		}
		newCfg := config.Config{
			Name:                           actualName,
			ANTHROPIC_API_KEY:             sc.APIKey,
			ANTHROPIC_BASE_URL:            sc.BaseURL,
			ANTHROPIC_MODEL:               sc.Model,
			ANTHROPIC_DEFAULT_OPUS_MODEL:  getEnvStringFromMap(sc.Env, "ANTHROPIC_DEFAULT_OPUS_MODEL"),
			ANTHROPIC_DEFAULT_SONNET_MODEL: getEnvStringFromMap(sc.Env, "ANTHROPIC_DEFAULT_SONNET_MODEL"),
			ANTHROPIC_DEFAULT_HAIKU_MODEL:  getEnvStringFromMap(sc.Env, "ANTHROPIC_DEFAULT_HAIKU_MODEL"),
			ANTHROPIC_SMALL_FAST_MODEL:     getEnvStringFromMap(sc.Env, "ANTHROPIC_SMALL_FAST_MODEL"),
		}
		if _, ok := cfg.Get(actualName); ok {
			var confirm bool
			qs := &survey.Confirm{Message: fmt.Sprintf("配置 '%s' 已存在，是否覆盖?", actualName), Default: true}
			if err := survey.AskOne(qs, &confirm); err != nil {
				continue
			}
			if !confirm {
				continue
			}
		}
		cfg.Set(actualName, newCfg)
	}
}

func findConfigByName(configs []importer.CCSwitchConfig, displayName string) int {
	for i, c := range configs {
		if c.Name == displayName || displayName == c.Name+" (缺失: APIKey)" || displayName == c.Name+" (缺失: BaseURL)" || displayName == c.Name+" (缺失: APIKey, BaseURL)" {
			return i
		}
	}
	return -1
}

func getEnvStringFromMap(m map[string]interface{}, key string) string {
	if m == nil {
		return ""
	}
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}