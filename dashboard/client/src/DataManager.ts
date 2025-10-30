// webhook$$bc-go;dashboard/client/src/DataManager.ts;grok$$
import { create } from 'zustand';
import { w3cwebsocket as WebSocket } from 'websocket';
import { produce } from "immer";

export interface Ticker {
  x: string; // exchange
  p: string; // pair
  b: number; // bid
  bs: number; // bid size
  a: number; // ask 
  as: number; // ask size
  et: number; // event time
}

export interface Alpha {
  i: number; // id
  p: string; // pair
  s: number; // spread
  sr: number; // spread ratio
  ss: number; // spread size
  ax: string; // ask exchange
  bx: string; // bid exchange
  ap: number; // ask price
  bp: number; // bid price
  as: number; // ask size
  bs: number; // bid size
  at: number; // ask time
  bt: number; // bid time
  af: number; // ask fee
  bf: number; // bid fee
  pl: number; // p/l
}

interface DashboardDataState {
  alpha: Alpha[];
  minSpreadPct: number;
  status: 'connected' | 'disconnected' | 'retry';
  reconnect: () => void;
  disconnect: () => void;
  setMinSpreadPct: (val: number) => void;
}

const API_HOST = process.env.NODE_ENV === 'production' ? document.location.host : 'localhost:8080';

export const useDataStore = create<DashboardDataState>((set, get) => {
  let ws: WebSocket | null = null;

  const setupWebSocket = () => {
    ws = new WebSocket(`ws://${API_HOST}/ws`);
    
    ws.onopen = () => {
      set(state => ({ ...state, status: 'connected' }));
    };
    
    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data.toString())
        set(produce((draft: DashboardDataState) => {
          draft.alpha.unshift(message);
          if (draft.alpha.length > 200) {
            draft.alpha.pop()
          }
        }));
      } catch (e) {
        console.error('WebSocket message parse error:', e);
      }
    };
    
    ws.onclose = (event) => {
      set(state => ({ ...state, status: 'disconnected' }));
      console.log('WebSocket disconnected:', event.code, event.reason);
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error)
      ws?.close();
    };
  };

  const reconnect = async () => {
    if (get().status === 'connected') return;

    set(state => ({ ...state, status: 'retry' }));
    try {
      ws?.close();
      setupWebSocket();
      await new Promise(resolve => setTimeout(resolve, 1000)); // Wait for connection
      if (get().status !== 'connected') {
        set(state => ({ ...state, status: 'disconnected' }));
      }
    } catch (error) {
      console.error('Reconnect failed:', error);
      set(state => ({ ...state, status: 'disconnected' }));
    }
  };

  const disconnect = async () => {
    ws?.close();
  };

  const setMinSpreadPct = (newMinSpreadPct: number) => {
    set(produce((draft: DashboardDataState) => {
      draft.minSpreadPct = newMinSpreadPct
    }));
  }

  setupWebSocket();

  return { 
    alpha: [],
    minSpreadPct: 0.025,
    status: 'disconnected',
    reconnect,
    disconnect,
    setMinSpreadPct,
  };
});