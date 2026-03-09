package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/EchoJamie/ccs/pkg/config"

	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename <old> <new>",
	Short: "重命名配置",
	Args:  cobra.ExactArgs(2),
	Run:   runRename,
}

func init() {
	rootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, args []string) {
	logger := slog.Default()
	oldName := args[0]
	newName := args[1]
	cfg, err := config.Load()
	if err != nil {
		logger.Error("加载配置失败", "error", err)
		os.Exit(1)
	}
	profile, ok := cfg.Get(oldName)
	if !ok {
		fmt.Printf("配置 '%s' 不存在\n", oldName)
		os.Exit(1)
	}
	if _, ok := cfg.Get(newName); ok {
		fmt.Printf("配置 '%s' 已存在\n", newName)
		os.Exit(1)
	}
	profile.Name = newName
	cfg.Delete(oldName)
	cfg.Set(newName, *profile)
	if err := cfg.Save(); err != nil {
		logger.Error("保存配置失败", "error", err)
		os.Exit(1)
	}
	fmt.Printf("配置 '%s' 已重命名为 '%s'\n", oldName, newName)
}