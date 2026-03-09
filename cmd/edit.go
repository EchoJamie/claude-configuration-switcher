package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"ccs/pkg/config"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "编辑配置文件",
	Run:   runEdit,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func runEdit(cmd *cobra.Command, args []string) {
	path := config.GetConfigPath()
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}
	execCmd := exec.Command(editor, path)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	if err := execCmd.Run(); err != nil {
		fmt.Printf("打开编辑器失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("配置文件已编辑:", path)
}