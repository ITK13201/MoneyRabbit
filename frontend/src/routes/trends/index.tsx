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
} from 'recharts'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import { useMonthlySummary } from '@/hooks/useSummary'
import type { MonthlySummary } from '@/types'

export const Route = createFileRoute('/trends/')({
  component: TrendsPage,
})

function jpy(amount: number) {
  return new Intl.NumberFormat('ja-JP', { style: 'currency', currency: 'JPY' }).format(amount)
}

function jpyShort(amount: number) {
  if (Math.abs(amount) >= 10000) {
    return `${Math.round(amount / 1000)}k`
  }
  return String(amount)
}

function buildChartData(months: MonthlySummary[], year: number) {
  return Array.from({ length: 12 }, (_, i) => {
    const m = i + 1
    const found = months.find(s => s.year === year && s.month === m)
    return {
      label: `${m}月`,
      income: found?.income ?? 0,
      expense: found?.expense ?? 0,
      net: found ? found.income - found.expense : 0,
    }
  })
}

function TrendsPage() {
  const now = new Date()
  const [year, setYear] = useState(now.getFullYear())

  const { data, isLoading } = useMonthlySummary(year)
  const months = data?.months ?? []
  const chartData = buildChartData(months, year)

  const totalIncome = months.reduce((s, m) => s + m.income, 0)
  const totalExpense = months.reduce((s, m) => s + m.expense, 0)
  const totalNet = totalIncome - totalExpense

  return (
    <div className="p-4 md:p-8 space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3">
        <button
          onClick={() => setYear(y => y - 1)}
          className="p-1 rounded hover:bg-zinc-100 text-zinc-500"
        >
          <ChevronLeft size={20} />
        </button>
        <h1 className="text-xl font-bold text-zinc-800 whitespace-nowrap">{year}年 月別収支</h1>
        <button
          onClick={() => setYear(y => y + 1)}
          disabled={year >= now.getFullYear()}
          className="p-1 rounded hover:bg-zinc-100 text-zinc-500 disabled:opacity-30 disabled:cursor-not-allowed"
        >
          <ChevronRight size={20} />
        </button>
      </div>

      {/* Annual summary cards */}
      <div className="grid grid-cols-3 gap-3">
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
            {totalNet >= 0 ? '+' : ''}{jpy(totalNet)}
          </p>
        </div>
      </div>

      {/* Chart */}
      <div className="bg-white rounded-lg border border-zinc-200 p-5">
        <h2 className="text-sm font-semibold text-zinc-700 mb-4">月別収支グラフ</h2>

        {isLoading ? (
          <div className="h-64 flex items-center justify-center text-zinc-400 text-sm">読み込み中…</div>
        ) : (
          /* Mobile: horizontally scrollable; Desktop: full width */
          <div className="overflow-x-auto -mx-1">
            <div className="min-w-[600px] md:min-w-0">
              <ResponsiveContainer width="100%" height={280}>
                <ComposedChart data={chartData} margin={{ top: 4, right: 8, left: 8, bottom: 0 }}>
                  <XAxis dataKey="label" tick={{ fontSize: 11 }} />
                  <YAxis tickFormatter={jpyShort} tick={{ fontSize: 11 }} width={48} />
                  <Tooltip
                    formatter={(value, name) => {
                      const labels: Record<string, string> = { income: '収入', expense: '支出', net: '収支差' }
                      return [jpy(Number(value)), labels[name as string] ?? name]
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
                  <Line
                    type="monotone"
                    dataKey="net"
                    stroke="#71717a"
                    strokeWidth={2}
                    dot={{ r: 3, fill: '#71717a' }}
                  />
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
            {chartData.map((row) => {
              const net = row.income - row.expense
              return (
                <tr key={row.label} className="hover:bg-zinc-50">
                  <td className="px-4 py-2.5 font-medium text-zinc-700">{row.label}</td>
                  <td className="px-4 py-2.5 text-right tabular-nums text-emerald-600">
                    {row.income > 0 ? jpy(row.income) : '—'}
                  </td>
                  <td className="px-4 py-2.5 text-right tabular-nums text-rose-600">
                    {row.expense > 0 ? jpy(row.expense) : '—'}
                  </td>
                  <td className={`px-4 py-2.5 text-right tabular-nums font-medium ${net > 0 ? 'text-emerald-600' : net < 0 ? 'text-rose-600' : 'text-zinc-400'}`}>
                    {row.income === 0 && row.expense === 0 ? '—' : `${net >= 0 ? '+' : ''}${jpy(net)}`}
                  </td>
                </tr>
              )
            })}
          </tbody>
        </table>
      </div>
    </div>
  )
}
