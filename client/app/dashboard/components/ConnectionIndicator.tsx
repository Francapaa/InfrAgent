interface ConnectionIndicatorProps {
  isConnected: boolean
  error: string | null
}

export function ConnectionIndicator({ isConnected, error }: ConnectionIndicatorProps) {
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
