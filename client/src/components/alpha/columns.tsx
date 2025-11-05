import type { ColumnDef } from "@tanstack/react-table";
import type { Alpha } from "./state";
import ExchangeIcon from "@/components/icon";
import clsx from "clsx";

const columnDefs: ColumnDef<Alpha,any>[] = [
  {
    header: 'Event Time',
    size: 15, 
    accessorFn: ({ at, bt }) => {
      if (at > 0) return new Date(at)
      else if (bt > 0) return new Date(bt)
      else return new Date()
    },
    cell: (info) => <span className="font-mono">{info.getValue().toLocaleString()}</span>,
  },
  { 
    header: 'Asset', 
    accessorKey: 'p',
    size: 10, 
    cell: (info) => <span className="font-normal">{info.getValue()}</span>,
  },
  { 
    header: 'Volume', 
    accessorKey: 'v',
    size: 10, 
    cell: (info) => <span className="font-light font-mono">{info.getValue().toFixed(3)}</span>,
  },
  { 
    header: 'Buy Ex.', 
    accessorKey: 'ax',
    size: 10, 
    cell: (info) => <ExchangeCell name={info.getValue()} />,
  },
  { 
    header: 'Buy', 
    accessorKey: 'ap',
    size: 10, 
    cell: (info) => <span className="font-light font-mono">{info.getValue().toFixed(3)}</span>,
  },
  { 
    header: 'Sell Ex.', 
    accessorKey: 'bx',
    size: 10, 
    cell: (info) => <ExchangeCell name={info.getValue()} />,
  },
  { 
    header: 'Sell', 
    accessorKey: 'bp',
    size: 10, 
    cell: (info) => <span className="font-light font-mono">{info.getValue().toFixed(3)}</span>,
  },
  { 
    header: 'Fee %', 
    accessorFn: ({af, bf}) => (af+bf).toFixed(2).toString(),
    size: 10, 
    cell: (info) => <span className="">{info.getValue()}%</span>,
  },
  { 
    header: 'P/L %', 
    size: 10, 
    accessorFn: ({s, bp, ap, af, bf}) => {
      const v1 = s / (bp + ap) * 2 * 100 // spread pct 
      const v2 = af + bf // fees
      return (v1 - v2).toFixed(2)
    },
    cell: (info) => <span className={"font-bold text-green-500"}>{info.getValue()}%</span>,
  },
  { 
    header: 'P/L $',
    size: 10, 
    accessorFn: ({pl}) => pl.toFixed(2),
    cell: (info) => {
      const val = info.getValue().toString()
      const className = [
        info.getValue() > 0 && 'font-bold text-green-500'
      ]
      return <span className={clsx(className)}>{val}$</span>
    } 
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
