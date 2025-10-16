// webhook$$bc-go;dashboard/client/src/components/nav-dashboard.tsx;grok$$
import React from 'react';
import {
  SidebarGroup,
  SidebarMenu,
  SidebarMenuItem,
  SidebarMenuButton,
  SidebarMenuSub,
  SidebarMenuSubItem,
  SidebarMenuSubButton,
} from '@/components/ui/sidebar';
import { Home, BarChart2, List } from 'lucide-react';

// Sample data for navigation
const navData = {
  navMain: [
    {
      title: 'Dashboard',
      url: '#',
      items: [
        { title: 'Analytics', url: '#analytics' },
        { title: 'Pairs', url: '#pairs' },
        { title: 'Trades', url: '#trades' },
      ],
    },
  ],
};

const NavDashboard: React.FC = () => {
  return (
    <SidebarGroup>
      <SidebarMenu>
        {navData.navMain.map((item) => (
          <SidebarMenuItem key={item.title}>
            <SidebarMenuButton asChild>
              <a href={item.url} className="font-medium">
                <Home className="mr-2 h-4 w-4" />
                <span>{item.title}</span>
              </a>
            </SidebarMenuButton>
            {item.items?.length ? (
              <SidebarMenuSub className="ml-3.5 border-l-1">
                {item.items.map((subItem) => (
                  <SidebarMenuSubItem key={subItem.title}>
                    <SidebarMenuSubButton asChild>
                      <a href={subItem.url}>
                        {subItem.title === 'Analytics' && <BarChart2 className="mr-2 h-4 w-4" />}
                        {subItem.title === 'Pairs' && <List className="mr-2 h-4 w-4" />}
                        {subItem.title === 'Trades' && <List className="mr-2 h-4 w-4" />}
                        <span>{subItem.title}</span>
                      </a>
                    </SidebarMenuSubButton>
                  </SidebarMenuSubItem>
                ))}
              </SidebarMenuSub>
            ) : null}
          </SidebarMenuItem>
        ))}
      </SidebarMenu>
    </SidebarGroup>
  );
};

export default NavDashboard;