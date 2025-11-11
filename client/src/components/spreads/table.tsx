import React, { useState } from "react"
import { useReactTable, getCoreRowModel, getGroupedRowModel, getSortedRowModel, flexRender, } from '@tanstack/react-table';
import type { GroupingState, SortingState } from '@tanstack/react-table';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Tooltip, TooltipContent, TooltipTrigger, } from "@/components/ui/tooltip";
import { Match } from "@/gen/v1/schema";
import { columnDefs } from "./columns";

export default function SpreadsTable({ data = [], }: { data: Match[]; }) {
  const [grouping, setGrouping] = React.useState<GroupingState>([
    "asset", 
    "sellExchange",
    "buyExchange",
  ])
  const [sorting, setSorting] = useState<SortingState>([
    {id: "asset", desc: false},
    {id: "sellExchange", desc: false},
    {id: "buyExchange", desc: false},
    {id: "time", desc: true},
  ]) 
  const table = useReactTable({
    data, 
    columns: columnDefs,
    getCoreRowModel: getCoreRowModel(),
    getGroupedRowModel: getGroupedRowModel(),
    onGroupingChange: setGrouping,
    getSortedRowModel: getSortedRowModel(),
    onSortingChange: setSorting,
    state: {
      grouping,
      sorting,
    },
    getRowId: ({ pair, askExchange, bidExchange }: Match) => `${pair}.${askExchange}.${bidExchange}`,
  });
  
  return (
    <Table>
      <TableHeader className="bg-muted">
        {table.getHeaderGroups().map((headerGroup) => (
          <TableRow key={headerGroup.id}>
            {headerGroup.headers.map((header) => {
              return (
                <TableHead key={header.id} 
                  colSpan={header.colSpan} 
                  style={{ width: `${header.column.columnDef.size}%` }}>
                  <Tooltip>
                    <TooltipTrigger>{flexRender(header.column.columnDef.header, header.getContext())}</TooltipTrigger>
                    <TooltipContent>{flexRender(header.column.columnDef.meta?.tooltip, header.getContext())}</TooltipContent>
                  </Tooltip>
                </TableHead>
                )
            })}
          </TableRow>
        ))}
      </TableHeader>
      <TableBody className="">
        {table.getRowModel().rows.map((assetRow, assetRowIdx) => 
          assetRow.subRows.map(buyRow => 
            buyRow.subRows.map(sellRow => 
              sellRow.subRows.map((spreadRow, spreadRowIdx) => {
                const asset = assetRow.getValue('asset')
                const buyExchange = buyRow.getValue('buyExchange')
                const sellExchange = sellRow.getValue('sellExchange')
                if (spreadRowIdx >= 1) {
                  return null
                }
                const keys = [asset, buyExchange, sellExchange, spreadRowIdx, spreadRow.getValue('id')]
                return (
                  <TableRow key={keys.join(".")} className={assetRowIdx % 2 === 1 ? "bg-muted/30" : "bg-muted/60"}>
                    {spreadRow.getAllCells().map((cell, cellIdx) => {
                      if (cellIdx < 4) {
                        return (
                          <TableCell key={[...keys, cellIdx].join(".")}>
                            {spreadRowIdx === 0 && flexRender(cell.column.columnDef.cell, cell.getContext())}
                          </TableCell>
                        )
                      } else {
                        return (
                          <TableCell key={[...keys, cellIdx].join(".")}>
                            {flexRender(cell.column.columnDef.cell, cell.getContext())}
                          </TableCell>
                        )
                      }
                    })}
                  </TableRow>
                )
              })
            )
          )
        )}
      </TableBody>
    </Table>
  );
};
