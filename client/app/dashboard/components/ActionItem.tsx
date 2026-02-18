import { AgentAction } from "@/app/hooks/use-agent-socket"

interface ActionItemProps {
  action: AgentAction
}

export function ActionItem({ action }: ActionItemProps) {
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
          <p className="mt-1 text-xs text-muted-foreground">Duración: {action.duration}ms</p>
        )}
      </div>
    </div>
  )
}
