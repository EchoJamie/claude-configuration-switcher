package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/EchoJamie/ccs/pkg/config"
	"github.com/EchoJamie/ccs/pkg/util"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化配置文件",
	Run:   runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	logger := slog.Default()
	path := config.GetConfigPath()
	if util.FileExists(path) {
		fmt.Println("配置文件已存在:", path)
		os.Exit(0)
	}
	cfg := config.ConfigStore{Profiles: make(map[string]config.Config)}
	if err := cfg.Save(); err != nil {
		logger.Error("创建配置文件失败", "error", err)
		os.Exit(1)
	}
	fmt.Println("配置文件已创建:", path)
}