import { useState, useEffect } from "react";
import { type ColumnDef } from '@tanstack/react-table';
import { Match } from "@/gen/v1/schema";
import ExchangeIcon from "@/components/icon";
import clsx from "clsx";

import '@tanstack/react-table' //or vue, svelte, solid, qwik, etc.

declare module '@tanstack/react-table' {
  interface ColumnMeta<TData, TValue> {
    tooltip?: string
  }
}

export const assetCol = (size: number): ColumnDef<Match, any> => ({
  id: 'asset',
  header: 'Asset',
  accessorFn: ({ pair }) => pair.split(":")[0],
  cell: info => info.getValue(),
  meta: {
    tooltip: 'Trade Asset',
  },
  size, 
})

export const sellExchangeCol = (size: number): ColumnDef<Match, any> => ({ 
  id: 'sellExchange',
  header: 'Sell Exchange', 
  accessorKey: 'bidExchange',
  cell: info => (
    <div className="flex flex-row gap-2 items-center">
      <ExchangeIcon exchange={info.getValue().toLowerCase()} />
      <span className="font-normal">{info.getValue()}</span>
    </div>
  ),
  meta: {
    tooltip: 'Highest bid ticker exchange',
  },
  size, 
})

export const buyExchangeCol = (size: number): ColumnDef<Match, any> => ({
  id: 'buyExchange',
  header: 'Buy Exchange', 
  accessorKey: 'askExchange',
  cell: info => (
    <div className="flex flex-row gap-2 items-center">
      <ExchangeIcon exchange={info.getValue().toLowerCase()} />
      <span className="font-normal">{info.getValue()}</span>
    </div>
  ),
  meta: {
    tooltip: 'Lowest ask ticker exchange',
  },
  size, 
})

export const sizeCol = (size: number): ColumnDef<Match, any> => ({ 
  id: 'size',
  header: 'Size', 
  accessorKey: 'calculations.volume',
  cell: info => <span className="">{info.getValue().toFixed(3)}</span>,
  meta: {
    tooltip: 'Arbitrage size',
  },
  size, 
})

export const sellPositionCol = (size: number): ColumnDef<Match, any> => ({ 
  id: 'sellPosition',
  header: 'Sell $', 
  accessorFn: ({ bidPrice, calculations }) => bidPrice * (calculations?.volume || 0),
  cell: info => new Intl.NumberFormat("us-US", { style: "currency", currency: "USD" }).format(info.getValue()),
  meta: {
    tooltip: 'Position sell price',
  },
  size, 
})

export const buyPositionCol = (size: number): ColumnDef<Match, any> => ({ 
  id: 'buyPosition',
  header: 'Buy $', 
  accessorFn: ({ askPrice, calculations }) => askPrice * (calculations?.volume || 0),
  cell: info => new Intl.NumberFormat("us-US", { style: "currency", currency: "USD" }).format(info.getValue()),
  meta: {
    tooltip: 'Position buy price',
  },
  size, 
})

export const spreadCol = (size: number): ColumnDef<Match, any> => ({ 
  id: 'spread',
  header: 'Spread $', 
  accessorKey: 'calculations.spread',
  cell: info => {
    const className = clsx([
      info.getValue() > 0.01 && "text-green-500",
      info.getValue() < 0 && "text-red-500",
    ])
    return <span className={className}>{new Intl.NumberFormat("us-US", { style: "currency", currency: "USD" }).format(info.getValue())}</span>
  },
  meta: {
    tooltip: 'Sell Price - Buy Price',
  },
  size, 
})

export const spreadPctCol = (size: number): ColumnDef<Match, any> => ({ 
  header: 'Spread %', 
  accessorKey: 'calculations.spreadPct',
  cell: info => {
    const className = clsx([
      info.getValue() > 0.01 && "text-green-500",
      info.getValue() < 0 && "text-red-500",
    ])
    return <span className={className}>{`${info.getValue().toFixed(2)}%`}</span>
  },
  meta: {
    tooltip: 'Spread / Average Price',
  },
  size, 
})

export const feesPctCol = (size: number): ColumnDef<Match, any> => ({ 
  id: 'fees',
  header: 'Fee %', 
  accessorFn: ({ askFeeRate, bidFeeRate }) => {
    return askFeeRate + bidFeeRate
  },
  cell: info => `${info.getValue().toFixed(2)}%`,
  meta: {
    tooltip: 'Combined exchanges fees',
  },
  size, 
})

export const profitLosstCol = (size: number): ColumnDef<Match, any> => ({ 
  id: 'profitLoss',
  header: 'P/L $', 
  accessorKey: 'calculations.profitLoss',
  cell: info => {
    const className = clsx([
      info.getValue() > 0 ? "text-green-500" : "text-red-500",
    ])
    return <span className={className}>{new Intl.NumberFormat("us-US", { style: "currency", currency: "USD" }).format(info.getValue())}</span>
  },
  meta: {
    tooltip: 'Spread - Fees',
  },
  size, 
})

export const eventTimeCol = (size: number): ColumnDef<Match, any> => ({
  id: 'time',
  header: 'Event Time',
  accessorFn: ({ askTime, bidTime }) => {
    let time = askTime <= bidTime ? askTime : bidTime
    if (Number(time) === 0) {
      return "-"
    } else {
      return new Date(Number(time)).toLocaleTimeString()
    }
  },
  cell: info => info.getValue(),
  meta: {
    tooltip: 'Oldest time of two prices if available',
  },
  size, 
})

export const ttlCol = (size: number): ColumnDef<Match, any> => ({ 
  header: 'TTL',
  accessorFn: ({ askTime, bidTime }) => {
    if (askTime > 0) return Number(askTime)
    else if (bidTime > 0) return Number(bidTime)
    else return Date.now()
  },
  cell: info => {
    const date = new Date(info.getValue());
    const ttl = 500;
    const [now, setNow] = useState(Date.now());
    const r = (now - new Date(date).getTime()) / ttl;
    const style = { maxWidth: r < 1 ? `${100 - Math.ceil(r*100)}%` : "0%" };
    useEffect(() => {
      const id = setInterval(() => {
        setNow(Date.now());
        if (Date.now() - date.getTime() >= ttl) {
          clearInterval(id);
        }; 
      }, 50);
      return () => clearInterval(id);
    }, []);
    return (
      <div className="mr-4 w-full h-full flex items-center">
        <div className="mr-6 w-full min-h-2 h-full flex items-center justify-start">
          <div className="w-full h-full min-h-2 bg-gray-600 z-2" style={style}></div>
        </div>
      </div>
    )
  },
  meta: {
    tooltip: 'TTL (Time To Live) 500 ms',
  },
  size, 
})

export const idCol = (size: number): ColumnDef<Match, any> => ({ 
  id: 'id',
  header: 'ID', 
  accessorFn: row => row.id.toString(),
  cell: info => info.getValue().length <= 6 ? info.getValue() : info.getValue().substring(info.getValue().length - 6, info.getValue().length),
  meta: {
    tooltip: 'Spread ID',
  },
  size, 
})

export const columnDefs: ColumnDef<Match,any>[] = [
  assetCol(6),
  sellExchangeCol(8.5),
  buyExchangeCol(8.5),
  sizeCol(7.5),
  sellPositionCol(8.5),
  buyPositionCol(8.5),
  spreadCol(7.5),
  spreadPctCol(7.5),
  feesPctCol(7.5),
  profitLosstCol(7.5),
  eventTimeCol(7.5),
  ttlCol(7.5),
  idCol(7.5),
];

