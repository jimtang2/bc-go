// webhook$$bc-go;dashboard/client/src/components/app-sidebar.tsx;grok$$
import React from 'react';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
} from '@/components/ui/sidebar';
import NavDashboard from './nav-dashboard';
import NavSecondary from './nav-secondary';

const AppSidebar: React.FC<{ variant?: 'sidebar' | 'floating' | 'inset' }> = ({ variant }) => {
  return (
    <Sidebar variant={variant}>
      <SidebarHeader>
        <div className="flex flex-col gap-2 p-2">
          <span className="text-lg font-semibold">Arbitrage</span>
        </div>
      </SidebarHeader>
      <SidebarContent className="overflow-x-hidden">
        <NavDashboard />
        <NavSecondary />
      </SidebarContent>
      <SidebarFooter>
        {/* <NavUser /> */}
      </SidebarFooter>
    </Sidebar>
  );
};

export default AppSidebar;