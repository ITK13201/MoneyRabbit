import { Link } from '@tanstack/react-router'
import { LayoutDashboard, Upload, List, Tag, Rabbit, TrendingUp } from 'lucide-react'
import { cn } from '@/lib/utils'

const navItems = [
  { to: '/', label: 'ダッシュボード', icon: LayoutDashboard },
  { to: '/transactions', label: '取引一覧', icon: List },
  { to: '/trends', label: 'トレンド', icon: TrendingUp },
  { to: '/import', label: 'インポート', icon: Upload },
  { to: '/categories', label: 'カテゴリ', icon: Tag },
] as const

export function Sidebar() {
  return (
    <aside className="hidden md:flex w-56 shrink-0 bg-zinc-900 text-zinc-100 flex-col h-screen sticky top-0">
      <div className="flex items-center gap-2 px-4 py-5 border-b border-zinc-700">
        <Rabbit className="text-emerald-400" size={22} />
        <span className="font-bold text-lg tracking-tight">MoneyRabbit</span>
      </div>
      <nav className="flex-1 py-3">
        {navItems.map(({ to, label, icon: Icon }) => (
          <Link
            key={to}
            to={to}
            className={cn(
              'flex items-center gap-3 px-4 py-2.5 text-sm transition-colors',
              'hover:bg-zinc-800 hover:text-white',
            )}
            activeProps={{ className: 'bg-zinc-800 text-white font-medium' }}
            activeOptions={{ exact: to === '/' }}
          >
            <Icon size={17} />
            {label}
          </Link>
        ))}
      </nav>
    </aside>
  )
}
