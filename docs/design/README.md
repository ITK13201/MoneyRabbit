# MoneyRabbit 設計ドキュメント

PWAで動作する個人用家計簿アプリ。銀行口座のCSVデータをインポートして収支を管理する。

## 対象プラットフォーム

- ブラウザ（PC / Android）
- iOS（ホーム画面追加によるPWA）

## ドキュメント一覧

| ファイル | 内容 |
|---|---|
| [tech-stack.md](./tech-stack.md) | 技術選定・採用理由 |
| [deployment.md](./deployment.md) | デプロイ構成（Docker Compose） |
| [observability.md](./observability.md) | ログ設計（slog / Promtail / Loki / Grafana） |
| [implementation-plan.md](./implementation-plan.md) | 実装計画（Step 1〜8） |
| [features.md](./features.md) | 機能設計・画面一覧 |
| [data-model.md](./data-model.md) | データモデル・テーブル定義 |
| [csv-import.md](./csv-import.md) | CSVインポート設計（ImportFormat・フロー・APIエンドポイント） |
| [manual-entry.md](./manual-entry.md) | 取引手動入力設計（スキーマ変更・API追加・フロントエンド実装） |

## 基本方針

- **個人利用**：認証なし（VPNによるアクセス制限）
- **バックエンドあり**：Next.js + PostgreSQL
- **CSVインポートのみ**：手動入力は行わない
- **複数口座対応**：銀行・口座を紐づけて管理
