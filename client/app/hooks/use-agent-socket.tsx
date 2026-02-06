"use client"

// custom hook para hacer la conexion WS
import { useEffect, useReducer, useCallback, useRef } from 'react';

// Tipo para datos JSON genéricos sin usar any
type JSONValue = string | number | boolean | null | JSONObject | JSONArray;
interface JSONObject {
  [key: string]: JSONValue;
}
interface JSONArray extends Array<JSONValue> {}

export type AgentStatus = 
  | 'idle' 
  | 'monitoring' 
  | 'analyzing' 
  | 'executing' 
  | 'error';

export type AgentAction = {
  id: string;
  type: string;
  description: string;
  status: 'pending' | 'running' | 'completed' | 'failed';
  timestamp: string;
  duration?: number;
};

type Agent = {
  id: string;
  client_id: string;
  state: AgentStatus;
  last_tick_at?: string;
  cooldown_until: string;
  created_at: string;
  updated_at: string;
};

type EventData = Record<string, JSONValue>;

type Event = {
  id: string;
  client_id: string;
  agent_id: string;
  type: string;
  service: string;
  severity: string;
  data: EventData;
  processed_at?: string;
  created_at: string;
};

type ActionResult = Record<string, JSONValue>;

type ActionParams = Record<string, JSONValue>;

type Action = {
  id: string;
  agent_id: string;
  client_id: string;
  type: string;
  target: string;
  params: ActionParams;
  reasoning: string;
  confidence: number;
  status: string;
  result: ActionResult;
  executed_at?: string;
  created_at: string;
};

type WSMessageData = Agent | Event | Action | JSONObject;

type WSMessage = {
  type: 'agent_update' | 'event_created' | 'action_taken' | 'agent_state_changed';
  data: WSMessageData;
};

// Estado del reducer
type State = {
  agents: Agent[];
  events: Event[];
  actions: AgentAction[];
  isConnected: boolean;
  error: string | null;
  status: AgentStatus;
  currentTask: string;
  metrics: {
    cpuUsage: number;
    memoryUsage: number;
    activeConnections: number;
    errorsDetected: number;
  };
};

// Acciones del reducer
type ReducerAction =
  | { type: 'SET_CONNECTED'; payload: boolean }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'UPDATE_AGENT'; payload: Agent }
  | { type: 'ADD_EVENT'; payload: Event }
  | { type: 'UPDATE_ACTION'; payload: Action }
  | { type: 'SET_STATUS'; payload: AgentStatus }
  | { type: 'SET_CURRENT_TASK'; payload: string }
  | { type: 'SET_METRICS'; payload: State['metrics'] }
  | { type: 'RESET' };

// Estado inicial
const initialState: State = {
  agents: [],
  events: [],
  actions: [],
  isConnected: false,
  error: null,
  status: 'idle',
  currentTask: '',
  metrics: {
    cpuUsage: 0,
    memoryUsage: 0,
    activeConnections: 0,
    errorsDetected: 0,
  },
};

// Helper para convertir Action backend a AgentAction frontend
function mapBackendActionToFrontend(action: Action): AgentAction {
  return {
    id: action.id,
    type: action.type,
    description: action.reasoning || `${action.type} on ${action.target}`,
    status: action.status as AgentAction['status'],
    timestamp: action.executed_at || action.created_at,
    duration: action.result?.duration as number | undefined,
  };
}

// Reducer
function websocketReducer(state: State, action: ReducerAction): State {
  switch (action.type) {
    case 'SET_CONNECTED':
      return {
        ...state,
        isConnected: action.payload,
        error: action.payload ? null : state.error,
      };

    case 'SET_ERROR':
      return {
        ...state,
        error: action.payload,
        isConnected: false,
      };

    case 'UPDATE_AGENT': {
      const index = state.agents.findIndex(a => a.id === action.payload.id);
      if (index >= 0) {
        const newAgents = [...state.agents];
        newAgents[index] = action.payload;
        return { 
          ...state, 
          agents: newAgents,
          status: action.payload.state,
        };
      }
      return { 
        ...state, 
        agents: [...state.agents, action.payload],
        status: action.payload.state,
      };
    }

    case 'ADD_EVENT':
      return {
        ...state,
        events: [action.payload, ...state.events].slice(0, 100), // Limitar a 100
        metrics: {
          ...state.metrics,
          errorsDetected: action.payload.severity === 'critical' 
            ? state.metrics.errorsDetected + 1 
            : state.metrics.errorsDetected,
        },
      };

    case 'UPDATE_ACTION': {
      const frontendAction = mapBackendActionToFrontend(action.payload);
      const index = state.actions.findIndex(a => a.id === frontendAction.id);
      if (index >= 0) {
        const newActions = [...state.actions];
        newActions[index] = frontendAction;
        return { ...state, actions: newActions };
      }
      return {
        ...state,
        actions: [frontendAction, ...state.actions].slice(0, 100),
      };
    }

    case 'SET_STATUS':
      return {
        ...state,
        status: action.payload,
      };

    case 'SET_CURRENT_TASK':
      return {
        ...state,
        currentTask: action.payload,
      };

    case 'SET_METRICS':
      return {
        ...state,
        metrics: action.payload,
      };

    case 'RESET':
      return initialState;

    default:
      return state;
  }
}

type CommandMessage = {
  type: 'command';
  command: 'pause' | 'resume';
};

type WebSocketMessage = CommandMessage | JSONObject;

export type UseWebSocketReturn = {
  state: {
    agents: Agent[];
    events: Event[];
    actions: AgentAction[];
    status: AgentStatus;
    currentTask: string;
    metrics: {
      cpuUsage: number;
      memoryUsage: number;
      activeConnections: number;
      errorsDetected: number;
    };
    lastUpdate?: string;
  };
  isConnected: boolean;
  connectionError: string | null;
  sendCommand: (command: 'pause' | 'resume') => void;
};

export function useWebSocket(url: string | undefined = process.env.NEXT_PUBLIC_WS_URL): UseWebSocketReturn {
  
  const actualUrl = url || '';

  const [state, dispatch] = useReducer(websocketReducer, initialState);
  const wsRef = useRef<WebSocket | null>(null);

  const sendMessage = useCallback((message: WebSocketMessage) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    }
  }, []);

  const sendCommand = useCallback((command: 'pause' | 'resume') => {
    sendMessage({ type: 'command', command });
    dispatch({ 
      type: 'SET_CURRENT_TASK', 
      payload: command === 'pause' ? 'Agente pausado' : 'Monitoreando servicios' 
    });
  }, [sendMessage]);

  // Si no hay URL, retornar estado por defecto sin conexión
  if (!actualUrl) {
    console.log("no tenemos la url del backend con ws");
    return {
      state: {
        agents: [],
        events: [],
        actions: [],
        status: 'idle' as AgentStatus,
        currentTask: '',
        metrics: {
          cpuUsage: 0,
          memoryUsage: 0,
          activeConnections: 0,
          errorsDetected: 0,
        },
        lastUpdate: new Date().toISOString(),
      },
      isConnected: false,
      connectionError: 'No se configuró la URL del WebSocket',
      sendCommand,
    };
  }

  useEffect(() => {
    const ws = new WebSocket(actualUrl);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log('WebSocket conectado');
      dispatch({ type: 'SET_CONNECTED', payload: true });
      dispatch({ type: 'SET_CURRENT_TASK', payload: 'Monitoreando servicios de producción' });
    };

    ws.onmessage = (event) => {
      try {
        const message: WSMessage = JSON.parse(event.data);
        
        switch (message.type) {
          case 'agent_update':
            dispatch({ type: 'UPDATE_AGENT', payload: message.data as Agent });
            break;
            
          case 'event_created':
            dispatch({ type: 'ADD_EVENT', payload: message.data as Event });
            break;
            
          case 'action_taken':
            dispatch({ type: 'UPDATE_ACTION', payload: message.data as Action });
            break;
            
          case 'agent_state_changed':
            console.log('Agent state changed:', message.data);
            if (typeof message.data === 'object' && message.data !== null) {
              const data = message.data as JSONObject;
              if ('state' in data && typeof data.state === 'string') {
                dispatch({ type: 'SET_STATUS', payload: data.state as AgentStatus });
              }
            }
            break;
        }
      } catch (err) {
        console.error('Error parsing message:', err);
      }
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
      dispatch({ type: 'SET_ERROR', payload: 'Error de conexión' });
    };

    ws.onclose = () => {
      console.log('WebSocket desconectado');
      dispatch({ type: 'SET_CONNECTED', payload: false });
    };

    return () => {
      ws.close();
      dispatch({ type: 'RESET' });
    };
  }, [url]);

  return {
    state: {
      agents: state.agents,
      events: state.events,
      actions: state.actions,
      status: state.status,
      currentTask: state.currentTask,
      metrics: state.metrics,
      lastUpdate: new Date().toISOString(),
    },
    isConnected: state.isConnected,
    connectionError: state.error,
    sendCommand,
  };
}
