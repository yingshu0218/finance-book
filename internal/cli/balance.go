package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"ledger/internal/app"
)

var balanceCmd = &cobra.Command{
	Use:   "balance [账本]",
	Short: "查看收支统计",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bookName := app.GetDefaultBook()
		if len(args) > 0 {
			bookName = args[0]
		}

		month, _ := cmd.Flags().GetString("month")

		result, err := app.GetBalance(bookName, month)
		if err != nil {
			fmt.Printf("查询失败: %v\n", err)
			return
		}

		fmt.Printf("收入: ¥%.2f\n", result.Income)
		fmt.Printf("支出: ¥%.2f\n", result.Expense)
		fmt.Printf("余额: ¥%.2f\n", result.Balance)
	},
}

func init() {
	balanceCmd.Flags().StringP("month", "m", "", "月份筛选 (YYYY-MM)")
}
