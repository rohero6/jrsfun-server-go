package router

import (
	"github.com/gin-gonic/gin"
	"jrsfun-server-go/controller/api"
	"net/http"
	"strings"
)

func StartRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	// 使用自定义的中间件
	router.Use(NormalizePathMiddleware())
	router.Any("//api/home/", api.Home)
	router.Any("/api/live", api.Live)
	router.Any("/api/channel", api.Channel)

	return router
}
func NormalizePathMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求的原始路径
		originalPath := c.Request.URL.Path

		// 规范化路径（将多个斜杠替换为单个斜杠）
		normalizedPath := strings.Replace(originalPath, "//", "/", -1)

		// 检查路径是否已经被规范化，如果不同，则重定向到规范化后的路径
		if originalPath != normalizedPath {
			c.Redirect(http.StatusMovedPermanently, normalizedPath)
			return
		}

		// 继续处理下一个处理程序
		c.Next()
	}
}
