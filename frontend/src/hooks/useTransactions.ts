import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { api, type TransactionListParams } from '@/lib/api'

export const txKeys = {
  all: ['transactions'] as const,
  list: (params: TransactionListParams) => [...txKeys.all, 'list', params] as const,
}

export function useTransactions(params: TransactionListParams = {}) {
  return useQuery({
    queryKey: txKeys.list(params),
    queryFn: () => api.transactions.list(params),
  })
}

export function useUpdateTransactionCategory() {
  const qc = useQueryClient()
  return useMutation({
    mutationFn: ({ id, categoryId }: { id: string; categoryId: string | null }) =>
      api.transactions.updateCategory(id, categoryId),
    onSuccess: () => qc.invalidateQueries({ queryKey: txKeys.all }),
  })
}
