// webhook$$bc-go;dashboard/client/src/components/site-header.tsx;grok$$
import { Separator } from '@/components/ui/separator';
import { SidebarTrigger } from '@/components/ui/sidebar';
import { Button } from "@/components/ui/button";
import { 
  // DropdownMenuLabel, DropdownMenuPortal, DropdownMenuSeparator, DropdownMenuShortcut, DropdownMenuSub, DropdownMenuSubContent, DropdownMenuSubTrigger, 
  DropdownMenu, DropdownMenuContent, DropdownMenuGroup, DropdownMenuItem,  DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Zap, ZapOff } from 'lucide-react';
import { useDataStore } from '../DataManager';

const StatusIndicator = () => {
  const { status, reconnect, disconnect } = useDataStore();
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        {status == "connected" ? 
        <Button variant="ghost" size="icon"><Zap /></Button> : 
        <Button variant="ghost" size="icon"><ZapOff /></Button>}
      </DropdownMenuTrigger>
      <DropdownMenuContent className="w-36" align="start">
        <DropdownMenuGroup>
          {status === "disconnected" ? 
          <DropdownMenuItem onClick={reconnect}>Connect</DropdownMenuItem> : 
          <DropdownMenuItem onClick={disconnect}>Disconnect</DropdownMenuItem>}
        </DropdownMenuGroup>
      </DropdownMenuContent>
    </DropdownMenu>    
  );
};

export function SiteHeader() {
  return (
    <header className="flex h-[var(--header-height)] shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-[var(--header-height)]">
      <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
        <SidebarTrigger className="-ml-1 cursor-pointer" />
        <Separator orientation="vertical" className="mx-2 data-[orientation=vertical]:h-4" />
        <h1 className="text-base font-medium flex-grow-1">Arbitrage</h1>
        <StatusIndicator />
      </div>
    </header>
  );
}

export default SiteHeader;