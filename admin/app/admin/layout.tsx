'use client'

import { useAuth } from '../providers/AuthProvider'
import Sidebar from '@/components/Sidebar'
import { usePathname } from 'next/navigation'

export default function AdminLayout({
  children,
}: {
  children: React.ReactNode
}) {
  const { user, loading } = useAuth()
  const pathname = usePathname()

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen bg-black">
        <span className="loading loading-spinner loading-lg text-primary"></span>
      </div>
    )
  }

  if (!user || user.role !== 'admin') {
    return null
  }

  return (
    <div className="flex min-h-screen bg-black">
      <Sidebar currentPath={pathname || ''} />
      <main className="flex-1 ml-0 lg:ml-64 p-4 lg:p-8">
        {children}
      </main>
    </div>
  )
}

