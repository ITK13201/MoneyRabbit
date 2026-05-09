# CSVインポート設計

## 概要

CSVインポートには以下の2種類がある。

| 種別 | 説明 | 例 |
|---|---|---|
| `bank_account` | 銀行の出入金履歴 | 三井住友銀行 入出金明細 |
| `credit_card` | クレジットカードの支払い明細 | SMBCカード 利用明細 |

両者はCSVの列構造が異なる。フォーマット設定はGoコードの定数として定義し、DBには持たない。

---

## ImportFormat の定義

対応フォーマットはGoコードの定数として定義する。DBテーブルは持たない。対応フォーマットを追加する場合はコード変更のみでよい。

```go
// backend/internal/service/csv/formats.go

type ImportFormatID string

const (
    ImportFormatSMBCBank ImportFormatID = "smbc_bank" // 三井住友銀行
    ImportFormatSMBCCard ImportFormatID = "smbc_card" // SMBCカード
)

type ImportFormat struct {
    ID          ImportFormatID
    Name        string // 表示名
    ImportType  string // "bank_account" | "credit_card"
    Encoding    string // "Shift_JIS" | "UTF-8"
    SkipRows    int
    ColDate     int
    ColDesc     int
    ColAmount   int // credit_card: 利用金額列（-1=未使用）
    ColWithdraw int // bank_account: 出金列（-1=未使用）
    ColDeposit  int // bank_account: 入金列（-1=未使用）
    ColBalance  int // 残高列（-1=未使用）
    DateFormat  string // Go形式
}

var Formats = map[ImportFormatID]ImportFormat{
    ImportFormatSMBCBank: { /* 後述 */ },
    ImportFormatSMBCCard: { /* 後述 */ },
}
```

### import_type による列の使い分け

| import_type | 使用する列 | amountの符号 |
|---|---|---|
| `bank_account` | `ColWithdraw`（出金）/ `ColDeposit`（入金） | 出金=負、入金=正 |
| `credit_card` | `ColAmount`（利用金額） | 常に負 |

`Transaction.amount` は符号付き整数（正=収入、負=支出）。

---

## 銀行別 ImportFormat 設定例

### 三井住友銀行（出入金履歴）

実際のCSVフォーマット（文字コード: Shift-JIS）:
```
年月日,お引出し,お預入れ,お取り扱い内容,残高,メモ,ラベル
2026/4/2,1480,,"V907163　券売機（東日本高速道路）／ｉＤ",4600092,"",
2026/4/10,,826,"振込　ｹｲﾃﾞｲﾃﾞｲｱｲｱｼﾞﾔｲﾙｶｲﾊﾂｾﾝﾀ-(ｶ",4600618,"",
2026/4/24,,359562,"給料振込　ｹ-ﾃﾞｲ-ﾃﾞｲ-ｱｲｱｼﾞﾔｲﾙｶｲﾊﾂｾﾝﾀ-ｶﾌﾞｼｷｶﾞｲｼﾔ",4957768,"",
```

| フィールド | 値 | 備考 |
|---|---|---|
| import_type | `bank_account` | |
| encoding | `Shift_JIS` | |
| skip_rows | `1` | ヘッダー行1行 |
| col_date | `0` | `年月日` 列 |
| col_description | `3` | `お取り扱い内容` 列 |
| col_withdrawal | `1` | `お引出し` 列 |
| col_deposit | `2` | `お預入れ` 列 |
| col_balance | `4` | `残高` 列（保存しない） |
| col_amount | `-1` | 未使用 |
| date_format | `"2006/1/2"` | 月・日のゼロパディングなし |

> **実装注意**: 日付が `2026/4/2`（ゼロパディングなし）のため、Go形式では `"2006/1/2"` を使用する。フロントエンド（JS）では `split('/')` で分割して `padStart(2,'0')` でISO 8601に変換する。

---

### SMBCカード（支払い明細）

実際のCSVフォーマット（文字コード: Shift-JIS）:
```
池田　匠　様,4980-09**-****-****,Ｏｌｉｖｅ／クレジット
2026/03/01,ヨドバシカメラ　通信販売（新経路）,166,１,１,166,
2026/03/06,ＧＯＯＧＬＥ　ＰＬＡＹ　ＪＡＰＡＮ,1280,１,１,1280,
2026/03/09,CLAUDE.AI SUBSCRIPTION (ANTHROPIC.COM),3281,１,１,3281,20.00　USD　164.096　03 10
,,,,,187313,
```

ヘッダー行はCSVに含まれない（Webサイト上での表示のみ）:
```
ご利用日  ご利用店名  ご利用金額  支払区分  今回回数  お支払い金額  （お支払い総額）  （内手数料）
```

| フィールド | 値 | 備考 |
|---|---|---|
| import_type | `credit_card` | |
| encoding | `Shift_JIS` | |
| skip_rows | `1` | 先頭の氏名・カード番号行 |
| col_date | `0` | `ご利用日` 列 |
| col_description | `1` | `ご利用店名` 列 |
| col_withdrawal | `-1` | 未使用 |
| col_deposit | `-1` | 未使用 |
| col_amount | `2` | `ご利用金額` 列（col5=お支払い金額ではなく、col2=ご利用金額を使う） |
| col_balance | `-1` | 未使用 |
| date_format | `"2006/01/02"` | ゼロパディングあり |

> **col2 vs col5**: `ご利用金額`（col2）は実際の利用総額、`お支払い金額`（col5）は分割払い時の今回分のみ。家計簿としては支出の実態を記録するためcol2を使用する。

> **最終行の合計行**: `,,,,,187313,` のようにcol0が空の行が末尾に存在する。col0（日付）が空の行はスキップする。

> **海外利用**: col7に `20.00　USD　164.096　03 10` のような通貨情報が付く場合があるが、col2（円換算済み金額）を使うため無視してよい。

---

## パース共通ロジック

### スキップ条件

| 条件 | 理由 |
|---|---|
| `row[col_date]` が空 | 合計行・空行 |
| `row[col_date]` が日付として不正 | メタデータ行 |
| `col_withdrawal` と `col_deposit` が両方空または0 | 残高のみ行 |

### 金額の正規化

```typescript
// カンマ区切り数値の除去 + 全角数字の変換
function parseAmount(raw: string): number {
  const normalized = raw
    .replace(/[０-９]/g, c => String.fromCharCode(c.charCodeAt(0) - 0xFEE0)) // 全角→半角
    .replace(/,/g, '')  // カンマ除去
  return parseInt(normalized, 10)
}
```

SMBCカードの支払区分（`１`）や今回回数（`１`）は全角数字なので、全角→半角変換が必要。

### 日付の正規化

```typescript
function parseDate(raw: string): string {
  const parts = raw.split('/')
  if (parts.length !== 3) throw new Error(`不正な日付: ${raw}`)
  const [y, m, d] = parts.map(Number)
  return `${y}-${String(m).padStart(2, '0')}-${String(d).padStart(2, '0')}`
}
// "2026/4/2"  → "2026-04-02"
// "2026/03/01" → "2026-03-01"
```

---

## インポートフロー

```
[フロントエンド]                              [バックエンド]
  1. インポート種別を選択（smbc_bank / smbc_card）
  2. CSVファイルを選択
  3. POST /api/import/preview ────────────→ import_format_id でフォーマット取得
    multipart: file + import_format_id       → Shift-JIS変換 + CSVパース
                                             → 列マッピング・金額正規化
                          ←────────────────  { rows: [...] }  ※保存しない
  4. プレビューテーブルを表示
  5. ユーザーが不要行を削除
  6. 「確定してインポート」
  7. POST /api/import/confirm ─────────────→ 重複検出 → カテゴリ分類 → DB保存
    JSON: { import_format_id, transactions: [...] } ←─  { imported, skipped, errors }
  8. 結果を表示
```

### バックエンドのパース責務（`internal/service/csv/`）

CSVパース処理はバックエンドで行う。ImportFormatの設定がサーバー側にあるため、パースロジックもサーバーに置くのが自然（将来のバッチ処理・自動インポートにも対応しやすい）。

1. `golang.org/x/text/encoding/japanese` でShift-JIS→UTF-8変換
2. `encoding/csv` でパース
3. `skip_rows` 行分をスキップしてデータ行を処理
4. `import_type` に応じて列マッピングを切り替え
5. 金額文字列の正規化（全角数字・カンマ除去）
6. 日付を `time.Time` にパース（`date_format` を使用）

---

## APIエンドポイント

### `GET /api/import-formats`

対応済みフォーマットの一覧を返す（読み取り専用）。口座作成時のセレクトボックスに使用する。

```json
[
  { "id": "smbc_bank", "name": "三井住友銀行", "import_type": "bank_account" },
  { "id": "smbc_card", "name": "SMBCカード",   "import_type": "credit_card" }
]
```

### `POST /api/import/preview`

CSVファイルをアップロードし、パース結果をプレビューとして返す（保存しない）。

**リクエスト**: `multipart/form-data`
- `file`: CSVファイル
- `import_format_id`: フォーマットID（`smbc_bank` or `smbc_card`）

**レスポンス**:
```json
{
  "rows": [
    {
      "date": "2026-04-02",
      "description": "V907163　券売機（東日本高速道路）／ｉＤ",
      "amount": -1480
    },
    {
      "date": "2026-04-24",
      "description": "給料振込　ｹ-ﾃﾞｲ-ﾃﾞｲ-ｱｲｱｼﾞﾔｲﾙｶｲﾊﾂｾﾝﾀ-ｶﾌﾞｼｷｶﾞｲｼﾔ",
      "amount": 359562
    }
  ],
  "skipped_rows": 1
}
```

- `amount`: 符号付き整数（負=支出、正=収入）
- `skipped_rows`: 合計行・空行などスキップした行数

### `POST /api/import/confirm`

プレビューで確認済みの取引データを保存する。ユーザーが削除した行は含めない。

**リクエスト**: `application/json`
```json
{
  "import_format_id": "smbc_bank",
  "transactions": [
    { "date": "2026-04-02", "description": "...", "amount": -1480 }
  ]
}
```

**レスポンス**:
```json
{
  "imported": 8,
  "skipped": 1,
  "errors": []
}
```

- `skipped`: 重複検出により除外された件数
- `errors`: バリデーションエラーの詳細（空配列が正常）

---

## バックエンド処理（`POST /api/import`）

```
受信した transactions[]
    ↓
バリデーション（date形式、amount != 0）
    ↓
重複検出
    同一 import_format_id + date + amount + description の既存レコードをスキップ
    ↓
カテゴリ自動分類（features.md 参照）
    1. キーワードルール（CategoryRule）でマッチ
    2. マッチなし → Claude API（claude-haiku-4-5）でバッチ推論
    3. APIエラー → category=null（未分類）のまま保存
    ↓
Transaction 一括保存（BulkCreate）
```

---

## 未対応事項（将来対応）

| 項目 | 理由 |
|---|---|
| 複数ファイルの一括インポート | 現在は1ファイルずつ |
| インポート履歴の管理 | どのファイルをいつインポートしたかの記録 |
| 文字コードの自動検出 | 現在はImportFormatの設定値に従う |
