import { Outlet } from '@tanstack/react-router'
import { Sidebar } from './Sidebar'
import { BottomNav } from './BottomNav'
import { useSwipeNavigation } from '@/hooks/useSwipeNavigation'

export function Layout() {
  useSwipeNavigation()

  return (
    <div className="flex min-h-screen bg-zinc-50">
      <Sidebar />
      <main className="flex-1 overflow-auto pb-[calc(4rem+env(safe-area-inset-bottom))] md:pb-0">
        <Outlet />
      </main>
      <BottomNav />
    </div>
  )
}
