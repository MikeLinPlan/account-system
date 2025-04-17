#!/bin/bash

# 進入後端目錄
cd backend

# 初始化 Go 模塊
go mod tidy

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

# 再次整理模塊
go mod tidy

echo "Go 模塊初始化完成！"
