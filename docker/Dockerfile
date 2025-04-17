FROM oven/bun:latest AS frontend-builder

WORKDIR /build
# 複製前端依賴文件
COPY frontend/package.json .
RUN bun install
# 複製前端源代碼
COPY ./frontend .
# 構建前端
RUN DISABLE_ESLINT_PLUGIN='true' bun run build

FROM golang:alpine AS backend-builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /build

# 安裝必要的工具
RUN apk add --no-cache git

# 複製後端源代碼
COPY ./backend .

# 初始化 Go 模块並下載依賴
RUN go mod init account-system \
    && go get -u github.com/google/uuid \
    && go get -u golang.org/x/crypto/bcrypt \
    && go get -u gorm.io/driver/mysql \
    && go get -u gorm.io/driver/sqlite \
    && go get -u gorm.io/gorm \
    && go get -u github.com/gin-contrib/sessions \
    && go get -u github.com/gin-contrib/sessions/cookie \
    && go get -u github.com/gin-gonic/gin \
    && go get -u golang.org/x/time/rate \
    && go get -u github.com/gin-contrib/gzip \
    && go get -u github.com/gin-contrib/static \
    && go get -u github.com/joho/godotenv \
    && go mod tidy

# 複製前端構建產物
COPY --from=frontend-builder /build/dist ./web/dist

# 構建後端
RUN go build -ldflags "-s -w" -o account-system

FROM alpine

RUN apk update \
    && apk upgrade \
    && apk add --no-cache ca-certificates tzdata \
    && update-ca-certificates

COPY --from=backend-builder /build/account-system /
EXPOSE 3000
WORKDIR /data
ENTRYPOINT ["/account-system"]
