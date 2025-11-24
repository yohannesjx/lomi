'use client'

import Link from 'next/link'
import { useAuth } from '../app/providers/AuthProvider'
import { useRouter } from 'next/navigation'

interface SidebarProps {
  currentPath: string
}

export default function Sidebar({ currentPath }: SidebarProps) {
  const { user, logout } = useAuth()
  const router = useRouter()

  const menuItems = [
    { path: '/admin', label: 'Dashboard', icon: '📊' },
    { path: '/admin/users', label: 'Users', icon: '👥' },
    { path: '/admin/photos', label: 'Photo Moderation', icon: '📸' },
    { path: '/admin/reports', label: 'Reports', icon: '🚨' },
    { path: '/admin/gifts', label: 'Gifts & Coins', icon: '🎁' },
    { path: '/admin/cashouts', label: 'Cashouts', icon: '💰' },
    { path: '/admin/broadcast', label: 'Broadcast', icon: '📢' },
    { path: '/admin/analytics', label: 'Analytics', icon: '📈' },
    { path: '/admin/settings', label: 'Settings', icon: '⚙️' },
  ]

  const handleLogout = () => {
    logout()
    router.push('/login')
  }

  return (
    <div className="fixed left-0 top-0 h-full w-64 bg-base-200 border-r border-base-300 z-50">
      <div className="flex flex-col h-full">
        {/* Header */}
        <div className="p-6 border-b border-base-300">
          <h1 className="text-2xl font-bold text-primary">Lomi Admin</h1>
          <p className="text-sm text-base-content/60 mt-1">
            {user?.name || 'Admin'}
          </p>
        </div>

        {/* Menu */}
        <nav className="flex-1 overflow-y-auto p-4">
          <ul className="space-y-2">
            {menuItems.map((item) => {
              const isActive = currentPath === item.path
              return (
                <li key={item.path}>
                  <Link
                    href={item.path}
                    className={`flex items-center gap-3 px-4 py-3 rounded-lg transition-colors ${
                      isActive
                        ? 'bg-primary text-primary-content font-semibold'
                        : 'hover:bg-base-300 text-base-content'
                    }`}
                  >
                    <span className="text-xl">{item.icon}</span>
                    <span>{item.label}</span>
                  </Link>
                </li>
              )
            })}
          </ul>
        </nav>

        {/* Footer */}
        <div className="p-4 border-t border-base-300">
          <button
            onClick={handleLogout}
            className="btn btn-error btn-sm w-full"
          >
            Logout
          </button>
        </div>
      </div>
    </div>
  )
}

