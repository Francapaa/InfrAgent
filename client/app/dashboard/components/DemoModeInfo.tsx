export function DemoModeInfo() {
  return (
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
  )
}
