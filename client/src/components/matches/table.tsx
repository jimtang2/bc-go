import React, { Fragment, useState } from "react"
import { useReactTable, getCoreRowModel, getExpandedRowModel, getGroupedRowModel, getSortedRowModel, flexRender, type GroupingState, type SortingState } from '@tanstack/react-table';
import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Pause, Play } from 'lucide-react';
import { Match } from "@/gen/v1/schema";
import { columnDefs } from "./columns";

export interface MatchTableProps {
  data: Match[];
  status: string;
  disconnect: () => void;
  reconnect: () => void;
};

export default function MatchTable({ data = [], status = "disconnected", disconnect, reconnect, }: MatchTableProps) {
  const [grouping, setGrouping] = React.useState<GroupingState>([
    "asset", 
    "buyExchange",
    "sellExchange",
  ])
  const [sorting, setSorting] = useState<SortingState>([
    {id: "asset", desc: false},
    {id: "buyExchange", desc: false},
    {id: "sellExchange", desc: false},
    {id: "time", desc: true},
  ]) 
  const table = useReactTable({
    data, 
    columns: columnDefs,
    getCoreRowModel: getCoreRowModel(),
    getGroupedRowModel: getGroupedRowModel(),
    getExpandedRowModel: getExpandedRowModel(),
    onGroupingChange: setGrouping,
    getSortedRowModel: getSortedRowModel(),
    onSortingChange: setSorting,
    state: {
      grouping,
      sorting,
    },
    getRowId: ({ pair, askExchange, bidExchange }: Match) => `${pair}.${askExchange}.${bidExchange}`,
    enableRowPinning: true,
  });
  return (
    <div className="mx-2 lg:mx-3 overflow-hidden rounded-lg border">
      <div className="mx-2 my-2 flex flex-row gap-2">
        <div className="flex-grow-1 flex flex-row items-center">
          {status === "connected" && 
          <Button variant="ghost" size="icon" onClick={disconnect}><Pause /></Button>}
          {status === "disconnected" && 
          <Button variant="ghost" size="icon" onClick={reconnect}><Play /></Button>}
        </div>
      </div>
      <Table>
        <TableHeader className="bg-muted">
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => {
                return (
                  <TableHead key={header.id} 
                    colSpan={header.colSpan} 
                    style={{ width: `${header.column.columnDef.size}%` }}>
                    {flexRender(header.column.columnDef.header, header.getContext())}
                  </TableHead>
                  )
              })}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody className="">
          {table.getRowModel().rows.map(assetRow => (
            <Fragment key={assetRow.getValue('asset')}>
              <TableRow className="bg-muted/30">
                <TableCell colSpan={assetRow.getAllCells().length}>{assetRow.getValue('asset')}</TableCell>
              </TableRow>
              {assetRow.subRows.map(buyRow => (
                <Fragment key={`${assetRow.getValue('asset')}.${buyRow.getValue('buyExchange')}`}>
                  {buyRow.subRows.map(sellRow => (
                    <Fragment key={`${assetRow.getValue('asset')}.${buyRow.getValue('buyExchange')}.${sellRow.getValue('sellExchange')}`}>
                      <TableRow>
                        {sellRow.getAllCells().map(cell => (
                          <TableCell key={`${assetRow.getValue('asset')}.${buyRow.getValue('buyExchange')}.${sellRow.getValue('sellExchange')}.${cell.column.id}.${sellRow.getValue('time')}`}>
                            {flexRender(cell.column.columnDef.cell, cell.getContext())}
                          </TableCell>
                        ))}
                      </TableRow>
                    </Fragment>
                  ))}
                </Fragment>
              ))}              
            </Fragment>)
          )}
        </TableBody>
      </Table>
    </div>
  );
};
