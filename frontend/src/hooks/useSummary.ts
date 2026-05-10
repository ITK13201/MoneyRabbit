import { useQuery } from '@tanstack/react-query'
import { api } from '@/lib/api'

export const summaryKeys = {
  all: ['summary'] as const,
  monthly: (year: number) => [...summaryKeys.all, 'monthly', year] as const,
}

export function useMonthlySummary(year: number) {
  return useQuery({
    queryKey: summaryKeys.monthly(year),
    queryFn: () => api.summary.monthly(year),
  })
}
