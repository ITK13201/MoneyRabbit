# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

個人用PWA家計簿アプリ。銀行口座のCSVをインポートして収支を管理する。認証なし（VPNでアクセス制限）。

## 開発コマンド

### 起動

```bash
# 全サービス起動
docker compose up

# 個別起動（開発時）
cd frontend && npm run dev       # http://localhost:5173
cd backend  && go run ./cmd/server  # http://localhost:8080
```

### フロントエンド

```bash
cd frontend
npm run dev        # 開発サーバー
npm run build      # プロダクションビルド
npm run lint       # ESLint
npm run typecheck  # 型チェック
```

### バックエンド

```bash
cd backend
go build ./...          # ビルド
go test ./...           # 全テスト
go test ./internal/usecase/transaction/...  # 単一パッケージのテスト
go vet ./...            # 静的解析
swag init -g cmd/server/main.go -o docs/swagger  # Swagger再生成
```

### スキーマ変更フロー（entスキーマを変更した場合）

```bash
cd backend
go generate ./ent/...                          # 1. entコード再生成
atlas migrate diff <name> \
  --dir "file://db/migrations" \
  --to "ent://ent/schema" \
  --dev-url "mysql://root:pass@localhost:3306/dev"  # 2. 差分SQL生成
goose -dir db/migrations mysql "${DATABASE_URL}" up  # 3. マイグレーション適用
```

### Swagger UI

`http://localhost:8080/swagger/index.html`

## アーキテクチャ

### 全体構成

```
frontend (Vite + React) → /api/* → backend (Go + Gin) → MySQL 8.4
```

NginxがSPAの静的配信と `/api/*` のリバースプロキシを兼ねる。

### バックエンドレイヤー構造

```
handler → usecase → domain ← service
                              (persistence / csv)
```

| レイヤー | パス | 責務 |
|---|---|---|
| handler | `internal/handler/` | HTTPリクエスト/レスポンス・Usecase interfaceを定義 |
| usecase | `internal/usecase/` | ビジネスロジック・Repository interfaceを定義 |
| domain | `internal/domain/entity/` | エンティティ定義のみ（外部依存なし） |
| service | `internal/service/persistence/` | Repository interfaceをentで実装 |

**インターフェースは使う側が定義する**（Goの慣習）。
- `handler/transaction.go` に usecase interface を定義
- `usecase/transaction/usecase.go` に Repository interface を定義

### テスト戦略

| 対象 | 手法 | モック対象 |
|---|---|---|
| `handler` | ユニットテスト | usecase interface |
| `usecase` | ユニットテスト | Repository interface |
| `service/persistence` | 統合テスト | 実DB（テスト用コンテナ） |

### フロントエンド構造

- `src/routes/` — TanStack Router ファイルベースルーティング（`routeTree.gen.ts` は自動生成）
- `src/hooks/` — TanStack Query カスタムフック（サーバー状態）
- `src/stores/` — Zustand ストア（クライアント状態）
- `src/lib/api.ts` — fetch wrapper（APIクライアント）

### マイグレーション管理

ent（スキーマ定義）→ Atlas（差分SQL生成）→ goose（適用・履歴管理）の順で管理する。gooseはバックエンド起動時に自動実行される。

## 環境変数

`.env.example` を参照。`.env` はGit管理外。

| 変数名 | 内容 |
|---|---|
| `DATABASE_URL` | `moneyrabbit:password@tcp(db:3306)/moneyrabbit?parseTime=true` |
| `MYSQL_PASSWORD` | MySQLユーザーパスワード |
| `MYSQL_ROOT_PASSWORD` | MySQLルートパスワード |
