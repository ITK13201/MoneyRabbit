# データモデル

MySQL 8 + ent（Go ORM）で管理する。  
entのスキーマ定義からGoコードとマイグレーションを自動生成する。

## ImportFormat（インポートフォーマット）

DBテーブルとして管理しない。対応している銀行・カードのフォーマット設定はGoコード内の定数として定義する。詳細は [csv-import.md](./csv-import.md) を参照。

```go
// backend/internal/service/csv/formats.go
type ImportFormatID string

const (
    ImportFormatSMBCBank ImportFormatID = "smbc_bank" // 三井住友銀行
    ImportFormatSMBCCard ImportFormatID = "smbc_card" // SMBCカード
)
```

## entスキーマ定義

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
        field.Enum("import_format_id").Values(             // インポート元
            "smbc_bank",
            "smbc_card",
        ),
        field.Time("imported_at").Default(time.Now),
    }
}

func (Transaction) Edges() []ent.Edge {
    return []ent.Edge{
        edge.From("category", Category.Type).Ref("transactions").Unique(),
    }
}
```

## ER図

```
Transaction ── Category ── CategoryRule
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
```
