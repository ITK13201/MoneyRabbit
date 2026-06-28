import { Link } from '@tanstack/react-router'
import { LayoutDashboard, Upload, List, Tag, TrendingUp } from 'lucide-react'
import { cn } from '@/lib/utils'

const navItems = [
  { to: '/', label: 'ダッシュボード', icon: LayoutDashboard },
  { to: '/transactions', label: '取引', icon: List },
  { to: '/trends', label: 'トレンド', icon: TrendingUp },
  { to: '/import', label: 'インポート', icon: Upload },
  { to: '/categories', label: 'カテゴリ', icon: Tag },
] as const

export function BottomNav() {
  return (
    <nav className="fixed bottom-0 inset-x-0 bg-zinc-900 border-t border-zinc-700 md:hidden z-50 pb-[env(safe-area-inset-bottom)]">
      <div className="relative flex">
        <span className="absolute right-1.5 top-1 text-[9px] text-zinc-600 pointer-events-none select-none">
          v{__APP_VERSION__}
        </span>
      {navItems.map(({ to, label, icon: Icon }) => (
        <Link
          key={to}
          to={to}
          className={cn(
            'flex-1 flex flex-col items-center justify-center gap-0.5 py-2 text-zinc-400 text-[10px] transition-colors',
            'hover:text-white',
          )}
          activeProps={{ className: 'text-emerald-400 font-medium' }}
          activeOptions={{ exact: to === '/' }}
        >
          <Icon size={20} />
          {label}
        </Link>
      ))}
      </div>
    </nav>
  )
}
