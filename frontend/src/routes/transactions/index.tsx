import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useTransactions, useUpdateTransactionCategory, useDeleteTransaction } from '@/hooks/useTransactions'
import { useCategories } from '@/hooks/useCategories'
import type { Category, Transaction } from '@/types'
import { Trash2 } from 'lucide-react'

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
  const { data: categories = [] } = useCategories()
  const updateCategory = useUpdateTransactionCategory()
  const deleteTransaction = useDeleteTransaction()

  const totalPages = data ? Math.ceil(data.total / pageSize) : 0

  function handleDelete(tx: Transaction) {
    if (confirm(`「${tx.description}」を削除しますか？`)) {
      deleteTransaction.mutate(tx.id)
    }
  }

  return (
    <div className="p-4 md:p-8 space-y-5">
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

      {/* モバイル: カードリスト */}
      <div className="md:hidden bg-white rounded-lg border border-zinc-200 divide-y divide-zinc-100">
        {isLoading ? (
          <div className="p-8 text-center text-zinc-400 text-sm">読み込み中…</div>
        ) : (data?.transactions ?? []).length === 0 ? (
          <div className="p-8 text-center text-zinc-400 text-sm">データがありません</div>
        ) : (data?.transactions ?? []).map(tx => (
          <TxCard
            key={tx.id}
            tx={tx}
            categories={categories}
            onCategoryChange={(catId) => updateCategory.mutate({ id: tx.id, categoryId: catId })}
            onDelete={() => handleDelete(tx)}
          />
        ))}
      </div>

      {/* デスクトップ: テーブル */}
      <div className="hidden md:block bg-white rounded-lg border border-zinc-200 overflow-x-auto">
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
                <th className="px-4 py-3 w-10"></th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-100">
              {(data?.transactions ?? []).map(tx => (
                <TxRow
                  key={tx.id}
                  tx={tx}
                  categories={categories}
                  onCategoryChange={(catId) => updateCategory.mutate({ id: tx.id, categoryId: catId })}
                  onDelete={() => handleDelete(tx)}
                />
              ))}
              {(data?.transactions ?? []).length === 0 && (
                <tr>
                  <td colSpan={5} className="px-4 py-8 text-center text-zinc-400">データがありません</td>
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

function TxCard({
  tx,
  categories,
  onCategoryChange,
  onDelete,
}: {
  tx: Transaction
  categories: Category[]
  onCategoryChange: (id: string | null) => void
  onDelete: () => void
}) {
  return (
    <div className="px-4 py-3 space-y-1.5">
      <div className="flex items-center justify-between gap-2">
        <span className="text-xs text-zinc-400">{tx.date.slice(0, 10)}</span>
        <div className="flex items-center gap-2">
          <span className={`text-sm font-medium tabular-nums ${tx.amount >= 0 ? 'text-emerald-600' : 'text-rose-600'}`}>
            {tx.amount >= 0 ? '+' : ''}{jpy(tx.amount)}
          </span>
          <button onClick={onDelete} className="text-zinc-300 hover:text-rose-500">
            <Trash2 size={14} />
          </button>
        </div>
      </div>
      <p className="text-sm text-zinc-800 truncate">{tx.description}</p>
      <select
        className="text-xs border border-zinc-200 rounded px-1.5 py-1 bg-white w-full"
        value={tx.category_id ?? ''}
        onChange={e => onCategoryChange(e.target.value || null)}
      >
        <option value="">未分類</option>
        {categories.map(cat => (
          <option key={cat.id} value={cat.id}>{cat.name}</option>
        ))}
      </select>
    </div>
  )
}

function TxRow({
  tx,
  categories,
  onCategoryChange,
  onDelete,
}: {
  tx: Transaction
  categories: Category[]
  onCategoryChange: (id: string | null) => void
  onDelete: () => void
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
          className="text-xs border border-zinc-200 rounded px-1.5 py-1 bg-white max-w-35"
          value={tx.category_id ?? ''}
          onChange={e => onCategoryChange(e.target.value || null)}
        >
          <option value="">未分類</option>
          {categories.map(cat => (
            <option key={cat.id} value={cat.id}>{cat.name}</option>
          ))}
        </select>
      </td>
      <td className="px-4 py-2.5 text-center">
        <button onClick={onDelete} className="text-zinc-300 hover:text-rose-500">
          <Trash2 size={14} />
        </button>
      </td>
    </tr>
  )
}
