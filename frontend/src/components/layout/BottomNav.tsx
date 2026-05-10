import { Link } from '@tanstack/react-router'
import { LayoutDashboard, Upload, List, Tag } from 'lucide-react'
import { cn } from '@/lib/utils'

const navItems = [
  { to: '/', label: 'ダッシュボード', icon: LayoutDashboard },
  { to: '/transactions', label: '取引', icon: List },
  { to: '/import', label: 'インポート', icon: Upload },
  { to: '/categories', label: 'カテゴリ', icon: Tag },
] as const

export function BottomNav() {
  return (
    <nav className="fixed bottom-0 inset-x-0 bg-zinc-900 border-t border-zinc-700 flex md:hidden z-50">
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
    </nav>
  )
}
