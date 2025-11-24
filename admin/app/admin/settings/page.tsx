'use client'

import { useEffect, useState } from 'react'
import { api } from '@/lib/api'

interface Settings {
  maintenance_mode: boolean
  registration_enabled: boolean
  gift_sending_enabled: boolean
  cashout_enabled: boolean
}

export default function SettingsPage() {
  const [settings, setSettings] = useState<Settings>({
    maintenance_mode: false,
    registration_enabled: true,
    gift_sending_enabled: true,
    cashout_enabled: true,
  })
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    loadSettings()
  }, [])

  const loadSettings = async () => {
    try {
      const response = await api.get('/admin/settings')
      setSettings(response.data)
    } catch (error) {
      console.error('Failed to load settings:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleToggle = async (key: keyof Settings) => {
    const newValue = !settings[key]
    setSaving(true)
    try {
      await api.put('/admin/settings', {
        [key]: newValue,
      })
      setSettings({ ...settings, [key]: newValue })
    } catch (error) {
      alert('Failed to update setting')
    } finally {
      setSaving(false)
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
      <h1 className="text-3xl font-bold mb-8 text-primary">Settings</h1>

      <div className="card bg-base-200 border border-base-300 max-w-2xl">
        <div className="card-body">
          <h2 className="card-title text-primary mb-6">Feature Toggles</h2>

          <div className="space-y-4">
            <div className="flex items-center justify-between p-4 bg-base-300 rounded-lg">
              <div>
                <h3 className="font-semibold">Maintenance Mode</h3>
                <p className="text-sm text-base-content/60">
                  Disable app access for all users
                </p>
              </div>
              <input
                type="checkbox"
                className="toggle toggle-error"
                checked={settings.maintenance_mode}
                onChange={() => handleToggle('maintenance_mode')}
                disabled={saving}
              />
            </div>

            <div className="flex items-center justify-between p-4 bg-base-300 rounded-lg">
              <div>
                <h3 className="font-semibold">Registration</h3>
                <p className="text-sm text-base-content/60">
                  Allow new user registration
                </p>
              </div>
              <input
                type="checkbox"
                className="toggle toggle-primary"
                checked={settings.registration_enabled}
                onChange={() => handleToggle('registration_enabled')}
                disabled={saving}
              />
            </div>

            <div className="flex items-center justify-between p-4 bg-base-300 rounded-lg">
              <div>
                <h3 className="font-semibold">Gift Sending</h3>
                <p className="text-sm text-base-content/60">
                  Allow users to send gifts
                </p>
              </div>
              <input
                type="checkbox"
                className="toggle toggle-primary"
                checked={settings.gift_sending_enabled}
                onChange={() => handleToggle('gift_sending_enabled')}
                disabled={saving}
              />
            </div>

            <div className="flex items-center justify-between p-4 bg-base-300 rounded-lg">
              <div>
                <h3 className="font-semibold">Cashouts</h3>
                <p className="text-sm text-base-content/60">
                  Allow users to request cashouts
                </p>
              </div>
              <input
                type="checkbox"
                className="toggle toggle-primary"
                checked={settings.cashout_enabled}
                onChange={() => handleToggle('cashout_enabled')}
                disabled={saving}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

