import type { ColumnDef } from "@tanstack/react-table";
import { Match } from "@/gen/v1/schema";
import { assetCol, sellExchangeCol, buyExchangeCol, sizeCol, sellPositionCol, buyPositionCol, spreadCol, spreadPctCol, feesPctCol, profitLosstCol, idCol, } from "../spreads/columns"

export const timestampCol = (size: number): ColumnDef<Match, any> => ({
  id: 'time',
  header: 'Timestamp',
  accessorFn: ({ timestamp }) => new Date(Number(timestamp)).toLocaleString(),
  cell: info => info.getValue(),
  meta: {
    tooltip: 'Match timestamp',
  },
  size, 
})

export const columnDefs: ColumnDef<Match,any>[] = [
  idCol(5),
  timestampCol(13),
  assetCol(5),
  sellExchangeCol(10),
  buyExchangeCol(10),
  sizeCol(9),
  sellPositionCol(8),
  buyPositionCol(8),
  spreadCol(8),
  spreadPctCol(8),
  feesPctCol(8),
  profitLosstCol(8),
];
