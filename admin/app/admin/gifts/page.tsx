'use client'

import { useEffect, useState } from 'react'
import { api } from '@/lib/api'

interface Transaction {
  id: string
  user_name: string
  type: string
  coin_amount: number
  birr_amount?: number
  gift_type?: string
  created_at: string
}

export default function GiftsPage() {
  const [transactions, setTransactions] = useState<Transaction[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    loadTransactions()
  }, [])

  const loadTransactions = async () => {
    try {
      const response = await api.get('/admin/transactions')
      setTransactions(response.data.transactions || [])
    } catch (error) {
      console.error('Failed to load transactions:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleRefund = async (transactionId: string) => {
    if (!confirm('Are you sure you want to refund this transaction?')) return

    try {
      await api.post(`/admin/transactions/${transactionId}/refund`)
      await loadTransactions()
      alert('Refund processed successfully')
    } catch (error) {
      alert('Failed to process refund')
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
      <h1 className="text-3xl font-bold mb-8 text-primary">Gifts & Coins</h1>

      <div className="card bg-base-200 border border-base-300">
        <div className="card-body p-0">
          <div className="overflow-x-auto">
            <table className="table">
              <thead>
                <tr>
                  <th>User</th>
                  <th>Type</th>
                  <th>Gift</th>
                  <th>Coins</th>
                  <th>Amount</th>
                  <th>Date</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {transactions.map((tx) => (
                  <tr key={tx.id}>
                    <td>{tx.user_name}</td>
                    <td>
                      <span className="badge badge-info">{tx.type}</span>
                    </td>
                    <td>{tx.gift_type || 'N/A'}</td>
                    <td>{tx.coin_amount.toLocaleString()} LC</td>
                    <td>
                      {tx.birr_amount
                        ? `${tx.birr_amount.toFixed(2)} ETB`
                        : 'N/A'}
                    </td>
                    <td>{new Date(tx.created_at).toLocaleString()}</td>
                    <td>
                      {tx.type === 'purchase' && (
                        <button
                          className="btn btn-sm btn-warning"
                          onClick={() => handleRefund(tx.id)}
                        >
                          Refund
                        </button>
                      )}
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

