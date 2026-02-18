import { ActionItem } from "./ActionItem"
import { AgentAction } from "@/app/hooks/use-agent-socket"

interface HistoryActionsProps {
  actions: AgentAction[]
  lastUpdate: string
}

export function HistoryActions({ actions, lastUpdate }: HistoryActionsProps) {
  return (
    <section>
      <div className="mb-4 flex items-center justify-between">
        <h2 className="text-lg font-semibold">Historial de Acciones</h2>
        {lastUpdate && (
          <span className="text-xs text-muted-foreground">
            Última actualización:{" "}
            {new Date(lastUpdate).toLocaleTimeString("es-ES")}
          </span>
        )}
      </div>
      <div className="rounded-lg border border-border bg-muted/20">
        {actions.length > 0 ? (
          <div className="max-h-96 overflow-y-auto px-4">
            {actions.map((action: AgentAction) => (
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
  )
}
