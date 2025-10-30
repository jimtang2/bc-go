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
};
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
};
interface AlphaDataState {
  alpha: Alpha[];
  status: 'connected' | 'disconnected';
  connect: () => void;
  reconnect: () => void;
  disconnect: () => void;
};
const API_HOST = process.env.NODE_ENV === 'production' ? document.location.host : 'localhost:8080';
export const useAlphaDataStore = create<AlphaDataState>((set, get) => {
  let ws: WebSocket | null = null;
  function connect() {
    // console.log(`DataManager.connect: start: readyState=${ws?.readyState}; status=${get().status}"`)
    if (ws?.readyState === 0 || ws?.readyState === 1) {
      return
    }
    // console.log("DataManager.connect: connecting")
    ws = new WebSocket(`ws://${API_HOST}/ws/alpha`);
    ws.onopen = () => {
      // console.log("DataManager.onopen")
      set(state => ({ ...state, status: 'connected' }));
    };
    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data.toString())
        set(produce((draft: AlphaDataState) => {
          draft.alpha.unshift(message);
          if (draft.alpha.length > 1000) {
            draft.alpha.pop()
          }
        }));
      } catch (e) {
        console.error('WebSocket message parse error:', e);
      }
    };
    ws.onclose = () => {
      set(state => ({ ...state, status: 'disconnected' }));
    };
    // ws.onerror = (error) => {
    //   console.log('ws.onerror:', error)
    // };
  };
  async function reconnect() {
    try {
      ws?.close();
      connect();
      await new Promise(resolve => setTimeout(resolve, 1000)); // Wait for connection
      if (get().status !== 'connected') {
        set(state => ({ ...state, status: 'disconnected' }));
      }
    } catch (error) {
      console.error('Reconnect failed:', error);
      set(state => ({ ...state, status: 'disconnected' }));
    }
  };
  async function disconnect() {
    // console.log("DataManager.disconnect")
    ws?.close();
    set(state => ({ ...state, status: 'disconnected' }));
  };
  return { 
    alpha: [],
    status: 'disconnected',
    connect,
    reconnect,
    disconnect,
  };
});
interface HistoryDataState {
  status: 'connected' | 'disconnected';
  reconnect: () => void;
  disconnect: () => void;
};
export const useHistoryDataStore = create<HistoryDataState>((set, get) => {
  let ws: WebSocket | null = null;
  function connect() {
    ws = new WebSocket(`ws://${API_HOST}/ws/history`);
    ws.onopen = () => {
      set(state => ({ ...state, status: 'connected' }));
    };
    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data.toString())
        console.log(message)
        set(produce((draft: HistoryDataState) => {
          console.log(draft)

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
  async function reconnect() {
    try {
      ws?.close();
      connect();
      await new Promise(resolve => setTimeout(resolve, 1000)); // Wait for connection
      if (get().status !== 'connected') {
        set(state => ({ ...state, status: 'disconnected' }));
      }
    } catch (error) {
      console.error('Reconnect failed:', error);
      set(state => ({ ...state, status: 'disconnected' }));
    }
  };
  async function disconnect() {
    ws?.close();
  };
  return { 
    status: 'disconnected',
    connect,
    reconnect,
    disconnect,
  };
});