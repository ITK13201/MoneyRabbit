# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 開発コマンド

```bash
pnpm dev        # 開発サーバー (http://localhost:3000)
pnpm build      # プロダクションビルド
pnpm lint       # ESLint
pnpm typecheck  # 型チェック
```

## 構造

- `src/routes/` — TanStack Router ファイルベースルーティング（`routeTree.gen.ts` は自動生成・コミット対象外）
  - `index.tsx` — ダッシュボード（月次収支サマリー・カテゴリ別棒グラフ・最近の取引）
  - `transactions/index.tsx` — 取引一覧（ページネーション・カテゴリ変更）
  - `import/index.tsx` — CSVインポート（preview → confirm の2ステップUI）
  - `categories/index.tsx` — カテゴリ・キーワードルールのCRUD
- `src/hooks/` — TanStack Query カスタムフック（`useTransactions`, `useCategories`）
- `src/lib/api.ts` — 型付きfetchラッパー（`api.categories.*`, `api.import.*`, `api.transactions.*`）
- `src/types/` — 共通型定義（APIレスポンス型）

新しいルートファイルを追加した後は `pnpm build` を実行して `routeTree.gen.ts` を再生成すること（`pnpm dev` では自動再生成されない場合がある）。

## 状態管理方針

- **サーバー状態**（APIデータ）→ TanStack Query（`src/hooks/` にカスタムフックを置く）
- **クライアント状態**（UIのみ）→ Zustand（`src/stores/` にストアを置く）
- React Context・Redux は使用しない

## TypeScript パスエイリアス

`@/*` → `./src/*`（`tsconfig.app.json` および `vite.config.ts` で設定済み）

## shadcn/ui コンポーネント追加

```bash
pnpm dlx shadcn@latest add <component-name>
```

生成されたコンポーネントは `src/components/ui/` に配置される。

## 開発時のAPI通信

`vite.config.ts` の proxy 設定により `/api/*` は `http://localhost:8080` に転送される。

