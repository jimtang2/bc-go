import { useState, useEffect } from 'react';
import Table from "./table"
import type { SpreadsTableProps } from "./table"
import { buildApiUrl } from "@/lib/utils"
import { Match } from "@/gen/v1/schema";

export default function MatchesPage() {
  const [ socket, setSocket ] = useState<WebSocket | null>(null)
  const [ status, setStatus ] = useState("disconnected")
  const [ data, setData ] = useState<Match[]>([])
  function connect() {
    const ws = new WebSocket(buildApiUrl("/ws", { websocket: true }));
    ws.binaryType = 'arraybuffer';
    setStatus("connecting");
    ws.onopen = () => {
      setStatus("connected")
      setSocket(ws)
    };
    ws.onmessage = (event) => {
      try {
        const match = Match.fromBinary(new Uint8Array(event.data))
        setData(prev => {
          if (!prev) {
            return data
          }
          const newData = [match, ...prev]
          if (newData.length > 5000) {
            newData.pop()
          }
          // console.log(newData.length)
          return newData
        })
      } catch (e) {
        console.error('WebSocket message parse error:', e);
      }
    };
    ws.onerror = () => {
      ws?.close()
    }
    ws.onclose = () => {
      setStatus("disconnected");
    };
    return () => ws.close()
  }
  useEffect(connect, [])
  const props: SpreadsTableProps = {
    data,
    status,
    disconnect: () => socket?.close(),
    reconnect: () => {socket?.close();connect();},
  }
  return <Table {...props} />
} 