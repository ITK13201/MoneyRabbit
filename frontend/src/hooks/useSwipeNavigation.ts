import { useEffect, useRef } from 'react'
import { useLocation } from '@tanstack/react-router'
import { usePeriodStore } from '@/stores/periodStore'

const SWIPE_THRESHOLD = 50

function isInsideHorizontalScroll(el: HTMLElement | null): boolean {
  while (el) {
    const { overflowX } = window.getComputedStyle(el)
    if ((overflowX === 'scroll' || overflowX === 'auto') && el.scrollWidth > el.clientWidth) {
      return true
    }
    el = el.parentElement
  }
  return false
}

export function useSwipeNavigation() {
  const { pathname } = useLocation()
  const { prevMonth, nextMonth, prevTrendsYear, nextTrendsYear } = usePeriodStore()
  const startX = useRef(0)
  const startY = useRef(0)
  const inScrollable = useRef(false)

  useEffect(() => {
    const isMonthPage = pathname === '/' || pathname === '/transactions'
    const isYearPage = pathname === '/trends'
    if (!isMonthPage && !isYearPage) return

    const goPrev = isMonthPage ? prevMonth : prevTrendsYear
    const goNext = isMonthPage ? nextMonth : nextTrendsYear

    const onTouchStart = (e: TouchEvent) => {
      startX.current = e.touches[0].clientX
      startY.current = e.touches[0].clientY
      inScrollable.current = isInsideHorizontalScroll(e.target as HTMLElement)
    }

    const onTouchEnd = (e: TouchEvent) => {
      if (inScrollable.current) return
      const dx = e.changedTouches[0].clientX - startX.current
      const dy = e.changedTouches[0].clientY - startY.current
      if (Math.abs(dx) < SWIPE_THRESHOLD || Math.abs(dx) < Math.abs(dy)) return
      dx < 0 ? goNext() : goPrev()
    }

    const onKeyDown = (e: KeyboardEvent) => {
      const t = e.target as HTMLElement
      if (t.tagName === 'INPUT' || t.tagName === 'TEXTAREA' || t.isContentEditable) return
      if (e.key === 'ArrowRight') goNext()
      if (e.key === 'ArrowLeft') goPrev()
    }

    document.addEventListener('touchstart', onTouchStart, { passive: true })
    document.addEventListener('touchend', onTouchEnd, { passive: true })
    document.addEventListener('keydown', onKeyDown)

    return () => {
      document.removeEventListener('touchstart', onTouchStart)
      document.removeEventListener('touchend', onTouchEnd)
      document.removeEventListener('keydown', onKeyDown)
    }
  }, [pathname, prevMonth, nextMonth, prevTrendsYear, nextTrendsYear])
}
