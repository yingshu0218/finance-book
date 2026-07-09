package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
