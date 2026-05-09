# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 開発コマンド

```bash
go run ./cmd/server             # 開発サーバー (http://localhost:8080)
go build ./...                  # ビルド
go test ./...                   # 全テスト
go test ./internal/usecase/...  # 単一パッケージのテスト
go vet ./...                    # 静的解析
swag init -g cmd/server/main.go -o docs/swagger  # Swagger再生成
```

## アーキテクチャ

```
handler → usecase → domain ← service
                              (persistence / csv)
```

| レイヤー | パス | 責務 |
|---|---|---|
| handler | `internal/handler/` | HTTPリクエスト/レスポンス・Usecase interfaceを定義 |
| usecase | `internal/usecase/` | ビジネスロジック・Repository interfaceを定義 |
| domain | `internal/domain/entity/` | エンティティ定義のみ（外部依存なし） |
| service/persistence | `internal/service/persistence/` | Repository interfaceをentで実装 |
| service/csv | `internal/service/csv/` | CSVパース処理（Shift-JIS変換・ImportFormat定数に基づく列マッピング） |
| service/classifier | `internal/service/classifier/` | 取引カテゴリ分類ロジック |

インターフェースは使う側が定義する（Goの慣習）。
- `handler/transaction.go` にusecase interfaceを定義
- `usecase/transaction/usecase.go` にRepository interfaceを定義

## 設計上の重要な決定

- **ImportFormat はDBテーブルではない**: 対応銀行・カードのCSVフォーマット設定は `internal/service/csv/formats.go` にGoコードの定数として定義する（三井住友銀行: `smbc_bank`、SMBCカード: `smbc_card`）
- **Accountテーブルなし**: 口座も固定のため、`Transaction` が `import_format_id` を直接持つ

## テスト戦略

| 対象 | 手法 | モック対象 |
|---|---|---|
| `handler` | ユニットテスト | usecase interface |
| `usecase` | ユニットテスト | Repository interface |
| `service/persistence` | 統合テスト | 実DB（テスト用コンテナ） |

## スキーマ変更フロー

entスキーマ（`ent/schema/`）を変更した場合：

```bash
go generate ./ent/...                          # 1. entコード再生成
atlas migrate diff <name> \
  --dir "file://db/migrations" \
  --to "ent://ent/schema" \
  --dev-url "mysql://root:pass@localhost:3306/dev"  # 2. 差分SQL生成
goose -dir db/migrations mysql "${DATABASE_URL}" up  # 3. マイグレーション適用
```

## 環境変数

`.env.example` を参照。

| 変数名 | 内容 |
|---|---|
| `DATABASE_URL` | `moneyrabbit:password@tcp(db:3306)/moneyrabbit?parseTime=true` |
| `MYSQL_PASSWORD` | MySQLユーザーパスワード |
| `MYSQL_ROOT_PASSWORD` | MySQLルートパスワード |
| `ANTHROPIC_API_KEY` | Claude APIキー |
