import { useState, useEffect } from "react";
import { useReactTable, getCoreRowModel, getPaginationRowModel, flexRender, type ColumnDef, } from '@tanstack/react-table';
import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { SelectLabel, Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";

import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
  PaginationNext,
  PaginationPrevious,
} from "@/components/ui/pagination"

import { ChevronFirst, ChevronLeft, ChevronRight, ChevronLast } from 'lucide-react';
import columnDefs from "./columns"
import type { Alpha, AlphaResponse } from "./state";

interface AlphaTableProps {
  response: AlphaResponse;
};

export default function AlphaTable({ response, }: StreamTableProps): React.FC {
  const { count, start, end, total } = response?.period || {}
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: 24,
  });
  const table = useReactTable({
    data: response?.items || [],
    columns: columnDefs,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    onPaginationChange: setPagination,
    state: { pagination },
    rowCount: count,
    getRowId: (row: Alpha) => row.id.toString(),
  });
  return (
    <div className="mx-2 lg:mx-3 overflow-hidden rounded-lg border">
      <div className="mx-2 my-2 flex flex-row gap-2">
        <PaginationButtons table={table} />
        <PaginationInfo table={table} response={response} />
        <PaginationMenu table={table} />
      </div>
      <Table>
        <TableHeader className="bg-muted">
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <TableHead key={header.id} style={{ width: `${header.column.columnDef.size}%` }}>{flexRender(header.column.columnDef.header, header.getContext())}</TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody className="font-mono">
          {table.getRowModel().rows.map(row => 
            <TableRow key={row.id}>
              {row.getVisibleCells().map((cell) => (
                <TableCell key={cell.id} style={{ width: `${cell.column.columnDef.size}%` }}>{flexRender(cell.column.columnDef.cell, cell.getContext())}</TableCell>
              ))}
            </TableRow>)}
        </TableBody>
      </Table>
    </div>
  );
};

function PaginationArrows({ table }) {
  return (
    <div className="flex flex-row items-center gap-1">
      <Button variant="ghost" 
        onClick={() => table.firstPage()} 
        disabled={!table.getCanPreviousPage()}>
        <ChevronFirst className="h-5 w-5" />
      </Button>
      <Button variant="ghost" 
        onClick={() => table.previousPage()} 
        disabled={!table.getCanPreviousPage()}>
        <ChevronLeft className="h-5 w-5" />
      </Button>
      <Button variant="ghost" 
        onClick={() => table.nextPage()} 
        disabled={!table.getCanNextPage()}>
        <ChevronRight className="h-5 w-5" />
      </Button>
      <Button variant="ghost" 
        onClick={() => table.lastPage()} 
        disabled={!table.getCanNextPage()}>
        <ChevronLast className="h-5 w-5" />
      </Button>
    </div>
  );
};

function PaginationInfo({ table, response }) {
  const pageCount = table.getPageCount();
  const rowCount = table.getRowCount();
  const since = new Date(response?.period.start).toLocaleDateString()
  return (
    <div className="flex flex-row items-center gap-1 text-sm">
      <span>Current Period: {response?.period.total.toFixed(2)}$ (since {since})</span>
    </div>
  );
};

export function PaginationButtons({ table }) {
  const state = table.getState().pagination
  const pageCount = table.getPageCount();
  const rowCount = table.getRowCount();
  const maxTabsCount = 7;
  const { pageIndex, pageSize } = state
  const [ pageNumbers, setPageNumbers ] = useState([])
  useEffect(() => {
    const pages = []
    for (let i = 0; i < pageCount; i++) {
      pages.push({ 
        index: i,
        label: (i+1).toString(),
        onClick: () => table.setPageIndex(i),
        isActive: i == pageIndex,
        isVisible: i == pageIndex,
      })
    }
    const visibleCount = pageCount > maxTabsCount ? maxTabsCount : pageCount;
    let i = -1
    while (true) {
      if (pages[pageIndex+i]) {
        pages[pageIndex+i].isVisible = true;
      }
      i = i < 0 ? i * (-1) : i * (-1) -1
      if (pages.filter(({isVisible}) => isVisible).length >= visibleCount) {
        break;
      }
    }
    setPageNumbers(pages)
  }, [state])

  return (
    <div className="flex flex-row flex-grow-1 items-center gap-1">
      <Pagination>
        <PaginationContent>
          <PaginationItem>
            <PaginationPrevious isActive={table.getCanPreviousPage()}
              onClick={() => table.getCanPreviousPage() && table.previousPage()} />
          </PaginationItem>

          {pageNumbers?.filter(({ isVisible}) => isVisible).map(({index, label, onClick, isActive}) => 
              <PaginationItem key={index}>
                <PaginationLink 
                  onClick={onClick} 
                  isActive={isActive}>{label}</PaginationLink>
              </PaginationItem>
          )}
          <PaginationItem>
            <PaginationNext isActive={table.getCanNextPage()}
              onClick={() => table.getCanNextPage() && table.nextPage()} />
          </PaginationItem>
        </PaginationContent>
      </Pagination>
    </div>
  )
}

export function PaginationMenu({ table }) {
  const pageCount = table.getPageCount();
  const rowCount = table.getRowCount();
  const { pageIndex, pageSize } = table.getState()
  return (
    <div className="flex flex-row items-center gap-1">
      <Select value={table.getState().pagination.pageSize} onValueChange={(val: string) => table.setPageSize(Number(val))}>
        <SelectTrigger className="w-18">
          <SelectValue placeholder="Display" />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            <SelectLabel>Display</SelectLabel>
            {[24, 48, 72, 96].map(pageSize => 
              <SelectItem key={pageSize} value={pageSize}>{pageSize}</SelectItem>
            )}
          </SelectGroup>
        </SelectContent>
      </Select>
    </div>

  )
}