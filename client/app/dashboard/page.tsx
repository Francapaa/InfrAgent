"use client"

import { useState, useEffect } from "react"
import Link from "next/link"
import { useRouter } from "next/navigation"
import { useWebSocketManager, type AgentAction, type AgentStatus } from "@/app/hooks/use-agent-socket"

const WS_URL = process.env.NEXT_PUBLIC_WS_URL;

console.log(WS_URL); 

function StatusBadge({ status }: { status: AgentStatus }) {
  const config: Record<AgentStatus, { label: string; color: string; pulse: boolean }> = {
    idle: { label: "Inactivo", color: "bg-muted-foreground", pulse: false },
    monitoring: { label: "Monitoreando", color: "bg-primary", pulse: true },
    analyzing: { label: "Analizando", color: "bg-blue-500", pulse: true },
    executing: { label: "Ejecutando", color: "bg-orange-500", pulse: true },
    error: { label: "Error", color: "bg-red-500", pulse: false },
  }

  const { label, color, pulse } = config[status]

  return (
    <div className="flex items-center gap-2">
      <span className="relative flex h-3 w-3">
        {pulse && (
          <span
            className={`absolute inline-flex h-full w-full animate-ping rounded-full opacity-75 ${color}`}
          />
        )}
        <span className={`relative inline-flex h-3 w-3 rounded-full ${color}`} />
      </span>
      <span className="text-sm font-medium">{label}</span>
    </div>
  )
}

function ConnectionIndicator({ isConnected, error }: { isConnected: boolean; error: string | null }) {
  return (
    <div className="flex items-center gap-2 text-sm">
      <span
        className={`h-2 w-2 rounded-full ${
          isConnected ? "bg-green-500" : error ? "bg-red-500" : "bg-muted-foreground"
        }`}
      />
      <span className="text-muted-foreground">
        {isConnected ? "Conectado" : error || "Desconectado"}
      </span>
    </div>
  )
}

function MetricCard({
  label,
  value,
  unit,
  warning,
}: {
  label: string
  value: number
  unit: string
  warning?: boolean
}) {
  return (
    <div className="rounded-lg border border-border bg-muted/30 p-4">
      <p className="text-sm text-muted-foreground">{label}</p>
      <p className={`mt-1 text-2xl font-bold ${warning ? "text-red-500" : "text-foreground"}`}>
        {value}
        <span className="text-sm font-normal text-muted-foreground">{unit}</span>
      </p>
    </div>
  )
}

function ActionItem({ action }: { action: AgentAction }) {
  const statusConfig: Record<AgentAction["status"], { icon: string; color: string }> = {
    pending: { icon: "○", color: "text-muted-foreground" },
    running: { icon: "◐", color: "text-primary" },
    completed: { icon: "●", color: "text-green-500" },
    failed: { icon: "✕", color: "text-red-500" },
  }

  const { icon, color } = statusConfig[action.status]
  const time = new Date(action.timestamp).toLocaleTimeString("es-ES", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  })

  return (
    <div className="flex items-start gap-3 border-b border-border py-3 last:border-0">
      <span className={`mt-0.5 font-mono text-lg ${color}`}>{icon}</span>
      <div className="flex-1 min-w-0">
        <div className="flex items-center justify-between gap-2">
          <span className="font-medium truncate">{action.type}</span>
          <span className="text-xs text-muted-foreground shrink-0">{time}</span>
        </div>
        <p className="mt-0.5 text-sm text-muted-foreground truncate">{action.description}</p>
        {action.duration && (
          <p className="mt-1 text-xs text-muted-foreground">
            Duración: {action.duration}ms
          </p>
        )}
      </div>
    </div>
  )
}

// Mock data for demo when not connected
function useMockData(isConnected: boolean) {
  const [mockState, setMockState] = useState({
    status: "monitoring" as AgentStatus,
    currentTask: "Monitoreando servicios de producción",
    actions: [
      {
        id: "1",
        type: "health_check",
        description: "Verificación de salud del servicio API",
        status: "completed" as const,
        timestamp: new Date(Date.now() - 60000).toISOString(),
        duration: 245,
      },
      {
        id: "2",
        type: "memory_alert",
        description: "Detectado uso elevado de memoria en worker-3",
        status: "completed" as const,
        timestamp: new Date(Date.now() - 120000).toISOString(),
        duration: 1200,
      },
      {
        id: "3",
        type: "auto_scale",
        description: "Escalado automático de réplicas: 3 → 5",
        status: "running" as const,
        timestamp: new Date(Date.now() - 180000).toISOString(),
      },
      {
        id: "4",
        type: "log_analysis",
        description: "Análisis de patrones en logs de errores",
        status: "pending" as const,
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

export default function DashboardPage() {
  const router = useRouter()
  if(!WS_URL){
    console.error("No existe la url de websockets, imposible de conectar")
    return; 
  }
  const { state, isConnected, connectionError, sendCommand } = useWebSocketManager(WS_URL)
  const mockState = useMockData(isConnected)
  const [isCheckingProfile, setIsCheckingProfile] = useState(true)
  
  const displayState = isConnected ? state : mockState

  // Verificar si el perfil está completo
  useEffect(() => {
    const checkProfile = async () => {

      try {
        const backendUrl = process.env.NEXT_PUBLIC_BACKEND_URL || "http://localhost:8080"
        const response = await fetch(`${backendUrl}/auth/me`, {
         credentials: 'include'
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

  // Mostrar loading mientras verifica el perfil
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
        {/* Status Section */}
        <section className="mb-8">
          <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h1 className="text-2xl font-bold">Estado del Agente</h1>
              {displayState.currentTask && (
                <p className="mt-1 text-muted-foreground">{displayState.currentTask}</p>
              )}
            </div>
            <div className="flex items-center gap-4">
              <StatusBadge status={displayState.status} />
              <div className="flex gap-2">
                <button
                  onClick={() => sendCommand("pause")}
                  className="rounded-md border border-border px-3 py-1.5 text-sm transition-colors hover:bg-muted"
                >
                  Pausar
                </button>
                <button
                  onClick={() => sendCommand("resume")}
                  className="rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
                >
                  Reanudar
                </button>
              </div>
            </div>
          </div>
        </section>

        {/* Metrics Grid */}
        <section className="mb-8">
          <h2 className="mb-4 text-lg font-semibold">Métricas en Tiempo Real</h2>
          <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
            <MetricCard
              label="CPU"
              value={Math.round(displayState.metrics.cpuUsage)}
              unit="%"
              warning={displayState.metrics.cpuUsage > 80}
            />
            <MetricCard
              label="Memoria"
              value={Math.round(displayState.metrics.memoryUsage)}
              unit="%"
              warning={displayState.metrics.memoryUsage > 85}
            />
            <MetricCard
              label="Conexiones"
              value={displayState.metrics.activeConnections}
              unit=""
            />
            <MetricCard
              label="Errores"
              value={displayState.metrics.errorsDetected}
              unit=""
              warning={displayState.metrics.errorsDetected > 0}
            />
          </div>
        </section>

        {/* Actions Log */}
        <section>
          <div className="mb-4 flex items-center justify-between">
            <h2 className="text-lg font-semibold">Historial de Acciones</h2>
            {displayState.lastUpdate && (
              <span className="text-xs text-muted-foreground">
                Última actualización:{" "}
                {new Date(displayState.lastUpdate).toLocaleTimeString("es-ES")}
              </span>
            )}
          </div>
          <div className="rounded-lg border border-border bg-muted/20">
            {displayState.actions.length > 0 ? (
              <div className="max-h-96 overflow-y-auto px-4">
                {displayState.actions.map((action: AgentAction) => (
                  <ActionItem key={action.id} action={action} />
                ))}
              </div>
            ) : (
              <div className="py-12 text-center text-muted-foreground">
                No hay acciones registradas
              </div>
            )}
          </div>
        </section>

        {/* WebSocket Info */}
        {!isConnected && (
          <section className="mt-8 rounded-lg border border-border bg-muted/20 p-6">
            <h3 className="mb-2 font-semibold">Modo Demo</h3>
            <p className="text-sm text-muted-foreground">
              Actualmente mostrando datos de demostración. Para conectar con tu backend en Go,
              configura la variable de entorno{" "}
              <code className="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">
                NEXT_PUBLIC_WS_URL
              </code>{" "}
              con la URL de tu servidor WebSocket.
            </p>
            <pre className="mt-4 overflow-x-auto rounded-md bg-background p-4 font-mono text-xs">
              <code className="text-muted-foreground">
{`// Ejemplo de mensaje esperado desde Go
{
  "type": "state_update",
  "payload": {
    "status": "monitoring",
    "currentTask": "Verificando servicios..."
  }
}

// Tipos de mensajes soportados:
// - state_update
// - action_added
// - action_updated  
// - metrics_update
// - status_change`}
              </code>
            </pre>
          </section>
        )}
      </main>
    </div>
  )
}
