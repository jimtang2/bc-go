// webhook$$bc-go;dashboard/client/src/components/site-header.tsx;grok$$
import { Separator } from '@/components/ui/separator';
import { SidebarTrigger } from '@/components/ui/sidebar';
import { useDataStore } from '../DataManager';
import { Loader2, RefreshCw } from 'lucide-react';

const RefreshButton = ({ status, onClick }: { status: 'connected' | 'disconnected' | 'retry'; onClick?: () => void }) => {
  if (status === 'connected') return null;

  if (status === 'retry') {
    return (
      <div className="flex flex-col items-center gap-2">
        <Loader2 className="animate-spin mx-auto" />
      </div>
    );
  }

  if (status === 'disconnected' && onClick) {
    return (
      <div className="flex flex-row items-center gap-1 cursor-pointer" onClick={onClick}>
        <RefreshCw className="h-4 w-4 text-accent-foreground" />
        <span className="text-accent-foreground">Refresh</span>
      </div>
    );
  }

  return null;
};

export function SiteHeader() {
  const { connectionStatus, reconnect } = useDataStore();

  return (
    <header className="flex h-[var(--header-height)] shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-[var(--header-height)]">
      <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
        <SidebarTrigger className="-ml-1 cursor-pointer" />
        <Separator orientation="vertical" className="mx-2 data-[orientation=vertical]:h-4" />
        <h1 className="text-base font-medium flex-grow-1">Dashboard</h1>
        {(connectionStatus === 'disconnected' || connectionStatus === 'retry') && (
          <RefreshButton status={connectionStatus} onClick={connectionStatus === 'disconnected' ? reconnect : undefined} />
        )}
      </div>
    </header>
  );
}

export default SiteHeader;