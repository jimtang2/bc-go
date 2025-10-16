// webhook$$bc-go;dashboard/client/src/components/table-trades.tsx;grok$$
"use client"

import React from 'react';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { useReactTable, getCoreRowModel, flexRender, type ColumnDef } from '@tanstack/react-table';
import { useDataStore, type Trade } from '../DataManager';

const TradesTable: React.FC = () => {
  const { trades } = useDataStore();

  const columns: ColumnDef<Trade>[] = [
    { 
      accessorKey: 'p', 
      header: 'Pair', 
      size: 25, 
      accessorFn: (row) => row.p,
   },
    { 
      accessorKey: 'x', 
      header: 'Exchanges', 
      size: 25, 
      accessorFn: (row) => row.x.join(', '), 
    },
    { 
      accessorKey: 'q', 
      header: 'Time', 
      size: 20, 
      accessorFn: (row) => new Intl.DateTimeFormat('en-US', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit', timeZoneName: 'short' }).format(new Date(row.q)), 
    },
    { 
      accessorKey: 'e', 
      header: 'Executed', 
      size: 20, 
      accessorFn: (row) => new Intl.DateTimeFormat('en-US', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', second: '2-digit', timeZoneName: 'short' }).format(new Date(row.e)), 
    },
    { 
      accessorKey: 'm', 
      header: 'Margin', 
      size: 10, 
      accessorFn: (row) => row.m, 
      cell: (info) => {
        const margin = (info.getValue() as number).toFixed(3);
        const isPositive = parseFloat(margin) > 0;
        return <span className={isPositive ? 'text-accent-foreground' : ''}>{margin}%</span>;
      }, 
    },
  ];

  const table = useReactTable({
    data: trades,
    columns,
    getCoreRowModel: getCoreRowModel(),
  });

  return (
    <div className="mx-4 lg:mx-6 overflow-hidden rounded-lg border">
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
            table.getRowModel().rows.map((row) => (
              <TableRow key={row.id}>
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id} style={{ width: `${cell.column.columnDef.size}%` }}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
          ) : (
            <TableRow>
              <TableCell colSpan={5} className="h-24 text-center">
                No trade data available
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
};

export default TradesTable;