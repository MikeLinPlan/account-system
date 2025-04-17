#!/bin/bash

# 進入後端目錄
cd backend

# 檢查 go.mod 文件是否存在
if [ -f "go.mod" ]; then
    echo "go.mod 文件已存在，嘗試下載依賴和整理模塊..."
    # 下載依賴並整理模塊
    go mod download
    go mod tidy
else
    echo "go.mod 文件不存在，初始化新的 Go 模塊..."
    # 初始化 Go 模塊
    go mod init account-system
    
    # 下載所有依賴
    go get -u github.com/google/uuid
    go get -u golang.org/x/crypto/bcrypt
    go get -u gorm.io/driver/mysql
    go get -u gorm.io/driver/sqlite
    go get -u gorm.io/gorm
    go get -u github.com/gin-contrib/sessions
    go get -u github.com/gin-contrib/sessions/cookie
    go get -u github.com/gin-gonic/gin
    go get -u golang.org/x/time/rate
    go get -u github.com/gin-contrib/gzip
    go get -u github.com/gin-contrib/static
    go get -u github.com/joho/godotenv
    
    # 整理模塊
    go mod tidy
fi

# 檢查 go.sum 文件是否為空
if [ ! -s "go.sum" ]; then
    echo "go.sum 文件為空，重新下載依賴..."
    # 下載所有依賴
    go get -u github.com/google/uuid
    go get -u golang.org/x/crypto/bcrypt
    go get -u gorm.io/driver/mysql
    go get -u gorm.io/driver/sqlite
    go get -u gorm.io/gorm
    go get -u github.com/gin-contrib/sessions
    go get -u github.com/gin-contrib/sessions/cookie
    go get -u github.com/gin-gonic/gin
    go get -u golang.org/x/time/rate
    go get -u github.com/gin-contrib/gzip
    go get -u github.com/gin-contrib/static
    go get -u github.com/joho/godotenv
    
    # 整理模塊
    go mod tidy
fi

echo "Go 模塊初始化完成！"
echo "go.mod 和 go.sum 文件已準備就緒！"
