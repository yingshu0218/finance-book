package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"ledger/internal/app"
)

var bookCmd = &cobra.Command{
	Use:   "book",
	Short: "账本管理",
}

var bookCreateCmd = &cobra.Command{
	Use:   "create <名称>",
	Short: "创建新账本",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		err := app.CreateBook(name)
		if err != nil {
			fmt.Printf("创建失败: %v\n", err)
			return
		}
		fmt.Printf("账本 '%s' 创建成功\n", name)
	},
}

var bookListCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有账本",
	Run: func(cmd *cobra.Command, args []string) {
		books := app.ListBooks()
		for _, b := range books {
			marker := ""
			if b.Name == app.GetDefaultBook() {
				marker = " (默认)"
			}
			fmt.Printf("- %s%s\n", b.Name, marker)
		}
	},
}

var bookDeleteCmd = &cobra.Command{
	Use:   "delete <名称>",
	Short: "删除账本",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]
		err := app.DeleteBook(name)
		if err != nil {
			fmt.Printf("删除失败: %v\n", err)
			return
		}
		fmt.Printf("账本 '%s' 删除成功\n", name)
	},
}

func init() {
	bookCmd.AddCommand(bookCreateCmd)
	bookCmd.AddCommand(bookListCmd)
	bookCmd.AddCommand(bookDeleteCmd)
}
