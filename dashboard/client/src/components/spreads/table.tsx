import { useReactTable, getCoreRowModel, flexRender, type ColumnDef, } from '@tanstack/react-table';
import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { SelectLabel, Select, SelectContent, SelectGroup, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
// import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { Pause, Play } from 'lucide-react';
import columnDefs from "./columns";
import type { Spread } from "./state";

interface SpreadTableProps {
  data: Spread[];
  status: string;
  disconnect: () => void;
  reconnect: () => void;
  displayCount: number;
  setDisplayCount: (string) => void;
  margin: number;
  setMargin: (string) => void;
  // latency: string;
};

export default function SpreadTable({ data = [], status = "disconnected", disconnect, reconnect, displayCount, setDisplayCount, margin, setMargin, latency = "0" }: SpreadTableProps): React.FC {
  const table = useReactTable({
    data,
    columns: columnDefs,
    getCoreRowModel: getCoreRowModel(),
    getRowId: (row: Spread) => row.i.toString(),
  });
  return (
    <div className="mx-2 lg:mx-3 overflow-hidden rounded-lg border">
      <div className="mx-2 my-2 flex flex-row gap-2">
        <div className="flex-grow-1 flex flex-row items-center">
          {status === "connected" ? 
            <Button variant="ghost" size="icon" onClick={disconnect}><Pause /></Button> : 
            <Button variant="ghost" size="icon" onClick={reconnect}><Play /></Button>
          }
{/*          <Tooltip>
            <TooltipTrigger asChild>
              <span className="text-sm">{latency} ms</span>
            </TooltipTrigger>
            <TooltipContent>
              <p>Latency from last received update event time</p>
            </TooltipContent>
          </Tooltip>*/}
        </div>
        <div className="flex flex-row items-center gap-2">
           <Select value={margin} onValueChange={(val: string) => setMargin(val)}>
            <SelectTrigger className="w-36">
              <SelectValue placeholder="Margin %" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>Minimum Margin %</SelectLabel>
                <SelectItem value="0.000">No minimum</SelectItem>
                <SelectItem value="0.025">0.025%</SelectItem>
                <SelectItem value="0.050">0.05%</SelectItem>
                <SelectItem value="0.100">0.10%</SelectItem>
                <SelectItem value="0.150">0.15%</SelectItem>
              </SelectGroup>
            </SelectContent>
          </Select>
          <Select value={displayCount} onValueChange={(val: string) => setDisplayCount(val)}>
            <SelectTrigger className="w-18">
              <SelectValue placeholder="Display" />
            </SelectTrigger>
            <SelectContent>
              <SelectGroup>
                <SelectLabel>Display</SelectLabel>
                <SelectItem value="24">24</SelectItem>
                <SelectItem value="48">48</SelectItem>
                <SelectItem value="72">72</SelectItem>
                <SelectItem value="96">96</SelectItem>
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
