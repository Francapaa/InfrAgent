import { AgentStatus } from "@/app/hooks/use-agent-socket"

interface StatusBadgeProps {
  status: AgentStatus
}

export function StatusBadge({ status }: StatusBadgeProps) {
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
