# 年間カテゴリ分布機能 設計書

## 概要

トレンドページ（`/trends`）に「カテゴリ分布」タブを追加し、選択年の支出を**カテゴリ別に集計・可視化**する。  
現在は月単位でしか見られないカテゴリ内訳を、一年を通した視点で把握できるようにする。

---

## 画面設計

### レイアウト

`/trends` ページにタブを追加し、既存の月別収支グラフとカテゴリ分布を切り替え可能にする。

```
[2026年]  ← →

[ 月別収支 ]  [ カテゴリ分布 ]   ← タブ

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
（タブに応じたコンテンツ）
```

### 「カテゴリ分布」タブのコンテンツ

既存のダッシュボードと同じく、**モバイル・デスクトップで別レイアウト**を使う。

#### デスクトップ（`hidden md:block`）

```
年間支出合計: ¥1,234,567

┌─────────────────────────────────────────────┐
│  [ドーナツチャート]  │  カテゴリ凡例リスト       │
│                     │  ● 食費      ¥432,000  35%│
│     中央に合計表示   │  ● 交通費    ¥246,913  20%│
│                     │  ● 娯楽      ¥185,185  15%│
│                     │  ...                      │
└─────────────────────────────────────────────┘
```

- `PieChart` + `Pie`（`innerRadius` 設定でドーナツ形）+ `Cell`
- チャート右または下に凡例リスト（カテゴリ名・金額・割合）
- カテゴリ数が多い場合（8件超）は上位8件 + 「その他」にまとめる

#### モバイル（`md:hidden`）

ドーナツチャートはモバイルで**ラベルが潰れて読めない**ため非表示。
ダッシュボードのカテゴリ内訳と同じランキングリスト形式を採用。

```
年間支出合計: ¥1,234,567

● 食費           ¥432,000   35%
━━━━━━━━━━━━━━━━━━━━━░░░░
● 交通費         ¥246,913   20%
━━━━━━━━━░░░░░░░░░░░░░░░░
● 娯楽           ¥185,185   15%
━━━━━━░░░░░░░░░░░░░░░░░░░
● 未分類         ¥123,456   10%
...
```

- カラードット + カテゴリ名 + 金額 + 割合バー（`h-1.5` の細いプログレスバー）
- BottomNav（約 64px）と被らないよう `pb-[calc(4rem+env(safe-area-inset-bottom))]` は `Layout` 側で確保済み

---

## バックエンド設計

### 新規エンドポイント

```
GET /api/summary/category-annual?year=2026
```

#### レスポンス

```json
{
  "year": 2026,
  "total_expense": 1234567,
  "categories": [
    {
      "category_id": "uuid-xxx",
      "category_name": "食費",
      "category_color": "#10b981",
      "total_expense": 432000,
      "percentage": 35.0
    },
    {
      "category_id": null,
      "category_name": "未分類",
      "category_color": "#94a3b8",
      "total_expense": 123456,
      "percentage": 10.0
    }
  ]
}
```

- `category_id` が `null` の行は未分類（`category_id IS NULL`）の取引をまとめたもの
- 支出（`amount < 0`）のみ集計。収入は含めない
- `percentage` はバックエンドで計算して返す（フロントの計算を省く）
- 金額が 0 のカテゴリは返さない

#### SQL

```sql
SELECT
    t.category_id,
    COALESCE(c.name,  '未分類')   AS category_name,
    COALESCE(c.color, '#94a3b8') AS category_color,
    SUM(-t.amount)               AS total_expense
FROM transactions t
LEFT JOIN categories c ON c.id = t.category_id
WHERE YEAR(t.date) = ?
  AND t.amount < 0
GROUP BY t.category_id, c.name, c.color
ORDER BY total_expense DESC
```

`percentage` はアプリ層（usecase）で total_expense の合計から算出する。

### 新規ファイル

| ファイル | 役割 |
|---|---|
| `internal/usecase/summary/category_annual.go` | `CategoryAnnualUsecase` 実装 |
| （`internal/usecase/summary/usecase.go` を更新） | `Repository` インターフェースに `CategoryAnnualSummary` メソッド追加 |
| （`internal/service/persistence/summary.go` を更新） | SQL 実装追加 |
| （`internal/handler/summary.go` を更新） | `CategoryAnnual` ハンドラ追加 |
| （`internal/handler/router.go` を更新） | `summary.GET("/category-annual", ...)` 追加 |

---

## フロントエンド設計

### 新規型定義（`src/types/index.ts`）

```ts
export interface CategoryAnnualItem {
  category_id: string | null
  category_name: string
  category_color: string
  total_expense: number
  percentage: number
}

export interface CategoryAnnualResult {
  year: number
  total_expense: number
  categories: CategoryAnnualItem[]
}
```

### API クライアント（`src/lib/api.ts`）

```ts
summary: {
  monthly: (year: number) => ...,
  categoryAnnual: (year: number) =>
    request<CategoryAnnualResult>(`/summary/category-annual?year=${year}`),
},
```

### カスタムフック（`src/hooks/useSummary.ts`）

```ts
export function useCategoryAnnual(year: number) {
  return useQuery({
    queryKey: summaryKeys.categoryAnnual(year),
    queryFn: () => api.summary.categoryAnnual(year),
  })
}
```

### コンポーネント構成

```
src/routes/trends/index.tsx
  └─ TrendsPage
       ├─ タブ切り替え UI（状態: 'monthly' | 'category'）
       ├─ MonthlyTab（既存の内容をそのまま切り出す）
       └─ CategoryAnnualTab（新規）
            ├─ ドーナツチャート（recharts の PieChart）
            └─ カテゴリ別テーブル
```

`CategoryAnnualTab` は `src/routes/trends/` 直下に `_CategoryAnnualTab.tsx` として切り出すか、同ファイル内の関数コンポーネントで実装する。ファイルが肥大化しなければ同ファイルで十分。

### チャート実装方針

- ライブラリ: 既存の `recharts`（追加インストール不要）
- コンポーネント: `PieChart` + `Pie` + `Cell` + `Tooltip`
- `innerRadius` を設定してドーナツ形にする
- モバイルでは `cx/cy` を小さめに、デスクトップでは `ResponsiveContainer` で拡大

---

## 実装手順

1. **バックエンド**
   1. `CategoryAnnualSummary` メソッドを `Repository` インターフェースに追加
   2. `persistence/summary.go` に SQL 実装
   3. `usecase/summary/category_annual.go` を作成
   4. `handler/summary.go` に `CategoryAnnual` ハンドラ追加
   5. `router.go` にルート追加

2. **フロントエンド**
   1. `src/types/index.ts` に型を追加
   2. `src/lib/api.ts` に `categoryAnnual` を追加
   3. `src/hooks/useSummary.ts` に `useCategoryAnnual` を追加
   4. `src/routes/trends/index.tsx` をタブ構造にリファクタし `CategoryAnnualTab` を実装

---

## スコープ外（将来検討）

- カテゴリ × 月のヒートマップ（どの月にどのカテゴリが多いか）
- 複数年比較（カテゴリ構成の経年変化）
- 収入側のカテゴリ分布
