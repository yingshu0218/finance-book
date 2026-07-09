package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"ledger/internal/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "启动 Web 服务",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetInt("port")
		host, _ := cmd.Flags().GetString("host")

		fmt.Printf("启动 Web 服务: http://%s:%d\n", host, port)
		server.Start(host, port)
	},
}

func init() {
	serveCmd.Flags().IntP("port", "p", 8080, "服务端口")
	serveCmd.Flags().String("host", "localhost", "服务地址")
}
