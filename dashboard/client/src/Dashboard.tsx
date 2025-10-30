import React from 'react';
import { Outlet } from 'react-router-dom';
import { Separator } from '@/components/ui/separator';
import { SidebarInset, SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar';
import AppSidebar from './components/sidebar';
const RootLayout: React.FC = () => {
  const props = {
    style: {
      '--sidebar-width': 'calc(var(--spacing) * 50)',
      '--header-height': 'calc(var(--spacing) * 12)',
    } as React.CSSProperties,
    defaultOpen: true
  }
  return (
    <SidebarProvider {...props}>
      <AppSidebar variant="inset" />
      <SidebarInset>
        <header className="flex h-[var(--header-height)] shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-[var(--header-height)]">
          <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
            <SidebarTrigger className="-ml-1 cursor-pointer" />
            <Separator orientation="vertical" className="mx-2 data-[orientation=vertical]:h-4" />
            <h1 className="text-base font-medium flex-grow-1">Arbitrage</h1>
          </div>
        </header>
        <div className="flex flex-1 flex-col">
          <div className="@container/main flex flex-1 flex-col gap-2">
            <div className="flex flex-col gap-4 py-4 md:gap-4 md:py-4">
              <Outlet />
            </div>
          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
};
export default RootLayout