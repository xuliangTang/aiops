package app

import (
	"aipos/pkg/app/controllers"
	"aipos/pkg/configurations"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/xuliangTang/athena/athena"
	"net/http"
)

func NewApiServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "apiserver",
		RunE: func(cmd *cobra.Command, args []string) error {
			port, _ := cmd.Flags().GetInt("port")
			return run(port)
		},
	}

	// 添加flag
	cmd.Flags().Int("port", 8080, "apiserver 8080")
	return cmd
}

func run(port int) error {
	server := athena.Ignite().Configuration(
		configurations.NewK8sConfig(),
	).Mount(
		"", nil,
		controllers.NewResourcesController(),
		controllers.NewPromptController(),
		controllers.NewShellController(),
	)

	server.StaticFS("/html", http.Dir("./asserts/html"))     // 网页
	server.StaticFS("/static", http.Dir("./asserts/static")) // 静态资源 包含脚本 样式
	server.Launch()
	return nil
}

func cross() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		if method != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization,X-Token")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
	}
}
