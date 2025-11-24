'use client'

import { useEffect, useState } from 'react'
import { api } from '@/lib/api'

interface Cashout {
  id: string
  user_name: string
  coins: number
  etb_amount: number
  net_amount: number
  payment_method: string
  payment_account: string
  status: string
  created_at: string
}

export default function CashoutsPage() {
  const [cashouts, setCashouts] = useState<Cashout[]>([])
  const [loading, setLoading] = useState(true)
  const [processing, setProcessing] = useState<string | null>(null)

  useEffect(() => {
    loadCashouts()
  }, [])

  const loadCashouts = async () => {
    try {
      const response = await api.get('/admin/payouts/pending')
      setCashouts(response.data.payouts || [])
    } catch (error) {
      console.error('Failed to load cashouts:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleApprove = async (cashoutId: string) => {
    setProcessing(cashoutId)
    try {
      await api.put(`/admin/payouts/${cashoutId}/process`, {
        status: 'processing',
      })
      await loadCashouts()
    } catch (error) {
      alert('Failed to approve cashout')
    } finally {
      setProcessing(null)
    }
  }

  const handleMarkPaid = async (cashoutId: string) => {
    setProcessing(cashoutId)
    try {
      await api.put(`/admin/payouts/${cashoutId}/process`, {
        status: 'completed',
      })
      await loadCashouts()
    } catch (error) {
      alert('Failed to mark as paid')
    } finally {
      setProcessing(null)
    }
  }

  const handleReject = async (cashoutId: string) => {
    const reason = prompt('Rejection reason:')
    if (!reason) return

    setProcessing(cashoutId)
    try {
      await api.put(`/admin/payouts/${cashoutId}/process`, {
        status: 'rejected',
        rejection_reason: reason,
      })
      await loadCashouts()
    } catch (error) {
      alert('Failed to reject cashout')
    } finally {
      setProcessing(null)
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
      <h1 className="text-3xl font-bold mb-8 text-primary">Cashouts</h1>

      <div className="card bg-base-200 border border-base-300">
        <div className="card-body p-0">
          <div className="overflow-x-auto">
            <table className="table">
              <thead>
                <tr>
                  <th>User</th>
                  <th>Coins</th>
                  <th>Amount</th>
                  <th>Net Amount</th>
                  <th>Payment Method</th>
                  <th>Account</th>
                  <th>Status</th>
                  <th>Date</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {cashouts.map((cashout) => (
                  <tr key={cashout.id}>
                    <td>{cashout.user_name}</td>
                    <td>{cashout.coins.toLocaleString()} LC</td>
                    <td>{cashout.etb_amount.toFixed(2)} ETB</td>
                    <td className="font-semibold text-primary">
                      {cashout.net_amount.toFixed(2)} ETB
                    </td>
                    <td>{cashout.payment_method}</td>
                    <td>{cashout.payment_account}</td>
                    <td>
                      <span
                        className={`badge ${
                          cashout.status === 'pending'
                            ? 'badge-warning'
                            : cashout.status === 'completed'
                            ? 'badge-success'
                            : 'badge-error'
                        }`}
                      >
                        {cashout.status}
                      </span>
                    </td>
                    <td>{new Date(cashout.created_at).toLocaleDateString()}</td>
                    <td>
                      <div className="flex gap-2">
                        {cashout.status === 'pending' && (
                          <>
                            <button
                              className="btn btn-sm btn-success"
                              onClick={() => handleApprove(cashout.id)}
                              disabled={processing === cashout.id}
                            >
                              Approve
                            </button>
                            <button
                              className="btn btn-sm btn-error"
                              onClick={() => handleReject(cashout.id)}
                              disabled={processing === cashout.id}
                            >
                              Reject
                            </button>
                          </>
                        )}
                        {cashout.status === 'processing' && (
                          <button
                            className="btn btn-sm btn-primary"
                            onClick={() => handleMarkPaid(cashout.id)}
                            disabled={processing === cashout.id}
                          >
                            {processing === cashout.id ? (
                              <span className="loading loading-spinner loading-xs"></span>
                            ) : (
                              'Mark Paid'
                            )}
                          </button>
                        )}
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

