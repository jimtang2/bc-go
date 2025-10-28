// webhook$$bc-go;dashboard/client/src/components/table-pairs.tsx;grok$$
"use client"

import React, { useMemo } from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { useReactTable, getCoreRowModel, flexRender, type ColumnDef, type CellContext } from '@tanstack/react-table';
import { useDataStore, type Alpha, type Ticker } from '../DataManager';
import { exchangesMap, calcProfit } from "@/lib/utils";
import clsx from "clsx";

const PairsTable: React.FC = () => {
  const { alpha: data } = useDataStore();
  const columns: ColumnDef<string, any>[] = [
    { 
      accessorKey: 'pair', 
      header: 'Pair', 
      size: 5, 
      cell: (info: CellContext<any, string>) => 
        <span className="font-bold">
          {info.getValue()}
        </span> 
    },
    { 
      header: 'Exchanges', 
      size: 5, 
      accessorFn: ({ bid, ask }) => `${exchangesMap[bid.x]}–${exchangesMap[ask.x]}`,
      cell: (info: CellContext<any, string>) => 
        <span className="">
          {info.getValue()}
        </span> 
    },
    { 
      accessorKey: 'spread_size', 
      header: 'Volume', 
      size: 5, 
      meta: { className: 'text-right font-mono' },
      cell: (info: CellContext<any, string>) => 
        <span className="">
          {info.getValue().toFixed(2)}
        </span> 
    },
    { 
      header: 'Spread (pre-fee)', 
      size: 5, 
      accessorFn: ({ spread, spread_size }) => spread * spread_size, 
      meta: { className: 'text-right font-mono' },
      cell: (info: CellContext<any, string>) => 
        <span className="">
          {info.getValue().toFixed(2)}$
        </span> 
    },
    { 
      header: 'Profit (post-fee)',
      size: 5, 
      accessorFn: ({ bid, ask, spread_size: size }) => calcProfit(bid, ask, size),
      meta: { className: 'text-right font-mono' },
      cell: (info: CellContext<any, string>) => 
        <span className="">
          {info.getValue().toFixed(2)}$
        </span> 
    },
    { 
      accessorKey: 'spread_ratio', 
      header: '% Margin', 
      size: 5, 
      meta: { align: 'right' },
      cell: (info: CellContext<any, string>) => 
        <span className="">
          {(info.getValue()*100).toFixed(2)}%
        </span> 
    },
  ];
  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    getRowId: (row: Alpha) => row.id.toString(),
  });

  return (
    <div className="mx-2 lg:mx-3 overflow-hidden rounded-lg border">
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
            table.getRowModel().rows.map((row) => {
              const { bid, ask, spread_size: size } = row.original
              const muted = calcProfit(bid, ask, size) < 0
              return <TableRow key={row.id} className={clsx(muted && "text-muted-foreground")}>
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id} style={{ width: `${cell.column.columnDef.size}%` }}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            })
          ) : (
            <TableRow>
              <TableCell colSpan={columns.length} className="h-24 text-center">
                No data available
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
};

export default PairsTable;