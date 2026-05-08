# データモデル

MySQL 8 + ent（Go ORM）で管理する。  
entのスキーマ定義からGoコードとマイグレーションを自動生成する。

## entスキーマ定義

### BankFormat（銀行CSVフォーマット定義）

```go
// backend/ent/schema/bankformat.go
type BankFormat struct{ ent.Schema }

func (BankFormat) Fields() []ent.Field {
    return []ent.Field{
        field.String("id"),
        field.String("name"),                               // 例: 三菱UFJ銀行
        field.Enum("encoding").Values("UTF-8", "Shift_JIS"),
        field.Int("skip_rows"),                             // ヘッダー行数
        field.Int("col_date"),                              // 列インデックス (0始まり)
        field.Int("col_description"),
        field.Int("col_withdrawal"),                        // 支出列
        field.Int("col_deposit"),                           // 収入列
        field.Int("col_balance").Optional().Nillable(),
        field.String("date_format"),                        // 例: 2006/01/02 (Go形式)
    }
}

func (BankFormat) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("accounts", Account.Type),
    }
}
```

### Account（口座）

```go
type Account struct{ ent.Schema }

func (Account) Fields() []ent.Field {
    return []ent.Field{
        field.String("id"),
        field.String("name"),                               // 例: 三菱UFJ 普通口座
        field.Time("created_at").Default(time.Now),
    }
}

func (Account) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("bank_format", BankFormat.Type).Ref("accounts").Unique(),
        edge.To("transactions", Transaction.Type),
    }
}
```

### Category（カテゴリ）

```go
type Category struct{ ent.Schema }

func (Category) Fields() []ent.Field {
    return []ent.Field{
        field.String("id"),
        field.String("name"),
        field.String("color"),                              // 例: #FF5733
        field.String("icon"),
        field.Enum("type").Values("income", "expense", "both"),
        field.Int("sort_order").Default(0),
    }
}

func (Category) Edges() []ent.Edge {
    return []ent.Edge{
        edge.To("rules", CategoryRule.Type),
        edge.To("transactions", Transaction.Type),
    }
}
```

### CategoryRule（自動カテゴリルール）

```go
type CategoryRule struct{ ent.Schema }

func (CategoryRule) Fields() []ent.Field {
    return []ent.Field{
        field.String("id"),
        field.String("keyword"),                            // 摘要への部分一致キーワード
        field.Int("priority").Default(0),                  // 大きいほど優先
    }
}

func (CategoryRule) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("category", Category.Type).Ref("rules").Unique().Required(),
    }
}
```

### Transaction（取引）

```go
type Transaction struct{ ent.Schema }

func (Transaction) Fields() []ent.Field {
    return []ent.Field{
        field.String("id"),
        field.Time("date"),
        field.String("description"),                        // 摘要
        field.Int("amount"),                               // 正: 収入, 負: 支出 (円)
        field.Int("balance").Optional().Nillable(),        // 残高
        field.Time("imported_at").Default(time.Now),
        field.String("source_file"),
    }
}

func (Transaction) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("account", Account.Type).Ref("transactions").Unique().Required(),
        edge.From("category", Category.Type).Ref("transactions").Unique(),
    }
}
```

## ER図

```
BankFormat
    │
    └── Account
            │
            └── Transaction ── Category ── CategoryRule
```

## マイグレーションフロー

entスキーマを変更するたびに以下の手順を実行する。

```bash
# 1. entスキーマからGoコードを再生成
go generate ./ent/...

# 2. Atlasでスキーマ差分のSQLを生成
atlas migrate diff <migration_name> \
  --dir "file://backend/db/migrations" \
  --to "ent://backend/ent/schema" \
  --dev-url "mysql://root:pass@localhost:3306/dev"

# 3. 生成されたSQLをgooseで適用
goose -dir backend/db/migrations mysql "${DATABASE_URL}" up
```

### ディレクトリ構成

```
backend/
├── ent/
│   ├── schema/          # entスキーマ定義（手書き）
│   └── (generated)      # go generateで自動生成
├── db/
│   └── migrations/      # Atlasが生成・gooseが適用するSQLファイル
│       ├── 20260507000001_init.sql
│       └── atlas.sum    # Atlasのチェックサムファイル
└── sqlc.yaml            # 不要（entを使用）
```
