# 故障排除指南

## 構建錯誤：缺少 go.sum 條目或 go.mod 已存在

如果您在構建 Docker 映像時遇到類似以下的錯誤：

```
missing go.sum entry for module providing package github.com/google/uuid (imported by account-system/common)
```

或者：

```
go: /build/go.mod already exists
```

這是因為 Go 模塊系統需要正確設置 `go.mod` 和 `go.sum` 文件。

### 解決方案 1：生成正確的 go.sum 文件

如果 `go.mod` 文件已存在，但 `go.sum` 文件為空或不完整，您可以按照以下步驟生成正確的 `go.sum` 文件：

1. 進入後端目錄：
   ```bash
   cd backend
   ```

2. 下載所有依賴並生成 go.sum：
   ```bash
   go mod download
   go mod tidy
   ```

3. 確認 go.sum 文件已經生成：
   ```bash
   cat go.sum
   ```

4. 返回上一級目錄：
   ```bash
   cd ..
   ```

5. 現在您可以構建 Docker 映像：
   ```bash
   docker-compose up -d
   ```

### 解決方案 2：修改 Dockerfile

我們已經修改了 Dockerfile，使其先複製 `go.mod` 和 `go.sum` 文件，然後下載依賴，最後才複製其他源代碼。這程序符合 Go 模塊的最佳實踐，並利用 Docker 的層級快取機制。

如果依然遇到問題，請確保您使用的是最新版本的 Dockerfile。

### 解決方案 3：完全重新初始化 Go 模塊

如果上述方法仍然不起作用，您可以完全重新初始化 Go 模塊：

1. 進入後端目錄：
   ```bash
   cd backend
   ```

2. 刪除現有的 go.mod 和 go.sum 文件：
   ```bash
   rm go.mod go.sum
   ```

3. 初始化新的 Go 模塊：
   ```bash
   go mod init account-system
   ```

4. 下載所有依賴：
   ```bash
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
   ```

5. 整理模塊：
   ```bash
   go mod tidy
   ```

6. 返回上一級目錄：
   ```bash
   cd ..
   ```

7. 現在您可以構建 Docker 映像：
   ```bash
   docker-compose up -d
   ```

### 解決方案 4：使用初始化腳本

我們提供了一個初始化腳本來自動執行上述步驟：

1. 確保腳本有執行權限：
   ```bash
   chmod +x init-go-modules.sh
   ```

2. 運行腳本：
   ```bash
   ./init-go-modules.sh
   ```

3. 構建 Docker 映像：
   ```bash
   docker-compose up -d
   ```

## 其他常見問題

### 1. 無法連接到數據庫

如果應用程序無法連接到數據庫，請檢查：

- `.env` 文件中的 `DB_PASSWORD` 是否與 Docker Compose 中的 `MYSQL_ROOT_PASSWORD` 一致
- 確保 MySQL 容器已經啟動並正常運行：
  ```bash
  docker ps | grep mysql
  ```
- 檢查 MySQL 容器的日誌：
  ```bash
  docker logs mysql
  ```

### 2. 無法連接到 Redis

如果應用程序無法連接到 Redis，請檢查：

- 確保 Redis 容器已經啟動並正常運行：
  ```bash
  docker ps | grep redis
  ```
- 檢查 Redis 容器的日誌：
  ```bash
  docker logs redis
  ```

### 3. 前端無法訪問 API

如果前端無法訪問 API，請檢查：

- 確保 API 服務已經啟動並正常運行：
  ```bash
  docker ps | grep account-service
  ```
- 檢查 API 服務的日誌：
  ```bash
  docker logs account-service
  ```
- 確保前端的 API 請求路徑正確
