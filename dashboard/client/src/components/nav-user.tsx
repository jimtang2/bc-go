// webhook$$bc-go;dashboard/client/src/components/nav-user.tsx;grok$$
import React from 'react';
import { SidebarMenu, SidebarMenuItem, SidebarMenuButton } from '@/components/ui/sidebar';
import { User } from 'lucide-react';

const NavUser: React.FC = () => {
  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <SidebarMenuButton>
          <User className="mr-2 h-4 w-4" />
          <span>User</span>
        </SidebarMenuButton>
      </SidebarMenuItem>
    </SidebarMenu>
  );
};

export default NavUser;