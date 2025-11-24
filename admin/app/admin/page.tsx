'use client'

import { useEffect, useState } from 'react'
import { api } from '@/lib/api'
import Link from 'next/link'

interface Stats {
  dau: number
  totalUsers: number
  pendingPhotos: number
  pendingCashouts: number
  todayRevenue: number
}

export default function DashboardPage() {
  const [stats, setStats] = useState<Stats>({
    dau: 0,
    totalUsers: 0,
    pendingPhotos: 0,
    pendingCashouts: 0,
    todayRevenue: 0,
  })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadStats()
  }, [])

  const loadStats = async () => {
    try {
      // TODO: Replace with actual admin stats endpoint
      // For now, using mock data structure
      const response = await api.get('/admin/stats')
      setStats(response.data)
    } catch (error) {
      console.error('Failed to load stats:', error)
      // Mock data for development
      setStats({
        dau: 1247,
        totalUsers: 15234,
        pendingPhotos: 23,
        pendingCashouts: 8,
        todayRevenue: 12500,
      })
    } finally {
      setLoading(false)
    }
  }

  const StatCard = ({
    title,
    value,
    icon,
    link,
    color = 'primary',
  }: {
    title: string
    value: number | string
    icon: string
    link?: string
    color?: string
  }) => {
    const content = (
      <div className={`stat bg-base-200 rounded-lg border border-base-300`}>
        <div className="stat-figure text-primary text-4xl">{icon}</div>
        <div className="stat-title text-base-content/70">{title}</div>
        <div className={`stat-value text-${color} text-3xl`}>
          {typeof value === 'number' ? value.toLocaleString() : value}
        </div>
      </div>
    )

    if (link) {
      return <Link href={link}>{content}</Link>
    }
    return content
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <span className="loading loading-spinner loading-lg text-primary"></span>
      </div>
    )
  }

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8 text-primary">Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        <StatCard
          title="Daily Active Users"
          value={stats.dau}
          icon="👥"
          color="primary"
        />
        <StatCard
          title="Total Users"
          value={stats.totalUsers}
          icon="📊"
          link="/admin/users"
        />
        <StatCard
          title="Pending Photos"
          value={stats.pendingPhotos}
          icon="📸"
          link="/admin/photos"
          color="warning"
        />
        <StatCard
          title="Pending Cashouts"
          value={stats.pendingCashouts}
          icon="💰"
          link="/admin/cashouts"
          color="warning"
        />
        <StatCard
          title="Today's Revenue"
          value={`${stats.todayRevenue.toLocaleString()} ETB`}
          icon="💵"
          color="success"
        />
      </div>

      {/* Quick Actions */}
      <div className="card bg-base-200 border border-base-300">
        <div className="card-body">
          <h2 className="card-title text-primary">Quick Actions</h2>
          <div className="flex flex-wrap gap-4 mt-4">
            <Link href="/admin/photos" className="btn btn-primary">
              Moderate Photos
            </Link>
            <Link href="/admin/cashouts" className="btn btn-warning">
              Review Cashouts
            </Link>
            <Link href="/admin/reports" className="btn btn-error">
              View Reports
            </Link>
            <Link href="/admin/broadcast" className="btn btn-info">
              Send Broadcast
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}

