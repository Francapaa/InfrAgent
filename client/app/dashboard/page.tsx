"use client"

import { useState, useEffect } from "react"
import { useRouter } from "next/navigation"
import { DashboardContainer } from "./components/DashboardContainer"

const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL || "http://localhost:8080"

export default function DashboardPage() {
  const router = useRouter()
  const [isCheckingProfile, setIsCheckingProfile] = useState(true)

  useEffect(() => {
    const checkProfile = async () => {
      try {
        const response = await fetch(`${BACKEND_URL}/auth/me`, {
          credentials: "include",
        })

        if (response.ok) {
          const userData = await response.json()

          // Si no tiene perfil completo, redirigir a onboarding
          if (!userData.company_name || !userData.webhook_url) {
            router.push("/onboarding")
            return
          }
        } else if (response.status === 401) {
          // Token inválido
          router.push("/login")
          return
        }
      } catch (err) {
        console.error("Error checking profile:", err)
      } finally {
        setIsCheckingProfile(false)
      }
    }

    checkProfile()
  }, [router])

  return <DashboardContainer isCheckingProfile={isCheckingProfile} />
}
