# 帳號管理系統

這是一個精簡版的帳號管理系統，從原始項目中提取出來，只保留了與帳號管理相關的功能。

## 功能

- 用戶註冊
- 用戶登入
- 個人資料管理
- 權限控制
- Token 管理

## 目錄結構

```
account-system/
├── backend/                # 後端 Go 代碼
│   ├── common/             # 通用工具和常量
│   ├── controller/         # 控制器
│   ├── middleware/         # 中間件
│   ├── model/              # 數據模型
│   └── router/             # 路由配置
├── frontend/               # 前端 React 代碼
│   ├── src/
│   │   ├── components/     # React 組件
│   │   ├── context/        # React 上下文
│   │   ├── pages/          # 頁面組件
│   │   └── helpers/        # 工具函數
├── docker/                 # Docker 配置
│   ├── Dockerfile          # 應用 Dockerfile
│   └── docker-compose.yml  # Docker Compose 配置
└── README.md               # 項目說明
```

## 技術棧

- **後端**：Go (Gin 框架)
- **前端**：React (Vite)
- **數據庫**：MySQL
- **緩存**：Redis
- **容器化**：Docker & Docker Compose

## 快速開始

1. 克隆本倉庫

2. 進入 Docker 配置目錄：
   ```bash
   cd docker
   ```

3. 複製環境變數檔案並修改：
   ```bash
   cp .env.example .env
   ```
   編輯 `.env` 檔案，修改 `DB_PASSWORD` 和 `SESSION_SECRET` 為安全的隨機字符串。

4. 運行 Docker Compose：
   ```bash
   docker-compose up -d
   ```

5. 訪問系統：
   ```
   http://localhost:3000
   ```

6. 使用默認管理員帳號登入：
   - 用戶名：`root`
   - 密碼：`123456`

   **注意：**首次登入後請立即修改默認密碼！

## 故障排除

如果您在構建或運行過程中遇到問題，請參考 [TROUBLESHOOTING.md](TROUBLESHOOTING.md) 文件。

## API 文檔

### 認證 API

- `POST /api/user/register` - 註冊新用戶
- `POST /api/user/login` - 用戶登入
- `GET /api/user/logout` - 用戶登出
- `GET /api/user/self` - 獲取當前用戶信息
- `PUT /api/user/self` - 更新當前用戶信息
- `DELETE /api/user/self` - 刪除當前用戶

### Token API

- `GET /api/user/token` - 生成訪問令牌
- `GET /api/token/` - 獲取所有令牌
- `POST /api/token/` - 創建新令牌
- `PUT /api/token/` - 更新令牌
- `DELETE /api/token/:id` - 刪除令牌

### 管理員 API

- `GET /api/user/` - 獲取所有用戶
- `POST /api/user/` - 創建用戶
- `PUT /api/user/` - 更新用戶
- `DELETE /api/user/:id` - 刪除用戶
