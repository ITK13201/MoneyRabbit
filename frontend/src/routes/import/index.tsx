import { createFileRoute } from '@tanstack/react-router'
import { useState, useRef } from 'react'
import { useQuery } from '@tanstack/react-query'
import { api } from '@/lib/api'
import type { PreviewRow, ImportFormat } from '@/types'
import { Upload, Trash2, CheckCircle } from 'lucide-react'

export const Route = createFileRoute('/import/')({
  component: ImportPage,
})

function jpy(amount: number) {
  return new Intl.NumberFormat('ja-JP', { style: 'currency', currency: 'JPY' }).format(Math.abs(amount))
}

function ImportPage() {
  const { data: formats = [] } = useQuery({
    queryKey: ['import-formats'],
    queryFn: () => api.importFormats.list(),
  })

  const [formatId, setFormatId] = useState('')
  const [preview, setPreview] = useState<PreviewRow[] | null>(null)
  const [skipped, setSkipped] = useState(0)
  const [isPreviewing, setIsPreviewing] = useState(false)
  const [isConfirming, setIsConfirming] = useState(false)
  const [result, setResult] = useState<{ imported: number; skipped: number } | null>(null)
  const [error, setError] = useState('')
  const fileRef = useRef<HTMLInputElement>(null)

  // selected format metadata
  const selectedFormat: ImportFormat | undefined = formats.find(f => f.id === formatId)

  async function handlePreview() {
    const file = fileRef.current?.files?.[0]
    if (!file || !formatId) return
    setError('')
    setResult(null)
    setIsPreviewing(true)
    try {
      const fd = new FormData()
      fd.append('file', file)
      fd.append('import_format_id', formatId)
      const res = await api.import.preview(fd)
      setPreview(res.rows)
      setSkipped(res.skipped_rows)
    } catch (e) {
      setError(String(e))
    } finally {
      setIsPreviewing(false)
    }
  }

  function removeRow(index: number) {
    setPreview(prev => prev?.filter((_, i) => i !== index) ?? null)
  }

  async function handleConfirm() {
    if (!preview || !formatId) return
    setIsConfirming(true)
    setError('')
    try {
      const res = await api.import.confirm({
        import_format_id: formatId,
        transactions: preview,
      })
      setResult({ imported: res.imported, skipped: res.skipped })
      setPreview(null)
      if (fileRef.current) fileRef.current.value = ''
    } catch (e) {
      setError(String(e))
    } finally {
      setIsConfirming(false)
    }
  }

  return (
    <div className="p-4 md:p-8 space-y-6 max-w-4xl">
      <h1 className="text-xl font-bold text-zinc-800">CSVインポート</h1>

      {result && (
        <div className="flex items-center gap-3 bg-emerald-50 border border-emerald-200 rounded-lg px-5 py-4 text-emerald-700">
          <CheckCircle size={18} />
          <span>{result.imported} 件をインポートしました（重複スキップ: {result.skipped} 件）</span>
        </div>
      )}

      {error && (
        <div className="bg-rose-50 border border-rose-200 rounded-lg px-5 py-4 text-rose-700 text-sm">{error}</div>
      )}

      {/* Step 1: Select format + file */}
      <div className="bg-white rounded-lg border border-zinc-200 p-5 space-y-4">
        <h2 className="text-sm font-semibold text-zinc-700">1. フォーマットとファイルを選択</h2>
        <div className="flex flex-col sm:flex-row gap-3">
          <select
            className="border border-zinc-200 rounded px-3 py-2 text-sm bg-white flex-1"
            value={formatId}
            onChange={e => { setFormatId(e.target.value); setPreview(null) }}
          >
            <option value="">フォーマットを選択…</option>
            {formats.map(f => (
              <option key={f.id} value={f.id}>{f.name}</option>
            ))}
          </select>
          <label className="flex items-center gap-2 px-4 py-2 rounded border border-zinc-200 text-sm cursor-pointer hover:bg-zinc-50">
            <Upload size={15} />
            CSVを選択
            <input ref={fileRef} type="file" accept=".csv" className="hidden" />
          </label>
          <button
            disabled={!formatId || isPreviewing}
            onClick={handlePreview}
            className="px-4 py-2 bg-zinc-800 text-white text-sm rounded disabled:opacity-40 hover:bg-zinc-700"
          >
            {isPreviewing ? '解析中…' : 'プレビュー'}
          </button>
        </div>
        {selectedFormat && (
          <p className="text-xs text-zinc-400">種別: {selectedFormat.import_type === 'bank_account' ? '銀行口座' : 'クレジットカード'}</p>
        )}
      </div>

      {/* Step 2: Preview */}
      {preview !== null && (
        <div className="bg-white rounded-lg border border-zinc-200 space-y-0 overflow-hidden">
          <div className="flex items-center justify-between px-5 py-3 border-b border-zinc-100">
            <h2 className="text-sm font-semibold text-zinc-700">
              2. プレビュー ({preview.length} 件{skipped > 0 ? `、スキップ ${skipped} 行` : ''})
            </h2>
            <button
              disabled={preview.length === 0 || isConfirming}
              onClick={handleConfirm}
              className="px-4 py-1.5 bg-emerald-600 text-white text-sm rounded disabled:opacity-40 hover:bg-emerald-700"
            >
              {isConfirming ? '保存中…' : '確定してインポート'}
            </button>
          </div>
          <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-zinc-50 text-zinc-600 text-xs uppercase">
              <tr>
                <th className="px-4 py-2 text-left">日付</th>
                <th className="px-4 py-2 text-left">摘要</th>
                <th className="px-4 py-2 text-right">金額</th>
                <th className="px-4 py-2 w-12"></th>
              </tr>
            </thead>
            <tbody className="divide-y divide-zinc-100">
              {preview.map((row, i) => (
                <tr key={i} className="hover:bg-zinc-50">
                  <td className="px-4 py-2 text-zinc-500 whitespace-nowrap">{row.date}</td>
                  <td className="px-4 py-2 text-zinc-800 truncate max-w-xs">{row.description}</td>
                  <td className={`px-4 py-2 text-right font-medium tabular-nums ${row.amount >= 0 ? 'text-emerald-600' : 'text-rose-600'}`}>
                    {row.amount >= 0 ? '+' : ''}{jpy(row.amount)}
                  </td>
                  <td className="px-4 py-2 text-center">
                    <button onClick={() => removeRow(i)} className="text-zinc-300 hover:text-rose-500">
                      <Trash2 size={14} />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          </div>
        </div>
      )}
    </div>
  )
}
