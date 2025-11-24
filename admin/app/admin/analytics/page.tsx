'use client'

import { useEffect, useState } from 'react'
import { api } from '@/lib/api'
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts'

interface AnalyticsData {
  users: { date: string; count: number }[]
  revenue: { date: string; amount: number }[]
  gifts: { date: string; count: number }[]
}

export default function AnalyticsPage() {
  const [data, setData] = useState<AnalyticsData>({
    users: [],
    revenue: [],
    gifts: [],
  })
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadAnalytics()
  }, [])

  const loadAnalytics = async () => {
    try {
      const response = await api.get('/admin/analytics')
      setData(response.data)
    } catch (error) {
      console.error('Failed to load analytics:', error)
      // Mock data
      const dates = Array.from({ length: 7 }, (_, i) => {
        const date = new Date()
        date.setDate(date.getDate() - (6 - i))
        return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
      })
      setData({
        users: dates.map((date) => ({ date, count: Math.floor(Math.random() * 100) + 50 })),
        revenue: dates.map((date) => ({ date, amount: Math.floor(Math.random() * 5000) + 2000 })),
        gifts: dates.map((date) => ({ date, count: Math.floor(Math.random() * 200) + 50 })),
      })
    } finally {
      setLoading(false)
    }
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
      <h1 className="text-3xl font-bold mb-8 text-primary">Analytics</h1>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Users Chart */}
        <div className="card bg-base-200 border border-base-300">
          <div className="card-body">
            <h2 className="card-title text-primary">New Users (7 Days)</h2>
            <ResponsiveContainer width="100%" height={300}>
              <LineChart data={data.users}>
                <CartesianGrid strokeDasharray="3 3" stroke="#333" />
                <XAxis dataKey="date" stroke="#999" />
                <YAxis stroke="#999" />
                <Tooltip contentStyle={{ backgroundColor: '#1E1E1E', border: '1px solid #333' }} />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="count"
                  stroke="#A7FF83"
                  strokeWidth={2}
                  name="Users"
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Revenue Chart */}
        <div className="card bg-base-200 border border-base-300">
          <div className="card-body">
            <h2 className="card-title text-primary">Revenue (7 Days)</h2>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={data.revenue}>
                <CartesianGrid strokeDasharray="3 3" stroke="#333" />
                <XAxis dataKey="date" stroke="#999" />
                <YAxis stroke="#999" />
                <Tooltip contentStyle={{ backgroundColor: '#1E1E1E', border: '1px solid #333' }} />
                <Legend />
                <Bar dataKey="amount" fill="#A7FF83" name="ETB" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Gifts Chart */}
        <div className="card bg-base-200 border border-base-300 lg:col-span-2">
          <div className="card-body">
            <h2 className="card-title text-primary">Gifts Sent (7 Days)</h2>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={data.gifts}>
                <CartesianGrid strokeDasharray="3 3" stroke="#333" />
                <XAxis dataKey="date" stroke="#999" />
                <YAxis stroke="#999" />
                <Tooltip contentStyle={{ backgroundColor: '#1E1E1E', border: '1px solid #333' }} />
                <Legend />
                <Bar dataKey="count" fill="#A7FF83" name="Gifts" />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>
    </div>
  )
}

