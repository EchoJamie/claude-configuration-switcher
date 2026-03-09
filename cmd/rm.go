package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/EchoJamie/ccs/pkg/config"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm <name>",
	Short: "删除配置",
	Args:  cobra.ExactArgs(1),
	Run:   runRm,
}

func init() {
	rootCmd.AddCommand(rmCmd)
}

func runRm(cmd *cobra.Command, args []string) {
	logger := slog.Default()
	name := args[0]
	cfg, err := config.Load()
	if err != nil {
		logger.Error("加载配置失败", "error", err)
		os.Exit(1)
	}
	if _, ok := cfg.Get(name); !ok {
		fmt.Printf("配置 '%s' 不存在\n", name)
		os.Exit(1)
	}
	var confirm bool
	qs := &survey.Confirm{Message: fmt.Sprintf("确定要删除配置 '%s' 吗?", name), Default: false}
	if err := survey.AskOne(qs, &confirm); err != nil {
		logger.Error("确认失败", "error", err)
		os.Exit(1)
	}
	if !confirm {
		fmt.Println("取消删除")
		os.Exit(0)
	}
	cfg.Delete(name)
	if err := cfg.Save(); err != nil {
		logger.Error("保存配置失败", "error", err)
		os.Exit(1)
	}
	fmt.Printf("配置 '%s' 已删除\n", name)
}