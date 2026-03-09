package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/EchoJamie/ccs/pkg/config"

	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <name>",
	Short: "切换 Claude 配置",
	Args:  cobra.ExactArgs(1),
	Run:   runUse,
}

func init() {
	rootCmd.AddCommand(useCmd)
}

func runUse(cmd *cobra.Command, args []string) {
	logger := slog.Default()
	name := args[0]
	cfg, err := config.Load()
	if err != nil {
		logger.Error("加载配置失败", "error", err)
		os.Exit(1)
	}
	profile, ok := cfg.Get(name)
	if !ok {
		fmt.Printf("配置 '%s' 不存在\n", name)
		os.Exit(1)
	}
	settingsPath := filepath.Join(".claude", "settings.local.json")
	settings := map[string]interface{}{}
	if data, err := os.ReadFile(settingsPath); err == nil {
		json.Unmarshal(data, &settings)
	}
	env := map[string]string{}
	if profile.ANTHROPIC_API_KEY != "" {
		env["ANTHROPIC_API_KEY"] = profile.ANTHROPIC_API_KEY
	}
	if profile.ANTHROPIC_BASE_URL != "" {
		env["ANTHROPIC_BASE_URL"] = profile.ANTHROPIC_BASE_URL
	}
	if profile.ANTHROPIC_MODEL != "" {
		env["ANTHROPIC_MODEL"] = profile.ANTHROPIC_MODEL
	}
	if profile.ANTHROPIC_DEFAULT_OPUS_MODEL != "" {
		env["ANTHROPIC_DEFAULT_OPUS_MODEL"] = profile.ANTHROPIC_DEFAULT_OPUS_MODEL
	}
	if profile.ANTHROPIC_DEFAULT_SONNET_MODEL != "" {
		env["ANTHROPIC_DEFAULT_SONNET_MODEL"] = profile.ANTHROPIC_DEFAULT_SONNET_MODEL
	}
	if profile.ANTHROPIC_DEFAULT_HAIKU_MODEL != "" {
		env["ANTHROPIC_DEFAULT_HAIKU_MODEL"] = profile.ANTHROPIC_DEFAULT_HAIKU_MODEL
	}
	if profile.ANTHROPIC_SMALL_FAST_MODEL != "" {
		env["ANTHROPIC_SMALL_FAST_MODEL"] = profile.ANTHROPIC_SMALL_FAST_MODEL
	}
	settings["env"] = env
	if err := os.MkdirAll(".claude", 0755); err != nil {
		logger.Error("创建目录失败", "error", err)
		os.Exit(1)
	}
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		logger.Error("序列化配置失败", "error", err)
		os.Exit(1)
	}
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		logger.Error("写入配置文件失败", "error", err)
		os.Exit(1)
	}
	fmt.Printf("已切换到配置 '%s'\n", name)
	fmt.Println("配置文件:", settingsPath)
}