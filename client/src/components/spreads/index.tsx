import { useState, useEffect, } from 'react';
import { buildApiUrl } from "@/lib/utils";
import { Match, } from "@/gen/v1/schema";
import MatchesTable from "./table";
import { MatchesChart } from "./chart";
import { ControlsOverlay, PlaybackOverlay } from "./controls";

export default function MatchesPage() {
  const [ matchesSocket, setMatchesSocket ] = useState<WebSocket | null>(null)
  const [ status, setStatus ] = useState("disconnected")
  const [ matchesData, setMatchesData ] = useState<Match[]>([])
  const [ chartData, setChartData ] = useState<Match[]>([])
  const [ showControls, setShowControls ] = useState<boolean>(false)

  function connectMatches() {
    const ws = new WebSocket(buildApiUrl("/ws/matches", { websocket: true }));
    ws.binaryType = 'arraybuffer';
    ws.onopen = () => {
      setStatus("connected")
      setMatchesSocket(ws)
    };
    ws.onmessage = (event) => {
      try {
        const match = Match.fromBinary(new Uint8Array(event.data))
        setMatchesData(prev => {
          if (!prev) {
            return matchesData
          }
          const newData = [match, ...prev]
          if (newData.length > 1200) {
            newData.pop()
          }
          return newData
        })
        setChartData(prev => {
          if (!prev) {
            return chartData
          }
          const newData = [...prev, match]
          if (newData.length > 300) {
            newData.shift()
          } else {
            do {
              newData.unshift({id:0n,pair:"",askExchange:"",askPrice:0,askSize:0,askTime:0n,askFeeRate:0,bidExchange:"",bidPrice:0,bidSize:0,bidTime:0n,bidFeeRate:0,timestamp:0n,calculations:{volume:0,priceDiff:0,priceAvg:0,spread:0,spreadPct:0,bidFee:0,askFee:0,profitLoss:0}})
            } while (newData.length < 300)
          }
          return newData
        })
      } catch (e) {
        console.error('WebSocket Match message parse error:', e);
      }
    };
    ws.onerror = () => ws?.close()
    ws.onclose = () => setStatus("disconnected");
    return () => ws?.close()
  }
  const toggleConnect = () => {
    if (status === "connected") {
      matchesSocket?.close()
    } else if (status === "disconnected") {
      connectMatches()
    }
  }
  useEffect(connectMatches, [])
  return (
    <div className="mx-2 lg:mx-3 overflow-hidden rounded-lg border">
      <MatchesChart data={chartData} toggleConnect={toggleConnect} setShowControls={setShowControls}>
        <ControlsOverlay show={showControls} toggleConnect={toggleConnect} status={status} />
        <PlaybackOverlay status={status} />
      </MatchesChart>
      <MatchesTable data={matchesData} />
    </div>
  )
} 