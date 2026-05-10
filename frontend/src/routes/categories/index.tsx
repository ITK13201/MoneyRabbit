import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import {
  useCategories,
  useCreateCategory,
  useDeleteCategory,
  useCreateRule,
  useDeleteRule,
} from '@/hooks/useCategories'
import type { Category, CategoryType } from '@/types'
import { Plus, Trash2, ChevronDown, ChevronRight } from 'lucide-react'

export const Route = createFileRoute('/categories/')({
  component: CategoriesPage,
})

const TYPE_LABELS: Record<CategoryType, string> = {
  income: '収入',
  expense: '支出',
  both: '両方',
}

function CategoriesPage() {
  const { data: categories = [], isLoading } = useCategories()
  const createCategory = useCreateCategory()
  const deleteCategory = useDeleteCategory()
  const [expandedId, setExpandedId] = useState<string | null>(null)
  const [showForm, setShowForm] = useState(false)

  if (isLoading) return <div className="p-8 text-zinc-400 text-sm">読み込み中…</div>

  return (
    <div className="p-4 md:p-8 space-y-5 max-w-3xl">
      <div className="flex items-center justify-between">
        <h1 className="text-xl font-bold text-zinc-800">カテゴリ管理</h1>
        <button
          onClick={() => setShowForm(v => !v)}
          className="flex items-center gap-1.5 px-3 py-1.5 bg-zinc-800 text-white text-sm rounded hover:bg-zinc-700"
        >
          <Plus size={15} />
          カテゴリを追加
        </button>
      </div>

      {showForm && (
        <CategoryForm
          onSubmit={async (data) => {
            await createCategory.mutateAsync(data)
            setShowForm(false)
          }}
          onCancel={() => setShowForm(false)}
        />
      )}

      <div className="space-y-2">
        {categories.map(cat => (
          <CategoryRow
            key={cat.id}
            cat={cat}
            expanded={expandedId === cat.id}
            onToggle={() => setExpandedId(id => id === cat.id ? null : cat.id)}
            onDelete={() => {
              if (confirm(`「${cat.name}」を削除しますか？`)) {
                deleteCategory.mutate(cat.id)
              }
            }}
          />
        ))}
        {categories.length === 0 && (
          <p className="text-sm text-zinc-400 py-4">カテゴリがありません。追加してください。</p>
        )}
      </div>
    </div>
  )
}

function CategoryRow({
  cat,
  expanded,
  onToggle,
  onDelete,
}: {
  cat: Category
  expanded: boolean
  onToggle: () => void
  onDelete: () => void
}) {
  const createRule = useCreateRule()
  const deleteRule = useDeleteRule()
  const [showRuleForm, setShowRuleForm] = useState(false)

  return (
    <div className="bg-white rounded-lg border border-zinc-200">
      <div className="flex items-center gap-3 px-4 py-3">
        <button onClick={onToggle} className="text-zinc-400 hover:text-zinc-600">
          {expanded ? <ChevronDown size={16} /> : <ChevronRight size={16} />}
        </button>
        <span
          className="w-4 h-4 rounded-full shrink-0"
          style={{ backgroundColor: cat.color }}
        />
        <span className="text-sm">{cat.icon}</span>
        <span className="font-medium text-sm text-zinc-800 flex-1">{cat.name}</span>
        <span className="text-xs text-zinc-400 bg-zinc-100 rounded px-2 py-0.5">{TYPE_LABELS[cat.type]}</span>
        <button
          onClick={onDelete}
          className="text-zinc-300 hover:text-rose-500 ml-2"
          title="削除"
        >
          <Trash2 size={14} />
        </button>
      </div>

      {expanded && (
        <div className="border-t border-zinc-100 px-4 py-3 space-y-2">
          <div className="flex items-center justify-between">
            <p className="text-xs text-zinc-500 font-medium">自動分類ルール</p>
            <button
              onClick={() => setShowRuleForm(v => !v)}
              className="text-xs text-zinc-500 hover:text-zinc-700 flex items-center gap-1"
            >
              <Plus size={12} />
              ルール追加
            </button>
          </div>

          {showRuleForm && (
            <RuleForm
              onSubmit={async ({ keyword, priority }) => {
                await createRule.mutateAsync({ category_id: cat.id, keyword, priority })
                setShowRuleForm(false)
              }}
              onCancel={() => setShowRuleForm(false)}
            />
          )}

          {(cat.rules ?? []).length === 0 ? (
            <p className="text-xs text-zinc-400">ルールなし</p>
          ) : (
            <ul className="space-y-1">
              {(cat.rules ?? []).map(rule => (
                <li key={rule.id} className="flex items-center gap-2 text-xs text-zinc-600">
                  <span className="font-mono bg-zinc-100 px-1.5 py-0.5 rounded">{rule.keyword}</span>
                  <span className="text-zinc-400">優先度: {rule.priority}</span>
                  <button
                    onClick={() => deleteRule.mutate(rule.id)}
                    className="text-zinc-300 hover:text-rose-500 ml-auto"
                  >
                    <Trash2 size={11} />
                  </button>
                </li>
              ))}
            </ul>
          )}
        </div>
      )}
    </div>
  )
}

function CategoryForm({
  onSubmit,
  onCancel,
}: {
  onSubmit: (data: import('@/types').CreateCategoryInput) => Promise<void>
  onCancel: () => void
}) {
  const [name, setName] = useState('')
  const [color, setColor] = useState('#6366f1')
  const [icon, setIcon] = useState('📦')
  const [type, setType] = useState<CategoryType>('expense')
  const [loading, setLoading] = useState(false)

  return (
    <form
      className="bg-zinc-50 border border-zinc-200 rounded-lg p-4 space-y-3"
      onSubmit={async e => {
        e.preventDefault()
        setLoading(true)
        await onSubmit({ name, color, icon, type, sort_order: 0 })
        setLoading(false)
      }}
    >
      <div className="grid grid-cols-2 gap-3">
        <div>
          <label className="text-xs text-zinc-600 block mb-1">名前</label>
          <input
            required
            value={name}
            onChange={e => setName(e.target.value)}
            className="w-full border border-zinc-200 rounded px-2.5 py-1.5 text-sm bg-white"
            placeholder="食費"
          />
        </div>
        <div>
          <label className="text-xs text-zinc-600 block mb-1">アイコン</label>
          <input
            value={icon}
            onChange={e => setIcon(e.target.value)}
            className="w-full border border-zinc-200 rounded px-2.5 py-1.5 text-sm bg-white"
            placeholder="🍔"
          />
        </div>
        <div>
          <label className="text-xs text-zinc-600 block mb-1">カラー</label>
          <input
            type="color"
            value={color}
            onChange={e => setColor(e.target.value)}
            className="h-9 w-full border border-zinc-200 rounded cursor-pointer"
          />
        </div>
        <div>
          <label className="text-xs text-zinc-600 block mb-1">種別</label>
          <select
            value={type}
            onChange={e => setType(e.target.value as CategoryType)}
            className="w-full border border-zinc-200 rounded px-2.5 py-1.5 text-sm bg-white"
          >
            <option value="expense">支出</option>
            <option value="income">収入</option>
            <option value="both">両方</option>
          </select>
        </div>
      </div>
      <div className="flex gap-2 justify-end">
        <button type="button" onClick={onCancel} className="px-3 py-1.5 text-sm text-zinc-600 border border-zinc-200 rounded hover:bg-zinc-100">
          キャンセル
        </button>
        <button type="submit" disabled={loading} className="px-3 py-1.5 text-sm bg-zinc-800 text-white rounded disabled:opacity-40 hover:bg-zinc-700">
          {loading ? '保存中…' : '保存'}
        </button>
      </div>
    </form>
  )
}

function RuleForm({
  onSubmit,
  onCancel,
}: {
  onSubmit: (data: { keyword: string; priority: number }) => Promise<void>
  onCancel: () => void
}) {
  const [keyword, setKeyword] = useState('')
  const [priority, setPriority] = useState(0)
  const [loading, setLoading] = useState(false)

  return (
    <form
      className="flex items-end gap-2"
      onSubmit={async e => {
        e.preventDefault()
        setLoading(true)
        await onSubmit({ keyword, priority })
        setLoading(false)
      }}
    >
      <div className="flex-1">
        <label className="text-xs text-zinc-500 block mb-1">キーワード</label>
        <input
          required
          value={keyword}
          onChange={e => setKeyword(e.target.value)}
          className="w-full border border-zinc-200 rounded px-2 py-1.5 text-xs bg-white"
          placeholder="スーパー"
        />
      </div>
      <div className="w-20">
        <label className="text-xs text-zinc-500 block mb-1">優先度</label>
        <input
          type="number"
          value={priority}
          onChange={e => setPriority(Number(e.target.value))}
          className="w-full border border-zinc-200 rounded px-2 py-1.5 text-xs bg-white"
        />
      </div>
      <button type="button" onClick={onCancel} className="px-2 py-1.5 text-xs text-zinc-500 border border-zinc-200 rounded">
        ×
      </button>
      <button type="submit" disabled={loading} className="px-2 py-1.5 text-xs bg-zinc-700 text-white rounded disabled:opacity-40">
        追加
      </button>
    </form>
  )
}
