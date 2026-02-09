"use client"

// Hook MÍNIMO para debuggear el error 1005
import { useEffect, useState, useCallback } from 'react';

export type AgentStatus = 'idle' | 'monitoring' | 'analyzing' | 'executing' | 'error';

export type AgentAction = {
  id: string;
  type: string;
  description: string;
  status: 'pending' | 'running' | 'completed' | 'failed';
  timestamp: string;
  duration?: number;
};

type Metrics = {
  cpuUsage: number;
  memoryUsage: number;
  activeConnections: number;
  errorsDetected: number;
};

export type UseWebSocketReturn = {
  state: {
    agents: never[];
    events: never[];
    actions: AgentAction[];
    status: AgentStatus;
    currentTask: string;
    metrics: Metrics;
    lastUpdate?: string;
  };
  isConnected: boolean;
  connectionError: string | null;
  reconnectAttempt: number;
  sendCommand: (command: 'pause' | 'resume') => void;
};

export function useWebSocket(url: string | undefined = process.env.NEXT_PUBLIC_WS_URL): UseWebSocketReturn {
  const [isConnected, setIsConnected] = useState(false);
  const [connectionError, setConnectionError] = useState<string | null>(null);
  const [reconnectAttempt, setReconnectAttempt] = useState(0);

  useEffect(() => {
    if (!url) {
      console.log('[WS] No URL provided');
      return;
    }

    console.log('[WS] Effect running, connecting to:', url);
    
    const ws = new WebSocket(url);
    
    ws.onopen = () => {
      console.log('[WS] Connected');
      setIsConnected(true);
      setConnectionError(null);
      setReconnectAttempt(0);
    };
    
    ws.onclose = (event) => {
      console.log(`[WS] Closed - Code: ${event.code}, Reason: "${event.reason}", Clean: ${event.wasClean}`);
      setIsConnected(false);
      
      if (event.code === 1005) {
        console.warn('[WS] Error 1005 - No status code received. This usually means the connection was closed unexpectedly.');
      }
    };
    
    ws.onerror = (error) => {
      console.error('[WS] Error:', error);
      setConnectionError('WebSocket error');
    };
    
    ws.onmessage = (event) => {
      console.log('[WS] Message received:', event.data);
    };

    return () => {
      console.log('[WS] Cleanup - closing connection');
      if (ws.readyState === WebSocket.OPEN || ws.readyState === WebSocket.CONNECTING) {
        ws.close(1000, 'Component unmounting');
      }
    };
  }, [url]); // Solo se ejecuta cuando cambia la URL

  const sendCommand = useCallback((command: 'pause' | 'resume') => {
    console.log('[WS] Send command:', command);
    // Implementar cuando sea necesario
  }, []);

  return {
    state: {
      agents: [],
      events: [],
      actions: [],
      status: 'idle',
      currentTask: '',
      metrics: {
        cpuUsage: 0,
        memoryUsage: 0,
        activeConnections: 0,
        errorsDetected: 0,
      },
      lastUpdate: new Date().toISOString(),
    },
    isConnected,
    connectionError,
    reconnectAttempt,
    sendCommand,
  };
}
