# 実装計画

## Step 1 — DBレイヤー

entスキーマ定義 → Atlas差分SQL生成 → gooseマイグレーション適用

対象エンティティ：`Category` / `CategoryRule` / `Transaction`

`ImportFormat` はDBテーブルとして持たず、Goコードの定数として定義する（`internal/service/csv/formats.go`）。

## Step 2 — バックエンドAPI（基本CRUD）

| エンドポイント群 | 内容 |
|---|---|
| `GET /api/import-formats` | 対応済みフォーマット一覧（Goコードから生成、読み取り専用） |
| `/api/categories` | カテゴリCRUD |
| `/api/category-rules` | 自動分類ルールCRUD |

## Step 3 — CSVインポート

バックエンド：CSVパース → 重複検出 → キーワードルール適用 → Claude API分類  
エンドポイント：`POST /api/import/preview`、`POST /api/import/confirm`

## Step 4 — 取引API

- `GET /api/transactions`（フィルタ・ページネーション）
- `PATCH /api/transactions/:id/category`（カテゴリ手動変更）

## Step 5 — フロントエンド基盤

- 共通レイアウト（サイドバー・ナビゲーション）
- APIクライアント（`src/lib/api.ts`）
- shadcn/ui 基本コンポーネント導入

## Step 6 — フロントエンド各画面

| 画面 | 依存Step |
|---|---|
| カテゴリ管理 | Step 2 |
| CSVインポート | Step 3 |
| 取引一覧 | Step 4 |

## Step 7 — ダッシュボード

集計API（月次サマリー・カテゴリ別） → Rechartsグラフ実装

## Step 8 — PWA対応

`vite-plugin-pwa` 設定・iOSホーム画面追加バナー
