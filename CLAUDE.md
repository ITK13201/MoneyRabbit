# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

個人用PWA家計簿アプリ。銀行口座のCSVをインポートして収支を管理する。認証なし（VPNでアクセス制限）。

## アーキテクチャ

```
frontend (Vite + React 19) → /api/* → backend (Go 1.26 + Gin) → MySQL 8.4（開発）/ MariaDB 11.4（Kubernetes）
```

設計ドキュメントは `docs/design/` に格納されている。

## セットアップ

Dockerが必要。

```bash
direnv allow           # .envrc を許可（初回のみ・nix dev shell + .env 自動ロード）
docker compose up      # 全サービス起動
```

dev shell（`flake.nix`）が go / nodejs / pnpm / goose / atlas を提供する。direnv が自動で `nix develop` を起動する。シークレットは `.env.tpl`（1Password）から自動注入される。

## CLAUDE.md の管理方針

**このファイルには最低限の情報のみ記載する。** 詳細は各ディレクトリの CLAUDE.md に分散して記載すること。

新しいディレクトリを追加した場合は、そのディレクトリ内に CLAUDE.md を作成し、以下のインデックスに追加する。

| ファイル | 内容 |
|---|---|
| [frontend/CLAUDE.md](./frontend/CLAUDE.md) | 開発コマンド・フロントエンド構造・状態管理方針 |
| [backend/CLAUDE.md](./backend/CLAUDE.md) | 開発コマンド・バックエンド構造・スキーマ変更フロー |
| [charts/moneyrabbit/CLAUDE.md](./charts/moneyrabbit/CLAUDE.md) | Helm Chart・Kubernetes デプロイ・リリースフロー |
