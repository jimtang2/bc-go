// webhook$$bc-go;dashboard/client/src/components/table-pairs.tsx;grok$$
"use client"
import React, { useState, useEffect } from 'react';
import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
// import { Label } from "@/components/ui/label";
// import { Input } from "@/components/ui/input";
import { SelectLabel, Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { 
  // type CellContext,
  useReactTable, getCoreRowModel, flexRender, type ColumnDef, } from '@tanstack/react-table';
import { useDataStore, type Alpha } from '../DataManager';
import clsx from "clsx";
import ExchangeIcon from "@/components/icon"
import { Pause, Play } from 'lucide-react';
const columns: ColumnDef<Alpha,any>[] = [
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
    header: 'TTL',
    size: 5, 
    accessorFn: ({ at, bt }) => {
      if (at > 0) return new Date(at)
      else if (bt > 0) return new Date(bt)
      else return new Date()
    },
    cell: (info) => <TTLCell date={info.getValue()} />,
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
    cell: (info) => {
      const name = info.getValue()
      return <div className="flex flex-row gap-2 items-center">
        <ExchangeIcon exchange={info.getValue().toLowerCase()} />
        <span className="font-normal">{name}</span>
      </div>
    },
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
    cell: (info) => {
      const name = info.getValue()
      return <div className="flex flex-row gap-2 items-center">
        <ExchangeIcon exchange={info.getValue().toLowerCase()} />
        <span className="font-normal">{name}</span>
      </div>
    },
  },
  { 
    header: 'Sell', 
    accessorKey: 'bp',
    size: 12, 
    cell: (info) => <span className="font-light font-mono">{info.getValue().toFixed(3)}</span>,
  },
  // { 
  //   header: 'Spread', 
  //   accessorKey: 's',
  //   size: 10, 
  //   cell: (info) => <span className="font-bold">{info.getValue().toFixed(2)}</span>,
  // },
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
      return <span className={clsx(["font-bold", sr>f ? "text-green-500" : "text-red-500"])}>{(sr*100).toFixed(2)}%</span>
    },
  },
  { 
    header: 'P/L $',
    size: 10, 
    accessorKey: 'pl',
    cell: (info) => {
      const pl = info.getValue()
      return <span className={clsx(["font-bold", pl > 0 ? 'text-green-500' : 'text-red-500'])}>{pl.toFixed(2)}</span> 
    }        
  },
];
const AlphaTable: React.FC = () => {
  const { 
    alpha,
    minSpreadPct,
    setMinSpreadPct,
  } = useDataStore();
  const [data, setData] = useState(alpha)  
  const [displayCount, setDisplayCount] = useState(25)
  useEffect(() => {
    // add new alpha + filter
    if (alpha.length > 0 && alpha[0].sr*100 > minSpreadPct) {      
      setData([alpha[0], ...data])
    }
  }, [alpha])
  useEffect(() => {
    // filter alpha buffer 
    setData(alpha.filter(({sr}) => sr*100 > minSpreadPct))
  }, [minSpreadPct])
  const table = useReactTable({
    data: data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getRowId: (row: Alpha) => row.i.toString(),
  });
  const handleChangeMinSpreadPct = (val: string) => {
    setMinSpreadPct(parseFloat(val))
  }
  const handleChangeDisplayCount = (val: string) => {
    setDisplayCount(parseInt(val))
  }
  return (
    <div className="mx-2 lg:mx-3 overflow-hidden rounded-lg border">
      <div className="mx-2 my-2 flex flex-row gap-2">
        <StreamControls className="flex-grow-1 flex flex-row items-center" data={data.length > 0 ? data[0] : null} />
        <div className="flex flex-row items-center gap-2">
          <Select value={`${displayCount.toString()}`} onValueChange={handleChangeDisplayCount}>
            <SelectTrigger className="w-24">
              <SelectValue placeholder="Display Count" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>Display Count</SelectLabel>
                <SelectItem value="25">25</SelectItem>
                <SelectItem value="50">50</SelectItem>
                <SelectItem value="75">75</SelectItem>
                <SelectItem value="100">100</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>        

           <Select value={`${minSpreadPct.toFixed(3).toString()}`} onValueChange={handleChangeMinSpreadPct}>
            <SelectTrigger className="w-36">
              <SelectValue placeholder="Minimum Spread %" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>Minimum Spread %</SelectLabel>
                <SelectItem value="0.025">0.025%</SelectItem>
                <SelectItem value="0.050">0.05%</SelectItem>
                <SelectItem value="0.100">0.10%</SelectItem>
                <SelectItem value="0.150">0.15%</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>        
        </div>
      </div>
      <Table>
        <TableHeader className="bg-muted">
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <TableHead key={header.id} style={{ width: `${header.column.columnDef.size}%` }}>
                  {flexRender(header.column.columnDef.header, header.getContext())}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody className="font-mono">
          {table.getRowModel().rows.length ? (
            table.getRowModel().rows.map((row, i) => 
              i < displayCount ?
              <TableRow key={row.id}>
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id} style={{ width: `${cell.column.columnDef.size}%` }}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow> : null
            )
          ) : (
              <TableRow>
                <TableCell></TableCell>
              </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
};
const TTLCell = ({ date }: { date: Date; }) => {
  const ttl = 1000
  const [now, setNow] = useState(Date.now());
  useEffect(() => {
    const interval = setInterval(() => {
      if (new Date().getTime() - date.getTime() < ttl) setNow(Date.now())  
    }, 50);
    return () => clearInterval(interval);
  }, []);
  const r = (now - new Date(date).getTime()) / ttl
  const style = {
    maxWidth: r < 1 ? `${100 - Math.ceil(r*100)}%` : "0%",
    height: "100%",
  }
  return <div className="mr-4 w-full h-full flex items-center">
    <div className="mr-6 w-full min-h-2 h-full flex items-center justify-start">
      <div className="w-full h-full min-h-2 bg-gray-600 z-2" style={style}></div>
    </div>
  </div>;
}
interface StreamControlsProps {
  className?: string;
  data: Alpha | null;
};

const StreamControls = ({ className = "", data }: StreamControlsProps) => {
  const [ latency, setLatency ] = useState(0)
  const { status, reconnect, disconnect } = useDataStore();
  useEffect(() => {
    if (!data) {
      return
    }
    const { at, bt } = data
    if (at > 0) {
      setLatency(Date.now() - at)
    } else if (bt > 0) {
      setLatency(Date.now() - bt)
    }
  }, [data])
  console.log(latency)
  return <div className={className}>
    {status === "connected" ? 
    <Button variant="ghost" size="icon" onClick={disconnect}><Pause /></Button> : 
    <Button variant="ghost" size="icon" onClick={reconnect}><Play /></Button>}
    <span className="text-sm">{latency.toString().padStart(4, " ")} ms</span>
  </div>

};
export default AlphaTable