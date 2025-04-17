package router

import (
	"embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

// SetRouter 設置所有路由
func SetRouter(router *gin.Engine, buildFS embed.FS, indexPage []byte) {
	// 設置 API 路由
	SetApiRouter(router)
	
	// 設置 Web 路由
	frontendBaseUrl := os.Getenv("FRONTEND_BASE_URL")
	if frontendBaseUrl == "" {
		// Explicitly unset FRONTEND_BASE_URL if it's empty
		os.Unsetenv("FRONTEND_BASE_URL")
		SetWebRouter(router, buildFS, indexPage)
	} else {
		frontendBaseUrl = strings.TrimSuffix(frontendBaseUrl, "/")
		router.NoRoute(func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s%s", frontendBaseUrl, c.Request.RequestURI))
		})
	}
}
