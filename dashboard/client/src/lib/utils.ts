import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"
import { type Ticker } from '../DataManager';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const feesMap = {
  "coinbase": 0.0060,
  "okx": 0.0010,
  "binance": 0.0010,
  "kraken": 0.0026,
  "bitfinex": 0.0020,
}

export const exchangesMap = {
  "okx": "OKX",
  "binance": "Binance",
  "bitfinex": "Bitfinex",
  "coinbase": "Coinbase",
  "kraken": "Kraken",
}

export function calcProfit(bid: Ticker, ask: Ticker, size: number) {
  const spread = (bid.b*size) - (ask.a*size); 
  const fee = (bid.b*size*feesMap[bid.x]) + (ask.a*size*feesMap[ask.x]);
  return spread - fee  
}
