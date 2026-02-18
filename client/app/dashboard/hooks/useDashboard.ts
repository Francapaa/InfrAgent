"use client"

import { useState, useEffect } from "react"
import { useWebSocketManager } from "@/app/hooks/use-agent-socket"

const WS_URL = process.env.NEXT_PUBLIC_WS_URL

interface DashboardState {
  status: string
  currentTask: string
  actions: Array<{
    id: string
    type: string
    description: string
    status: "pending" | "running" | "completed" | "failed"
    timestamp: string
    duration?: number
  }>
  metrics: {
    cpuUsage: number
    memoryUsage: number
    activeConnections: number
    errorsDetected: number
  }
  lastUpdate: string
}

interface UseDashboardReturn {
  state: DashboardState
  isConnected: boolean
  connectionError: string | null
  wsUrlError: boolean
  sendCommand: (command: "pause" | "resume") => void
}

export function useDashboard(): UseDashboardReturn {
  const [wsUrlError, setWsUrlError] = useState(false)

  // Verificar que WS_URL existe
  useEffect(() => {
    if (!WS_URL) {
      console.error("No existe la url de websockets, imposible de conectar")
      setWsUrlError(true)
    }
  }, [])

  // Inicializar WebSocket solo si WS_URL existe
  const wsManager = WS_URL
    ? useWebSocketManager(WS_URL)
    : {
        state: {
          status: "idle",
          currentTask: "",
          actions: [],
          metrics: { cpuUsage: 0, memoryUsage: 0, activeConnections: 0, errorsDetected: 0 },
          lastUpdate: new Date().toISOString(),
        },
        isConnected: false,
        connectionError: "WebSocket URL no configurada",
        sendCommand: () => {},
      }

  const displayState = {
    status: wsManager.state.status,
    currentTask: wsManager.state.currentTask,
    actions: wsManager.state.actions,
    metrics: wsManager.state.metrics,
    lastUpdate: wsManager.state.lastUpdate,
  }

  return {
    state: displayState,
    isConnected: wsManager.isConnected,
    connectionError: wsManager.connectionError,
    wsUrlError,
    sendCommand: wsManager.sendCommand,
  }
}
