import { createFileRoute } from '@tanstack/react-router'
import { useMemo, useState } from 'react'
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell } from 'recharts'
import { useTransactions } from '@/hooks/useTransactions'
import { TrendingUp, TrendingDown, Minus, ChevronLeft, ChevronRight } from 'lucide-react'

export const Route = createFileRoute('/')({
  component: Dashboard,
})

function jpy(amount: number) {
  return new Intl.NumberFormat('ja-JP', { style: 'currency', currency: 'JPY' }).format(Math.abs(amount))
}

function Dashboard() {
  const now = new Date()
  const [year, setYear] = useState(now.getFullYear())
  const [month, setMonth] = useState(now.getMonth() + 1)

  const isCurrentMonth = year === now.getFullYear() && month === now.getMonth() + 1

  function prevMonth() {
    if (month === 1) { setYear(y => y - 1); setMonth(12) }
    else setMonth(m => m - 1)
  }
  function nextMonth() {
    if (month === 12) { setYear(y => y + 1); setMonth(1) }
    else setMonth(m => m + 1)
  }

  const { data, isLoading } = useTransactions({ year, month, page_size: 500 })

  const stats = useMemo(() => {
    const txs = data?.transactions ?? []
    const income = txs.filter(t => t.amount > 0).reduce((s, t) => s + t.amount, 0)
    const expense = txs.filter(t => t.amount < 0).reduce((s, t) => s + t.amount, 0)

    const byCategory: Record<string, { name: string; color: string; total: number }> = {}
    for (const tx of txs) {
      if (tx.amount >= 0) continue
      const key = tx.category?.id ?? '__none__'
      if (!byCategory[key]) {
        byCategory[key] = { name: tx.category?.name ?? '未分類', color: tx.category?.color ?? '#94a3b8', total: 0 }
      }
      byCategory[key].total += Math.abs(tx.amount)
    }

    return {
      income,
      expense,
      net: income + expense,
      breakdown: Object.values(byCategory).sort((a, b) => b.total - a.total).slice(0, 8),
      recent: txs.slice(0, 10),
    }
  }, [data])

  if (isLoading) return <div className="p-8 text-zinc-500">読み込み中…</div>

  return (
    <div className="p-4 md:p-8 space-y-6">
      <div className="flex items-center gap-3">
        <button onClick={prevMonth} className="p-1 rounded hover:bg-zinc-100 text-zinc-500">
          <ChevronLeft size={20} />
        </button>
        <h1 className="text-xl font-bold text-zinc-800 whitespace-nowrap text-center">{year}年{month}月 の収支</h1>
        <button onClick={nextMonth} disabled={isCurrentMonth} className="p-1 rounded hover:bg-zinc-100 text-zinc-500 disabled:opacity-30 disabled:cursor-not-allowed">
          <ChevronRight size={20} />
        </button>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <SummaryCard label="収入" amount={stats.income} positive />
        <SummaryCard label="支出" amount={stats.expense} />
        <SummaryCard label="収支差" amount={stats.net} positive={stats.net >= 0} />
      </div>

      {stats.breakdown.length > 0 && (
        <div className="bg-white rounded-lg border border-zinc-200 p-5">
          <h2 className="text-sm font-semibold text-zinc-700 mb-4">カテゴリ別支出</h2>

          {/* モバイル: ランキングリスト */}
          <div className="md:hidden space-y-2">
            {stats.breakdown.map((e) => (
              <div key={e.name}>
                <div className="flex items-center justify-between mb-1">
                  <div className="flex items-center gap-2 min-w-0">
                    <span className="w-2.5 h-2.5 rounded-full shrink-0" style={{ backgroundColor: e.color }} />
                    <span className="text-sm text-zinc-700 truncate">{e.name}</span>
                  </div>
                  <span className="text-sm font-medium tabular-nums text-zinc-800 ml-2 shrink-0">{jpy(e.total)}</span>
                </div>
                <div className="h-1.5 rounded-full bg-zinc-100 overflow-hidden">
                  <div
                    className="h-full rounded-full"
                    style={{
                      width: `${Math.round((e.total / stats.breakdown[0].total) * 100)}%`,
                      backgroundColor: e.color,
                    }}
                  />
                </div>
              </div>
            ))}
          </div>

          {/* デスクトップ: 横棒グラフ（高さをアイテム数に応じて調整） */}
          <div className="hidden md:block">
            <ResponsiveContainer width="100%" height={stats.breakdown.length * 36 + 20}>
              <BarChart data={stats.breakdown} layout="vertical">
                <XAxis type="number" tickFormatter={v => `${Math.round((v as number) / 1000)}k`} tick={{ fontSize: 11 }} />
                <YAxis type="category" dataKey="name" width={80} tick={{ fontSize: 12 }} />
                <Tooltip formatter={(v) => jpy(Number(v))} />
                <Bar dataKey="total" radius={[0, 4, 4, 0]}>
                  {stats.breakdown.map((e, i) => <Cell key={i} fill={e.color} />)}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
      )}

      {stats.recent.length > 0 && (
        <div className="bg-white rounded-lg border border-zinc-200">
          <h2 className="text-sm font-semibold text-zinc-700 px-5 py-3 border-b border-zinc-100">最近の取引</h2>
          <ul className="divide-y divide-zinc-100">
            {stats.recent.map(tx => (
              <li key={tx.id} className="flex items-center justify-between gap-3 px-5 py-2.5">
                <div className="min-w-0">
                  <p className="text-sm text-zinc-800 truncate">{tx.description}</p>
                  <p className="text-xs text-zinc-400">{tx.date.slice(0, 10)}</p>
                </div>
                <span className={`text-sm font-medium tabular-nums shrink-0 ${tx.amount >= 0 ? 'text-emerald-600' : 'text-rose-600'}`}>
                  {tx.amount >= 0 ? '+' : ''}{jpy(tx.amount)}
                </span>
              </li>
            ))}
          </ul>
        </div>
      )}

      {stats.recent.length === 0 && (
        <p className="text-sm text-zinc-400">今月の取引データがありません。CSVをインポートしてください。</p>
      )}
    </div>
  )
}

function SummaryCard({ label, amount, positive }: { label: string; amount: number; positive?: boolean }) {
  const color = positive ? 'text-emerald-600' : 'text-rose-600'
  const Icon = positive ? TrendingUp : amount === 0 ? Minus : TrendingDown
  return (
    <div className="bg-white rounded-lg border border-zinc-200 p-5">
      <div className="flex items-center gap-2 text-zinc-500 text-xs mb-2">
        <Icon size={14} />{label}
      </div>
      <p className={`text-2xl font-bold tabular-nums ${color}`}>
        {amount >= 0 ? '' : '-'}{new Intl.NumberFormat('ja-JP', { style: 'currency', currency: 'JPY' }).format(Math.abs(amount))}
      </p>
    </div>
  )
}
