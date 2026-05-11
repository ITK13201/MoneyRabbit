import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import {
  useTransactions,
  useUpdateTransactionCategory,
  useDeleteTransaction,
  useCreateTransaction,
  useUpdateTransaction,
} from '@/hooks/useTransactions'
import { useCategories } from '@/hooks/useCategories'
import type { Category, Transaction } from '@/types'
import { Trash2, ChevronLeft, ChevronRight, Pencil, Plus, X } from 'lucide-react'

export const Route = createFileRoute('/transactions/')({
  component: TransactionsPage,
})

function jpy(amount: number) {
  return new Intl.NumberFormat('ja-JP', { style: 'currency', currency: 'JPY' }).format(
    Math.abs(amount),
  )
}

// ---- Transaction Form Dialog ----

interface TxFormState {
  date: string
  description: string
  absAmount: string
  isExpense: boolean
  categoryId: string
}

function emptyForm(): TxFormState {
  return {
    date: new Date().toISOString().slice(0, 10),
    description: '',
    absAmount: '',
    isExpense: true,
    categoryId: '',
  }
}

function txToForm(tx: Transaction): TxFormState {
  return {
    date: tx.date.slice(0, 10),
    description: tx.description,
    absAmount: String(Math.abs(tx.amount)),
    isExpense: tx.amount < 0,
    categoryId: tx.category_id ?? '',
  }
}

interface TxDialogProps {
  mode: 'create' | 'edit'
  tx?: Transaction
  categories: Category[]
  onClose: () => void
}

function TxDialog({ mode, tx, categories, onClose }: TxDialogProps) {
  const [form, setForm] = useState<TxFormState>(mode === 'edit' && tx ? txToForm(tx) : emptyForm())
  const [error, setError] = useState('')
  const createTx = useCreateTransaction()
  const updateTx = useUpdateTransaction()

  const isPending = createTx.isPending || updateTx.isPending

  function set<K extends keyof TxFormState>(key: K, value: TxFormState[K]) {
    setForm((f) => ({ ...f, [key]: value }))
  }

  function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault()
    setError('')

    const absVal = parseInt(form.absAmount, 10)
    if (!form.date || !form.description.trim()) {
      setError('日付と摘要は必須です')
      return
    }
    if (!form.absAmount || isNaN(absVal) || absVal <= 0) {
      setError('金額は1以上の整数を入力してください')
      return
    }

    const amount = form.isExpense ? -absVal : absVal
    const input = {
      date: form.date,
      description: form.description.trim(),
      amount,
      category_id: form.categoryId || null,
    }

    if (mode === 'create') {
      createTx.mutate(input, { onSuccess: onClose, onError: (e) => setError(e.message) })
    } else {
      updateTx.mutate(
        { id: tx!.id, input },
        { onSuccess: onClose, onError: (e) => setError(e.message) },
      )
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/40">
      <div className="bg-white rounded-xl shadow-xl w-full max-w-md">
        <div className="flex items-center justify-between px-5 py-4 border-b border-zinc-200">
          <h2 className="text-base font-semibold text-zinc-800">
            {mode === 'create' ? '取引を追加' : '取引を編集'}
          </h2>
          <button
            onClick={onClose}
            className="text-zinc-400 hover:text-zinc-600"
            aria-label="閉じる"
          >
            <X size={18} />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="px-5 py-4 space-y-4">
          {/* 日付 */}
          <div className="space-y-1">
            <label className="text-xs font-medium text-zinc-600">日付</label>
            <input
              type="date"
              value={form.date}
              onChange={(e) => set('date', e.target.value)}
              className="w-full border border-zinc-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>

          {/* 摘要 */}
          <div className="space-y-1">
            <label className="text-xs font-medium text-zinc-600">摘要</label>
            <input
              type="text"
              value={form.description}
              onChange={(e) => set('description', e.target.value)}
              placeholder="例: スーパーABC"
              className="w-full border border-zinc-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              required
            />
          </div>

          {/* 金額 */}
          <div className="space-y-1">
            <label className="text-xs font-medium text-zinc-600">金額</label>
            <div className="flex gap-2">
              <div className="flex rounded-lg border border-zinc-200 overflow-hidden text-sm shrink-0">
                <button
                  type="button"
                  onClick={() => set('isExpense', true)}
                  className={`px-3 py-2 transition-colors ${form.isExpense ? 'bg-rose-500 text-white' : 'bg-white text-zinc-500 hover:bg-zinc-50'}`}
                >
                  支出
                </button>
                <button
                  type="button"
                  onClick={() => set('isExpense', false)}
                  className={`px-3 py-2 transition-colors ${!form.isExpense ? 'bg-emerald-500 text-white' : 'bg-white text-zinc-500 hover:bg-zinc-50'}`}
                >
                  収入
                </button>
              </div>
              <input
                type="number"
                min="1"
                value={form.absAmount}
                onChange={(e) => set('absAmount', e.target.value)}
                placeholder="金額（円）"
                className="flex-1 border border-zinc-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
              />
            </div>
          </div>

          {/* カテゴリ */}
          <div className="space-y-1">
            <label className="text-xs font-medium text-zinc-600">カテゴリ（任意）</label>
            <select
              value={form.categoryId}
              onChange={(e) => set('categoryId', e.target.value)}
              className="w-full border border-zinc-200 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 bg-white"
            >
              <option value="">未分類</option>
              {categories.map((cat) => (
                <option key={cat.id} value={cat.id}>
                  {cat.name}
                </option>
              ))}
            </select>
          </div>

          {error && <p className="text-xs text-rose-500">{error}</p>}

          <div className="flex gap-2 pt-1">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 px-4 py-2 rounded-lg border border-zinc-200 text-sm text-zinc-600 hover:bg-zinc-50"
            >
              キャンセル
            </button>
            <button
              type="submit"
              disabled={isPending}
              className="flex-1 px-4 py-2 rounded-lg bg-blue-600 text-white text-sm font-medium hover:bg-blue-700 disabled:opacity-50"
            >
              {isPending ? '保存中…' : mode === 'create' ? '追加' : '保存'}
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

// ---- Page ----

type DialogState =
  | { open: false }
  | { open: true; mode: 'create' }
  | { open: true; mode: 'edit'; tx: Transaction }

function TransactionsPage() {
  const now = new Date()
  const [year, setYear] = useState(now.getFullYear())
  const [month, setMonth] = useState(now.getMonth() + 1)
  const [page, setPage] = useState(0)
  const [dialog, setDialog] = useState<DialogState>({ open: false })
  const pageSize = 50

  const { data, isLoading } = useTransactions({ year, month, page, page_size: pageSize })
  const { data: categories = [] } = useCategories()
  const updateCategory = useUpdateTransactionCategory()
  const deleteTransaction = useDeleteTransaction()

  const totalPages = data ? Math.ceil(data.total / pageSize) : 0
  const isCurrentMonth = year === now.getFullYear() && month === now.getMonth() + 1

  function prevMonth() {
    if (month === 1) {
      setYear((y) => y - 1)
      setMonth(12)
    } else setMonth((m) => m - 1)
    setPage(0)
  }
  function nextMonth() {
    if (month === 12) {
      setYear((y) => y + 1)
      setMonth(1)
    } else setMonth((m) => m + 1)
    setPage(0)
  }

  function handleDelete(tx: Transaction) {
    if (confirm(`「${tx.description}」を削除しますか？`)) {
      deleteTransaction.mutate(tx.id)
    }
  }

  return (
    <div className="p-4 md:p-8 space-y-5">
      {dialog.open && (
        <TxDialog
          mode={dialog.mode}
          tx={dialog.mode === 'edit' ? dialog.tx : undefined}
          categories={categories}
          onClose={() => setDialog({ open: false })}
        />
      )}

      <div className="flex items-center justify-between gap-2">
        <div className="flex items-center gap-2">
          <button onClick={prevMonth} className="p-1 rounded hover:bg-zinc-100 text-zinc-500">
            <ChevronLeft size={20} />
          </button>
          <h1 className="text-xl font-bold text-zinc-800 whitespace-nowrap">
            {year}年{month}月 の取引
          </h1>
          <button
            onClick={nextMonth}
            disabled={isCurrentMonth}
            className="p-1 rounded hover:bg-zinc-100 text-zinc-500 disabled:opacity-30 disabled:cursor-not-allowed"
          >
            <ChevronRight size={20} />
          </button>
        </div>
        <div className="flex items-center gap-3">
          <span className="text-sm text-zinc-400 shrink-0">{data?.total ?? 0} 件</span>
          <button
            onClick={() => setDialog({ open: true, mode: 'create' })}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-blue-600 text-white text-sm font-medium hover:bg-blue-700"
          >
            <Plus size={14} />
            追加
          </button>
        </div>
      </div>

      {/* モバイル: カードリスト */}
      <div className="md:hidden bg-white rounded-lg border border-zinc-200 divide-y divide-zinc-100">
        {isLoading ? (
          <div className="p-8 text-center text-zinc-400 text-sm">読み込み中…</div>
        ) : (data?.transactions ?? []).length === 0 ? (
          <div className="p-8 text-center text-zinc-400 text-sm">データがありません</div>
        ) : (
          (data?.transactions ?? []).map((tx) => (
            <TxCard
              key={tx.id}
              tx={tx}
              categories={categories}
              onCategoryChange={(catId) =>
                updateCategory.mutate({ id: tx.id, categoryId: catId })
              }
              onEdit={() => setDialog({ open: true, mode: 'edit', tx })}
              onDelete={() => handleDelete(tx)}
            />
          ))
        )}
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
                <th className="px-4 py-3 w-16"></th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-100">
              {(data?.transactions ?? []).map((tx) => (
                <TxRow
                  key={tx.id}
                  tx={tx}
                  categories={categories}
                  onCategoryChange={(catId) =>
                    updateCategory.mutate({ id: tx.id, categoryId: catId })
                  }
                  onEdit={() => setDialog({ open: true, mode: 'edit', tx })}
                  onDelete={() => handleDelete(tx)}
                />
              ))}
              {(data?.transactions ?? []).length === 0 && (
                <tr>
                  <td colSpan={5} className="px-4 py-8 text-center text-zinc-400">
                    データがありません
                  </td>
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
            onClick={() => setPage((p) => Math.max(0, p - 1))}
            className="px-3 py-1.5 rounded border border-zinc-200 text-sm disabled:opacity-40 hover:bg-zinc-50"
          >
            ← 前
          </button>
          <span className="text-sm text-zinc-600">
            {page + 1} / {totalPages}
          </span>
          <button
            disabled={page >= totalPages - 1}
            onClick={() => setPage((p) => p + 1)}
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
  onEdit,
  onDelete,
}: {
  tx: Transaction
  categories: Category[]
  onCategoryChange: (id: string | null) => void
  onEdit: () => void
  onDelete: () => void
}) {
  return (
    <div className="px-4 py-3 space-y-1.5">
      <div className="flex items-center justify-between gap-2">
        <span className="text-xs text-zinc-400">{tx.date.slice(0, 10)}</span>
        <div className="flex items-center gap-2">
          <span
            className={`text-sm font-medium tabular-nums ${tx.amount >= 0 ? 'text-emerald-600' : 'text-rose-600'}`}
          >
            {tx.amount >= 0 ? '+' : ''}
            {jpy(tx.amount)}
          </span>
          <button onClick={onEdit} className="text-zinc-300 hover:text-blue-500">
            <Pencil size={14} />
          </button>
          <button onClick={onDelete} className="text-zinc-300 hover:text-rose-500">
            <Trash2 size={14} />
          </button>
        </div>
      </div>
      <p className="text-sm text-zinc-800 truncate">{tx.description}</p>
      <select
        className="text-xs border border-zinc-200 rounded px-1.5 py-1 bg-white w-full"
        value={tx.category_id ?? ''}
        onChange={(e) => onCategoryChange(e.target.value || null)}
      >
        <option value="">未分類</option>
        {categories.map((cat) => (
          <option key={cat.id} value={cat.id}>
            {cat.name}
          </option>
        ))}
      </select>
    </div>
  )
}

function TxRow({
  tx,
  categories,
  onCategoryChange,
  onEdit,
  onDelete,
}: {
  tx: Transaction
  categories: Category[]
  onCategoryChange: (id: string | null) => void
  onEdit: () => void
  onDelete: () => void
}) {
  return (
    <tr className="hover:bg-zinc-50">
      <td className="px-4 py-2.5 text-zinc-500 text-xs whitespace-nowrap">
        {tx.date.slice(0, 10)}
      </td>
      <td className="px-4 py-2.5 text-zinc-800 max-w-xs truncate">{tx.description}</td>
      <td
        className={`px-4 py-2.5 text-right font-medium tabular-nums ${tx.amount >= 0 ? 'text-emerald-600' : 'text-rose-600'}`}
      >
        {tx.amount >= 0 ? '+' : ''}
        {jpy(tx.amount)}
      </td>
      <td className="px-4 py-2.5">
        <select
          className="text-xs border border-zinc-200 rounded px-1.5 py-1 bg-white max-w-35"
          value={tx.category_id ?? ''}
          onChange={(e) => onCategoryChange(e.target.value || null)}
        >
          <option value="">未分類</option>
          {categories.map((cat) => (
            <option key={cat.id} value={cat.id}>
              {cat.name}
            </option>
          ))}
        </select>
      </td>
      <td className="px-4 py-2.5 text-center">
        <div className="flex items-center justify-center gap-2">
          <button onClick={onEdit} className="text-zinc-300 hover:text-blue-500">
            <Pencil size={14} />
          </button>
          <button onClick={onDelete} className="text-zinc-300 hover:text-rose-500">
            <Trash2 size={14} />
          </button>
        </div>
      </td>
    </tr>
  )
}
