# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

個人用PWA家計簿アプリ。銀行口座のCSVをインポートして収支を管理する。認証なし（VPNでアクセス制限）。

## 構成

```
frontend (Vite + React) → /api/* → backend (Go + Gin) → MySQL 8.4
```

## 起動

```bash
# 全サービス起動
docker compose up

# 個別起動（開発時）
cd frontend && pnpm dev          # http://localhost:3000
cd backend  && go run ./cmd/server  # http://localhost:8080
```

## CLAUDE.md の管理方針

**このファイルには最低限の情報のみ記載する。** 詳細は各ディレクトリの CLAUDE.md に分散して記載すること。

新しいディレクトリを追加した場合は、そのディレクトリ内に CLAUDE.md を作成し、以下のインデックスに追加する。

| ファイル | 内容 |
|---|---|
| [frontend/CLAUDE.md](./frontend/CLAUDE.md) | 開発コマンド・フロントエンド構造 |
| [backend/CLAUDE.md](./backend/CLAUDE.md) | 開発コマンド・バックエンド構造・スキーマ変更フロー |
