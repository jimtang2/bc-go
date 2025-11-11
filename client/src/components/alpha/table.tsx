import { useEffect, useRef } from "react";
import { useReactTable, getCoreRowModel, flexRender } from '@tanstack/react-table';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { columnDefs } from "./columns"
import { Match } from "@/gen/v1/schema";

export interface ProfitableMatchesTableProps {
  profit: string;
  profitSince: string;
  data: Match[];
  limit: number;
  offset: number;
  setOffset: (updater: (prev: number) => number) => void;
  count: number;
};

export default function ProfitableMatchesTable({ data, profit, profitSince, limit, offset, setOffset, count }: ProfitableMatchesTableProps) {
  const table = useReactTable({
    data,
    columns: columnDefs,
    getCoreRowModel: getCoreRowModel(),
    getRowId: (row: Match) => row.id.toString(),
  });
  const lastRowRef = useRef<HTMLTableRowElement>(null);
  const rows = table.getRowModel().rows;
  useEffect(() => {
    if (!lastRowRef.current || data.length === 0) return;
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          if (offset + limit < count) {
            setOffset(prev => prev + limit);
          }
        }
      },
      { root: null, threshold: 0.1 }
    );
    observer.observe(lastRowRef.current);
    return () => observer.disconnect();
  }, [data, setOffset]);
  return (
    <div className="mx-2 lg:mx-3 overflow-hidden rounded-lg border">
      <div className="mx-2 my-2 flex flex-row justify-end gap-2">
        <div className="flex flex-row items-center gap-2">
          <span className="text-sm">Alpha since {profitSince} ({new Intl.NumberFormat("us-US").format(count)} trades):</span>
          <span className="font-bold text-base">{profit}</span>
        </div>
      </div>
      <Table>
        <TableHeader className="bg-muted">
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <TableHead key={header.id} style={{ width: `${header.column.columnDef.size}%` }}>
                  <Tooltip>
                    <TooltipTrigger>{flexRender(header.column.columnDef.header, header.getContext())}</TooltipTrigger>
                    <TooltipContent>{flexRender(header.column.columnDef.meta?.tooltip, header.getContext())}</TooltipContent>
                  </Tooltip>
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody className="">
          {rows.map((row, rowIdx) => {
            const isLastRow = rowIdx === rows.length - 1;
            return (
              <TableRow
                key={row.id}
                ref={isLastRow ? lastRowRef : null}
                className={rowIdx % 2 === 1 ? "bg-muted/30" : "bg-muted/60"}
              >
                {row.getVisibleCells().map((cell) => (
                  <TableCell key={cell.id} style={{ width: `${cell.column.columnDef.size}%` }}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </TableCell>
                ))}
              </TableRow>
            );
          })}
        </TableBody>
      </Table>
    </div>
  );
};