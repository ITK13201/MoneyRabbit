# 取引手動入力 設計

CSVインポートに加え、フォームから取引を直接登録・編集できる機能（Phase 2）の設計。

## 概要

- 取引一覧ページのダイアログから新規作成・編集を行う
- CSVインポートと同一の `Transaction` モデルを使用し、APIを拡張するだけで実現する

---

## データモデル変更

### Transaction スキーマ

`import_format_id` を nullable に変更する。`null` = 手動入力、非null = CSVインポート。

```go
// ent/schema/transaction.go
func (Transaction) Fields() []ent.Field {
    return []ent.Field{
        field.UUID("id", uuid.UUID{}).Default(uuid.New).Immutable(),
        field.Time("date"),
        field.String("description"),
        field.Int("amount"),                                // 正: 収入, 負: 支出 (円)
        field.Enum("import_format_id").
            Values("smbc_bank", "smbc_card").
            Optional().
            Nillable(),                                     // 手動入力時は nil
        field.Time("imported_at").Default(time.Now).Immutable(),
    }
}
```

### entity.Transaction

```go
// internal/domain/entity/transaction.go
type Transaction struct {
    ID             uuid.UUID  `json:"id"`
    Date           time.Time  `json:"date"`
    Description    string     `json:"description"`
    Amount         int        `json:"amount"`
    ImportFormatID *string    `json:"import_format_id"` // *string に変更（null = 手動入力）
    ImportedAt     time.Time  `json:"imported_at"`
    CategoryID     *uuid.UUID `json:"category_id"`
    Category       *Category  `json:"category,omitempty"`
}
```

### マイグレーション

`import_format_id` カラムを `NOT NULL` → nullable に変更するだけ（データ消失なし、破壊的変更なし）。

---

## APIエンドポイント

### 追加エンドポイント

| Method | Path | 説明 |
|---|---|---|
| `POST /api/transactions` | 取引を手動で1件作成 |
| `PUT /api/transactions/:id` | 取引の日付・摘要・金額・カテゴリを更新 |

既存の `DELETE /api/transactions/:id` はそのまま使用する。

### POST `/api/transactions`

**リクエスト:**

```json
{
  "date": "2026-05-01",
  "description": "スーパーABC",
  "amount": -3200,
  "category_id": "uuid-or-null"
}
```

- `amount`: 正 = 収入、負 = 支出（円）
- `category_id`: 省略または null = 未分類

**レスポンス:** `201 Created` + `entity.Transaction`

### PUT `/api/transactions/:id`

手動入力・CSVインポート問わず全取引に適用可能（CSVの摘要や金額を修正したいケースもあるため）。

**リクエスト:**

```json
{
  "date": "2026-05-01",
  "description": "スーパーABC",
  "amount": -3200,
  "category_id": "uuid-or-null"
}
```

**レスポンス:** `200 OK` + `entity.Transaction`

---

## バックエンド実装

### レイヤー構成（追加分）

```
handler/transaction.go     createUsecase / updateUsecase インターフェース追加
usecase/transaction/
  create.go                CreateManual(ctx, CreateManualInput) (*entity.Transaction, error)
  update.go                Update(ctx, uuid.UUID, UpdateInput) (*entity.Transaction, error)
  usecase.go               CreateManualTransaction / UpdateTransaction をRepository interfaceに追加
service/persistence/
  transaction.go           CreateManualTransaction / UpdateTransaction を実装
```

### usecase の型

```go
// usecase/transaction/usecase.go に追加

type CreateManualInput struct {
    Date        time.Time
    Description string
    Amount      int
    CategoryID  *uuid.UUID
}

type UpdateInput struct {
    Date        time.Time
    Description string
    Amount      int
    CategoryID  *uuid.UUID
}
```

### Repository interface 追加メソッド

```go
CreateManualTransaction(ctx context.Context, input CreateManualInput) (*entity.Transaction, error)
UpdateTransaction(ctx context.Context, id uuid.UUID, input UpdateInput) (*entity.Transaction, error)
```

---

## フロントエンド実装

### UX方針

新規ページ（`/transactions/new`）ではなく、**取引一覧ページ上のダイアログ**で実装する。
PWA（モバイル）でページ遷移なしに操作できるため UX が良い。

### 変更ファイル

| ファイル | 変更内容 |
|---|---|
| `src/lib/api.ts` | `api.transactions.create` / `api.transactions.update` を追加 |
| `src/hooks/useTransactions.ts` | `useCreateTransaction` / `useUpdateTransaction` mutation を追加 |
| `src/routes/transactions/index.tsx` | ダイアログ統合・「＋ 手動追加」ボタン・編集ボタン追加 |
| `src/types/` | `CreateTransactionRequest` / `UpdateTransactionRequest` 型を追加 |

### ダイアログ仕様

shadcn/ui の `Dialog` コンポーネントを使用する。新規作成・編集で同一ダイアログを使いまわす（`mode: "create" | "edit"`）。

**フォーム項目:**

| 項目 | 入力形式 | バリデーション |
|---|---|---|
| 日付 | `<input type="date">` | 必須 |
| 摘要 | テキスト | 必須、最大255文字 |
| 金額 | 数値（円） | 必須、0以外 |
| カテゴリ | セレクト | 任意（未選択 = 未分類） |

**金額の入力UX:** 支出/収入のトグルボタン + 絶対値入力にして、内部で符号変換する（負数をユーザーに入力させない）。

### 取引一覧の変更点

- ヘッダーに「＋ 手動追加」ボタンを追加
- 各取引行に編集アイコン（鉛筆）と削除アイコン（ゴミ箱）を追加
- 削除は確認ダイアログを挟む（既存の削除APIをそのまま使用）

---

## 実装順序

1. entスキーマ変更（`import_format_id` を Optional/Nillable に）
2. `go generate ./ent/...`
3. Atlas でマイグレーション差分SQL生成
4. `entity.Transaction.ImportFormatID` を `*string` に変更
5. `usecase/transaction/usecase.go` — `CreateManualInput` / `UpdateInput` / Repository interfaceメソッド追加
6. `usecase/transaction/create.go` — `CreateManual` ユースケース実装
7. `usecase/transaction/update.go` — `Update` ユースケース実装
8. `service/persistence/transaction.go` — Repository実装追加
9. `handler/transaction.go` — `Create` / `Update` ハンドラー追加、interfaceに追加
10. `handler/router.go` — ルーティング追加
11. Swagger再生成
12. フロントエンド: `api.ts` → `hooks` → ダイアログコンポーネント → 取引一覧統合
