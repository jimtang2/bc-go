import { useState, useEffect } from "react";
import Table from "./table"
import { buildApiUrl } from "@/lib/utils";
import { ProfitableMatchesResponse } from "@/gen/v1/schema";

export default () => {
  const [ response, setResponse ] = useState<ProfitableMatchesResponse | null>(null);
  const [ limit ] = useState(0)
  const [ offset ] = useState(0)

  useEffect(() => {
    fetch(buildApiUrl(`/api/alpha?l=${limit}&o=${offset}`))
      .then(resp => resp.arrayBuffer())
      .then(arrayBuffer => {
        const uint8 = new Uint8Array(arrayBuffer);
        const pbMsg = ProfitableMatchesResponse.fromBinary(uint8);
        console.log(pbMsg)
        setResponse(pbMsg);
      }) 
  }, [])
  const props = {
    response,
  }
  return <Table {...props} />
} 