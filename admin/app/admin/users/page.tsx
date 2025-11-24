'use client'

import { useEffect, useState } from 'react'
import { api } from '@/lib/api'

interface User {
  id: string
  name: string
  email?: string
  telegram_id?: string
  city?: string
  is_active: boolean
  created_at: string
}

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [banning, setBanning] = useState<string | null>(null)

  useEffect(() => {
    loadUsers()
  }, [])

  const loadUsers = async () => {
    try {
      const response = await api.get('/users')
      setUsers(response.data.users || [])
    } catch (error) {
      console.error('Failed to load users:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleBan = async (userId: string) => {
    if (!confirm('Are you sure you want to ban this user?')) return

    setBanning(userId)
    try {
      await api.post(`/admin/users/${userId}/ban`)
      await loadUsers()
    } catch (error) {
      alert('Failed to ban user')
    } finally {
      setBanning(null)
    }
  }

  const handleUnban = async (userId: string) => {
    setBanning(userId)
    try {
      await api.post(`/admin/users/${userId}/unban`)
      await loadUsers()
    } catch (error) {
      alert('Failed to unban user')
    } finally {
      setBanning(null)
    }
  }

  const filteredUsers = users.filter(
    (user) =>
      user.name.toLowerCase().includes(search.toLowerCase()) ||
      user.email?.toLowerCase().includes(search.toLowerCase()) ||
      user.telegram_id?.toString().includes(search)
  )

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <span className="loading loading-spinner loading-lg text-primary"></span>
      </div>
    )
  }

  return (
    <div>
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold text-primary">Users</h1>
        <div className="form-control w-full max-w-xs">
          <input
            type="text"
            placeholder="Search users..."
            className="input input-bordered w-full"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>
      </div>

      <div className="card bg-base-200 border border-base-300">
        <div className="card-body p-0">
          <div className="overflow-x-auto">
            <table className="table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Email/Telegram</th>
                  <th>City</th>
                  <th>Status</th>
                  <th>Joined</th>
                  <th>Actions</th>
                </tr>
              </thead>
              <tbody>
                {filteredUsers.map((user) => (
                  <tr key={user.id}>
                    <td>
                      <div className="font-semibold">{user.name}</div>
                    </td>
                    <td>
                      {user.email || `@${user.telegram_id}` || 'N/A'}
                    </td>
                    <td>{user.city || 'N/A'}</td>
                    <td>
                      {user.is_active ? (
                        <span className="badge badge-success">Active</span>
                      ) : (
                        <span className="badge badge-error">Banned</span>
                      )}
                    </td>
                    <td>
                      {new Date(user.created_at).toLocaleDateString()}
                    </td>
                    <td>
                      <div className="flex gap-2">
                        <button
                          className="btn btn-sm btn-info"
                          onClick={() =>
                            window.open(`/admin/users/${user.id}`, '_blank')
                          }
                        >
                          View
                        </button>
                        {user.is_active ? (
                          <button
                            className="btn btn-sm btn-error"
                            onClick={() => handleBan(user.id)}
                            disabled={banning === user.id}
                          >
                            {banning === user.id ? (
                              <span className="loading loading-spinner loading-xs"></span>
                            ) : (
                              'Ban'
                            )}
                          </button>
                        ) : (
                          <button
                            className="btn btn-sm btn-success"
                            onClick={() => handleUnban(user.id)}
                            disabled={banning === user.id}
                          >
                            {banning === user.id ? (
                              <span className="loading loading-spinner loading-xs"></span>
                            ) : (
                              'Unban'
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

      <div className="mt-4 text-sm text-base-content/60">
        Showing {filteredUsers.length} of {users.length} users
      </div>
    </div>
  )
}

