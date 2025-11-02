import React, { useState, useEffect } from 'react';
import { useStore } from './state';
import Table from "./table"

export default () => {
  const { response, fetchData } = useStore();
  useEffect(() => {
    fetchData()
  }, [])
  const props = {
    response,
  }
  return <Table {...props} />
} 