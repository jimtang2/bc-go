// webhook$$bc-go;dashboard/client/src/components/table-pairs.tsx;grok$$
"use client"
import React, { useState, useEffect } from 'react';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
// import { Label } from "@/components/ui/label";
// import { Input } from "@/components/ui/input";
import { SelectLabel, Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { 
  // type CellContext,
  useReactTable, getCoreRowModel, flexRender, type ColumnDef, } from '@tanstack/react-table';
import { useDataStore, type Alpha } from '../DataManager';
import clsx from "clsx";
import exchangeIcons from "@/components/icon"

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
  const columns: ColumnDef<Alpha,any>[] = [
    { 
      header: 'Event Time',
      size: 5, 
      accessorFn: ({ at, bt }) => {
        if (at > 0) return new Date(at)
        else if (bt > 0) return new Date(bt)
        else return new Date()
      },
      cell: (info) => <span>{info.getValue().toLocaleTimeString()}</span>,
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
      size: 10, 
      cell: (info) => <span className="font-light">{info.getValue().toFixed(5)}</span>,
    },
    { 
      header: 'Buy Ex.', 
      accessorKey: 'ax',
      size: 15, 
      cell: (info) => {
        const name = info.getValue()
        return <div className="flex flex-row gap-2 items-center">
          {exchangeIcons[info.getValue().toLowerCase()]()}
          <span className="font-normal">{name}</span>
        </div>
      },
    },
    { 
      header: 'Sell Ex.', 
      accessorKey: 'bx',
      size: 15, 
      cell: (info) => {
        const name = info.getValue()
        return <div className="flex flex-row gap-2 items-center">
          {exchangeIcons[info.getValue().toLowerCase()]()}
          <span className="font-normal">{name}</span>
        </div>
      },
    },
    { 
      header: 'Buy', 
      accessorKey: 'ap',
      size: 10, 
      cell: (info) => <span className="font-light">{info.getValue().toFixed(4)}</span>,
    },
    { 
      header: 'Sell', 
      accessorKey: 'bp',
      size: 10, 
      cell: (info) => <span className="font-light">{info.getValue().toFixed(4)}</span>,
    },
    // { 
    //   header: 'Spread', 
    //   accessorKey: 's',
    //   size: 10, 
    //   cell: (info) => <span className="font-bold">{info.getValue().toFixed(2)}</span>,
    // },
    { 
      header: 'Spread %', 
      accessorKey: 'sr',
      size: 10, 
      cell: (info) => <span className="font-bold">{(info.getValue()*100).toFixed(3)}%</span>,
    },
    { 
      header: 'Fee %', 
      accessorFn: ({af, bf}) => (af+bf).toFixed(2).toString(),
      size: 10, 
      cell: (info) => <span className="">{info.getValue()}%</span>,
    },
    { 
      header: 'P/L $',
      size: 10, 
      accessorKey: 'pr',
      cell: (info) => {
        const p = info.getValue()
        return <span className={clsx(["font-bold", p > 0 ? 'text-green-500' : 'text-red-500'])}>{p.toFixed(2)}</span> 
      }        
    },
  ];
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
        <div className="flex-grow-1"></div>
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
        <TableBody>
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
export default AlphaTable