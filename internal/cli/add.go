package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"ledger/internal/app"
)

var addCmd = &cobra.Command{
	Use:   "add [账本]",
	Short: "添加记账记录",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bookName := app.GetDefaultBook()
		if len(args) > 0 {
			bookName = args[0]
		}

		amount, _ := cmd.Flags().GetFloat64("amount")
		category, _ := cmd.Flags().GetString("category")
		date, _ := cmd.Flags().GetString("date")
		note, _ := cmd.Flags().GetString("note")

		err := app.AddEntry(bookName, amount, category, date, note)
		if err != nil {
			fmt.Printf("添加失败: %v\n", err)
			return
		}

		fmt.Println("添加成功")
	},
}

func init() {
	addCmd.Flags().Float64P("amount", "a", 0, "金额（正数收入，负数支出）")
	addCmd.Flags().StringP("category", "c", "", "分类")
	addCmd.Flags().StringP("date", "d", "", "日期 (YYYY-MM-DD)")
	addCmd.Flags().StringP("note", "n", "", "备注")

	addCmd.MarkFlagRequired("amount")
	addCmd.MarkFlagRequired("category")
}
