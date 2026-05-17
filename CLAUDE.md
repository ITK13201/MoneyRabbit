# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

個人用PWA家計簿アプリ。銀行口座のCSVをインポートして収支を管理する。認証なし（VPNでアクセス制限）。

## 構成

```
frontend (Vite + React 19) → /api/* → backend (Go 1.26 + Gin) → MySQL 8.4
```

- **ORM**: ent（スキーマ → Goコード自動生成）、マイグレーションは Atlas で差分SQL生成 → goose でサーバー起動時に自動適用
- **カテゴリ分類**: キーワードルール（DB）→ Claude API（claude-sonnet-4-6）→ 未分類 の優先順位でAI分類
- **フロント状態管理**: TanStack Query（サーバー状態）+ Zustand（UIのみ）。React Context / Redux は使用しない

設計ドキュメントは [`docs/design/`](./docs/design/README.md) に格納されている。

| ファイル | 内容 |
|---|---|
| `data-model.md` | テーブル定義・エンティティ関係 |
| `csv-import.md` | ImportFormat 設計・インポートフロー・APIエンドポイント |
| `manual-entry.md` | 手動入力設計（`import_format_id = nil` で表現） |
| `features.md` | 機能設計・画面一覧 |
| `tech-stack.md` | 技術選定・採用理由 |
| `deployment.md` | デプロイ構成（Docker Compose） |

## セットアップ・起動

Dockerが必要（DB・バックエンドはコンテナで動作）。

```bash
cp .env.example .env   # 初回のみ：パスワード等を設定
docker compose up      # 全サービス起動（フロント・バックエンド・MySQL）

# 個別起動（開発時）
cd frontend && pnpm dev             # http://localhost:3000
cd backend  && go run ./cmd/server  # http://localhost:8080
```

Swagger UI: http://localhost:8080/docs/swagger/index.html

## CLAUDE.md の管理方針

**このファイルには最低限の情報のみ記載する。** 詳細は各ディレクトリのCLAUDE.mdに分散して記載すること。

新しいディレクトリを追加した場合は、そのディレクトリ内にCLAUDE.mdを作成し、以下のインデックスに追加する。

| ファイル | 内容 |
|---|---|
| [frontend/CLAUDE.md](./frontend/CLAUDE.md) | 開発コマンド・フロントエンド構造 |
| [backend/CLAUDE.md](./backend/CLAUDE.md) | 開発コマンド・バックエンド構造・スキーマ変更フロー |
