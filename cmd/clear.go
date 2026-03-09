package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "清理当前配置",
	Run:   runClear,
}

func init() {
	rootCmd.AddCommand(clearCmd)
}

func runClear(cmd *cobra.Command, args []string) {
	settingsPath := filepath.Join(".claude", "settings.local.json")
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		fmt.Println("未找到配置文件")
		return
	}
	settings := map[string]interface{}{}
	if data, err := os.ReadFile(settingsPath); err == nil {
		json.Unmarshal(data, &settings)
	}
	delete(settings, "env")
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		fmt.Println("序列化配置失败")
		os.Exit(1)
	}
	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		fmt.Println("写入配置文件失败")
		os.Exit(1)
	}
	fmt.Println("已清理当前配置")
}