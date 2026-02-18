"use client"

import Link from "next/link"
import { ConnectionIndicator } from "./ConnectionIndicator"
import { StatusSection } from "./StatusSection"
import { MetricsView } from "./metricsView"
import { HistoryActions } from "./historyActions"
import { DemoModeInfo } from "./DemoModeInfo"
import { useDashboard } from "../hooks/useDashboard"
import { useMockData } from "../hooks/useMockData"
import { AgentAction, AgentStatus } from "@/app/hooks/use-agent-socket"

// Estado combinado que puede venir de WebSocket o de datos mock
interface DisplayState {
  status: AgentStatus
  currentTask: string
  actions: AgentAction[]
  metrics: {
    cpuUsage: number
    memoryUsage: number
    activeConnections: number
    errorsDetected: number
  }
  lastUpdate: string
}

interface DashboardContainerProps {
  isCheckingProfile: boolean
}

export function DashboardContainer({ isCheckingProfile }: DashboardContainerProps) {
  const {
    state: wsState,
    isConnected,
    connectionError,
    sendCommand,
  } = useDashboard()

  // Usar datos de mock cuando no estamos conectados
  const mockState = useMockData(isConnected)

  // Determinar qué estado mostrar
  const displayState: DisplayState = isConnected
    ? {
        status: wsState.status as AgentStatus,
        currentTask: wsState.currentTask,
        actions: wsState.actions,
        metrics: wsState.metrics,
        lastUpdate: wsState.lastUpdate,
      }
    : {
        status: mockState.status,
        currentTask: mockState.currentTask,
        actions: mockState.actions as AgentAction[],
        metrics: mockState.metrics,
        lastUpdate: mockState.lastUpdate,
      }

  if (isCheckingProfile) {
    return (
      <div className="min-h-screen bg-black flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-yellow-400 mx-auto mb-4"></div>
          <p className="text-gray-400">Verificando perfil...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <header className="border-b border-border">
        <div className="mx-auto flex max-w-6xl items-center justify-between px-6 py-4">
          <Link href="/" className="flex items-center gap-2">
            <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary">
              <span className="text-lg font-bold text-primary-foreground">A</span>
            </div>
            <span className="text-lg font-bold">InfraAgent</span>
          </Link>
          <ConnectionIndicator isConnected={isConnected} error={connectionError} />
        </div>
      </header>

      <main className="mx-auto max-w-6xl px-6 py-8">
        <StatusSection
          status={displayState.status}
          currentTask={displayState.currentTask}
          onPause={() => sendCommand("pause")}
          onResume={() => sendCommand("resume")}
        />

        <MetricsView metrics={displayState.metrics} />

        <HistoryActions
          actions={displayState.actions}
          lastUpdate={displayState.lastUpdate}
        />

        {!isConnected && <DemoModeInfo />}
      </main>
    </div>
  )
}
