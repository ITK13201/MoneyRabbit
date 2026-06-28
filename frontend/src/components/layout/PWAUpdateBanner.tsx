import { useState, useEffect, useRef } from 'react'
import { RefreshCw } from 'lucide-react'

export function PWAUpdateBanner() {
  const [needRefresh, setNeedRefresh] = useState(false)
  const waitingWorkerRef = useRef<ServiceWorker | null>(null)

  useEffect(() => {
    if (!('serviceWorker' in navigator)) return

    const onUpdateFound = (registration: ServiceWorkerRegistration) => {
      const newWorker = registration.installing
      if (!newWorker) return

      newWorker.addEventListener('statechange', () => {
        if (newWorker.state === 'installed' && navigator.serviceWorker.controller) {
          waitingWorkerRef.current = newWorker
          setNeedRefresh(true)
        }
      })
    }

    navigator.serviceWorker.ready.then((registration) => {
      if (registration.waiting && navigator.serviceWorker.controller) {
        waitingWorkerRef.current = registration.waiting
        setNeedRefresh(true)
      }
      registration.addEventListener('updatefound', () => onUpdateFound(registration))
    })

    let refreshing = false
    navigator.serviceWorker.addEventListener('controllerchange', () => {
      if (!refreshing) {
        refreshing = true
        window.location.reload()
      }
    })
  }, [])

  const handleUpdate = () => {
    waitingWorkerRef.current?.postMessage({ type: 'SKIP_WAITING' })
  }

  if (!needRefresh) return null

  return (
    <div className="fixed top-4 left-1/2 -translate-x-1/2 z-50 flex items-center gap-3 px-4 py-3 rounded-xl bg-zinc-800 text-zinc-100 shadow-lg border border-zinc-700 text-sm">
      <span>新しいバージョンが利用可能です</span>
      <button
        onClick={handleUpdate}
        className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-emerald-500 hover:bg-emerald-400 text-white text-xs font-medium transition-colors"
      >
        <RefreshCw size={13} />
        更新
      </button>
    </div>
  )
}
