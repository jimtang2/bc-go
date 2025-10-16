// webhook$$bc-go;dashboard/client/src/components/section-charts.tsx;grok$$
import React from 'react';
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from '@/components/ui/chart';
import { AreaChart, Area, CartesianGrid, XAxis } from 'recharts';

const chartData = [
  { date: new Date().toISOString().split('T')[0], portfolio: 1000, avgMargin: 2.5 },
  { date: new Date(Date.now() - 86400000).toISOString().split('T')[0], portfolio: 1050, avgMargin: 2.7 },
  { date: new Date(Date.now() - 172800000).toISOString().split('T')[0], portfolio: 980, avgMargin: 2.3 },
  { date: new Date(Date.now() - 259200000).toISOString().split('T')[0], portfolio: 1020, avgMargin: 2.6 },
  { date: new Date(Date.now() - 345600000).toISOString().split('T')[0], portfolio: 1100, avgMargin: 2.8 },
];

interface ChartConfig {
  [key: string]: { label: string };
}

const chartConfig: ChartConfig = {
  portfolio: { label: 'Portfolio' },
  avgMargin: { label: 'Avg Margin' },
};

export const PortfolioChart: React.FC = () => {
  return (
    <ChartContainer config={chartConfig} className="aspect-auto h-[100px] w-full">
      <AreaChart data={chartData}>
        <defs>
          <linearGradient id="fillPortfolio" x1="0" y1="0" x2="0" y2="1">
            <stop offset="5%" stopColor="var(--chart-1)" stopOpacity={0.8} />
            <stop offset="95%" stopColor="var(--chart-1)" stopOpacity={0.1} />
          </linearGradient>
        </defs>
        <CartesianGrid vertical={false} />
        <XAxis
          dataKey="date"
          tickLine={false}
          axisLine={false}
          tickMargin={8}
          tickFormatter={(value) => value.split('-').slice(1).join('-')}
        />
        <ChartTooltip
          cursor={false}
          content={<ChartTooltipContent labelFormatter={(value) => value} />}
        />
        <Area dataKey="portfolio" type="natural" fill="url(#fillPortfolio)" stroke="var(--chart-1)" />
      </AreaChart>
    </ChartContainer>
  );
};

export const AverageMarginChart: React.FC = () => {
  return (
    <ChartContainer config={chartConfig} className="aspect-auto h-[100px] w-full">
      <AreaChart data={chartData}>
        <defs>
          <linearGradient id="fillAvgMargin" x1="0" y1="0" x2="0" y2="1">
            <stop offset="5%" stopColor="var(--chart-1)" stopOpacity={0.8} />
            <stop offset="95%" stopColor="var(--chart-1)" stopOpacity={0.1} />
          </linearGradient>
        </defs>
        <CartesianGrid vertical={false} />
        <XAxis
          dataKey="date"
          tickLine={false}
          axisLine={false}
          tickMargin={8}
          tickFormatter={(value) => value.split('-').slice(1).join('-')}
        />
        <ChartTooltip
          cursor={false}
          content={<ChartTooltipContent labelFormatter={(value) => value} />}
        />
        <Area dataKey="avgMargin" type="natural" fill="url(#fillAvgMargin)" stroke="var(--chart-1)" />
      </AreaChart>
    </ChartContainer>
  );
};