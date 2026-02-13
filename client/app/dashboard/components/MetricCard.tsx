interface MetricCardProps {
  label: string
  value: number
  unit: string
  warning?: boolean
}

export function MetricCard({ label, value, unit, warning }: MetricCardProps) {
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
