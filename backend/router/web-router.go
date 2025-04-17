package router

import (
	"account-system/middleware"
	"embed"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"io/fs" // 新增導入
	"net/http"
	"strings"
)

// SetWebRouter 設置 Web 路由
func SetWebRouter(router *gin.Engine, buildFS embed.FS, indexPage []byte) {
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.GlobalWebRateLimit())

	// 提取 "web/dist" 子文件系統
	subFS, err := fs.Sub(buildFS, "web/dist")
	if err != nil {
		// 處理錯誤，例如 panic 或 log
		panic("無法獲取嵌入式文件系統的子目錄: " + err.Error())
	}

	// 提供靜態文件，使用 http.FS 適配器
	router.Use(static.Serve("/", static.ServeFileSystem(http.FS(subFS))))

	// 處理所有未匹配的路由，返回前端 index.html
	router.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, "/api") || strings.HasPrefix(c.Request.RequestURI, "/assets") {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "API 路徑不存在",
			})
			return
		}
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexPage)
	})
}
