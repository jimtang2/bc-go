import { useReactTable, getCoreRowModel, flexRender, type ColumnDef, } from '@tanstack/react-table';
import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { SelectLabel, Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
import columnDefs from "./columns"

interface AlphaTableProps {
  data: Alpha[];
  status: string;
  disconnect: () => void;
  reconnect: () => void;
  displayCount: number;
  setDisplayCount: (string) => void;
  minSpread: number;
  setMinSpread: (string) => void;
};

export default function AlphaTable({ data = [], status = "disconnected", disconnect, reconnect, displayCount, setDisplayCount, minSpread, setMinSpread }: StreamTableProps): React.FC {
  const table = useReactTable({
    data,
    columns: columnDefs,
    getCoreRowModel: getCoreRowModel(),
    getRowId: (row: Alpha) => row.i.toString(),
  });
  return (
    <div className="mx-2 lg:mx-3 overflow-hidden rounded-lg border">
      <div className="mx-2 my-2 flex flex-row gap-2">
        <div className="flex-grow-1"></div>
        <div className="flex flex-row items-center gap-2">          
          <Select value={displayCount} onValueChange={(val: string) => setDisplayCount(val)}>
            <SelectTrigger className="w-24">
              <SelectValue placeholder="Display" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>Display</SelectLabel>
                <SelectItem value="25">25</SelectItem>
                <SelectItem value="50">50</SelectItem>
                <SelectItem value="75">75</SelectItem>
                <SelectItem value="100">100</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
           <Select value={minSpread} onValueChange={(val: string) => setMinSpread(val)}>
            <SelectTrigger className="w-36">
              <SelectValue placeholder="Minimum Spread %" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>Minimum Spread %</SelectLabel>
                <SelectItem value="0.000">No minimum</SelectItem>
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
                <TableHead key={header.id} style={{ width: `${header.column.columnDef.size}%` }}>{flexRender(header.column.columnDef.header, header.getContext())}</TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody className="font-mono">
          {table.getRowModel().rows.filter((row, i) => i < parseInt(displayCount)).map(row => 
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
