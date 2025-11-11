import { useState, useEffect } from "react"
import { Button } from "@/components/ui/button";
import { Tooltip, TooltipContent, TooltipTrigger, } from "@/components/ui/tooltip";
import { Pause, Play, Info } from 'lucide-react';
import clsx from "clsx"

interface ControlsOverlayProps {
	show: boolean;
	toggleConnect: () => void;
	status: string;
}

export function ControlsOverlay({ show = false, toggleConnect, status, }: ControlsOverlayProps) {
	const cls = {
		container: [
			"absolute w-full h-12.5",
			show ? "opacity-100 -translate-y-12.5" : "opacity-0 translate-y-0",
			"animate-all duration-150",
		],
	}
	const styles = {
		container: {
			backgroundColor: "#6661",
		}
	}
	return (
		<div className={clsx(cls.container)} style={styles.container}>
			<div className="mx-2 my-2 flex flex-row gap-2 ">
			  <div className="flex-grow-1 flex flex-row items-center">
			    <Button variant="ghost" size="icon" onClick={toggleConnect}>
			      {status === "connected" ? <Pause /> : <Play />}
			    </Button>
			  </div>
			  <div className="flex flex-row items-center mx-2 text-sm">
			    <Tooltip>
			      <TooltipTrigger>
			        <Info />
			      </TooltipTrigger>
			      <TooltipContent>
			        <p>Real-time tickers match within &lt; 100ms</p>
			      </TooltipContent>
			    </Tooltip>
			  </div>
			</div>			
		</div>
		)
}

interface PlaybackOverlayProps {
	status: string;
}

export function PlaybackOverlay({ status }: PlaybackOverlayProps) {
	const [ show, setShow ] = useState(false)
	const cls = {
		container: [
			"absolute w-full h-12.5",
			"flex items-center justify-center",
			"-translate-y-40",
			"pointer-events-none",
		],
		indicator: [
			"w-16 h-16 rounded-full",
			"flex items-center justify-center",
			"animate-all duration-150",
			show ? "opacity-100" : "opacity-0",
		],
	}
	const styles = {
		indicator: {
			backgroundColor: "#6663",
		}
	}
	useEffect(() => {
		setShow(true)
		setTimeout(() => {
			setShow(false)
		}, 800)
	}, [status])
	return (
		<div className={clsx(cls.container)}>
			<div className={clsx(cls.indicator)} style={styles.indicator}>
				{status === "disconnected" ? <Pause /> : <Play />}
			</div>
		</div>
		)
}