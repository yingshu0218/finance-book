package cli

import (
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"ledger/internal/app"
)

var listCmd = &cobra.Command{
	Use:   "list [账本]",
	Short: "列出记账记录",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		bookName := app.GetDefaultBook()
		if len(args) > 0 {
			bookName = args[0]
		}

		month, _ := cmd.Flags().GetString("month")
		category, _ := cmd.Flags().GetString("category")

		entries, err := app.ListEntries(bookName, month, category)
		if err != nil {
			fmt.Printf("查询失败: %v\n", err)
			return
		}

		w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\t日期\t分类\t金额\t备注")
		for _, e := range entries {
			sign := ""
			if e.Amount > 0 {
				sign = "+"
			}
			fmt.Fprintf(w, "%d\t%s\t%s\t%s%.2f\t%s\n",
				e.ID, e.Date, e.Category, sign, e.Amount, e.Note)
		}
		w.Flush()
	},
}

func init() {
	listCmd.Flags().StringP("month", "m", "", "月份筛选 (YYYY-MM)")
	listCmd.Flags().StringP("category", "c", "", "分类筛选")
}
