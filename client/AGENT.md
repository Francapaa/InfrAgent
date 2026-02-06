# DEVOPS AGENT


-EL DEVOPS AGENT ES UN AGENTE EL CUAL PUEDE EJECUTAR ACCIONES EN LOS PROYECTOS DE LAS PERSONAS. COMO? CON UN SDK

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
│  │      A Monitorea la app cada 30s                       │ │
│  │      B Expone webhook para recibir acciones            │ │
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


# Que ocurre en la aplicacion? 

-El frontend contiene la parte del login (/login) en la cual los usuarios pueden loguearse para asi poder usar nuestras funcionalidades
-Cuando el usuario/cliente se registra en nuestra plataforma (por primera vez) se le da: Client_id , api_key (tiene que guardarlo) y webHookSecret(tiene que guardarlo)
-En nuestra BD postgreSQL creo el usuario y ademas el agente; ALGO MUY IMPORTANTE 1 AGENTE CADA 1 USUARIO
-Luego de registrarse y obtener los datos que le damos (api_key y webHookSecret) el usuario tiene que hacer esto (ejemplo en GO):

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

# La pregunta, que es el SDK? 

Basicamente el sdk cada 30s corre un tick el cual detecta un problema (GRACIAS A LA ARQUITECTURA EVENT DRIVEN)

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


  EN POSTGRESQL AHORA HAY:
  SELECT * FROM events WHERE processed_at IS NULL;

id         | agent_id     | type     | service | severity | processed_at
-----------|--------------|----------|---------|----------|-------------
event-123  | agent-abc    | app_down | api     | critical | NULL


EL AGENTE NO SABE QUE HAY UN PROBLEMA HASTA CORRER EL TICK 

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


# AHORA LO MAS IMPORTANTE, QUE TIENE QUE HACER EL FRONTEND? 


1) Permitir el login/registro gracias a las herramientas de GOOGLE
2) Permitir el Registro en la plataforma de manera LOCAL
3) LO MAS IMPORTANTE: Poder integrar el dashboard (/dashboard) con la arquitectura webSockets para asi poder ver los datos en tiempo REAL. Ademas , poder ver las decisiones que tomó el AGENTE. 
4) Ademas tiene que haber (esquina superior de la derecha) un logo del usuario para que pueda ver los datos MAS IMPORTANTES DE SU PERFIL





# BUENAS PRÁCTICAS QUE TENÉS QUE SEGUIR (OBLIGATORIAS)

1) TODO EL CÓDIGO TIENE QUE ESTAR TIPADO.  
   SIEMPRE, SIN EXCEPCIÓN ALGUNA.

2) NUNCA USAR EL TIPO `any`.  
   SI NO SABÉS EL TIPO, CREÁ UNO CORRECTO O INFERILO.

4) EVITAR COMPONENTES MONOLÍTICOS.  
   EL CÓDIGO DE UNA RUTA DEBE ESTAR COMPUESTO POR COMPONENTES MÁS CHICOS PARA MEJOR MANTENIBILIDAD Y LEGIBILIDAD.

5) NO INSTALAR NINGUNA LIBRERÍA SIN CONSULTAR PREVIAMENTE AL PROGRAMADOR RESPONSABLE.

6) EVITAR EL USO EXCESIVO DE `useState`.  
   - NO CREAR UN `useState` POR CADA PIEZA DE LÓGICA.  
   - AGRUPAR ESTADO Y LÓGICA RELACIONADA EN OBJETOS.  
   - CUANDO EL ESTADO SEA COMPLEJO O TENGA MUCHAS TRANSICIONES, USAR `useReducer`.

7) SEPARACIÓN DE RESPONSABILIDADES (OBLIGATORIO).
   - LAS LLAMADAS AL BACKEND NO DEBEN ESTAR EN `page.tsx` NI EN COMPONENTES DE UI.  
   - CREAR ARCHIVOS/SERVICIOS DEDICADOS PARA FETCH, API CLIENTS O ACTIONS.
   - LA LÓGICA DE COOKIES Y AUTENTICACIÓN (MIDDLEWARE DE NEXT) DEBE VIVIR EN ARCHIVOS SEPARADOS, NUNCA MEZCLADA CON UI.

8) MANEJO DE AUTENTICACIÓN Y SEGURIDAD.
   - EL TOKEN JWT SE GUARDA ÚNICAMENTE EN COOKIES.  
   - NUNCA USAR `localStorage` PARA TOKENS DE AUTENTICACIÓN.  
   - ESTO ES UNA DECISIÓN DE SEGURIDAD Y NO ES OPCIONAL.

9) PRIORIZAR CÓDIGO CLARO, PREDECIBLE Y ESCALABLE.
   - LA LÓGICA RELACIONADA DEBE VIVIR JUNTA.  
   - EVITAR SIDE EFFECTS INESPERADOS.  
   - SI ALGO SE PUEDE HACER SIMPLE, HACERLO SIMPLE.


CUALQUIER DUDA SOBRE IMPLEMENTACIÓN, ARQUITECTURA O LIBRERÍAS  
DEBE SER CONSULTADA ANTES DE ESCRIBIR CÓDIGO.



# STACK A UTILIZAR A LA HORA DE GENERAR CODIGO

- FRONTEND: 
1) NEXT.JS PARA GENERAR RUTAS Y COMPONENTES
2) PARA GENERAR ESTILOS UTILIZAR TAILWINDCSS 





