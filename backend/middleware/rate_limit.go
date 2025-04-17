package middleware

import (
	"account-system/common"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

// 速率限制器映射
var (
	apiLimiters    = make(map[string]*rate.Limiter)
	apiLimitersMux sync.Mutex

	webLimiters    = make(map[string]*rate.Limiter)
	webLimitersMux sync.Mutex

	criticalLimiters    = make(map[string]*rate.Limiter)
	criticalLimitersMux sync.Mutex
)

// 獲取 API 速率限制器
func getApiLimiter(key string) *rate.Limiter {
	apiLimitersMux.Lock()
	defer apiLimitersMux.Unlock()

	limiter, exists := apiLimiters[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(float64(common.GlobalApiRateLimitNum)/float64(common.GlobalApiRateLimitDuration)), common.GlobalApiRateLimitNum)
		apiLimiters[key] = limiter
	}

	return limiter
}

// 獲取 Web 速率限制器
func getWebLimiter(key string) *rate.Limiter {
	webLimitersMux.Lock()
	defer webLimitersMux.Unlock()

	limiter, exists := webLimiters[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(float64(common.GlobalWebRateLimitNum)/float64(common.GlobalWebRateLimitDuration)), common.GlobalWebRateLimitNum)
		webLimiters[key] = limiter
	}

	return limiter
}

// 獲取關鍵操作速率限制器
func getCriticalLimiter(key string) *rate.Limiter {
	criticalLimitersMux.Lock()
	defer criticalLimitersMux.Unlock()

	limiter, exists := criticalLimiters[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Limit(float64(common.CriticalRateLimitNum)/float64(common.CriticalRateLimitDuration)), common.CriticalRateLimitNum)
		criticalLimiters[key] = limiter
	}

	return limiter
}

// 清理過期的限制器
func cleanupLimiters() {
	for {
		time.Sleep(time.Hour)

		apiLimitersMux.Lock()
		apiLimiters = make(map[string]*rate.Limiter)
		apiLimitersMux.Unlock()

		webLimitersMux.Lock()
		webLimiters = make(map[string]*rate.Limiter)
		webLimitersMux.Unlock()

		criticalLimitersMux.Lock()
		criticalLimiters = make(map[string]*rate.Limiter)
		criticalLimitersMux.Unlock()
	}
}

// GlobalAPIRateLimit 全局 API 速率限制
func GlobalAPIRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !common.GlobalApiRateLimitEnable {
			c.Next()
			return
		}

		ip := c.ClientIP()
		limiter := getApiLimiter(ip)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "請求過於頻繁，請稍後再試",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// GlobalWebRateLimit 全局 Web 速率限制
func GlobalWebRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !common.GlobalWebRateLimitEnable {
			c.Next()
			return
		}

		ip := c.ClientIP()
		limiter := getWebLimiter(ip)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "請求過於頻繁，請稍後再試",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// CriticalRateLimit 關鍵操作速率限制
func CriticalRateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getCriticalLimiter(ip)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "關鍵操作請求過於頻繁，請稍後再試",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// 初始化速率限制器清理
func init() {
	go cleanupLimiters()
}
