import React, { useMemo } from "react";
import { LineChart, CartesianGrid, YAxis, Line, Tooltip, type TooltipContentProps, ResponsiveContainer, } from "recharts";
import { Match, } from "@/gen/v1/schema";

interface MatchesChartProps {
	data: Match[],
	toggleConnect: () => void;
	setShowControls: (showControls: boolean) => void;
	children?: React.ReactNode;
}

export function MatchesChart({ data = [], toggleConnect, setShowControls, children, }: MatchesChartProps) {
  const pctFmt = useMemo(() => (p: number) => p.toFixed(2), [])
  const tooltipProps = useMemo(() => ({
	  cursor: { fill: "#222" },
	  content: CustomTooltip,
  }), [])
  return (
  	<div className="relative" 
  		onMouseEnter={() => setShowControls(true)} 
  		onMouseLeave={() => setShowControls(false)}>
	    <ResponsiveContainer 
	    	width="100%" 
	    	height={240}>
	      <LineChart 
	      	data={data} 
	      	margin={{ bottom: 50 }} 
	      	responsive 
	      	onClick={toggleConnect}>
	        <CartesianGrid 
	        	stroke="#222"
	        	vertical={false} />
	        <YAxis
	          tickFormatter={pctFmt}
	          tick={{fontSize: 10}}
	          ticks={[-0.10, 0, 0.10, 0.20, 0.30]}
	          unit="%"
	          width={40}
	          angle={-45}
	          stroke="#555" />
					<Line
					  dataKey="calculations.spreadPct"
					  stroke="#356c60"
					  dot={false}
					  isAnimationActive={false} />
					<Tooltip {...tooltipProps} />
	      </LineChart>
	    </ResponsiveContainer>
  		{children}
  	</div>
  );
}

function CustomTooltip({ active, payload }: TooltipContentProps<string | number, string>) {
	if (!active || !payload || !payload[0] || !payload[0].payload) {
		return null
	}
	const { bidExchange, bidPrice, bidFeeRate, askExchange, askPrice, askFeeRate, pair, calculations: { spreadPct = 0 } = {}} = payload[0].payload
	if (!bidExchange || !bidPrice || !askExchange || !askPrice || !pair || !bidFeeRate || !askFeeRate || !spreadPct) {
		return null
	}
	const containerStyle = {
		background: 'linear-gradient(to bottom, #222, #111)',
		borderColor: '#333',
		fontSize: "small",		
	}
	const fields: Record<string, number | string> = {
		'Pair': pair,
		'Sell Exchange': bidExchange,
		'Sell Price': new Intl.NumberFormat("us-US", { style: "currency", currency: "USD" }).format(bidPrice),
		'Buy Exchange': askExchange,
		'Buy Price': new Intl.NumberFormat("us-US", { style: "currency", currency: "USD" }).format(askPrice),
		'Spread': `${spreadPct.toFixed(2)}%`,
		'Fees': `-${(bidFeeRate + askFeeRate).toFixed(2)}%`,
	}
	return <table style={containerStyle} className="m-2 py-2 px-4 rounded-md border-1">
		<tbody>
			{Object.keys(fields).map(k =>
				<tr key={k}>
					<td className="flex-grow-1">{k}:</td>
					<td className="pl-4">{fields[k]}</td>
				</tr>)}
		</tbody>
	</table>
}
