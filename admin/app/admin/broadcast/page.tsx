'use client'

import { useState } from 'react'
import { api } from '@/lib/api'

export default function BroadcastPage() {
  const [title, setTitle] = useState('')
  const [message, setMessage] = useState('')
  const [sending, setSending] = useState(false)

  const handleSend = async () => {
    if (!title || !message) {
      alert('Please fill in both title and message')
      return
    }

    if (!confirm('Send broadcast to all users?')) return

    setSending(true)
    try {
      await api.post('/admin/broadcast', {
        title,
        message,
      })
      alert('Broadcast sent successfully!')
      setTitle('')
      setMessage('')
    } catch (error) {
      alert('Failed to send broadcast')
    } finally {
      setSending(false)
    }
  }

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8 text-primary">Broadcast</h1>

      <div className="card bg-base-200 border border-base-300 max-w-2xl">
        <div className="card-body">
          <h2 className="card-title text-primary">Send Push Notification</h2>

          <div className="form-control w-full mb-4">
            <label className="label">
              <span className="label-text">Title</span>
            </label>
            <input
              type="text"
              placeholder="Notification title"
              className="input input-bordered w-full"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
            />
          </div>

          <div className="form-control w-full mb-6">
            <label className="label">
              <span className="label-text">Message</span>
            </label>
            <textarea
              className="textarea textarea-bordered h-32"
              placeholder="Notification message"
              value={message}
              onChange={(e) => setMessage(e.target.value)}
            ></textarea>
          </div>

          <div className="card-actions">
            <button
              className="btn btn-primary"
              onClick={handleSend}
              disabled={sending || !title || !message}
            >
              {sending ? (
                <>
                  <span className="loading loading-spinner"></span>
                  Sending...
                </>
              ) : (
                '📢 Send Broadcast'
              )}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}

