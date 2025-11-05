import type { ColumnDef } from "@tanstack/react-table";
import ExchangeIcon from "@/components/icon";
import { Match } from "@/gen/v1/schema";

const columnDefs: ColumnDef<Match,any>[] = [
  { 
    id: 'asset',
    header: 'Asset', 
    accessorFn: ({ pair }) => pair.split(":")[0],
    cell: info => info.getValue(),
    size: 7.5,
  },
  { 
    id: 'buyExchange',
    header: 'Buy Exchange', 
    accessorKey: 'askExchange',
    cell: (info) => <ExchangeCell name={info.getValue()} />,
    size: 10, 
  },
  { 
    id: 'sellExchange',
    header: 'Sell Exchange', 
    accessorKey: 'bidExchange',
    cell: (info) => <ExchangeCell name={info.getValue()} />,
    size: 10, 
  },
  { 
    header: 'Size', 
    accessorKey: 'calculations.volume',
    size: 10, 
    cell: (info) => info.getValue(),
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
    size: 7.5,
  },
  { 
    id: 'fees',
    header: 'Fee %', 
    accessorFn: ({ askFeeRate, bidFeeRate }) => {
      return askFeeRate + bidFeeRate
    },
    cell: info => `${info.getValue().toFixed(2)}%`,
    size: 7.5,
  },
  { 
    id: 'profitLoss',
    header: 'P/L $', 
    accessorKey: 'calculations.profitLoss',
    cell: info => {
      return <span className={info.getValue() > 0 ? "text-green-500" : "text-red-500"}>{new Intl.NumberFormat("us-US", { style: "currency", currency: "USD" }).format(info.getValue())}</span>
    },
    size: 7.5,
  },
  {
    id: 'time',
    header: 'Event Time',
    accessorFn: ({ askTime, bidTime, timestamp }) => {
      if (askTime > 0) {
        return askTime 
      } else if (bidTime > 0) {
        return bidTime
      } else {
        return timestamp
      }
    },
    cell: info => new Date(Number(info.getValue())).toLocaleString("en-US"),
    size: 10,
  },
];

function ExchangeCell({ name }: { name: string; }) {
  return (
    <div className="flex flex-row gap-2 items-center">
      <ExchangeIcon exchange={name.toLowerCase()} />
      <span className="font-normal">{name}</span>
    </div>
  );
};

export default columnDefs;
