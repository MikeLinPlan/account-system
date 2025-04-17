package model

import (
	"account-system/common"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"time"
)

var DB *gorm.DB

// InitDB 初始化數據庫
func InitDB() error {
	// 根據環境變數決定使用 MySQL 還是 SQLite
	dsn := os.Getenv("SQL_DSN")
	var db *gorm.DB
	var err error

	if dsn == "" {
		// 使用 SQLite
		sqlitePath := os.Getenv("SQLITE_PATH")
		if sqlitePath == "" {
			sqlitePath = "data/account-system.db"
		}
		// 確保目錄存在
		dir := "data"
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return fmt.Errorf("failed to create directory: %v", err)
			}
		}
		db, err = gorm.Open(sqlite.Open(sqlitePath), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("failed to connect to SQLite: %v", err)
		}
		common.SysLog("using SQLite: " + sqlitePath)
	} else {
		// 使用 MySQL
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("failed to connect to MySQL: %v", err)
		}
		common.SysLog("using MySQL: " + dsn)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get DB instance: %v", err)
	}

	// 設置連接池
	maxIdleConns := common.GetIntEnv("SQL_MAX_IDLE_CONNS", 10)
	maxOpenConns := common.GetIntEnv("SQL_MAX_OPEN_CONNS", 100)
	maxLifetime := common.GetIntEnv("SQL_MAX_LIFETIME", 60)

	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(maxLifetime) * time.Second)

	DB = db

	// 自動遷移數據表結構
	err = db.AutoMigrate(&User{}, &Token{})
	if err != nil {
		return fmt.Errorf("failed to migrate database: %v", err)
	}

	// 創建根用戶帳號（如果需要）
	err = createRootAccountIfNeed()
	if err != nil {
		return fmt.Errorf("failed to create root account: %v", err)
	}

	return nil
}

// CloseDB 關閉數據庫連接
func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
