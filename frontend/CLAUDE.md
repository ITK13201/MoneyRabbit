# CLAUDE.md

## 開発コマンド

```bash
pnpm dev        # 開発サーバー (http://localhost:3000)
pnpm build      # プロダクションビルド
pnpm lint       # ESLint
pnpm typecheck  # 型チェック
```

## 構造

- `src/routes/` — TanStack Router ファイルベースルーティング（`routeTree.gen.ts` は自動生成・コミット対象外）
- `src/hooks/` — TanStack Query カスタムフック（サーバー状態）
- `src/stores/` — Zustand ストア（クライアント状態）
- `src/components/ui/` — shadcn/ui コンポーネント
- `src/lib/api.ts` — fetch wrapper（APIクライアント）

## 開発時のAPI通信

`vite.config.ts` の proxy 設定により `/api/*` は `http://localhost:8080` に転送される。
