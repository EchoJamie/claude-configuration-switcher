package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"ccs/pkg/config"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "添加新的 Claude 配置",
	Run:   runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) {
	logger := slog.Default()
	cfg, err := config.Load()
	if err != nil {
		logger.Error("加载配置失败", "error", err)
		os.Exit(1)
	}
	var name string
	qs := &survey.Input{Message: "请输入配置名称:"}
	if err := survey.AskOne(qs, &name); err != nil {
		logger.Error("输入失败", "error", err)
		os.Exit(1)
	}
	if _, ok := cfg.Get(name); ok {
		fmt.Printf("配置 '%s' 已存在，是否覆盖? (y/N): ", name)
		var confirm string
		fmt.Scanln(&confirm)
		if confirm != "y" && confirm != "Y" {
			fmt.Println("取消添加")
			os.Exit(0)
		}
	}
	var apiKey, baseURL, defaultModel string
	survey.AskOne(&survey.Input{Message: "请输入 API Key:", Help: "格式: sk-xxx"}, &apiKey)
	if apiKey == "" {
		fmt.Println("API Key 不能为空，取消添加")
		os.Exit(1)
	}
	survey.AskOne(&survey.Input{Message: "请输入 Base URL:"}, &baseURL)
	if baseURL == "" {
		fmt.Println("Base URL 不能为空，取消添加")
		os.Exit(1)
	}
	survey.AskOne(&survey.Input{Message: "请输入默认模型 (opus/sonnet/haiku):"}, &defaultModel)
	newCfg := config.Config{
		Name:                       name,
		ANTHROPIC_API_KEY:         apiKey,
		ANTHROPIC_BASE_URL:        baseURL,
		ANTHROPIC_MODEL:           defaultModel,
	}
	survey.AskOne(&survey.Input{Message: "请输入 Opus 细分模型:"}, &newCfg.ANTHROPIC_DEFAULT_OPUS_MODEL)
	survey.AskOne(&survey.Input{Message: "请输入 Sonnet 细分模型:"}, &newCfg.ANTHROPIC_DEFAULT_SONNET_MODEL)
	survey.AskOne(&survey.Input{Message: "请输入 Haiku 细分模型:"}, &newCfg.ANTHROPIC_DEFAULT_HAIKU_MODEL)
	survey.AskOne(&survey.Input{Message: "请输入快速模型:"}, &newCfg.ANTHROPIC_SMALL_FAST_MODEL)

	cfg.Set(name, newCfg)
	if err := cfg.Save(); err != nil {
		logger.Error("保存配置失败", "error", err)
		os.Exit(1)
	}
	fmt.Printf("配置 '%s' 添加成功\n", name)
}