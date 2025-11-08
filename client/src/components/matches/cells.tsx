import { useState, useEffect } from "react";
import ExchangeIcon from "@/components/icon";

export function TTLCell({ date }: { date: Date; }) {
  const ttl = 1000;
  const [now, setNow] = useState(Date.now());
  const r = (now - new Date(date).getTime()) / ttl;
  const style = { maxWidth: r < 1 ? `${100 - Math.ceil(r*100)}%` : "0%" };
  useEffect(() => {
    const id = setInterval(() => {
      setNow(Date.now());
      if (Date.now() - date.getTime() >= ttl) {
        clearInterval(id);
      }; 
    }, 50);
    return () => clearInterval(id);
  }, []);
  return (
    <div className="mr-4 w-full h-full flex items-center">
      <div className="mr-6 w-full min-h-2 h-full flex items-center justify-start">
        <div className="w-full h-full min-h-2 bg-gray-600 z-2" style={style}></div>
      </div>
    </div>
  );
};

export function ExchangeCell({ name }: { name: string; }) {
  return (
    <div className="flex flex-row gap-2 items-center">
      <ExchangeIcon exchange={name?.toLowerCase()} />
      <span className="font-normal">{name}</span>
    </div>
  );
};
