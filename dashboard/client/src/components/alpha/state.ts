import { create } from 'zustand';
import { w3cwebsocket as WebSocket } from 'websocket';
import { produce } from "immer";
import { getAPIHost } from "@/lib/utils";

export interface Alpha {
  i:   number;  // id
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
};

export interface AlphaResponse {
  items: Alpha[];
  period: {
    count: number;
    start: number;
    total: number;
  };
};

interface State {
  response:   AlphaResponse;
  fetchData:  () => void;
};

export const useStore = create<State>((set, get) => {
  async function fetchData() {
    fetch(`http://${getAPIHost()}/alpha`)
      .then(resp => resp.json())
      .then(json => set({ response: json }))
  }
  
  return { 
    response: null,
    fetchData: fetchData,
  };
});