import { useState, useEffect } from "react";
import { useReactTable, getCoreRowModel, getPaginationRowModel, flexRender, type Table as TTable } from '@tanstack/react-table';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { SelectLabel, Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import { Pagination, PaginationContent, PaginationItem, PaginationLink, PaginationNext, PaginationPrevious } from "@/components/ui/pagination"
import columnDefs from "./columns"
import { Match, ProfitableMatchesResponse } from "@/gen/v1/schema";

interface ProfitableMatchesTableProps {
  response: ProfitableMatchesResponse | null;
};

export default function ProfitableMatchesTable({ response, }: ProfitableMatchesTableProps) {
  const [pagination, setPagination] = useState({
    pageIndex: 0,
    pageSize: 24,
  });
  const table = useReactTable({
    data: response?.matches || [],
    columns: columnDefs,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    onPaginationChange: setPagination,
    state: { pagination },
    getRowId: (row: Match) => row.id.toString(),
  });
  return (
    <div className="mx-2 lg:mx-3 overflow-hidden rounded-lg border">
      <div className="mx-2 my-2 flex flex-row gap-2">
        <PaginationButtons table={table} />
        <PaginationInfo response={response} />
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
        <TableBody className="">
          {table.getRowModel().rows.map((row, rowIdx) => 
            <TableRow key={row.id} className={rowIdx % 2 === 1 ? "bg-muted/30" : "bg-muted/60"}>
              {row.getVisibleCells().map((cell) => (
                <TableCell key={cell.id} style={{ width: `${cell.column.columnDef.size}%` }}>{flexRender(cell.column.columnDef.cell, cell.getContext())}</TableCell>
              ))}
            </TableRow>)}
        </TableBody>
      </Table>
    </div>
  );
};

function PaginationInfo({ response }: { response: ProfitableMatchesResponse | null; }) {
  const since = new Date(Number(response?.periodStart || 0)).toLocaleDateString()
  const profit = new Intl.NumberFormat("us-US", { style: "currency", currency: "USD" }).format(response?.periodProfit || 0)
  return (
    <div className="flex flex-row items-center gap-1 text-sm">
      <span>Current Period: {profit} (since {since})</span>
    </div>
  );
};

interface PaginationItemProps {
  index: number;
  label: string;
  onClick: () => void;
  isActive: boolean;
  isVisible: boolean;
}

export function PaginationButtons({ table }: { table: TTable<Match>; }) {
  const state = table.getState().pagination
  const pageCount = table.getPageCount();
  const maxTabsCount = 7;
  const { pageIndex } = state
  const [ pageNumbers, setPageNumbers ] = useState<PaginationItemProps[]>([])
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

export function PaginationMenu({ table }: { table: TTable<Match>; }) {
  const { pageSize: pageSizeNumber } = table.getState().pagination
  const pageSize = pageSizeNumber.toString()
  return (
    <div className="flex flex-row items-center gap-1">
      <Select value={pageSize} onValueChange={(val: string) => table.setPageSize(Number(val))}>
        <SelectTrigger className="w-18">
          <SelectValue placeholder="Display" />
        </SelectTrigger>
        <SelectContent>
          <SelectGroup>
            <SelectLabel>Display</SelectLabel>
            {[24, 48, 72, 96].map(pageSize => 
              <SelectItem key={pageSize} value={pageSize.toString()}>{pageSize}</SelectItem>
            )}
          </SelectGroup>
        </SelectContent>
      </Select>
    </div>

  )
}