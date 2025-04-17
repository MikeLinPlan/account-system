package router

import (
	"account-system/middleware"
	"embed"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// EmbedFolder 嵌入式文件夾
type EmbedFolder struct {
	fs   embed.FS
	root string
}

// Open 打開文件
func (f *EmbedFolder) Open(name string) (http.File, error) {
	return f.fs.Open(f.root + name)
}

// SetWebRouter 設置 Web 路由
func SetWebRouter(router *gin.Engine, buildFS embed.FS, indexPage []byte) {
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.GlobalWebRateLimit())
	
	// 提供靜態文件
	router.Use(static.Serve("/", &EmbedFolder{fs: buildFS, root: "web/dist"}))
	
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
