package router

import (
	"account-system/middleware"
	"embed"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"net/http"
	"path" // 需要 path 來處理路徑
	"strings"
)

// EmbedFolder 實現 static.ServeFileSystem 接口
type EmbedFolder struct {
	fs   embed.FS
	root string
}

// Open 實現 http.FileSystem 接口的 Open 方法
func (f *EmbedFolder) Open(name string) (http.File, error) {
	// 將請求的路徑與根目錄結合
	fullName := path.Join(f.root, name)
	// 使用 http.FS 適配器來打開文件，確保返回 http.File
	return http.FS(f.fs).Open(fullName)
}

// Exists 實現 static.ServeFileSystem 接口的 Exists 方法
func (f *EmbedFolder) Exists(prefix string, filepath string) bool {
	// 組合完整路徑
	p := path.Join(f.root, filepath)
	// 嘗試打開文件或目錄
	if _, err := f.fs.Open(p); err != nil {
		return false // 打開失敗，表示不存在
	}
	return true // 打開成功，表示存在
}

// SetWebRouter 設置 Web 路由
func SetWebRouter(router *gin.Engine, buildFS embed.FS, indexPage []byte) {
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middleware.GlobalWebRateLimit())

	// 提供靜態文件，使用自定義的 EmbedFolder
	// 注意：這裡的 root 是相對於 buildFS 的根目錄
	router.Use(static.Serve("/", &EmbedFolder{fs: buildFS, root: "web/dist"}))

	// 處理所有未匹配的路由，返回前端 index.html
	router.NoRoute(func(c *gin.Context) {
		// 檢查是否是 API 或靜態資源路徑
		if strings.HasPrefix(c.Request.RequestURI, "/api") || strings.HasPrefix(c.Request.RequestURI, "/assets") {
			// 如果是靜態資源，讓 static.Serve 處理，這裡不應該攔截
			// 但如果 static.Serve 找不到，最終還是會到這裡，所以需要區分
			// 這裡假設 /assets 是前端資源的固定前綴
			if strings.HasPrefix(c.Request.RequestURI, "/assets") {
				// 嘗試讓 static handler 處理，如果到這裡表示找不到
				c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "靜態資源未找到"})
				return
			}
			// 如果是 API 路徑
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "API 路徑不存在"})
			return
		}
		// 其他所有路徑都返回 index.html
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexPage)
	})
}
