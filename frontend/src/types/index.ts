export type CategoryType = 'income' | 'expense' | 'both'

export interface CategoryRule {
  id: string
  keyword: string
  priority: number
  category_id: string
}

export interface Category {
  id: string
  name: string
  color: string
  icon: string
  type: CategoryType
  sort_order: number
  rules: CategoryRule[]
}

export interface Transaction {
  id: string
  date: string
  description: string
  amount: number
  import_format_id: string
  imported_at: string
  category_id: string | null
  category: Category | null
}

export interface ImportFormat {
  id: string
  name: string
  import_type: string
}

export interface PreviewRow {
  date: string
  description: string
  amount: number
}

export interface PreviewResult {
  rows: PreviewRow[]
  skipped_rows: number
}

export interface ImportResult {
  imported: number
  skipped: number
  errors: string[]
}

export interface TransactionListResult {
  transactions: Transaction[]
  total: number
}

export interface CreateCategoryInput {
  name: string
  color: string
  icon: string
  type: CategoryType
  sort_order: number
}

export interface CreateRuleInput {
  category_id: string
  keyword: string
  priority: number
}
