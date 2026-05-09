import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useTransactions, useUpdateTransactionCategory } from '@/hooks/useTransactions'
import { useCategories } from '@/hooks/useCategories'
import type { Transaction } from '@/types'

export const Route = createFileRoute('/transactions/')({
  component: TransactionsPage,
})

function jpy(amount: number) {
  return new Intl.NumberFormat('ja-JP', { style: 'currency', currency: 'JPY' }).format(Math.abs(amount))
}

function TransactionsPage() {
  const now = new Date()
  const [year, setYear] = useState(now.getFullYear())
  const [month, setMonth] = useState(now.getMonth() + 1)
  const [page, setPage] = useState(0)
  const pageSize = 50

  const { data, isLoading } = useTransactions({ year, month, page, page_size: pageSize })
  const { data: categories } = useCategories()
  const updateCategory = useUpdateTransactionCategory()

  const totalPages = data ? Math.ceil(data.total / pageSize) : 0

  return (
    <div className="p-8 space-y-5">
      <h1 className="text-xl font-bold text-zinc-800">取引一覧</h1>

      {/* Filters */}
      <div className="flex items-center gap-3">
        <select
          className="border border-zinc-200 rounded px-2 py-1.5 text-sm bg-white"
          value={year}
          onChange={e => { setYear(Number(e.target.value)); setPage(0) }}
        >
          {[2024, 2025, 2026, 2027].map(y => (
            <option key={y} value={y}>{y}年</option>
          ))}
        </select>
        <select
          className="border border-zinc-200 rounded px-2 py-1.5 text-sm bg-white"
          value={month}
          onChange={e => { setMonth(Number(e.target.value)); setPage(0) }}
        >
          {Array.from({ length: 12 }, (_, i) => i + 1).map(m => (
            <option key={m} value={m}>{m}月</option>
          ))}
        </select>
        <span className="text-sm text-zinc-400">{data?.total ?? 0} 件</span>
      </div>

      {/* Table */}
      <div className="bg-white rounded-lg border border-zinc-200 overflow-hidden">
        {isLoading ? (
          <div className="p-8 text-center text-zinc-400 text-sm">読み込み中…</div>
        ) : (
          <table className="w-full text-sm">
            <thead className="bg-zinc-50 text-zinc-600 text-xs uppercase tracking-wide">
              <tr>
                <th className="px-4 py-3 text-left">日付</th>
                <th className="px-4 py-3 text-left">摘要</th>
                <th className="px-4 py-3 text-right">金額</th>
                <th className="px-4 py-3 text-left">カテゴリ</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-100">
              {(data?.transactions ?? []).map(tx => (
                <TxRow
                  key={tx.id}
                  tx={tx}
                  categories={categories ?? []}
                  onCategoryChange={(catId) =>
                    updateCategory.mutate({ id: tx.id, categoryId: catId })
                  }
                />
              ))}
              {(data?.transactions ?? []).length === 0 && (
                <tr>
                  <td colSpan={4} className="px-4 py-8 text-center text-zinc-400">データがありません</td>
                </tr>
              )}
            </tbody>
          </table>
        )}
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex items-center gap-2 justify-center">
          <button
            disabled={page === 0}
            onClick={() => setPage(p => Math.max(0, p - 1))}
            className="px-3 py-1.5 rounded border border-zinc-200 text-sm disabled:opacity-40 hover:bg-zinc-50"
          >
            ← 前
          </button>
          <span className="text-sm text-zinc-600">{page + 1} / {totalPages}</span>
          <button
            disabled={page >= totalPages - 1}
            onClick={() => setPage(p => p + 1)}
            className="px-3 py-1.5 rounded border border-zinc-200 text-sm disabled:opacity-40 hover:bg-zinc-50"
          >
            次 →
          </button>
        </div>
      )}
    </div>
  )
}

function TxRow({
  tx,
  categories,
  onCategoryChange,
}: {
  tx: Transaction
  categories: import('@/types').Category[]
  onCategoryChange: (id: string | null) => void
}) {
  return (
    <tr className="hover:bg-zinc-50">
      <td className="px-4 py-2.5 text-zinc-500 text-xs whitespace-nowrap">{tx.date.slice(0, 10)}</td>
      <td className="px-4 py-2.5 text-zinc-800 max-w-xs truncate">{tx.description}</td>
      <td className={`px-4 py-2.5 text-right font-medium tabular-nums ${tx.amount >= 0 ? 'text-emerald-600' : 'text-rose-600'}`}>
        {tx.amount >= 0 ? '+' : ''}{jpy(tx.amount)}
      </td>
      <td className="px-4 py-2.5">
        <select
          className="text-xs border border-zinc-200 rounded px-1.5 py-1 bg-white max-w-[140px]"
          value={tx.category_id ?? ''}
          onChange={e => onCategoryChange(e.target.value || null)}
        >
          <option value="">未分類</option>
          {categories.map(cat => (
            <option key={cat.id} value={cat.id}>{cat.name}</option>
          ))}
        </select>
      </td>
    </tr>
  )
}
