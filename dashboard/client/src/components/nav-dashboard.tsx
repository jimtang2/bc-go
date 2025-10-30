// webhook$$bc-go;dashboard/client/src/components/nav-dashboard.tsx;grok$$
import React from 'react';
import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  // SidebarMenuSub,
  // SidebarMenuSubItem,
  // SidebarMenuSubButton,
} from '@/components/ui/sidebar';
import { Blend, History } from 'lucide-react';

const NavDashboard: React.FC = () => {
  return (
    <SidebarGroup>
      <SidebarMenu>
        <SidebarMenuItem>
          <SidebarMenuButton asChild>
            <a href="#" className="font-medium">
              <Blend className="mr-2 h-4 w-4" />
              <span>Alpha</span>
            </a>
          </SidebarMenuButton>
        </SidebarMenuItem>
        <SidebarMenuItem>
          <SidebarMenuButton asChild>
            <a href="#" className="font-medium">
              <History className="mr-2 h-4 w-4" />
              <span>History</span>
            </a>
          </SidebarMenuButton>
        </SidebarMenuItem>
      </SidebarMenu>
    </SidebarGroup>
  );
};

export default NavDashboard;