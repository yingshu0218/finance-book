package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Start(host string, port int) {
	r := gin.Default()

	setupAPI(r)
	setupStatic(r)

	addr := fmt.Sprintf("%s:%d", host, port)
	if err := r.Run(addr); err != nil {
		fmt.Printf("服务启动失败: %v\n", err)
	}
}

func setupStatic(r *gin.Engine) {
	r.StaticFile("/", "./web/dist/index.html")
	r.Static("/assets", "./web/dist/assets")
}
