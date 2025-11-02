import React, { useState, useEffect } from 'react';
import { useStore, type Spread } from './state';
import { useDebounce } from "use-debounce";
import Table from "./table"

export default () => {
  const { spreads, status, connect, reconnect, disconnect } = useStore();
  const [ data, setData ] = useState(spreads)
  const [ displayCount, setDisplayCount ] = useState("24")
  const [ margin, setMargin ] = useState("0.000")
  const [ latency, setLatency ] = useState(0)
  const [ throttledData, setThrottledData ] = useDebounce(data, 30)

  useEffect(() => {
    connect()
    return () => disconnect()
  }, [])

  useEffect(() => {
    if (spreads.length > 0) {
      const {s, ap, bp} = spreads[0]
      if (s/(ap+bp)*2*100>=margin) {
        setData([spreads[0], ...data])  
      }
    } 
  }, [spreads])
  
  useEffect(() => {
    setData(spreads.filter(({ s, ap, bp }) => s/(ap+bp)*2*100 >= margin))
  }, [margin])

  // when window is throttled the data gets backed up and the latency (duration between last received message timestamp and current time) increases; this fix disconnects the websocket connection once the latency reaches 20+ seconds
  useEffect(() => {
    if (data?.length == 0) {
      return
    };
    if (data[0].at > 0) {
      setLatency(Date.now() - data[0].at);
    } else if (data[0].bt > 0) {
      setLatency(Date.now() - data[0].bt);
    };
  }, [data]);

  useEffect(() => {
    if (latency > 20000) {
      disconnect();
    };
  }, [latency]);

  const props = {
    data: throttledData,
    status,
    disconnect,
    reconnect,
    displayCount,
    setDisplayCount,
    margin,
    setMargin,
    // latency,
  }

  return <Table {...props} />
} 