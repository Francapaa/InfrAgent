"use client"

import { useState, useEffect } from "react"
import { AgentStatus } from "@/app/hooks/use-agent-socket"

interface MockAction {
  id: string
  type: string
  description: string
  status: "pending" | "running" | "completed" | "failed"
  timestamp: string
  duration?: number
}

interface MockState {
  status: AgentStatus
  currentTask: string
  actions: MockAction[]
  metrics: {
    cpuUsage: number
    memoryUsage: number
    activeConnections: number
    errorsDetected: number
  }
  lastUpdate: string
}

export function useMockData(isConnected: boolean): MockState {
  const [mockState, setMockState] = useState<MockState>({
    status: "monitoring",
    currentTask: "Monitoreando servicios de producción",
    actions: [
      {
        id: "1",
        type: "health_check",
        description: "Verificación de salud del servicio API",
        status: "completed",
        timestamp: new Date(Date.now() - 60000).toISOString(),
        duration: 245,
      },
      {
        id: "2",
        type: "memory_alert",
        description: "Detectado uso elevado de memoria en worker-3",
        status: "completed",
        timestamp: new Date(Date.now() - 120000).toISOString(),
        duration: 1200,
      },
      {
        id: "3",
        type: "auto_scale",
        description: "Escalado automático de réplicas: 3 → 5",
        status: "running",
        timestamp: new Date(Date.now() - 180000).toISOString(),
      },
      {
        id: "4",
        type: "log_analysis",
        description: "Análisis de patrones en logs de errores",
        status: "pending",
        timestamp: new Date(Date.now() - 240000).toISOString(),
      },
    ],
    metrics: {
      cpuUsage: 42,
      memoryUsage: 67,
      activeConnections: 1284,
      errorsDetected: 3,
    },
    lastUpdate: new Date().toISOString(),
  })

  useEffect(() => {
    if (isConnected) return

    const interval = setInterval(() => {
      setMockState((prev) => ({
        ...prev,
        metrics: {
          cpuUsage: Math.min(100, Math.max(0, prev.metrics.cpuUsage + (Math.random() - 0.5) * 10)),
          memoryUsage: Math.min(100, Math.max(0, prev.metrics.memoryUsage + (Math.random() - 0.5) * 5)),
          activeConnections: Math.max(0, prev.metrics.activeConnections + Math.floor((Math.random() - 0.5) * 100)),
          errorsDetected: prev.metrics.errorsDetected,
        },
        lastUpdate: new Date().toISOString(),
      }))
    }, 2000)

    return () => clearInterval(interval)
  }, [isConnected])

  return mockState
}
