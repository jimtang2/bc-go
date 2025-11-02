import { create } from 'zustand';
import { buildApiUrl } from "@/lib/utils";

export interface Alpha {
  id:   number;  // id
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
  response:   AlphaResponse | null;
  fetchData:  () => void;
};

export const useStore = create<State>((set) => {
  async function fetchData() {
    fetch(buildApiUrl("/api/alpha"))
      .then(resp => resp.json())
      .then(json => set(state => ({ ...state, response: json })))
  }
  
  return { 
    response: null,
    fetchData: fetchData,
  };
});