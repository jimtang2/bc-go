import { useState, useEffect } from "react";
import { type ColumnDef } from '@tanstack/react-table';
import { Match } from "@/gen/v1/schema";
import ExchangeIcon from "@/components/icon";
import clsx from "clsx";

export const columnDefs: ColumnDef<Match,any>[] = [
  {
    id: 'asset',
    header: 'Asset', 
    accessorFn: ({ pair }) => pair.split(":")[0],
    cell: info => info.getValue(),
    size: 5,
  },
  {
    id: 'buyExchange',
    header: 'Buy Exchange', 
    accessorKey: 'askExchange',
    cell: info => (
      <div className="flex flex-row gap-2 items-center">
        <ExchangeIcon exchange={info.getValue().toLowerCase()} />
        <span className="font-normal">{info.getValue()}</span>
      </div>
    ),
    size: 10, 
  },
  { 
    id: 'sellExchange',
    header: 'Sell Exchange', 
    accessorKey: 'bidExchange',
    cell: info => (
      <div className="flex flex-row gap-2 items-center">
        <ExchangeIcon exchange={info.getValue().toLowerCase()} />
        <span className="font-normal">{info.getValue()}</span>
      </div>
    ),
    size: 10, 
  },
  { 
    id: 'fees',
    header: 'Fee %', 
    accessorFn: ({ askFeeRate, bidFeeRate }) => {
      return askFeeRate + bidFeeRate
    },
    cell: info => `${info.getValue().toFixed(2)}%`,
    size: 5, 
  },
  { 
    id: 'size',
    header: 'Size', 
    accessorKey: 'calculations.volume',
    cell: info => <span className="">{info.getValue().toFixed(3)}</span>,
    size: 10, 
  },
  { 
    id: 'buyPosition',
    header: 'Buy $', 
    accessorFn: ({ bidPrice, calculations }) => bidPrice * (calculations?.volume || 0),
    cell: info => new Intl.NumberFormat("us-US", { style: "currency", currency: "USD" }).format(info.getValue()),
    size: 10, 
  },
  { 
    id: 'sellPosition',
    header: 'Sell $', 
    accessorFn: ({ askPrice, calculations }) => askPrice * (calculations?.volume || 0),
    cell: info => new Intl.NumberFormat("us-US", { style: "currency", currency: "USD" }).format(info.getValue()),
    size: 10, 
  },
  { 
    id: 'spread',
    header: 'Spread', 
    accessorKey: 'calculations.spread',
    cell: info => new Intl.NumberFormat("us-US", { style: "currency", currency: "USD" }).format(info.getValue()),
    size: 10, 
  },
  { 
    header: 'Spread %', 
    accessorKey: 'calculations.spreadPct',
    cell: info => `${info.getValue().toFixed(2)}%`,
    size: 10, 
  },
  { 
    id: 'profitLoss',
    header: 'P/L $', 
    accessorKey: 'calculations.profitLoss',
    cell: info => {
      const className = clsx([
        "font-bold", 
        info.getValue() > 0 ? "text-green-500" : "text-red-500",
      ])
      return <span className={className}>{new Intl.NumberFormat("us-US", { style: "currency", currency: "USD" }).format(info.getValue())}</span>
    },
    size: 10, 
  },
  {
    id: 'time',
    header: 'Event Time',
    accessorFn: ({ askTime, bidTime }) => {
      if (askTime > 0) return Number(askTime)
      else if (bidTime > 0) return Number(bidTime)
      else return Date.now()
    },
    cell: info => new Date(info.getValue()).toLocaleTimeString(),
    size: 5, 
  },
  { 
    header: 'TTL',
    accessorFn: ({ askTime, bidTime }) => {
      if (askTime > 0) return Number(askTime)
      else if (bidTime > 0) return Number(bidTime)
      else return Date.now()
    },
    cell: info => {
      const date = new Date(info.getValue());
      const ttl = 1000;
      const [now, setNow] = useState(Date.now());
      const r = (now - new Date(date).getTime()) / ttl;
      const style = { maxWidth: r < 1 ? `${100 - Math.ceil(r*100)}%` : "0%" };
      useEffect(() => {
        const id = setInterval(() => {
          setNow(Date.now());
          if (Date.now() - date.getTime() >= ttl) {
            clearInterval(id);
          }; 
        }, 100);
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
    size: 5, 
  },
  { 
    id: 'id',
    header: 'ID', 
    accessorKey: 'id',
    cell: info => info.getValue(),
    size: 5, 
  },
];

