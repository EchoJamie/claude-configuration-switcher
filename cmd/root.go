package cmd

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

var verbose bool
var configPath string

var rootCmd = &cobra.Command{
	Use:   "ccs",
	Short: "Claude Configuration Switcher",
	Long:  "CCS 是一个用于管理 Claude API 配置的 CLI 工具",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if cmd.Name() != "init" {
			initLogger()
		}
		return nil
	},
}

func initLogger() {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))
	slog.SetDefault(logger)
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "显示详细日志")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "自定义配置文件路径")
}