import React, { useState, useEffect } from 'react';
import { useStore, type Spread } from './state';
import Table from "./table"

export default () => {
  const { spreads, status, connect, reconnect, disconnect } = useStore();
  const [ data, setData ] = useState(spreads)
  const [ displayCount, setDisplayCount ] = useState("25")
  const [ margin, setMargin ] = useState("0.000")
  const [ latency, setLatency ] = useState("0")

  useEffect(() => {
    connect()
    return () => disconnect()
  }, [])

  useEffect(() => {
    if (spreads.length > 0 && spreads[0].m*100>=margin) {
      setData([spreads[0], ...data])
    } 
  }, [spreads])

  useEffect(() => {
    if (data?.length > 0) {
      if (data[0].at > 0) {
        setLatency((Date.now() - data[0].at).toString());
      } else if (data[0].bt > 0) {
        setLatency((Date.now() - data[0].bt).toString());
      };
    };
  }, [data]);
  
  useEffect(() => {
    setData(spreads.filter(({ m }) => m*100 >= margin))
  }, [margin])

  // when window is throttled the data gets backed up and the latency (duration between last received message timestamp and current time) increases; this fix disconnects the websocket connection once the latency reaches 20+ seconds
  useEffect(() => {
    if (latency > 20000) {
      disconnect();
    };
  }, [latency]);

  const props = {
    data,
    status,
    disconnect,
    reconnect,
    displayCount,
    setDisplayCount,
    margin,
    setMargin,
    latency,
  }

  return <Table {...props} />
} 