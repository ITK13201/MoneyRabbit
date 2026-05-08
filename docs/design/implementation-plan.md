# 実装計画

## Step 1 — DBレイヤー

entスキーマ定義 → Atlas差分SQL生成 → gooseマイグレーション適用

対象エンティティ：`BankFormat` / `Account` / `Category` / `CategoryRule` / `Transaction`

## Step 2 — バックエンドAPI（基本CRUD）

| エンドポイント群 | 内容 |
|---|---|
| `/api/bank-formats` | 銀行フォーマットCRUD |
| `/api/accounts` | 口座CRUD |
| `/api/categories` | カテゴリCRUD |
| `/api/category-rules` | 自動分類ルールCRUD |

## Step 3 — CSVインポート

バックエンド：CSVパース → 重複検出 → キーワードルール適用 → Claude API分類  
エンドポイント：`POST /api/import`

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
| 銀行フォーマット設定 | Step 2 |
| 口座管理 | Step 2 |
| カテゴリ管理 | Step 2 |
| CSVインポート | Step 3 |
| 取引一覧 | Step 4 |

## Step 7 — ダッシュボード

集計API（月次サマリー・カテゴリ別） → Rechartsグラフ実装

## Step 8 — PWA対応

`vite-plugin-pwa` 設定・iOSホーム画面追加バナー
