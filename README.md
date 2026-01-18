# InfrAgent
# Arquitectura Completa del Sistema
#AGENTE DE INFRAESTRUCTURA
## 🏗️ Los 3 Componentes Físicos

```
┌─────────────────────────────────────────────────────────────┐
│                    SERVIDOR DEL USUARIO                      │
│                  (Infraestructura del cliente)               │
│                                                              │
│  ┌────────────────────────────────────────────────────────┐ │
│  │ 1. SU APLICACIÓN (payments-api)                         │ │
│  │    - Puerto 8080                                        │ │
│  │    - Es su negocio (e-commerce, fintech, etc)          │ │
│  └────────────────────────────────────────────────────────┘ │
│                           ▲                                  │
│                           │ monitorea                        │
│  ┌────────────────────────┴───────────────────────────────┐ │
│  │ 2. SDK (que la plataforma da)                            │ │
│  │    - Corre en el servidor del cliente                   │ │
│  │    - Hace 2 cosas:                                      │ │
│  │      A) Monitorea la app cada 30s                       │ │
│  │      B) Expone webhook para recibir acciones            │ │
│  └─────────────────────────────────────────────────────────┘ │
│           │                                     ▲             │
│           │ reporta eventos                     │             │
│           │ (cuando detecta problema)           │             │
│           │                                     │ recibe      │
│           │                                     │ webhooks    │
└───────────┼─────────────────────────────────────┼─────────────┘
            │                                     │
            │ INTERNET                            │
            │                                     │
┌───────────▼─────────────────────────────────────┼─────────────┐
│                    TU SERVIDOR                  │             │
│                (Tu plataforma SaaS)             │             │
│                                                 │             │
│  ┌──────────────────────────────────────────┐  │             │
│  │ 3A. INGEST API (Gin)                      │  │             │
│  │     - Puerto 8080                         │  │             │
│  │     - Recibe eventos del SDK ◄────────────┼──┘             │
│  │     - Guarda en PostgreSQL                │                │
│  └──────────────────┬───────────────────────┘                │
│                     │                                         │
│                     ▼                                         │
│  ┌──────────────────────────────────────────┐                │
│  │ PostgreSQL                                │                │
│  │ - events (pending = NULL)                 │                │
│  │ - agents                                  │                │
│  │ - clients                                 │                │
│  └──────────────────┬───────────────────────┘                │
│                     │                                         │
│                     │ lee                                     │
│                     ▼                                         │
│  ┌──────────────────────────────────────────┐                │
│  │ 3B. AGENT (loop)                          │                │
│  │     - Cada 30s lee eventos pending        │                │
│  │     - Pregunta a Gemini qué hacer         │                │
│  │     - Ejecuta acción ──────────────────────┼───────────────┘
│  └──────────────────────────────────────────┘   llama webhook
│                                                  del cliente
└─────────────────────────────────────────────────────────────┘
```

---

## 📊 Flujo Temporal Completo (Paso a Paso)

### MOMENTO 1: Setup inicial (una sola vez)

**El cliente se registra en la plataforma:**

```bash
# El cliente ejecuta (o hace desde un formulario web):
curl -X POST https://tu-plataforma.com/api/clients/register \
  -d '{
    "email": "admin@cliente.com",
    "company_name": "Cliente Corp",
    "webhook_url": "https://servidor-cliente.com:9000/webhooks/agent"
  }'

# Respuesta:
{
  "client_id": "client-abc123",
  "api_key": "agent_key_xyz789",           ← Guardar
  "webhook_secret": "whsec_secret456"      ← Guardar
}
```

**¿Qué pasó internamente?**

```
Cliente                         Servidor
  │                                 │
  │ POST /api/clients/register      │
  ├────────────────────────────────>│
  │                                 │
  │                                 ▼
  │                    cmd/api/main.go (Gin escuchando)
  │                                 │
  │                                 ▼
  │                    internal/api/ingest.go
  │                    RegisterClient()
  │                                 │
  │                                 ├─ Generar api_key
  │                                 ├─ Generar webhook_secret
  │                                 ├─ Hash api_key (bcrypt)
  │                                 │
  │                                 ▼
  │                    PostgreSQL INSERT INTO clients
  │                    (trigger crea agent automáticamente)
  │                                 │
  │                                 ▼
  │                    PostgreSQL INSERT INTO agents
  │                    (1 agent por cliente)
  │                                 │
  │  JSON con credenciales          │
  │<────────────────────────────────┤
  │                                 │
```

---

### MOMENTO 2: Cliente instala SDK (una sola vez)

**En el servidor del cliente:**

```go
// main.go (en servidor del cliente)
package main

import "tu-sdk"

func main() {
    sdk := New(
        "https://tu-plataforma.com",  
        "agent_key_xyz789",            // De la respuesta del registro
        "whsec_secret456",             // De la respuesta del registro
    )
    
    // TAREA 1: Exponer webhook
    r := gin.Default()
    r.POST("/webhooks/agent", sdk.WebhookHandler())
    go r.Run(":9000")
    
    // TAREA 2: Monitorear
    sdk.MonitorAndReport()
}
```

**Esto hace que el servidor del cliente tenga:**
- ✅ Un webhook escuchando en `:9000/webhooks/agent`
- ✅ Un loop monitoreando su propia app cada 30s

---

### MOMENTO 3: SDK detecta problema (cada 30s automáticamente)

**Loop del SDK (corre en servidor del cliente):**

```go
// Esto corre en el servidor DEL CLIENTE
for {
    // Chequea SU PROPIA app
    resp, err := http.Get("http://localhost:8080/health")
    
    if err != nil || resp.StatusCode != 200 {
        // ¡Problema detectado!
        sdk.ReportEvent("app_down", "api", "critical", data)
    }
    
    sleep(30 * time.Second)
}
```

**¿Qué hace `ReportEvent()`?**

```go
func (sdk *SDK) ReportEvent(tipo, service, severity string, data map[string]interface{}) {
    // Construye JSON
    payload := {
        "type": tipo,
        "service": service,
        "severity": severity,
        "data": data,
    }
    
    // Hace HTTP POST a TU servidor
    POST https://tu-plataforma.com/api/events
    Headers:
        Authorization: Bearer agent_key_xyz789
    Body:
        {"type":"app_down","service":"api","severity":"critical"}
}
```

**Flujo:**

```
Servidor del Cliente                    Servidor
  │                                         │
  │ SDK detecta: app no responde            │
  │         ↓                                │
  │ POST /api/events                         │
  │ Authorization: Bearer agent_key_xyz789   │
  │ {"type":"app_down","service":"api"}      │
  ├─────────────────────────────────────────>│
  │                                         │
  │                                         ▼
  │                            cmd/api/main.go (Gin)
  │                                         │
  │                                         ▼
  │                            internal/api/ingest.go
  │                            CreateEvent()
  │                                         │
  │                                         ├─ Valida API key
  │                                         ├─ Obtiene agent_id del cliente
  │                                         │
  │                                         ▼
  │                            PostgreSQL INSERT INTO events
  │                            (processed_at = NULL)
  │                                         │
  │ {"event_id":"event-123","status":"received"}
  │<─────────────────────────────────────────┤
  │                                         │
```

**En PostgreSQL ahora hay:**

```sql
SELECT * FROM events WHERE processed_at IS NULL;

id         | agent_id     | type     | service | severity | processed_at
-----------|--------------|----------|---------|----------|-------------
event-123  | agent-abc    | app_down | api     | critical | NULL
```

---

### MOMENTO 4: Agent despierta (cada 30s, automáticamente)

**El agente NO sabe que hay un evento hasta que hace su tick:**

```
Tu Servidor (Agent corriendo en background)

┌─────────────────────────────────┐
│ cmd/agent/main.go               │
│                                 │
│ func main() {                   │
│   ticker := time.NewTicker(30s) │
│                                 │
│   for {                         │
│     <-ticker.C                  │ ← Cada 30 segundos
│     agent.tick()                │
│   }                             │
│ }                               │
└─────────────────────────────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│ internal/agent/agent.go                 │
│                                         │
│ func tick() {                           │
│   // 1. Chequear cooldown               │
│   if now < cooldown_until {             │
│     return  // Skip, en cooldown        │
│   }                                     │
│                                         │
│   // 2. LEER eventos pending            │
│   events = db.GetPendingEvents() ───────┼──┐
│                                         │  │
│   if len(events) == 0 {                 │  │
│     return  // Nada que hacer           │  │
│   }                                     │  │
│                                         │  │
│   // 3. Construir contexto              │  │
│   context = buildContext(events) ───────┼──┼──┐
│                                         │  │  │
│   // 4. Preguntar a Gemini              │  │  │
│   decision = gemini.Decide(context) ────┼──┼──┼──┐
│                                         │  │  │  │
│   // 5. Ejecutar acción                 │  │  │  │
│   executor.Execute(decision) ───────────┼──┼──┼──┼──┐
│ }                                       │  │  │  │  │
└─────────────────────────────────────────┘  │  │  │  │
                                             │  │  │  │
    Lee de PostgreSQL                        │  │  │  │
                         ◄───────────────────┘  │  │  │
                                                │  │  │
    Obtiene historial, config                   │  │  │
                         ◄──────────────────────┘  │  │
                                                   │  │
    Llama a Gemini API                             │  │
                         ◄─────────────────────────┘  │
                                                      │
    Llama webhook del cliente                         │
                         ◄────────────────────────────┘
```

---

### MOMENTO 5: Agent ejecuta acción (llama webhook del cliente)

**Agent decide hacer restart:**

```
 Servidor (Agent)                  Servidor del Cliente (SDK webhook)
  │                                         │
  │ Gemini decidió: "restart api"           │
  │         ↓                                │
  │ internal/agent/executor.go               │
  │ Execute(decision)                        │
  │         │                                │
  │         ├─ Obtener client webhook_url    │
  │         ├─ Obtener client webhook_secret │
  │         │                                │
  │         ├─ Construir payload:            │
  │         │   {                            │
  │         │     "action": "restart",       │
  │         │     "target": "api",           │
  │         │     "timestamp": 1704649200    │
  │         │   }                            │
  │         │                                │
  │         ├─ Calcular HMAC:                │
  │         │   signature = HMAC-SHA256(     │
  │         │       payload,                 │
  │         │       webhook_secret           │
  │         │   )                            │
  │         │                                │
  │         ▼                                │
  │ POST /webhooks/agent                     │
  │ X-Agent-Signature: abc123def...          │
  │ {"action":"restart","target":"api"}      │
  ├─────────────────────────────────────────>│
  │                                         │
  │                                         ▼
  │                          sdk/go/example_client.go
  │                          WebhookHandler()
  │                                         │
  │                                         ├─ Leer body
  │                                         ├─ Obtener signature del header
  │                                         │
  │                                         ├─ Calcular HMAC:
  │                                         │   expected = HMAC-SHA256(
  │                                         │       body,
  │                                         │       webhook_secret
  │                                         │   )
  │                                         │
  │                                         ├─ Validar:
  │                                         │   if signature != expected {
  │                                         │       return 401
  │                                         │   }
  │                                         │
  │                                         ├─ Parsear action
  │                                         │
  │                                         ▼
  │                          executeAction("restart", "api")
  │                                         │
  │                                         ▼
  │                          exec.Command("systemctl", "restart", "my-api")
  │                                         │
  │                                         ├─ La app se reinicia
  │                                         │
  │ {"ok": true}                            │
  │<─────────────────────────────────────────┤
  │                                         │
  │ Guardar resultado en DB                 │
  │ UPDATE events SET processed_at = NOW()  │
  │                                         │
```



---

### ✅ LO QUE SÍ PASA:

```
SDK del cliente ──(monitorea)──> Su propia API (localhost)
      │
      │ (si detecta problema)
      │
      └──(HTTP POST)──>  Ingest API
                       "Reporto: mi API está caída"
```

**El cliente SE MONITOREA A SÍ MISMO** y te reporta problemas.

---

## 📋 Resumen de Responsabilidades

| Componente | Ubicación | Responsabilidad | Cuándo actúa |
|------------|-----------|-----------------|--------------|
| **SDK (Monitor)** | Servidor del cliente | Chequear su propia app cada 30s | Automático (loop) |
| **SDK (Webhook)** | Servidor del cliente | Recibir y ejecutar acciones | Cuando tu agente lo llama |
| **Ingest API** |  servidor | Recibir eventos y guardar en DB | Cuando SDK reporta |
| **Agent** |  servidor | Leer eventos, decidir, ejecutar | Cada 30s (loop) |
| **Gemini** | API de Google | Analizar y decidir acción | Cuando agent pregunta |
| **PostgreSQL** |  servidor | Almacenar todo | Siempre |

---

## 🔄 Ciclo de Vida de un Evento

```
1. [T+0s] SDK: "Mi API no responde"
   └─> POST /api/events → PostgreSQL (processed_at = NULL)

2. [T+5s] Agent tick: "Hay eventos pending?"
   └─> SELECT * FROM events WHERE processed_at IS NULL
   └─> SÍ: event-123

3. [T+6s] Agent: "¿Qué hago?"
   └─> Gemini: "Analiza esto"
   └─> Gemini: "Respuesta: restart api"

4. [T+7s] Agent: "Ejecuto restart"
   └─> POST cliente.com/webhooks/agent
   └─> Cliente: "OK, reiniciando"

5. [T+8s] Agent: "Marco evento como procesado"
   └─> UPDATE events SET processed_at = NOW()

6. [T+9s] Agent: "Entro en cooldown 5 min"
   └─> UPDATE agents SET cooldown_until = NOW() + 5min
```

---

## ❓ Preguntas Frecuentes

### 1. ¿El agente sabe INMEDIATAMENTE cuando pasa algo?

**NO.** El agente solo se entera cada 30 segundos cuando hace su tick.

**Timeline:**
```
15:30:00 - Cliente detecta problema, reporta evento
15:30:05 - (agente dormido)
15:30:10 - (agente dormido)
15:30:15 - (agente dormido)
15:30:20 - (agente dormido)
15:30:25 - (agente dormido)
15:30:30 - ¡Agent tick! Lee evento, actúa
```

**Delay máximo:** 30 segundos
(SE PUEDE MIGRAR HACIA NOTIFY/LISTEN, REDIS, SQS O KAFKA)
---

### 2. ¿Por qué el cliente no llama directamente al agente?

**Porque NO QUEREMOS que el agente esté esperando requests TODO el tiempo.**

**Arquitectura basada en eventos:**
- Cliente → Reporta evento → DB (asíncrono)
- Agent → Lee DB cuando está listo → Actúa

**Ventajas:**
- ✅ Agent puede estar offline temporalmente
- ✅ Eventos se acumulan en DB
- ✅ Agent procesa en batch
- ✅ Más escalable

---

### 4. ¿Qué pasa si el agent está apagado?

**Los eventos se acumulan:**

```sql
-- Events sin procesar
SELECT * FROM events WHERE processed_at IS NULL;

id         | created_at           | type
-----------|----------------------|----------
event-123  | 2026-01-09 15:30:00  | app_down
event-124  | 2026-01-09 15:31:00  | app_down
event-125  | 2026-01-09 15:32:00  | app_down
```

Cuando el agent vuelve a encenderse, los procesa todos.

---

### 5. ¿El SDK chequea TODA la infraestructura del cliente?

**Depende de cómo lo configure el cliente.**

**Ejemplo básico:**
```go
// Solo chequea la API
if !isHealthy("http://localhost:8080/health") {
    sdk.ReportEvent("app_down", "api", "critical")
}
```

**Ejemplo avanzado:**
```go
// Chequea múltiples servicios
for _, service := range []string{"api", "worker", "cron"} {
    if !isHealthy(service) {
        sdk.ReportEvent("app_down", service, "critical")
    }
}

// Chequea DB
if !isDatabaseHealthy() {
    sdk.ReportEvent("db_down", "postgres", "critical")
}

// Chequea métricas
if getCPU() > 90 {
    sdk.ReportEvent("high_cpu", "api", "warning", {"cpu": 95})
}
```

**El cliente decide qué monitorear.**

---

## 🎯 Arquitectura en UNA Imagen

```
┌──────────────────────────────────────────────────────────────┐
│                  SERVIDOR DEL CLIENTE                         │
│                                                               │
│  App (su negocio)  ◄──monitorea── SDK ──reporta──┐          │
│      :8080                          :9000          │          │
│                                       ▲            │          │
│                                       │            │          │
│                                  recibe webhooks   │          │
│                                  (para ejecutar    │          │
│                                   acciones)        │          │
└───────────────────────────────────────┼────────────┼──────────┘
                                        │            │
                                        │            │ HTTP POST
                                        │            │ /api/events
                                        │            │
┌───────────────────────────────────────┼────────────▼──────────┐
│                     SERVIDOR                                 │
│                                       │                        │
│  ┌──────────────────────────────┐    │                        │
│  │ Ingest API (Gin)              │◄───┘                        │
│  │ Recibe eventos                │                            │
│  └──────────┬───────────────────┘                            │
│             │ guarda                                          │
│             ▼                                                  │
│  ┌─────────────────────────┐                                  │
│  │ PostgreSQL               │                                  │
│  │ events (pending)         │                                  │
│  └──────────┬──────────────┘                                  │
│             │ lee (cada 30s)                                   │
│             ▼                                                  │
│  ┌─────────────────────────┐                                  │
│  │ Agent Loop               │                                  │
│  │ 1. Lee eventos           │                                  │
│  │ 2. Pregunta a Gemini ────┼──> Gemini API                   │
│  │ 3. Ejecuta acción ───────┼──────────────────────┐          │
│  └─────────────────────────┘                       │          │
└────────────────────────────────────────────────────┼──────────┘
                                                     │
                                         llama webhook
                                                     │
                                                     └─ HTTP POST
                                                        /webhooks/agent
```

---



**El flujo es:**

1. Cliente se monitorea a sí mismo (SDK)
2. Cliente reporta problemas a tu API (Ingest)
3.  API guarda en DB
4.  Agent lee DB cada 30s
5. Agent pregunta a LLM qué hacer
6. Agent ejecuta acción llamando webhook del cliente
7. Cliente ejecuta la acción localmente


