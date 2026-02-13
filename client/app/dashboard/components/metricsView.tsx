import { MetricCard } from "./MetricCard"

interface MetricsViewProps {
  metrics: {
    cpuUsage: number
    memoryUsage: number
    activeConnections: number
    errorsDetected: number
  }
}

export function MetricsView({ metrics }: MetricsViewProps) {
  return (
    <section className="mb-8">
      <h2 className="mb-4 text-lg font-semibold">Métricas en Tiempo Real</h2>
      <div className="grid grid-cols-2 gap-4 lg:grid-cols-4">
        <MetricCard
          label="CPU"
          value={Math.round(metrics.cpuUsage)}
          unit="%"
          warning={metrics.cpuUsage > 80}
        />
        <MetricCard
          label="Memoria"
          value={Math.round(metrics.memoryUsage)}
          unit="%"
          warning={metrics.memoryUsage > 85}
        />
        <MetricCard label="Conexiones" value={metrics.activeConnections} unit="" />
        <MetricCard
          label="Errores"
          value={metrics.errorsDetected}
          unit=""
          warning={metrics.errorsDetected > 0}
        />
      </div>
    </section>
  )
}
