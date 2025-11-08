import type { ColumnDef } from "@tanstack/react-table";
import ExchangeIcon from "@/components/icon";
import clsx from "clsx";
import { Match } from "@/gen/v1/schema";

const columnDefs: ColumnDef<Match,any>[] = [
  {
    header: 'Event Time',
    size: 5, 
    accessorFn: ({ askTime, bidTime }) => {
      if (askTime > 0) return new Date(Number(askTime)).toLocaleTimeString()
      else if (bidTime > 0) return new Date(Number(bidTime)).toLocaleTimeString()
      else return ""
    },
    cell: (info) => <span className="font-mono">{info.getValue()}</span>,
  },
  { 
    header: 'Asset', 
    accessorKey: 'pair',
    size: 5, 
    cell: (info) => <span className="font-normal">{info.getValue().split(":")[0]}</span>,
  },
  { 
    header: 'Buy Ex.', 
    accessorKey: 'askExchange',
    size: 12.5, 
    cell: (info) => <ExchangeCell name={info.getValue()} />,
  },
  { 
    header: 'Sell Ex.', 
    accessorKey: 'bidExchange',
    size: 12.5, 
    cell: (info) => <ExchangeCell name={info.getValue()} />,
  },
  { 
    header: 'Volume', 
    accessorKey: 'calculations.volume',
    size: 7.5, 
    cell: (info) => <span className="font-light font-mono">{info.getValue()}</span>,
  },
  // { 
  //   header: 'Price Spread', 
  //   accessorKey: 'calculations.spread',
  //   size: 5, 
  //   cell: (info) => <span className="font-light font-mono">{info.getValue().toFixed(2)}</span>,
  // },
  // { 
  //   header: 'Spread %', 
  //   accessorKey: 'calculations.spreadPct',
  //   size: 7.5, 
  //   cell: (info) => <span className="font-light font-mono">{info.getValue().toFixed(3)}%</span>,
  // },
  // { 
  //   header: 'Buy Position $', 
  //   accessorFn: ({ bidPrice, calculations: { volume }}) => bidPrice * volume,
  //   size: 10, 
  //   cell: (info) => <span className="font-light font-mono">{info.getValue().toFixed(2).toLocaleString()}</span>,
  // },
  { 
    header: 'Sell Position $', 
    accessorFn: ({ askPrice, calculations}) => askPrice * (calculations?.volume || 0),
    size: 10, 
    cell: (info) => <span className="font-light font-mono">{info.getValue().toFixed(2).toLocaleString()}</span>,
  },
  { 
    header: 'P/L $', 
    accessorKey: 'calculations.profitLoss',
    size: 7.5, 
    cell: (info) => {
      const className = clsx(["font-bold", info.getValue() > 0 ? "text-green-500" : "text-red-500"])
      return <span className={className}>{info.getValue().toFixed(2)}</span>
    },
  },
  // { 
  //   header: 'Fee %', 
  //   accessorFn: ({ askFeeRate, bidFeeRate }) => (askFeeRate + bidFeeRate).toFixed(2).toString(),
  //   size: 7.5, 
  //   cell: (info) => <span className="">{info.getValue()}%</span>,
  // },
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
