import { StatusBadge } from "./StatusBadge"
import { AgentStatus } from "@/app/hooks/use-agent-socket"

interface StatusSectionProps {
  status: AgentStatus
  currentTask: string
  onPause: () => void
  onResume: () => void
}

export function StatusSection({ status, currentTask, onPause, onResume }: StatusSectionProps) {
  return (
    <section className="mb-8">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold">Estado del Agente</h1>
          {currentTask && <p className="mt-1 text-muted-foreground">{currentTask}</p>}
        </div>
        <div className="flex items-center gap-4">
          <StatusBadge status={status} />
          <div className="flex gap-2">
            <button
              onClick={onPause}
              className="rounded-md border border-border px-3 py-1.5 text-sm transition-colors hover:bg-muted"
            >
              Pausar
            </button>
            <button
              onClick={onResume}
              className="rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground transition-colors hover:bg-primary/90"
            >
              Reanudar
            </button>
          </div>
        </div>
      </div>
    </section>
  )
}
