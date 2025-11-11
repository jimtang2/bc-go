import { useState, useEffect } from "react";
import Table, { type ProfitableMatchesTableProps } from "./table"
import { buildApiUrl } from "@/lib/utils";
import { ProfitableMatchesResponse } from "@/gen/v1/schema";
import { Match } from "@/gen/v1/schema";

export default () => {
  const limit = 100
  const [ offset, setOffset ] = useState<number>(0)
  const [ data, setData ] = useState<Match[]>([])
  const [ profit, setProfit ] = useState<string>("")
  const [ profitSince, setProfitSince ] = useState<string>("")
  const [ count, setCount ] = useState<number>(0)

  useEffect(() => {
    fetch(buildApiUrl(`/api/alpha?limit=${limit}&offset=${offset}`))
      .then(resp => resp.arrayBuffer())
      .then(arrayBuffer => {
        const bytes = new Uint8Array(arrayBuffer);
        const { matches, count, periodStart, periodProfit } = ProfitableMatchesResponse.fromBinary(bytes);
        setProfitSince(new Date(Number(periodStart || 0)).toLocaleDateString())
        setProfit(new Intl.NumberFormat("us-US", { style: "currency", currency: "USD" }).format(periodProfit || 0))
        setData([...data, ...matches])
        setCount(count)
      }) 
  }, [offset])

  const props: ProfitableMatchesTableProps = {
    profit,
    profitSince,
    data,
    limit,
    offset,
    setOffset,
    count,
  }

  return <Table {...props} />
} 