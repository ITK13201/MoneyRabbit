# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 開発コマンド

```bash
go run ./cmd/server                                    # 開発サーバー (http://localhost:8080)
go build ./...                                         # ビルド
go test ./...                                          # 全テスト
go test ./internal/usecase/category/... -run TestList  # 単一テスト実行
go vet ./...                                           # 静的解析
go mod tidy                                            # 依存関係整理（パッケージ追加後に実行）
swag init -g cmd/server/main.go -o docs/swagger && mv docs/swagger/docs.go docs/swagger/swagger.go  # Swagger docs 再生成（ハンドラ変更時・swag は nixpkgs 未収録なので go install github.com/swaggo/swag/cmd/swag@latest で導入）
```

Swagger UI は起動後 http://localhost:8080/docs/swagger/index.html で確認できる。`docs/` はコード生成物のためコミット対象だが、手動編集しないこと。

## APIエンドポイント

| メソッド | パス | ハンドラ |
|---|---|---|
| GET | `/api/health` | ヘルスチェック |
| GET | `/api/import-formats` | フォーマット一覧（DB不使用・定数返却） |
| POST | `/api/import/preview` | CSVパース（DB書き込みなし） |
| POST | `/api/import/confirm` | 分類＋保存 |
| GET | `/api/categories` | カテゴリ一覧 |
| POST | `/api/categories` | カテゴリ作成 |
| GET | `/api/categories/:id` | カテゴリ取得 |
| PUT | `/api/categories/:id` | カテゴリ更新 |
| DELETE | `/api/categories/:id` | カテゴリ削除 |
| POST | `/api/category-rules` | キーワードルール作成 |
| PUT | `/api/category-rules/:id` | キーワードルール更新 |
| DELETE | `/api/category-rules/:id` | キーワードルール削除 |
| GET | `/api/transactions` | 取引一覧 |
| POST | `/api/transactions` | 手動取引作成 |
| PUT | `/api/transactions/:id` | 取引更新 |
| PATCH | `/api/transactions/:id/category` | カテゴリのみ更新 |
| DELETE | `/api/transactions/:id` | 取引削除 |
| GET | `/api/summary/monthly` | 月次サマリー |

## アーキテクチャ

```
handler → usecase → domain ← service
                              (persistence / csv / classifier)
```

| レイヤー | パス | 責務 |
|---|---|---|
| handler | `internal/handler/` | HTTPリクエスト/レスポンス・Usecase interfaceを定義 |
| usecase | `internal/usecase/` | ビジネスロジック・Repository interfaceを定義 |
| domain | `internal/domain/entity/` | エンティティ定義のみ（外部依存なし） |
| service/persistence | `internal/service/persistence/` | Repository interfaceをentで実装 |
| service/csv | `internal/service/csv/` | CSVパース処理（Shift-JIS変換・ImportFormat定数に基づく列マッピング） |
| service/classifier | `internal/service/classifier/` | 取引カテゴリ分類ロジック（キーワードルール照合 → Claude API呼び出し） |

インターフェースは使う側が定義する（Goの慣習）。
- `handler/transaction.go` にusecase interfaceを定義
- `usecase/transaction/usecase.go` にRepository interfaceを定義

## 設計上の重要な決定

- **ImportFormat はDBテーブルではない**: 対応銀行・カードのCSVフォーマット設定は `internal/service/csv/formats.go` にGoコードの定数として定義する（三井住友銀行: `smbc_bank`、SMBCカード: `smbc_card`）
- **Accountテーブルなし**: 口座も固定のため、`Transaction` が `import_format_id` を直接持つ
- **手動入力は `import_format_id = nil` で表現**: CSVインポート時は非null、`POST /api/transactions` での手動作成時はnull。`entity.Transaction.ImportFormatID` は `*string`
- **CSVインポートは2ステップ**: `POST /api/import/preview`（パースのみ、DB書き込みなし）→ユーザー確認 → `POST /api/import/confirm`（分類＋保存）
- **カテゴリ分類の優先順位**: キーワードルール（DB）が完全一致で優先 → Claude API（claude-sonnet-4-6）でAI分類 → いずれも失敗時は `CategoryID = nil`（未分類）

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
  --dev-url "mysql://root:pass@localhost:3306/dev"  # 2. 差分SQL生成（atlas は nix dev shell で提供）
```

マイグレーションは **サーバー起動時に自動適用**される（`db/migrate.go` が goose で `db/migrations/*.sql` を embed して実行）。手動適用は不要。

`dev` データベースはatlas用の作業領域（空のDBが必要）。`docker compose up -d db` で起動後に `mysql -uroot -p... -e "CREATE DATABASE IF NOT EXISTS dev;"` で作成する。

> **注意**: マイグレーションファイルは goose 形式（`-- +goose Up` / `-- +goose Down`）で記述されている。Atlas が差分生成のためにこれらを dev DB へ適用する際、`-- +goose Down` セクションの DROP 文も実行してしまうため Atlas `migrate diff` が失敗するケースがある。その場合はマイグレーション SQL を手書きし、`atlas migrate hash --dir "file://db/migrations"` でチェックサムを再生成すること。

コンテナ内でもgoose・atlasを直接実行できる（Dockerfileに同梱）。シェルは`/busybox/sh`：

```bash
# マイグレーション状態確認
docker compose exec backend /busybox/sh -c "goose -dir /db/migrations mysql \"\$DATABASE_URL\" status"
# 全ロールバック（リセット）
docker compose exec backend /busybox/sh -c "goose -dir /db/migrations mysql \"\$DATABASE_URL\" reset"
# スキーマ差分生成
docker compose exec backend /busybox/sh -c "atlas migrate diff <name> --dir file:///db/migrations --to ent://ent/schema --dev-url mysql://root:pass@db:3306/dev"
```

## 依存関係の自動更新

Renovate が依存関係の更新 PR を自動作成する。`Dockerfile` の `ARG GOOSE_VERSION` は `renovate.json` のカスタム正規表現で管理されているため、手動で変更しても Renovate に上書きされる。

## 環境変数

| 変数名 | 内容 |
|---|---|
| `DATABASE_URL` | `moneyrabbit:password@tcp(db:3306)/moneyrabbit?parseTime=true` |
| `MYSQL_PASSWORD` | MySQLユーザーパスワード |
| `MYSQL_ROOT_PASSWORD` | MySQLルートパスワード |
| `ANTHROPIC_API_KEY` | Claude APIキー |
| `PORT` | バックエンドポート（デフォルト: `8080`） |
