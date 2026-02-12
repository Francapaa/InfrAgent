"use client";

import { useState } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';

// --- TIPOS (Sin cambios) ---
export type AgentStatus = 'idle' | 'monitoring' | 'analyzing' | 'executing' | 'error';
export type AgentAction = {
  id: string; type: string; description: string;
  status: 'pending' | 'running' | 'completed' | 'failed';
  timestamp: string; duration?: number;
};
type Metrics = { cpuUsage: number; memoryUsage: number; activeConnections: number; errorsDetected: number; };

interface BackendData {
  agents: any[]; events: any[]; actions: AgentAction[];
  status: AgentStatus; currentTask: string; metrics: Metrics; timestamp?: string;
}

export function useWebSocketManager(url: string) {
  const [data, setData] = useState<BackendData | null>(null);
  const [specificError, setSpecificError] = useState<string | null>(null);

  const { sendJsonMessage, readyState } = useWebSocket(url, {
    share: true,
    shouldReconnect: () => true,
    reconnectAttempts: 10,
    reconnectInterval: 5000, // Aumentamos a 5s para no saturar al reintentar
    heartbeat: {
      message: JSON.stringify({ type: 'ping' }),
      interval: 20000, // 20 segundos
      timeout: 60000,  // Si en 60s no hay respuesta, considera conexión muerta
    },
    onOpen: () => {
      console.log('%c[WS] ✅ Conectado!', 'color: green; font-weight: bold');
      setSpecificError(null);
    },

    onClose: (event) => {
      let reason = "Desconocido";
      
      switch (event.code) {
        case 1000: reason = "Cierre normal (el servidor terminó la tarea)"; break;
        case 1001: reason = "El servidor se está apagando (Going Away)"; break;
        case 1002: reason = "Error de protocolo (Go no entendió algo)"; break;
        case 1003: reason = "Tipo de dato no aceptado (Go esperaba texto y enviaste binario)"; break;
        case 1005: reason = "No se recibió código de estado (Cierre abrupto)"; break;
        case 1006: reason = "Conexión cerrada anormalmente (Posible Firewall o Crash)"; break;
        case 1011: reason = "Error interno del servidor (Go lanzó un Panic)"; break;
        default: reason = `Código ${event.code}: ${event.reason || 'Sin razón específica'}`;
      }

      console.warn(`%c[WS] ❌ Desconectado: ${reason}`, 'color: orange; font-weight: bold');
      setSpecificError(reason);
    },

    onError: (event) => {
      console.error('%c[WS] 🔥 Error de red/socket:', 'color: red;', event);
      setSpecificError("Error de red o servidor no alcanzable");
    },

    onMessage: (event) => {
      try {
        console.log(event.data)
        const payload = JSON.parse(event.data);
        console.log(payload); 
        setData(payload);
      } catch (e) {
        console.error("[WS] Mensaje no es JSON válido:", event.data);
      }
    },
  });

  const connectionStatus = {
    [ReadyState.CONNECTING]: 'Conectando...',
    [ReadyState.OPEN]: 'Conectado',
    [ReadyState.CLOSING]: 'Cerrando...',
    [ReadyState.CLOSED]: 'Desconectado',
    [ReadyState.UNINSTANTIATED]: 'No iniciado',
  }[readyState];

  return {
    state: {
      agents: data?.agents || [],
      events: data?.events || [],
      actions: data?.actions || [],
      status: data?.status || 'idle',
      currentTask: data?.currentTask || '',
      metrics: data?.metrics || { cpuUsage: 0, memoryUsage: 0, activeConnections: 0, errorsDetected: 0 },
      lastUpdate: data?.timestamp || new Date().toISOString(),
    },
    isConnected: readyState === ReadyState.OPEN,
    connectionStatus,
    connectionError: specificError, 
    sendCommand: (command: 'pause' | 'resume') => {
      sendJsonMessage({ action: command, timestamp: new Date().toISOString() });
    },
  };
}