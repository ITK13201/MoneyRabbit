import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import {
  ComposedChart,
  Bar,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
  ResponsiveContainer,
  PieChart,
  Pie,
  Cell,
} from 'recharts'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import { usePeriodStore } from '@/stores/periodStore'
import { useMonthlySummary, useCategoryAnnual } from '@/hooks/useSummary'
import { cn } from '@/lib/utils'
import type { MonthlySummary, CategoryAnnualItem } from '@/types'

export const Route = createFileRoute('/trends/')({
  component: TrendsPage,
})

type Tab = 'monthly' | 'category'

function jpy(amount: number) {
  return new Intl.NumberFormat('ja-JP', { style: 'currency', currency: 'JPY' }).format(amount)
}

function jpyShort(amount: number) {
  if (Math.abs(amount) >= 10000) return `${Math.round(amount / 1000)}k`
  return String(amount)
}

// 上位 MAX_CATEGORIES 件に絞り、残りを「その他」にまとめる
const MAX_CATEGORIES = 8

function mergeSmallCategories(categories: CategoryAnnualItem[]): CategoryAnnualItem[] {
  if (categories.length <= MAX_CATEGORIES) return categories
  const top = categories.slice(0, MAX_CATEGORIES - 1)
  const rest = categories.slice(MAX_CATEGORIES - 1)
  const restTotal = rest.reduce((s, c) => s + c.total_expense, 0)
  const totalAll = categories.reduce((s, c) => s + c.total_expense, 0)
  return [
    ...top,
    {
      category_id: null,
      category_name: 'その他',
      category_color: '#d4d4d8',
      total_expense: restTotal,
      percentage: totalAll > 0 ? (restTotal / totalAll) * 100 : 0,
    },
  ]
}

// ──────────────────────────────────────────────
// 月別収支タブ
// ──────────────────────────────────────────────

function buildChartData(months: MonthlySummary[], year: number) {
  return Array.from({ length: 12 }, (_, i) => {
    const m = i + 1
    const found = months.find(s => s.year === year && s.month === m)
    const income = found?.income ?? 0
    const expense = found?.expense ?? 0
    return { label: `${m}月`, income, expense: -expense, net: income - expense }
  })
}

function MonthlyTab({ year }: { year: number }) {
  const { data, isLoading } = useMonthlySummary(year)
  const months = data?.months ?? []
  const chartData = buildChartData(months, year)

  const totalIncome = months.reduce((s, m) => s + m.income, 0)
  const totalExpense = months.reduce((s, m) => s + m.expense, 0)
  const totalNet = totalIncome - totalExpense

  return (
    <div className="space-y-6">
      {/* Annual summary cards */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div className="bg-white rounded-lg border border-zinc-200 p-4">
          <p className="text-xs text-zinc-500 mb-1">年間収入</p>
          <p className="text-lg font-bold tabular-nums text-emerald-600">{jpy(totalIncome)}</p>
        </div>
        <div className="bg-white rounded-lg border border-zinc-200 p-4">
          <p className="text-xs text-zinc-500 mb-1">年間支出</p>
          <p className="text-lg font-bold tabular-nums text-rose-600">{jpy(totalExpense)}</p>
        </div>
        <div className="bg-white rounded-lg border border-zinc-200 p-4">
          <p className="text-xs text-zinc-500 mb-1">年間収支</p>
          <p className={`text-lg font-bold tabular-nums ${totalNet >= 0 ? 'text-emerald-600' : 'text-rose-600'}`}>
            {(totalNet >= 0 ? '+' : '') + jpy(totalNet)}
          </p>
        </div>
      </div>

      {/* Chart */}
      <div className="bg-white rounded-lg border border-zinc-200 p-5">
        <h2 className="text-sm font-semibold text-zinc-700 mb-4">月別収支グラフ</h2>
        {isLoading ? (
          <div className="h-64 flex items-center justify-center text-zinc-400 text-sm">読み込み中…</div>
        ) : (
          <div className="overflow-x-auto">
            <div style={{ minWidth: 560 }}>
              <ResponsiveContainer width="100%" height={280}>
                <ComposedChart data={chartData} margin={{ top: 4, right: 8, left: 8, bottom: 0 }}>
                  <XAxis dataKey="label" tick={{ fontSize: 11 }} />
                  <YAxis tickFormatter={jpyShort} tick={{ fontSize: 11 }} width={48} />
                  <Tooltip
                    formatter={(value, name) => {
                      const labels: Record<string, string> = { income: '収入', expense: '支出', net: '収支差' }
                      return [jpy(Math.abs(Number(value))), labels[name as string] ?? name]
                    }}
                  />
                  <Legend
                    formatter={(value) => {
                      const labels: Record<string, string> = { income: '収入', expense: '支出', net: '収支差' }
                      return labels[value] ?? value
                    }}
                  />
                  <Bar dataKey="income" fill="#10b981" radius={[2, 2, 0, 0]} maxBarSize={28} />
                  <Bar dataKey="expense" fill="#f43f5e" radius={[2, 2, 0, 0]} maxBarSize={28} />
                  <Line type="monotone" dataKey="net" stroke="#71717a" strokeWidth={2} dot={{ r: 3, fill: '#71717a' }} />
                </ComposedChart>
              </ResponsiveContainer>
            </div>
          </div>
        )}
      </div>

      {/* Monthly table */}
      <div className="bg-white rounded-lg border border-zinc-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-zinc-50 text-zinc-600 text-xs uppercase tracking-wide">
            <tr>
              <th className="px-4 py-3 text-left">月</th>
              <th className="px-4 py-3 text-right">収入</th>
              <th className="px-4 py-3 text-right">支出</th>
              <th className="px-4 py-3 text-right">収支差</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-zinc-100">
            {chartData.map((row) => (
              <tr key={row.label} className="hover:bg-zinc-50">
                <td className="px-4 py-2.5 font-medium text-zinc-700 whitespace-nowrap">{row.label}</td>
                <td className="px-4 py-2.5 text-right tabular-nums text-emerald-600 whitespace-nowrap">
                  {row.income > 0 ? jpy(row.income) : '—'}
                </td>
                <td className="px-4 py-2.5 text-right tabular-nums text-rose-600 whitespace-nowrap">
                  {row.expense < 0 ? jpy(Math.abs(row.expense)) : '—'}
                </td>
                <td className={`px-4 py-2.5 text-right tabular-nums font-medium whitespace-nowrap ${row.net > 0 ? 'text-emerald-600' : row.net < 0 ? 'text-rose-600' : 'text-zinc-400'}`}>
                  {row.income === 0 && row.expense === 0 ? '—' : `${row.net >= 0 ? '+' : ''}${jpy(row.net)}`}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

// ──────────────────────────────────────────────
// 年間カテゴリ分布タブ
// ──────────────────────────────────────────────

function CategoryAnnualTab({ year }: { year: number }) {
  const { data, isLoading } = useCategoryAnnual(year)
  const categories = mergeSmallCategories(data?.categories ?? [])
  const totalExpense = data?.total_expense ?? 0

  if (isLoading) {
    return <div className="py-16 text-center text-zinc-400 text-sm">読み込み中…</div>
  }

  if (categories.length === 0) {
    return <div className="py-16 text-center text-zinc-400 text-sm">この年の支出データがありません。</div>
  }

  return (
    <div className="space-y-4">
      {/* 年間支出合計 */}
      <div className="bg-white rounded-lg border border-zinc-200 p-4">
        <p className="text-xs text-zinc-500 mb-1">年間支出合計</p>
        <p className="text-lg font-bold tabular-nums text-rose-600">{jpy(totalExpense)}</p>
      </div>

      {/* モバイル: ランキングリスト */}
      <div className="md:hidden bg-white rounded-lg border border-zinc-200 p-5 space-y-3">
        {categories.map((c) => (
          <div key={c.category_id ?? '__none__'}>
            <div className="flex items-center justify-between mb-1">
              <div className="flex items-center gap-2 min-w-0">
                <span className="w-2.5 h-2.5 rounded-full shrink-0" style={{ backgroundColor: c.category_color }} />
                <span className="text-sm text-zinc-700 truncate">{c.category_name}</span>
              </div>
              <div className="flex items-center gap-2 shrink-0 ml-2">
                <span className="text-xs text-zinc-400 tabular-nums">{c.percentage.toFixed(1)}%</span>
                <span className="text-sm font-medium tabular-nums text-zinc-800">{jpy(c.total_expense)}</span>
              </div>
            </div>
            <div className="h-1.5 rounded-full bg-zinc-100 overflow-hidden">
              <div
                className="h-full rounded-full"
                style={{ width: `${c.percentage}%`, backgroundColor: c.category_color }}
              />
            </div>
          </div>
        ))}
      </div>

      {/* デスクトップ: ドーナツチャート + 凡例リスト */}
      <div className="hidden md:block bg-white rounded-lg border border-zinc-200 p-6">
        <div className="flex items-center gap-8">
          {/* ドーナツチャート */}
          <div className="relative shrink-0">
            <PieChart width={220} height={220}>
              <Pie
                data={categories}
                cx={110}
                cy={110}
                innerRadius={65}
                outerRadius={105}
                dataKey="total_expense"
                strokeWidth={0}
              >
                {categories.map((c) => (
                  <Cell key={c.category_id ?? '__none__'} fill={c.category_color} />
                ))}
              </Pie>
              <Tooltip
                formatter={(value, _name, props) => [
                  jpy(Number(value)),
                  (props.payload as CategoryAnnualItem).category_name,
                ]}
              />
            </PieChart>
            {/* 中央ラベル */}
            <div className="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
              <span className="text-xs text-zinc-400">支出合計</span>
              <span className="text-sm font-bold text-zinc-700 tabular-nums">{jpy(totalExpense)}</span>
            </div>
          </div>

          {/* 凡例リスト */}
          <div className="flex-1 space-y-2.5">
            {categories.map((c) => (
              <div key={c.category_id ?? '__none__'} className="flex items-center gap-2">
                <span className="w-3 h-3 rounded-sm shrink-0" style={{ backgroundColor: c.category_color }} />
                <span className="text-sm text-zinc-700 flex-1 truncate">{c.category_name}</span>
                <span className="text-xs text-zinc-400 tabular-nums w-12 text-right">{c.percentage.toFixed(1)}%</span>
                <span className="text-sm font-medium tabular-nums text-zinc-800 w-28 text-right">{jpy(c.total_expense)}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  )
}

// ──────────────────────────────────────────────
// ページ本体
// ──────────────────────────────────────────────

function TrendsPage() {
  const now = new Date()
  const { trendsYear: year, prevTrendsYear, nextTrendsYear } = usePeriodStore()
  const [tab, setTab] = useState<Tab>('monthly')

  return (
    <div className="p-4 md:p-8 space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3">
        <button onClick={prevTrendsYear} className="p-1 rounded hover:bg-zinc-100 text-zinc-500">
          <ChevronLeft size={20} />
        </button>
        <h1 className="text-xl font-bold text-zinc-800 whitespace-nowrap">{year}年</h1>
        <button
          onClick={nextTrendsYear}
          disabled={year >= now.getFullYear()}
          className="p-1 rounded hover:bg-zinc-100 text-zinc-500 disabled:opacity-30 disabled:cursor-not-allowed"
        >
          <ChevronRight size={20} />
        </button>
      </div>

      {/* Tab switcher */}
      <div className="flex gap-1 bg-zinc-100 rounded-lg p-1 w-fit">
        {(['monthly', 'category'] as const).map((t) => (
          <button
            key={t}
            onClick={() => setTab(t)}
            className={cn(
              'px-3 py-1.5 text-sm rounded-md transition-colors',
              tab === t
                ? 'bg-white shadow-sm text-zinc-900 font-medium'
                : 'text-zinc-500 hover:text-zinc-700',
            )}
          >
            {t === 'monthly' ? '月別収支' : 'カテゴリ分布'}
          </button>
        ))}
      </div>

      {/* Content */}
      {tab === 'monthly' ? <MonthlyTab year={year} /> : <CategoryAnnualTab year={year} />}
    </div>
  )
}
