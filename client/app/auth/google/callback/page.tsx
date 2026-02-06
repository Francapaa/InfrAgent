"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"

export const dynamic = 'force-dynamic'

export default function GoogleCallback() {
  const router = useRouter()
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const handleCallback = async () => {
      try {
        const params = new URLSearchParams(window.location.search)
        const code = params.get("code")
        
        if (!code) {
          setError("No se recibió código de autorización")
          setLoading(false)
          return
        }

        const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/auth/google/callback?${window.location.search}`)
        const data = await response.json()

        if (response.ok && data.Token) {
          localStorage.setItem("token", data.Token)
          localStorage.setItem("user", JSON.stringify({
            email: data.Email || "",
            name: data.Name || "",
          }))
          router.push("/dashboard")
        } else {
          setError(data.Error || "Error al procesar la respuesta")
          setLoading(false)
        }
      } catch (err) {
        setError("Error al conectar con el servidor")
        setLoading(false)
      }
    }

    handleCallback()
  }, [router])

  if (error) {
    return (
      <div className="min-h-screen bg-black flex items-center justify-center">
        <div className="text-center">
          <p className="text-red-500 text-xl mb-4">{error}</p>
          <a href="/login" className="text-yellow-400 hover:text-yellow-300">
            Volver al login
          </a>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-black flex items-center justify-center">
      <div className="text-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-yellow-400 mx-auto mb-4"></div>
        <p className="text-gray-400">Procesando inicio de sesión...</p>
      </div>
    </div>
  )
}
