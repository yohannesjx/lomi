'use client'

import { useEffect, useState } from 'react'
import { api } from '@/lib/api'

interface Report {
  id: string
  reporter_name: string
  reported_user_name: string
  reason: string
  description?: string
  status: string
  created_at: string
}

export default function ReportsPage() {
  const [reports, setReports] = useState<Report[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadReports()
  }, [])

  const loadReports = async () => {
    try {
      const response = await api.get('/admin/reports/pending')
      setReports(response.data.reports || [])
    } catch (error) {
      console.error('Failed to load reports:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleAction = async (reportId: string, action: string) => {
    try {
      await api.put(`/admin/reports/${reportId}/review`, { action })
      await loadReports()
    } catch (error) {
      alert('Failed to process report')
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
      <h1 className="text-3xl font-bold mb-8 text-primary">Reports</h1>

      <div className="card bg-base-200 border border-base-300">
        <div className="card-body p-0">
          <div className="overflow-x-auto">
            <table className="table">
              <thead>
                <tr>
                  <th>Reporter</th>
                  <th>Reported User</th>
                  <th>Reason</th>
                  <th>Description</th>
                  <th>Date</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {reports.map((report) => (
                  <tr key={report.id}>
                    <td>{report.reporter_name}</td>
                    <td>{report.reported_user_name}</td>
                    <td>
                      <span className="badge badge-error">{report.reason}</span>
                    </td>
                    <td className="max-w-xs truncate">
                      {report.description || 'N/A'}
                    </td>
                    <td>{new Date(report.created_at).toLocaleDateString()}</td>
                    <td>
                      <div className="flex gap-2">
                        <button
                          className="btn btn-sm btn-info"
                          onClick={() =>
                            window.open(`/admin/reports/${report.id}`, '_blank')
                          }
                        >
                          View
                        </button>
                        <button
                          className="btn btn-sm btn-warning"
                          onClick={() => handleAction(report.id, 'warn')}
                        >
                          Warn
                        </button>
                        <button
                          className="btn btn-sm btn-error"
                          onClick={() => handleAction(report.id, 'ban')}
                        >
                          Ban
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  )
}

