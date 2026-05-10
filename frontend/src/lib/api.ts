import type {
  Category,
  CategoryRule,
  CreateCategoryInput,
  CreateRuleInput,
  ImportFormat,
  ImportResult,
  MonthlySummaryResult,
  PreviewResult,
  PreviewRow,
  Transaction,
  TransactionListResult,
} from '@/types'

const BASE = '/api'

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, {
    headers: {
      'Content-Type': 'application/json',
      ...init?.headers,
    },
    ...init,
  })
  if (!res.ok) {
    const text = await res.text()
    throw new Error(`${res.status}: ${text}`)
  }
  if (res.status === 204) return undefined as T
  return res.json() as Promise<T>
}

export interface TransactionListParams {
  year?: number
  month?: number
  category_id?: string
  page?: number
  page_size?: number
}

export const api = {
  categories: {
    list: () => request<Category[]>('/categories'),
    create: (body: CreateCategoryInput) =>
      request<Category>('/categories', {
        method: 'POST',
        body: JSON.stringify(body),
      }),
    update: (id: string, body: Partial<CreateCategoryInput>) =>
      request<Category>(`/categories/${id}`, {
        method: 'PUT',
        body: JSON.stringify(body),
      }),
    delete: (id: string) => request<void>(`/categories/${id}`, { method: 'DELETE' }),
    createRule: (body: CreateRuleInput) =>
      request<CategoryRule>('/category-rules', {
        method: 'POST',
        body: JSON.stringify(body),
      }),
    updateRule: (id: string, body: Partial<CreateRuleInput>) =>
      request<CategoryRule>(`/category-rules/${id}`, {
        method: 'PUT',
        body: JSON.stringify(body),
      }),
    deleteRule: (id: string) => request<void>(`/category-rules/${id}`, { method: 'DELETE' }),
  },

  importFormats: {
    list: () => request<ImportFormat[]>('/import-formats'),
  },

  import: {
    preview: async (formData: FormData): Promise<PreviewResult> => {
      const res = await fetch(`${BASE}/import/preview`, {
        method: 'POST',
        body: formData,
      })
      if (!res.ok) throw new Error(`${res.status}: ${await res.text()}`)
      return res.json()
    },
    confirm: (body: { import_format_id: string; transactions: PreviewRow[] }) =>
      request<ImportResult>('/import/confirm', {
        method: 'POST',
        body: JSON.stringify(body),
      }),
  },

  summary: {
    monthly: (year: number) =>
      request<MonthlySummaryResult>(`/summary/monthly?year=${year}`),
  },

  transactions: {
    list: (params: TransactionListParams = {}) => {
      const q = new URLSearchParams()
      if (params.year != null) q.set('year', String(params.year))
      if (params.month != null) q.set('month', String(params.month))
      if (params.category_id) q.set('category_id', params.category_id)
      if (params.page != null) q.set('page', String(params.page))
      if (params.page_size != null) q.set('page_size', String(params.page_size))
      const qs = q.toString()
      return request<TransactionListResult>(`/transactions${qs ? '?' + qs : ''}`)
    },
    updateCategory: (id: string, categoryId: string | null) =>
      request<Transaction>(`/transactions/${id}/category`, {
        method: 'PATCH',
        body: JSON.stringify({ category_id: categoryId }),
      }),
    delete: (id: string) => request<void>(`/transactions/${id}`, { method: 'DELETE' }),
  },
}

