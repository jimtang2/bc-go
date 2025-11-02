import { create } from 'zustand';
import { w3cwebsocket as WebSocket } from 'websocket';
import { produce } from "immer";
import { getAPIHost } from "@/lib/utils";

export interface Spread {
  i:   number; // offset id
  p:   string;  // pair
  s:   number;  // spread
  v:   number;  // volume
  ax:  string;  // ask exchange
  bx:  string;  // bid exchange
  ap:  number;  // ask price
  bp:  number;  // bid price
  as:  number;  // ask size
  bs:  number;  // bid size
  at:  number;  // ask time
  bt:  number;  // bid time
  af:  number;  // ask fee
  bf:  number;  // bid fee
  pl:  number;  // p/l
  ts:  number;  // timestamp
  m:   number;  // margin, computed client side onmessage
};

interface State {
  spreads:     Spread[];
  bufferSize:  number;
  status:      'connected' | 'disconnected';
  connect:     () => void;
  reconnect:   () => void;
  disconnect:  () => void;
};

export const useStore = create<State>((set, get) => {
  let ws: WebSocket | null = null;
  function connect() {
    if (ws?.readyState === 0 || ws?.readyState === 1) {
      return
    }
    ws = new WebSocket(`ws://${getAPIHost()}/ws`);
    ws.onopen = () => {
      set(state => ({ ...state, status: 'connected' }));
    };
    ws.onmessage = (event) => {
      try {
        const message = JSON.parse(event.data.toString())
        set(produce((draft: State) => {
          draft.spreads.unshift(message);
          if (draft.spreads.length > get().bufferSize) {
            draft.spreads.pop()
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
      await new Promise(resolve => setTimeout(resolve, 1000));
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
    set(state => ({ ...state, status: 'disconnected' }));
  };
  return { 
    spreads: [],
    bufferSize: 1000,
    status: 'disconnected',
    connect,
    reconnect,
    disconnect,
  };
});