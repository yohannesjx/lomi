'use client'

import { createContext, useContext, useEffect, useState } from 'react'
import { useRouter, usePathname } from 'next/navigation'
import { api } from '@/lib/api'

interface User {
  id: string
  name: string
  email?: string
  role?: string
}

interface AuthContextType {
  user: User | null
  loading: boolean
  login: (token: string) => Promise<void>
  logout: () => void
  isAdmin: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [loading, setLoading] = useState(true)
  const router = useRouter()
  const pathname = usePathname()

  useEffect(() => {
    checkAuth()
  }, [])

  useEffect(() => {
    // Protect admin routes
    if (!loading && pathname?.startsWith('/admin')) {
      if (!user) {
        router.push('/login')
      } else if (user.role !== 'admin') {
        router.push('/')
      }
    }
  }, [user, loading, pathname, router])

  const checkAuth = async () => {
    try {
      const token = localStorage.getItem('admin_token')
      if (!token) {
        setLoading(false)
        return
      }

      const response = await api.get('/users/me')
      const userData = response.data

      if (userData.role !== 'admin') {
        localStorage.removeItem('admin_token')
        setUser(null)
        router.push('/')
        return
      }

      setUser(userData)
    } catch (error) {
      console.error('Auth check failed:', error)
      localStorage.removeItem('admin_token')
      setUser(null)
    } finally {
      setLoading(false)
    }
  }

  const login = async (token: string) => {
    localStorage.setItem('admin_token', token)
    await checkAuth()
  }

  const logout = () => {
    localStorage.removeItem('admin_token')
    setUser(null)
    router.push('/login')
  }

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        login,
        logout,
        isAdmin: user?.role === 'admin',
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return context
}

