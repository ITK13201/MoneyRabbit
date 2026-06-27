import { create } from 'zustand'

const now = new Date()
const THIS_YEAR = now.getFullYear()
const THIS_MONTH = now.getMonth() + 1

interface PeriodState {
  // ダッシュボード・取引画面で共有
  year: number
  month: number
  // トレンド画面
  trendsYear: number

  prevMonth: () => void
  nextMonth: () => void
  prevTrendsYear: () => void
  nextTrendsYear: () => void
}

export const usePeriodStore = create<PeriodState>((set) => ({
  year: THIS_YEAR,
  month: THIS_MONTH,
  trendsYear: THIS_YEAR,

  prevMonth: () =>
    set((s) =>
      s.month === 1 ? { year: s.year - 1, month: 12 } : { month: s.month - 1 },
    ),

  nextMonth: () =>
    set((s) => {
      if (s.year === THIS_YEAR && s.month === THIS_MONTH) return s
      return s.month === 12 ? { year: s.year + 1, month: 1 } : { month: s.month + 1 }
    }),

  prevTrendsYear: () => set((s) => ({ trendsYear: s.trendsYear - 1 })),

  nextTrendsYear: () =>
    set((s) => (s.trendsYear >= THIS_YEAR ? s : { trendsYear: s.trendsYear + 1 })),
}))
