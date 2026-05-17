# MoneyRabbit

個人用PWA家計簿アプリ。銀行口座のCSVをインポートして収支を管理する。

- **認証なし**（VPNによるアクセス制限）
- **AI自動分類**：取引摘要をキーワードルール → Claude API の順で自動カテゴリ分類
- **PWA対応**：iOS・Androidのホーム画面に追加してネイティブアプリ風に使用可能

## スクリーンショット

| ダッシュボード | 取引一覧 | CSVインポート | カテゴリ管理 |
|---|---|---|---|
| 月次サマリー・グラフ | フィルタ・カテゴリ編集 | preview → confirm の2ステップ | キーワードルールCRUD |

## 技術スタック

```
frontend (Vite + React 19) → /api/* → backend (Go + Gin) → MySQL 8.4
```

| レイヤー | 技術 |
|---|---|
| フロントエンド | React 19 / TanStack Router / TanStack Query / Zustand / Tailwind CSS v4 / shadcn/ui |
| バックエンド | Go / Gin / ent (ORM) / Atlas (スキーマ管理) / goose (マイグレーション) |
| AI分類 | Claude API (`claude-haiku-4-5`) |
| インフラ | Docker Compose / Nginx / MySQL 8.4 |

## セットアップ

Docker が必要。

```bash
cp .env.example .env   # パスワード・APIキーを設定
docker compose up      # 全サービス起動
```

| サービス | URL |
|---|---|
| フロントエンド | http://localhost |
| バックエンド API | http://localhost:8080 |
| Swagger UI | http://localhost:8080/docs/swagger/index.html |

## 開発

```bash
# フロントエンド
cd frontend
pnpm dev        # 開発サーバー (http://localhost:3000)
pnpm lint
pnpm typecheck

# バックエンド
cd backend
go run ./cmd/server   # http://localhost:8080
go test ./...
go vet ./...
```

フロントエンドの `/api/*` は `vite.config.ts` の proxy 設定でバックエンドに転送される。

## 環境変数

`.env.example` を参照。

| 変数名 | 内容 |
|---|---|
| `DATABASE_URL` | MySQL接続文字列 |
| `MYSQL_PASSWORD` | MySQLユーザーパスワード |
| `MYSQL_ROOT_PASSWORD` | MySQLルートパスワード |
| `ANTHROPIC_API_KEY` | Claude APIキー |

## ドキュメント

設計ドキュメントは [`docs/design/`](./docs/design/README.md) に格納されている。
