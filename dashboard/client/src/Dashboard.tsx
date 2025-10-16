// webhook$$bc-go;dashboard/client/src/Dashboard.tsx;grok$$
import React from 'react';
import AppSidebar from './components/app-sidebar';
import SiteHeader from './components/site-header';
import SectionCards from './components/section-cards';
import TradesTable from './components/table-trades';
import PairsTable from './components/table-pairs';
import { SidebarInset, SidebarProvider } from '@/components/ui/sidebar';

const Dashboard: React.FC = () => {
  return (
    <SidebarProvider
      style={
        {
          '--sidebar-width': 'calc(var(--spacing) * 65)',
          '--header-height': 'calc(var(--spacing) * 12)',
        } as React.CSSProperties
      }
      defaultOpen={false}
    >
      <AppSidebar variant="inset" />
      <SidebarInset>
        <SiteHeader />
        <div className="flex flex-1 flex-col">
          <div className="@container/main flex flex-1 flex-col gap-2">
            <div className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
              {/*<SectionCards />*/}
              <PairsTable />
              {/*<TradesTable />*/}
            </div>
          </div>
        </div>
      </SidebarInset>
    </SidebarProvider>
  );
};

export default Dashboard;