// webhook$$bc-go;dashboard/client/src/DataManager.ts;grok$$
import { create } from 'zustand';
import { w3cwebsocket as WebSocket } from 'websocket';
import type { IMessageEvent, ICloseEvent } from 'websocket';
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
  id: number;
  ask: Ticker;
  bid: Ticker;
  p: string; // pair
  s: number; // spread
  sr: number; // spread ratio
  ss: number; // spread size
}

interface DashboardDataState {
  alpha: Alpha[];
  exchanges: string[];
  pairs: string[];
  status: 'connected' | 'disconnected' | 'retry';
  reconnect: () => void;
}

const API_HOST = process.env.NODE_ENV === 'production' ? document.location.host : 'localhost:8080';

export const useDataStore = create<DashboardDataState>((set, get) => {
  let ws: WebSocket | null = null;

  const setupWebSocket = () => {
    ws = new WebSocket(`ws://${API_HOST}/ws`);
    
    ws.onopen = () => {
      set(state => ({ ...state, status: 'connected' }));
    };
    
    ws.onmessage = (event: IMessageEvent) => {
      try {
        const message = JSON.parse(event.data)
        set(produce((draft: DashboardDataState) => {
          if (!draft.alpha) {
            draft.alpha = []
          }
          draft.alpha.unshift(message as Alpha);
          if (draft.alpha.length > 20) {
            draft.alpha.pop()
          }
        }));
      } catch (e) {
        console.error('WebSocket message parse error:', e);
      }
    };
    
    ws.onclose = (event: ICloseEvent) => {
      set(state => ({ ...state, status: 'disconnected' }));
      console.log('WebSocket disconnected:', event.code, event.reason);
    };

    ws.onerror = (error: Error) => {
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

  setupWebSocket();

  return { 
    alpha: [],
    exchanges: [],
    pairs: [], 
    status: 'disconnected',
    reconnect,
  };
});