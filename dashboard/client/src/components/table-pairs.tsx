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
import { useDataStore, type Quote } from '../DataManager';

const PairsTable: React.FC = () => {
  const { pairs: data, exchanges } = useDataStore();
  const columns: ColumnDef<string, any>[] = useMemo(() => ([
    { 
      accessorKey: 'p', 
      header: 'Pair', 
      size: 15, 
      accessorFn: (row: string) => row, 
      cell: (info: CellContext<any, string>) => <span className="font-bold">{info.getValue().toUpperCase()}</span> 
    },
    ...exchanges.map(x => ({
      accessorFn: (row: string) => row + "." + x,
      header: x.charAt(0).toUpperCase() + x.slice(1),
      size: (85 / (exchanges.length + 1)),
      cell: (info: CellContext<any, string>) => info ? <QuoteCell info={info} /> : null,
    }))
  ]), [exchanges]);
  const table = useReactTable({
    data,
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
              <TableRow key={row.original}>
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id} style={{ width: `${cell.column.columnDef.size}%` }}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            ))
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

const QuoteCell = ({ info }: { info: CellContext<any, any> }) => {
  const k = info.getValue()
  const q: Quote = useDataStore(({ quotes }) => quotes[k])
  const text = q?.p | -1

  return (
    <span>{text}</span>
  );
};

export default PairsTable;