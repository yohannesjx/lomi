'use client'

import { useEffect, useState } from 'react'
import { api } from '@/lib/api'

interface Photo {
  id: string
  user_id: string
  user_name: string
  url: string
  thumbnail_url?: string
  is_approved: boolean
  moderation_status: string
  created_at: string
}

export default function PhotosPage() {
  const [photos, setPhotos] = useState<Photo[]>([])
  const [loading, setLoading] = useState(true)
  const [processing, setProcessing] = useState<string | null>(null)
  const [rejectReason, setRejectReason] = useState('')
  const [selectedPhoto, setSelectedPhoto] = useState<Photo | null>(null)

  useEffect(() => {
    loadPhotos()
  }, [])

  const loadPhotos = async () => {
    try {
      const response = await api.get('/admin/moderation/pending')
      setPhotos(response.data.photos || [])
    } catch (error) {
      console.error('Failed to load photos:', error)
      // Mock data for development
      setPhotos([])
    } finally {
      setLoading(false)
    }
  }

  const handleApprove = async (photoId: string) => {
    setProcessing(photoId)
    try {
      await api.put(`/admin/moderation/${photoId}/approve`)
      await loadPhotos()
    } catch (error) {
      alert('Failed to approve photo')
    } finally {
      setProcessing(null)
    }
  }

  const handleReject = async (photoId: string) => {
    if (!rejectReason) {
      alert('Please select a rejection reason')
      return
    }

    setProcessing(photoId)
    try {
      await api.put(`/admin/moderation/${photoId}/reject`, {
        reason: rejectReason,
      })
      setSelectedPhoto(null)
      setRejectReason('')
      await loadPhotos()
    } catch (error) {
      alert('Failed to reject photo')
    } finally {
      setProcessing(null)
    }
  }

  const rejectReasons = [
    'Inappropriate content',
    'Not a real person',
    'Blurry/low quality',
    'Contains text/watermark',
    'Violates community guidelines',
    'Other',
  ]

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <span className="loading loading-spinner loading-lg text-primary"></span>
      </div>
    )
  }

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8 text-primary">
        Photo Moderation
      </h1>

      {photos.length === 0 ? (
        <div className="card bg-base-200 border border-base-300">
          <div className="card-body text-center py-16">
            <p className="text-xl text-base-content/60">
              No pending photos to moderate 🎉
            </p>
          </div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {photos.map((photo) => (
            <div
              key={photo.id}
              className="card bg-base-200 border border-base-300"
            >
              <figure className="aspect-square bg-base-300">
                <img
                  src={photo.thumbnail_url || photo.url}
                  alt="Photo to moderate"
                  className="w-full h-full object-cover"
                />
              </figure>
              <div className="card-body">
                <h3 className="card-title text-sm">{photo.user_name}</h3>
                <p className="text-xs text-base-content/60">
                  {new Date(photo.created_at).toLocaleString()}
                </p>
                <div className="card-actions justify-end mt-4">
                  <button
                    className="btn btn-success btn-sm"
                    onClick={() => handleApprove(photo.id)}
                    disabled={processing === photo.id}
                  >
                    {processing === photo.id ? (
                      <span className="loading loading-spinner loading-xs"></span>
                    ) : (
                      '✓ Approve'
                    )}
                  </button>
                  <button
                    className="btn btn-error btn-sm"
                    onClick={() => setSelectedPhoto(photo)}
                    disabled={processing === photo.id}
                  >
                    ✗ Reject
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Reject Modal */}
      {selectedPhoto && (
        <dialog className="modal modal-open">
          <div className="modal-box">
            <h3 className="font-bold text-lg mb-4">Reject Photo</h3>
            <p className="mb-4">Select a reason for rejection:</p>
            <select
              className="select select-bordered w-full mb-4"
              value={rejectReason}
              onChange={(e) => setRejectReason(e.target.value)}
            >
              <option value="">Choose reason...</option>
              {rejectReasons.map((reason) => (
                <option key={reason} value={reason}>
                  {reason}
                </option>
              ))}
            </select>
            <div className="modal-action">
              <button
                className="btn"
                onClick={() => {
                  setSelectedPhoto(null)
                  setRejectReason('')
                }}
              >
                Cancel
              </button>
              <button
                className="btn btn-error"
                onClick={() => handleReject(selectedPhoto.id)}
                disabled={!rejectReason || processing === selectedPhoto.id}
              >
                {processing === selectedPhoto.id ? (
                  <span className="loading loading-spinner loading-xs"></span>
                ) : (
                  'Reject'
                )}
              </button>
            </div>
          </div>
        </dialog>
      )}
    </div>
  )
}

