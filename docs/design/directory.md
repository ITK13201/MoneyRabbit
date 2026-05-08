# ディレクトリ設計

## 全体構成

```
MoneyRabbit/
├── frontend/
├── backend/
├── docker-compose.yml
├── .env
├── .env.example
└── docs/
```

---

## フロントエンド（Vite + React）

```
frontend/
├── src/
│   ├── routes/                  # TanStack Router ファイルベースルーティング
│   │   ├── __root.tsx           # ルートレイアウト
│   │   ├── index.tsx            # / ダッシュボード
│   │   ├── transactions/
│   │   │   └── index.tsx        # /transactions
│   │   ├── import/
│   │   │   └── index.tsx        # /import
│   │   ├── accounts/
│   │   │   └── index.tsx        # /accounts
│   │   ├── categories/
│   │   │   └── index.tsx        # /categories
│   │   └── settings/
│   │       └── index.tsx        # /settings
│   ├── components/
│   │   ├── ui/                  # shadcn/ui コンポーネント
│   │   └── layout/              # Header, Sidebar 等
│   ├── hooks/                   # TanStack Query カスタムフック
│   ├── stores/                  # Zustand ストア
│   ├── lib/
│   │   ├── api.ts               # APIクライアント (fetch wrapper)
│   │   └── utils.ts             # shadcn/ui ユーティリティ (cn関数等)
│   ├── types/                   # 共通型定義
│   ├── main.tsx
│   └── routeTree.gen.ts         # TanStack Router 自動生成（コミット対象外）
├── public/
│   └── icons/                   # PWAアイコン
├── index.html
├── vite.config.ts
├── tsconfig.json
├── tailwind.config.ts
├── components.json              # shadcn/ui 設定
├── nginx.conf
└── Dockerfile
```

---

## バックエンド（Go + Gin）

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # エントリーポイント
├── docs/                        # ドキュメント類
│   └── swagger/                 # swag generate で自動生成（コミット対象）
│       ├── docs.go
│       ├── swagger.json
│       └── swagger.yaml
├── internal/
│   ├── domain/                  # ドメイン層：エンティティのみ（外部依存なし）
│   │   └── entity/
│   │       ├── transaction.go
│   │       ├── account.go
│   │       ├── category.go
│   │       └── bank_format.go
│   ├── usecase/                 # ユースケース層：アプリケーションビジネスロジック
│   │   ├── transaction/
│   │   │   ├── usecase.go       # Repository interface + Classifier interface + Usecase struct を定義
│   │   │   ├── list.go
│   │   │   ├── import_csv.go
│   │   │   └── update_category.go
│   │   ├── account/
│   │   │   └── usecase.go       # Repository interface + Usecase struct を定義
│   │   └── category/
│   │       └── usecase.go       # Repository interface + Usecase struct を定義
│   ├── service/                 # サービス層：usecaseのRepository interfaceを実装
│   │   ├── persistence/         # entを使ったDBアクセス実装
│   │   │   ├── transaction.go
│   │   │   ├── account.go
│   │   │   └── category.go
│   │   ├── csv/                 # CSVパース・銀行フォーマット変換
│   │   │   └── parser.go
│   │   └── classifier/          # Claude APIによるカテゴリ自動分類
│   │       └── claude.go        # usecase/transaction の Classifier interface を実装
│   ├── handler/                 # ハンドラー層：HTTPリクエスト/レスポンス（Gin）
│   │   ├── transaction.go       # Usecase interface・request/response型をこのファイルに定義
│   │   ├── account.go
│   │   ├── category.go
│   │   ├── import.go
│   │   └── router.go            # ルーティング定義・/swagger/*でSwagger UI配信
│   └── middleware/              # Ginミドルウェア
│       └── cors.go
├── ent/                         # ent 自動生成コード（go generate）
│   └── schema/                  # スキーマ定義（手書き）
│       ├── transaction.go
│       ├── account.go
│       ├── category.go
│       ├── categoryrule.go
│       └── bankformat.go
├── db/
│   └── migrations/              # goose マイグレーションファイル
│       └── 20260507000001_init.sql
├── atlas.hcl                    # Atlas 設定
├── go.mod
├── go.sum
└── Dockerfile

# Swagger UI: http://localhost:8080/swagger/index.html
# 再生成: swag init -g cmd/server/main.go -o docs/swagger
```

### レイヤー責務と依存方向

```
handler → usecase → domain ← service
                              (persistence / csv)
```

| レイヤー | 責務 | 外部依存 |
|---|---|---|
| `domain` | エンティティ定義のみ | なし |
| `usecase` | ビジネスロジック・Repository interfaceを定義 | domain のみ |
| `service` | usecase の Repository interface を実装（ent） | ent・domain |
| `handler` | HTTPリクエスト/レスポンス・Usecase interfaceを定義 | usecase |

### テスト戦略

| 対象 | テスト手法 | モック対象 |
|---|---|---|
| `handler` | ユニットテスト | usecase interface をモック |
| `usecase` | ユニットテスト | Repository interface をモック |
| `service/persistence` | 統合テスト | 実DB（テスト用コンテナ） |

インターフェースは**使う側**（handler・usecase）が定義する（Goの慣習）。

---

## Docker Compose

```
MoneyRabbit/
├── docker-compose.yml
├── .env                         # Git管理外
└── .env.example                 # Git管理対象
```
