package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"ccs/pkg/config"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有配置",
	Run:   runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) {
	logger := slog.Default()
	cfg, err := config.Load()
	if err != nil {
		logger.Error("加载配置失败", "error", err)
		os.Exit(1)
	}
	names := cfg.List()
	if len(names) == 0 {
		fmt.Println("暂无配置，请先使用 ccs add 添加")
		return
	}
	fmt.Println("可用配置:")
	for _, name := range names {
		fmt.Printf("  - %s\n", name)
	}
}