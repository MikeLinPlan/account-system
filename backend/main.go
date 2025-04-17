package main

import (
	"account-system/common"
	"account-system/model"
	"account-system/router"
	"embed"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

//go:embed web/dist
var buildFS embed.FS

//go:embed web/dist/index.html
var indexPage []byte

func main() {
	// 加載 .env 文件
	err := godotenv.Load(".env")
	if err != nil {
		common.SysLog("Support for .env file is disabled: " + err.Error())
	}

	// 加載環境變數
	common.LoadEnv()

	// 設置日誌
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	common.SysLog("Account System started")

	// 設置 Gin 模式
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化數據庫
	err = model.InitDB()
	if err != nil {
		common.FatalLog("failed to initialize database: " + err.Error())
	}
	defer func() {
		err := model.CloseDB()
		if err != nil {
			common.FatalLog("failed to close database: " + err.Error())
		}
	}()

	// 初始化 HTTP 服務器
	server := gin.New()
	server.Use(gin.Logger())
	server.Use(gin.Recovery())

	// 初始化會話存儲
	store := cookie.NewStore([]byte(common.SessionSecret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   2592000, // 30 天
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
	server.Use(sessions.Sessions("session", store))

	// 設置路由
	router.SetRouter(server, buildFS, indexPage)

	// 獲取端口
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// 啟動服務器
	common.SysLog(fmt.Sprintf("Server is running on port %s", port))
	err = server.Run(":" + port)
	if err != nil {
		common.FatalLog("failed to start HTTP server: " + err.Error())
	}
}
