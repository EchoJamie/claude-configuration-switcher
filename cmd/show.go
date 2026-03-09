package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"ccs/pkg/config"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "显示当前配置",
	Run:   runShow,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func runShow(cmd *cobra.Command, args []string) {
	logger := slog.Default()
	settingsPath := filepath.Join(".claude", "settings.local.json")
	settings := map[string]interface{}{}
	if data, err := os.ReadFile(settingsPath); err != nil {
		fmt.Println("未找到配置文件，请先使用 ccs use <name> 切换配置")
		os.Exit(0)
	} else {
		json.Unmarshal(data, &settings)
	}
	env, ok := settings["env"].(map[string]interface{})
	if !ok {
		fmt.Println("配置文件中没有 env 字段")
		os.Exit(0)
	}
	cfg, err := config.Load()
	if err != nil {
		logger.Error("加载配置失败", "error", err)
		os.Exit(1)
	}
	apiKey, _ := env["ANTHROPIC_API_KEY"].(string)
	if apiKey == "" {
		fmt.Println("当前配置中没有 ANTHROPIC_API_KEY")
		os.Exit(0)
	}
	var currentName string
	for _, c := range cfg.Profiles {
		if c.ANTHROPIC_API_KEY == apiKey {
			currentName = c.Name
			break
		}
	}
	if currentName == "" {
		fmt.Println("未找到匹配的配置")
		os.Exit(0)
	}
	fmt.Printf("当前配置: %s\n", currentName)
}