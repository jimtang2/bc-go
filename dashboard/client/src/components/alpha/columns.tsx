import type { ColumnDef } from "@tanstack/react-table";
import type { Alpha } from "./state";
import ExchangeIcon from "@/components/icon";
import clsx from "clsx";

const columnDefs: ColumnDef<Alpha,any>[] = [
  {
    header: 'Event Time',
    size: 5, 
    accessorFn: ({ at, bt }) => {
      if (at > 0) return new Date(at)
      else if (bt > 0) return new Date(bt)
      else return new Date()
    },
    cell: (info) => <span className="font-mono">{info.getValue().toLocaleTimeString()}</span>,
  },
  { 
    header: 'Asset', 
    accessorKey: 'p',
    size: 5, 
    cell: (info) => <span className="font-normal">{info.getValue().split(":")[0]}</span>,
  },
  { 
    header: 'Volume', 
    accessorKey: 'ss',
    size: 12, 
    cell: (info) => <span className="font-light font-mono">{info.getValue().toFixed(3)}</span>,
  },
  { 
    header: 'Buy Ex.', 
    accessorKey: 'ax',
    size: 12, 
    cell: (info) => <ExchangeCell name={info.getValue()} />,
  },
  { 
    header: 'Buy', 
    accessorKey: 'ap',
    size: 12, 
    cell: (info) => <span className="font-light font-mono">{info.getValue().toFixed(3)}</span>,
  },
  { 
    header: 'Sell Ex.', 
    accessorKey: 'bx',
    size: 12, 
    cell: (info) => <ExchangeCell name={info.getValue()} />,
  },
  { 
    header: 'Sell', 
    accessorKey: 'bp',
    size: 12, 
    cell: (info) => <span className="font-light font-mono">{info.getValue().toFixed(3)}</span>,
  },
  { 
    header: 'Fee %', 
    accessorFn: ({af, bf}) => (af+bf).toFixed(2).toString(),
    size: 7.5, 
    cell: (info) => <span className="">{info.getValue()}%</span>,
  },
  { 
    header: 'Spread %', 
    accessorFn: ({sr, af, bf}) => [sr, (af+bf).toFixed(2)],
    size: 7.5, 
    cell: (info) => {
      const [sr, f] = info.getValue()
      const className = clsx(["font-bold", sr>f ? "text-green-500" : "text-red-500"])
      const val = (sr*100).toFixed(2)
      return <span className={className}>{val}%</span>
    },
  },
  { 
    header: 'P/L $',
    size: 10, 
    accessorKey: 'pl',
    cell: (info) => {
      const val = info.getValue().toFixed(2)
      const className = clsx(["font-bold", val > 0 ? 'text-green-500' : 'text-red-500'])
      return <span className={className}>{val.toString()}</span> 
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
