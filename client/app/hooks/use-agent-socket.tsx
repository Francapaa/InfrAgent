"use client"

// custom hook para hacer la conexion WS
import { useEffect, useReducer, useCallback, useRef } from 'react';

type Agent = {
  id: string;
  client_id: string;
  state: string;
  last_tick_at?: string;
  cooldown_until: string;
  created_at: string;
  updated_at: string;
};

type Event = {
  id: string;
  client_id: string;
  agent_id: string;
  type: string;
  service: string;
  severity: string;
  data: Record<string, any>;
  processed_at?: string;
  created_at: string;
};

type Action = {
  id: string;
  agent_id: string;
  client_id: string;
  type: string;
  target: string;
  params: Record<string, any>;
  reasoning: string;
  confidence: number;
  status: string;
  result: Record<string, any>;
  executed_at?: string;
  created_at: string;
};

type WSMessage = {
  type: 'agent_update' | 'event_created' | 'action_taken' | 'agent_state_changed';
  data: any;
};

// Estado del reducer
type State = {
  agents: Agent[];
  events: Event[];
  actions: Action[];
  isConnected: boolean;
  error: string | null;
};

// Acciones del reducer
type ReducerAction =
  | { type: 'SET_CONNECTED'; payload: boolean }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'UPDATE_AGENT'; payload: Agent }
  | { type: 'ADD_EVENT'; payload: Event }
  | { type: 'UPDATE_ACTION'; payload: Action }
  | { type: 'RESET' };

// Estado inicial
const initialState: State = {
  agents: [],
  events: [],
  actions: [],
  isConnected: false,
  error: null,
};

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
        return { ...state, agents: newAgents };
      }
      return { ...state, agents: [...state.agents, action.payload] };
    }

    case 'ADD_EVENT':
      return {
        ...state,
        events: [action.payload, ...state.events].slice(0, 100), // Limitar a 100
      };

    case 'UPDATE_ACTION': {
      const index = state.actions.findIndex(a => a.id === action.payload.id);
      if (index >= 0) {
        const newActions = [...state.actions];
        newActions[index] = action.payload;
        return { ...state, actions: newActions };
      }
      return {
        ...state,
        actions: [action.payload, ...state.actions].slice(0, 100),
      };
    }

    case 'RESET':
      return initialState;

    default:
      return state;
  }
}

type UseWebSocketReturn = {
  agents: Agent[];
  events: Event[];
  actions: Action[];
  isConnected: boolean;
  error: string | null;
  sendMessage: (message: any) => void;
};

export function useWebSocket(url: string | undefined = process.env.NEXT_PUBLIC_WS_URL): UseWebSocketReturn | undefined {
  
  if (!url){
    console.log("no tenemos la url del backend con ws")
    return
  }

    const [state, dispatch] = useReducer(websocketReducer, initialState);
  const wsRef = useRef<WebSocket | null>(null);

  const sendMessage = useCallback((message: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify(message));
    }
  }, []);

  useEffect(() => {
    const ws = new WebSocket(url);
    wsRef.current = ws;

    ws.onopen = () => {
      console.log('WebSocket conectado');
      dispatch({ type: 'SET_CONNECTED', payload: true });
    };

    ws.onmessage = (event) => {
      try {
        const message: WSMessage = JSON.parse(event.data); //data que llega del back
        
        switch (message.type) {
          case 'agent_update':
            dispatch({ type: 'UPDATE_AGENT', payload: message.data });
            break;
            
          case 'event_created':
            dispatch({ type: 'ADD_EVENT', payload: message.data });
            break;
            
          case 'action_taken':
            dispatch({ type: 'UPDATE_ACTION', payload: message.data });
            break;
            
          case 'agent_state_changed':
            console.log('Agent state changed:', message.data);
            break;
        }//por cada caso cambiamos el estado con la nueva data
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
    agents: state.agents,
    events: state.events,
    actions: state.actions,
    isConnected: state.isConnected,
    error: state.error,
    sendMessage,
  };
}