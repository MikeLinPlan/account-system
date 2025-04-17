package common

import (
	"crypto/rand"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
)

// GetUUID 生成 UUID
func GetUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

// GetRandomString 生成隨機字符串
func GetRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			log.Println("failed to generate random string:", err)
			return GetUUID()[:length]
		}
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

// Password2Hash 將密碼轉換為哈希值
func Password2Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// ValidatePasswordAndHash 驗證密碼與哈希值
func ValidatePasswordAndHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetIntEnv 獲取整數環境變數
func GetIntEnv(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		SysError(fmt.Sprintf("failed to parse %s: %v, using default value %d", key, err, defaultValue))
		return defaultValue
	}
	return value
}

// GetBoolEnv 獲取布爾環境變數
func GetBoolEnv(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		SysError(fmt.Sprintf("failed to parse %s: %v, using default value %v", key, err, defaultValue))
		return defaultValue
	}
	return value
}

// SysLog 系統日誌
func SysLog(message string) {
	log.Printf("[INFO] %s\n", message)
}

// SysError 系統錯誤日誌
func SysError(message string) {
	log.Printf("[ERROR] %s\n", message)
}

// FatalLog 致命錯誤日誌
func FatalLog(message string) {
	log.Fatalf("[FATAL] %s\n", message)
}

// LoadEnv 加載環境變數
func LoadEnv() {
	// 從環境變數加載配置
	PasswordLoginEnabled = GetBoolEnv("PASSWORD_LOGIN_ENABLED", true)
	PasswordRegisterEnabled = GetBoolEnv("PASSWORD_REGISTER_ENABLED", true)
	EmailVerificationEnabled = GetBoolEnv("EMAIL_VERIFICATION_ENABLED", false)
	RegisterEnabled = GetBoolEnv("REGISTER_ENABLED", true)
	EmailDomainRestrictionEnabled = GetBoolEnv("EMAIL_DOMAIN_RESTRICTION_ENABLED", false)
	EmailAliasRestrictionEnabled = GetBoolEnv("EMAIL_ALIAS_RESTRICTION_ENABLED", false)

	// 加載速率限制配置
	GlobalApiRateLimitEnable = GetBoolEnv("GLOBAL_API_RATE_LIMIT_ENABLE", false)
	GlobalApiRateLimitNum = GetIntEnv("GLOBAL_API_RATE_LIMIT_NUM", 60)
	GlobalApiRateLimitDuration = int64(GetIntEnv("GLOBAL_API_RATE_LIMIT_DURATION", 60))

	GlobalWebRateLimitEnable = GetBoolEnv("GLOBAL_WEB_RATE_LIMIT_ENABLE", false)
	GlobalWebRateLimitNum = GetIntEnv("GLOBAL_WEB_RATE_LIMIT_NUM", 60)
	GlobalWebRateLimitDuration = int64(GetIntEnv("GLOBAL_WEB_RATE_LIMIT_DURATION", 60))

	CriticalRateLimitNum = GetIntEnv("CRITICAL_RATE_LIMIT_NUM", 20)
	CriticalRateLimitDuration = int64(GetIntEnv("CRITICAL_RATE_LIMIT_DURATION", 1200))

	// 如果環境變數中有 SESSION_SECRET，則使用它
	envSessionSecret := os.Getenv("SESSION_SECRET")
	if envSessionSecret != "" {
		SessionSecret = envSessionSecret
	}
}
