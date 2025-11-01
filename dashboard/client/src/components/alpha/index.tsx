import React, { useState, useEffect } from 'react';
import { useStore } from './state';
import Table from "./table"

export default () => {
  const { alpha, status, connect, reconnect, disconnect } = useStore();
  const [ data, setData ] = useState(alpha)
  const [ displayCount, setDisplayCount ] = useState("25")
  const [ minSpread, setMinSpread ] = useState("0.000")

  useEffect(() => {
    connect()
    return () => disconnect()
  }, [])

  useEffect(() => {
    if (alpha.length > 0 && alpha[0].sr*100>=minSpread) {
      setData([alpha[0], ...data])
    }
  }, [alpha])
  
  useEffect(() => {
    setData(alpha.filter(({ sr }) => sr*100 >= minSpread))
  }, [minSpread])

  const props = {
    data,
    status,
    disconnect,
    reconnect,
    displayCount,
    setDisplayCount,
    minSpread,
    setMinSpread,
  }

  return <Table {...props} />
} 