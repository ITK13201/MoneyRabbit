# 技術選定

## フロントエンド

| カテゴリ | 技術 | 理由 |
|---|---|---|
| ビルドツール | Vite | 高速・軽量・PWA対応（vite-plugin-pwa） |
| フレームワーク | React 19 + TypeScript | 2026年もエコシステム最大・採用率1位 |
| ルーター | TanStack Router | 完全型安全・ファイルベースルーティング |
| スタイリング | Tailwind CSS v4 + shadcn/ui | shadcnは2026年で75k+ stars・デファクトスタンダード |
| クライアント状態 | Zustand | 軽量3KB・シンプルなストア管理 |
| サーバー状態 | TanStack Query v6 | 2026年のサーバー状態管理の標準・キャッシュ・再フェッチ自動管理 |
| グラフ | Recharts | React親和性が高い・週360万DL・小〜中規模データセット向け |
| CSVパース | PapaParse | ブラウザ向けCSVライブラリのデファクト |
| PWA | vite-plugin-pwa | Viteエコシステムの標準PWAソリューション・Workboxベース |

## バックエンド

| カテゴリ | 技術 | 理由 |
|---|---|---|
| 言語 | Go | 軽量・高速・シングルバイナリ |
| フレームワーク | Gin | 2026年Goフレームワーク採用率48%・最大コミュニティ・安定 |
| API仕様 | swaggo/swag | GoアノテーションからSwagger UIを自動生成 |
| カテゴリ分類 | Claude API (claude-haiku-4-5) | 摘要テキストからカテゴリを自動推論 |
| ORM | ent | コードファースト・型安全・グラフベースのリレーション管理 |
| スキーマ管理 | Atlas | entスキーマ差分からSQLを自動生成・トランザクション対応 |
| マイグレーション | goose | バージョン管理・適用履歴の記録・Go製マイグレーションも記述可 |

## オブザーバビリティ

| カテゴリ | 技術 | 理由 |
|---|---|---|
| ロガー | log/slog (Go stdlib) | Go 1.21標準・外部依存なし・構造化JSON出力 |
| ログ収集 | Promtail | Dockerコンテナログを自動収集 |
| ログ集約 | Loki | Grafana Labs製・軽量・ラベルベース検索 |
| 可視化 | Grafana | Lokiのデファクトフロントエンド |

## データベース

| カテゴリ | 技術 |
|---|---|
| RDBMS | MySQL 8.4 LTS（2032年までサポート） |

## デプロイ

Docker Composeで動作する。詳細は [deployment.md](./deployment.md) を参照。

```
Vite (frontend) ──HTTP──▶ Gin (backend) ──▶ PostgreSQL
```

## 方針

- **個人利用**：認証なし（VPNによるネットワークレベルのアクセス制限）
- **フロント/バック分離**：フロントは静的ファイル配信、バックはJSON API
- **Type-safe**：entのスキーマ定義からGoコードを生成し、型安全なDB操作
- **マイグレーション**：ent → Atlas（SQL自動生成）→ goose（適用・履歴管理）
- **コンテナ**：GoはCGO無効のシングルバイナリ、distrolessイメージで軽量化

## iOS PWAの注意点

- ホーム画面追加を促すバナーを独自実装（iOSはWeb App Install Promptに非対応）
- `display: standalone` でネイティブアプリ風UI
