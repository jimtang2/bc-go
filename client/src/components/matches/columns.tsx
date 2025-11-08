import { type ColumnDef } from '@tanstack/react-table';
import { ExchangeCell, TTLCell } from "./cells";
import { Match } from "@/gen/v1/schema";
import clsx from "clsx";

export const columnDefs: ColumnDef<Match,any>[] = [
  {
    id: 'asset',
    header: 'Asset', 
    accessorFn: ({ pair }) => pair.split(":")[0],
    cell: () => <></>,
    size: 2.5,
  },
  {
    id: 'buyExchange',
    header: 'Buy Ex.', 
    accessorKey: 'askExchange',
    cell: info => <ExchangeCell name={info.getValue()} />,    
    size: 10,
  },
  { 
    id: 'sellExchange',
    header: 'Sell Ex.', 
    accessorKey: 'bidExchange',
    cell: info => <ExchangeCell name={info.getValue()} />,
    size: 10,
  },
  { 
    id: 'fees',
    header: 'Fee %', 
    accessorFn: ({ askFeeRate, bidFeeRate }) => {
      return askFeeRate + bidFeeRate
    },
    cell: info => `${info.getValue().toFixed(2)}%`,
    aggregationFn: (columnId, _, childRows) => {
      const lastRow = childRows[0]
      return lastRow.getValue(columnId)
    },
    size: 7.5, 
  },
  { 
    id: 'volume',
    header: 'Volume', 
    accessorKey: 'calculations.volume',
    cell: info => <span className="">{info.getValue().toFixed(3)}</span>,
    aggregationFn: (columnId, _, childRows) => {
      const lastRow = childRows[0]
      return lastRow.getValue(columnId)
    },
    size: 7.5, 
  },
  { 
    id: 'buyPosition',
    header: 'Buy $', 
    accessorFn: ({ bidPrice, calculations }) => bidPrice * (calculations?.volume || 0),
    cell: info => <span className="">{new Intl.NumberFormat("us-US", { style: "currency", currency: "USD" }).format(info.getValue())}</span>,
    aggregationFn: (columnId, _, childRows) => {
      const lastRow = childRows[0]
      return lastRow.getValue(columnId)
    },
    size: 7.5, 
  },
  { 
    id: 'sellPosition',
    header: 'Sell $', 
    accessorFn: ({ askPrice, calculations }) => askPrice * (calculations?.volume || 0),
    cell: info => <span className="">{new Intl.NumberFormat("us-US", { style: "currency", currency: "USD" }).format(info.getValue())}</span>,
    aggregationFn: (columnId, _, childRows) => {
      const lastRow = childRows[0]
      return lastRow.getValue(columnId)
    },
    size: 7.5, 
  },
  { 
    id: 'spread',
    header: 'Spread', 
    accessorKey: 'calculations.spread',
    cell: info => <span className="">{new Intl.NumberFormat("us-US", { style: "currency", currency: "USD" }).format(info.getValue())}</span>,
    aggregationFn: (columnId, _, childRows) => {
      const lastRow = childRows[0]
      return lastRow.getValue(columnId)
    },
    size: 7.5, 
  },
  { 
    header: 'Spread %', 
    accessorKey: 'calculations.spreadPct',
    cell: info => `${info.getValue().toFixed(2)}%`,
    aggregationFn: (columnId, _, childRows) => {
      const lastRow = childRows[0]
      return lastRow.getValue(columnId)
    },
    size: 7.5, 
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
    aggregationFn: (columnId, _, childRows) => {
      const lastRow = childRows[0]
      return lastRow.getValue(columnId)
    },
    size: 7.5, 
  },
  {
    id: 'time',
    header: 'Event Time',
    accessorFn: ({ askTime, bidTime }) => {
      if (askTime > 0) return Number(askTime)
      else if (bidTime > 0) return Number(bidTime)
      else return Date.now()
    },
    cell: info => <span className="">{new Date(info.getValue()).toLocaleTimeString()}</span>,
    aggregationFn: (columnId, _, childRows) => {
      const lastRow = childRows[0]
      return lastRow.getValue(columnId)
    },
    size: 5, 
  },
  { 
    header: 'TTL',
    accessorFn: ({ askTime, bidTime }) => {
      if (askTime > 0) return Number(askTime)
      else if (bidTime > 0) return Number(bidTime)
      else return Date.now()
    },
    cell: info => <TTLCell date={new Date(info.getValue())} />,
    aggregationFn: (columnId, _, childRows) => {
      const lastRow = childRows[0]
      return lastRow.getValue(columnId)
    },
    size: 5, 
  },
];

