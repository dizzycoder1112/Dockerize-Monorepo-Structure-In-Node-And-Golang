# Go Counter Server

Rails Counter Server 的 Go 重寫版本，使用 Gin 框架。

## 📁 專案結構

```
go-counter-server/
├── cmd/
│   └── main.go                 # 應用程式入口
├── internal/                   # 私有應用程式代碼
│   ├── config/                 # 配置管理
│   │   └── config.go
│   ├── middleware/             # 中介軟體
│   │   ├── cors.go            # CORS 設定
│   │   └── logger.go          # 請求日誌
│   ├── handlers/               # HTTP 處理器 (Controller)
│   │   ├── health.go          # Health check
│   │   ├── authentication.go  # 認證相關 (TODO)
│   │   ├── girls.go           # Girls API (TODO)
│   │   └── ...                # 其他 controllers
│   ├── models/                 # 資料模型 (Phase 2)
│   ├── services/               # 業務邏輯層 (Phase 3)
│   ├── repository/             # 資料存取層 (Phase 2)
│   └── router/                 # 路由設定
│       └── router.go
├── pkg/                        # 公開函式庫
│   └── response/               # 統一回應格式
│       └── response.go
├── mock/                       # Mock 資料 (Phase 1)
├── .env
├── .env.example
├── go.mod
└── go.sum
```

## 🚀 快速開始

### 1. 安裝依賴
```bash
go mod download
```

### 2. 設定環境變數
```bash
cp .env.example .env
# 編輯 .env 檔案設定你的環境變數
```

### 3. 執行開發伺服器
```bash
# 使用 air (hot reload)
air -c .air.toml

# 或直接執行
go run cmd/main.go
```

### 4. 測試 API
```bash
# Health check
curl http://localhost:8080/health

# 根路徑
curl http://localhost:8080/
```

## 📋 開發階段

### Phase 1: Mock API (進行中)
- [x] 基礎架構建立
- [x] Health check endpoint
- [ ] 15 個 GET APIs with mock data

### Phase 2: 資料庫整合 (待開發)
- [ ] GORM + PostgreSQL
- [ ] Models 定義
- [ ] Repository 層
- [ ] 真實資料查詢

### Phase 3: 完整 CRUD (待開發)
- [ ] POST APIs
- [ ] PATCH APIs
- [ ] DELETE APIs
- [ ] WebSocket 支援

## 🏗️ 架構說明

### Clean Architecture 分層

1. **Handlers** (Presentation Layer)
   - 處理 HTTP 請求和回應
   - 參數驗證
   - 呼叫 Service 層

2. **Services** (Business Logic Layer)
   - 核心業務邏輯
   - 資料轉換
   - 呼叫 Repository 層

3. **Repository** (Data Access Layer)
   - 資料庫操作
   - CRUD 操作封裝

4. **Models** (Domain Layer)
   - 資料結構定義
   - 業務實體

### Middleware

- **CORS**: 跨域請求處理
- **Logger**: 請求日誌記錄
- **Auth**: JWT 認證 (Phase 2)

### Response 格式

統一的 API 回應格式：

```json
{
  "success": true,
  "message": "optional message",
  "data": { ... }
}
```

錯誤回應：
```json
{
  "success": false,
  "error": "error message"
}
```

## 🔧 技術棧

- **框架**: Gin
- **ORM**: GORM (Phase 2)
- **資料庫**: PostgreSQL (Phase 2)
- **認證**: JWT (Phase 2)
- **環境變數**: godotenv
- **Hot Reload**: Air

## 📚 API 文件

詳細 API 文件請參考 [CLAUDE.md](../../CLAUDE.md)

## 🤝 貢獻

此專案是 rails-counter-server 的 Go 重寫版本。