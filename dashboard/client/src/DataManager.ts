// webhook$$bc-go;dashboard/client/src/DataManager.ts;grok$$
import { create } from 'zustand';
import { w3cwebsocket as WebSocket } from 'websocket';
import type { IMessageEvent, ICloseEvent } from 'websocket';
import { produce } from "immer";

export interface Ticker {
  e: string; // exchange
  s: string; // symbol
  p: number; // average price
}

interface DashboardDataState {
  pairs: Record<string,boolean>;
  exchanges: Record<string,boolean>;
  tickers: Record<string,Record<string,Ticker>>;
  status: 'connected' | 'disconnected' | 'retry';
  reconnect: () => void;
}

const API_HOST = process.env.NODE_ENV === 'production' ? document.location.host : 'localhost:8080';

export const useDataStore = create<DashboardDataState>((set, get) => {
  let ws: WebSocket | null = null;

  const fetchLists = async () => {
    const { exchanges, pairs } = await fetch(`http://${API_HOST}/lists`).then(resp => resp.json());
    set(state => ({ ...state, pairs, exchanges }));
  };

  const setupWebSocket = () => {
    ws = new WebSocket(`ws://${API_HOST}/ws`);
    
    ws.onopen = () => {
      set(state => ({ ...state, status: 'connected' }));
    };
    
    ws.onmessage = (event: IMessageEvent) => {
      try {
        const message = JSON.parse(event.data)
        set(produce((draft: DashboardDataState) => {
        //   if (isQuoteMsg(data)) {
        //     const q = data as Quote;
        //     const k = `${q.p}.${q.x}`
        //     draft.tickers[k] = q
        //   } else if (isTradeMsg(data)) {
        //     draft.trades.unshift(data as Trade)
        //     if (draft.trades.length > 20) draft.trades.pop()            
        //   } else if (isMetricMsg(data)) {
        //     draft.metrics.unshift(data as Metric)
        //     if (draft.metrics.length > 20) draft.metrics.pop()
        //   }
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
    pairs: {},
    exchanges: {},
    tickers: {}, 
    status: 'disconnected',
    reconnect,
  };
});