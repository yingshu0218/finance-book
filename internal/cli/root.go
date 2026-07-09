package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"ledger/internal/app"
)

var rootCmd = &cobra.Command{
	Use:   "ledger",
	Short: "轻量级记账工具",
	Long:  "一个基于 Go 的轻量级记账软件，支持 CLI 和 Web 界面",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	if err := app.Init(); err != nil {
		fmt.Printf("初始化失败: %v\n", err)
		os.Exit(1)
	}

	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(balanceCmd)
	rootCmd.AddCommand(bookCmd)
	rootCmd.AddCommand(serveCmd)
}
