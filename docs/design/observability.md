# オブザーバビリティ設計

## 構成概要

```
backend (Go) ──stdout JSON──▶ Promtail ──▶ Loki ──▶ Grafana
```

バックエンドは構造化JSONログを標準出力に書き出すのみ。  
Promtail・Loki・Grafanaは本番（Kubernetes）側で管理する。

## 環境別構成

| 環境 | ログ確認方法 |
|---|---|
| 開発（Docker Compose） | `docker compose logs -f backend` |
| 本番（Kubernetes） | Promtail → Loki → Grafana |

## Goロガー設計

Go 1.21以降標準の `log/slog` を使用する。外部依存なしで構造化JSONログを出力できる。

### 初期化（cmd/server/main.go）

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))
slog.SetDefault(logger)
```

### Ginロギングミドルウェア（internal/middleware/logger.go）

リクエストごとに以下のフィールドをJSON出力する。

```go
slog.InfoContext(c.Request.Context(), "request",
    slog.String("method",     c.Request.Method),
    slog.String("path",       c.Request.URL.Path),
    slog.Int("status",        c.Writer.Status()),
    slog.Duration("latency",  latency),
    slog.String("request_id", c.GetString("request_id")),
)
```

### ログフィールド規則

| フィールド | 型 | 内容 |
|---|---|---|
| `time` | string | RFC3339形式 |
| `level` | string | `INFO` / `WARN` / `ERROR` |
| `msg` | string | ログメッセージ |
| `request_id` | string | リクエスト追跡ID（UUIDv4） |
| `method` | string | HTTPメソッド |
| `path` | string | リクエストパス |
| `status` | int | HTTPステータスコード |
| `latency` | string | レスポンスタイム |

### リクエストID（internal/middleware/requestid.go）

リクエスト間の追跡のため、全リクエストにUUIDv4のIDを付与する。

```go
id := c.GetHeader("X-Request-ID")
if id == "" {
    id = uuid.NewString()
}
c.Set("request_id", id)
c.Header("X-Request-ID", id)
```

## ディレクトリ

```
backend/
└── internal/
    └── middleware/
        ├── logger.go        # Ginロギングミドルウェア
        └── requestid.go     # リクエストIDミドルウェア
```
